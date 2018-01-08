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

// Connects *receiver* to signal events sent by *sender*
func (s *Sima) Connect(receiver ReceiverType, sender *Topic) ReceiverType {
	if sender == nil {
		sender = s.topicFactory.GetByName(ANY)
	}

	receiverId := HashValue(receiver)
	var senderId uint64

	if sender == s.topicFactory.GetByName(ANY) {
		senderId = ANY_ID
	} else {
		senderId = HashValue(sender)
	}

	s.receivers.Set(receiverId, receiver)
	if v, ok := s.bySender.GetOK(senderId); ok {
		v.(mapset.Set).Add(receiverId)
	} else {
		s.bySender.Set(senderId, mapset.NewSet())
		s.bySender.Get(senderId).(mapset.Set).Add(receiverId)
	}
	if v, ok := s.byReceiver.GetOK(receiverId); ok {
		v.(mapset.Set).Add(senderId)
	} else {
		s.bySender.Set(receiverId, mapset.NewSet())
		s.bySender.Get(receiverId).(mapset.Set).Add(senderId)
	}

	return receiver
}

// True if there is a receiver for *sender* at the time of the call
func (s *Sima) HasReceiversFor(sender *Topic) bool {
	if s.receivers.Len() == 0 {
		return false
	}

	if s.bySender.Has(ANY_ID) == true {
		return true
	}

	if sender == s.topicFactory.GetByName(ANY) {
		return false
	}

	key := HashValue(sender)
	return s.bySender.Has(key)
}

// Emit this signal on behalf of *sender*, passing on Context.
// Returns a list of senders that accepted the signal
func (s *Sima) GetReceiversFor(sender *Topic) []ReceiverType {
	if s.receivers.Len() == 0 {
		return []ReceiverType{}
	}

	senderId := HashValue(sender)
	var ids []interface{}

	if s.bySender.Has(senderId) {
		ids = s.bySender.Get(senderId).(mapset.Set).ToSlice()
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

// Emit this signal on behalf of *sender*, passing on Context.
// Returns a list of results from the receivers
func (s *Sima) Dispatch(context context.Context, sender *Topic) []interface{} {
	if s.receivers.Len() == 0 {
		return []interface{}{}
	}

	var result []interface{}
	if sender == nil {
		sender = s.topicFactory.GetByName(ANY)
	}

	for _, receiver := range s.GetReceiversFor(sender) {
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
