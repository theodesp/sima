package sima

import (
	"context"
	"fmt"
	"github.com/OneOfOne/cmap"
	"github.com/deckarep/golang-set"
	"runtime"
	"sync/atomic"
	"unsafe"
)

const ANY_ID = iota

type simaInternal struct {
	receivers    *cmap.CMap
	byReceiver   *cmap.CMap
	bySender     *cmap.CMap
	topicFactory *TopicFactory
	closed       uint64
}

type Sima struct {
	*simaInternal
	pad [128 - unsafe.Sizeof(simaInternal{})%128]byte
}

func NewSima(topicFactory *TopicFactory) *Sima {
	s := &Sima{
		simaInternal: &simaInternal{
			receivers:    cmap.New(),
			byReceiver:   cmap.New(),
			bySender:     cmap.New(),
			topicFactory: topicFactory,
			closed:       uint64(0),
		},
	}
	runtime.SetFinalizer(s, (*Sima).clearState)
	return s
}

// Connects *receiver* of signal events sent by *sender*
func (s *Sima) Connect(receiver ReceiverType, senderName string) ReceiverType {
	receiverKey := HashValue(receiver)
	_, senderKey := getSenderKeyValue(senderName, s)

	s.receivers.Set(receiverKey, receiver)
	if v, ok := s.bySender.GetOK(senderKey); ok {
		v.(mapset.Set).Add(receiverKey)
	} else {
		s.bySender.Set(senderKey, mapset.NewSet())
		s.bySender.Get(senderKey).(mapset.Set).Add(receiverKey)
	}
	if v, ok := s.byReceiver.GetOK(receiverKey); ok {
		v.(mapset.Set).Add(senderKey)
	} else {
		s.bySender.Set(receiverKey, mapset.NewSet())
		s.bySender.Get(receiverKey).(mapset.Set).Add(senderKey)
	}

	return receiver
}

// True if there is a receiver for *sender* at the time of the call
func (s *Sima) HasReceiversFor(senderName string) bool {
	if s.receivers.Len() == 0 {
		return false
	}

	if s.bySender.Has(ANY_ID) == true {
		return true
	}

	_, senderKey := getSenderKeyValue(senderName, s)
	if v, ok := s.bySender.GetOK(senderKey); ok {
		return len(v.(mapset.Set).ToSlice()) > 0
	} else {
		return false
	}
}

// Emit this signal on behalf of *sender*, passing on Context.
// Returns a list of senders that accepted the signal
func (s *Sima) GetReceiversFor(senderName string) []ReceiverType {
	if s.receivers.Len() == 0 {
		return []ReceiverType{}
	}

	_, senderKey := getSenderKeyValue(senderName, s)
	var ids []interface{}

	if s.bySender.Has(senderKey) {
		ids = s.bySender.Get(senderKey).(mapset.Set).ToSlice()
	} else {
		ids = []interface{}{}
	}

	var result []ReceiverType
	for _, receiverId := range ids {
		receiver, ok := s.receivers.GetOK(receiverId)
		if !ok {
			continue
		}
		result = append(result, receiver.(ReceiverType))
	}

	return result
}

// Disconnect *receiver* from this signal's events.
// Returns true if successful or false otherwise
func (s *Sima) Disconnect(receiver ReceiverType, senderName string) bool {
	if s.receivers.Len() == 0 {
		return false
	}

	_, senderKey := getSenderKeyValue(senderName, s)
	receiverKey := HashValue(receiver)

	if senderName == "" {
		return s.disconnectAllByReceiver(senderKey, receiverKey)
	} else  {
		return s.disconnect(senderKey, receiverKey)
	}
}

func (s *Sima) disconnectAllByReceiver(senderKey uint64, receiverKey uint64) bool  {
	var isMissingKey bool
	if v := s.byReceiver.DeleteAndGet(receiverKey); v != nil  {
		v.(mapset.Set).Clear()

		s.bySender.ForEach(func(key interface{}, v interface{}) bool {
			v.(mapset.Set).Remove(receiverKey)
			return true
		})
	} else {
		isMissingKey = true
	}

	if !isMissingKey {
		s.receivers.Delete(receiverKey)
		return true
	} else {
		return false
	}
}

func (s *Sima) disconnect(senderKey uint64, receiverKey uint64) bool  {
	var isMissingKey bool
	if v, ok := s.bySender.GetOK(senderKey); ok {
		v.(mapset.Set).Remove(receiverKey)
	} else {
		isMissingKey = true
	}

	if v, ok := s.byReceiver.GetOK(receiverKey); ok && !isMissingKey {
		v.(mapset.Set).Remove(senderKey)
		return true
	} else {
		return false
	}
}

// Emit this signal on behalf of *sender*, passing on Context.
// Returns a list of results from the receivers
func (s *Sima) Dispatch(context context.Context, senderName string) []interface{} {
	if s.receivers.Len() == 0 {
		return []interface{}{}
	}

	var result []interface{}
	sender, _ := getSenderKeyValue(senderName, s)
	for _, receiver := range s.GetReceiversFor(senderName) {
		result = append(result, receiver(context, sender))
	}

	return result
}

// Clean up remaining state
func (p *Sima) clearState() {
	if p != nil {
		for _, key := range p.byReceiver.Keys() {
			p.byReceiver.DeleteAndGet(key).(mapset.Set).Clear()
		}
		for _, key := range p.bySender.Keys() {
			p.bySender.DeleteAndGet(key).(mapset.Set).Clear()
		}
		for _, key := range p.receivers.Keys() {
			p.receivers.Delete(key)
		}
		// Set closed to true atomically
		atomic.StoreUint64(&(p.closed), uint64(1))
		fmt.Println("state cleared")
	}
}

func getSenderKeyValue(senderName string, s *Sima) (*Topic, uint64) {
	var key uint64
	var sender *Topic

	if senderName == "" {
		sender = s.topicFactory.GetByName(ANY)
		key = ANY_ID
	} else {
		sender = s.topicFactory.GetByName(senderName)
		key = HashValue(sender)
	}

	return sender, key
}
