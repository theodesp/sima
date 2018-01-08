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
	receivers     *cmap.CMap
	by_receiver   *cmap.CMap
	by_sender     *cmap.CMap
	symbolFactory *SymbolFactory
	closed        uint64
}

type Sima struct {
	*simaInternal
	pad [128 - unsafe.Sizeof(simaInternal{})%128]byte
}

func New() *Sima {
	s := &Sima{
		simaInternal: &simaInternal{},
	}
	runtime.SetFinalizer(s, (*Sima).clearState)
	return s
}

// Connects *receiver* to signal events sent by *sender*
func (s *Sima) Connect(receiver ReceiverType, sender *Symbol) ReceiverType {
	if sender == nil {
		sender = s.symbolFactory.GetNamed(ANY)
	}

	receiverId := HashValue(receiver)
	var senderId uint64
	if sender == s.symbolFactory.GetNamed(ANY) {
		senderId = ANY_ID
	} else {
		senderId = HashValue(sender)
	}

	s.receivers.Set(receiverId, receiver)
	s.by_sender.Get(senderId).(mapset.Set).Add(receiverId)
	s.by_receiver.Get(receiverId).(mapset.Set).Add(senderId)

	return receiver
}

// True if there is a receiver for *sender* at the time of the call
func (s *Sima) HasReceiversFor(sender *Symbol) bool {
	if s.receivers.Len() == 0 {
		return false
	}

	if s.by_sender.Has(ANY_ID) == true {
		return true
	}

	if sender == s.symbolFactory.GetNamed(ANY) {
		return false
	}

	key := HashValue(sender)
	return s.by_sender.Has(key)
}

// Emit this signal on behalf of *sender*, passing on Context.
// Returns a list of senders that accepted the signal
func (s *Sima) GetReceiversFor(sender *Symbol) []ReceiverType {
	if s.receivers.Len() == 0 {
		return []ReceiverType{}
	}

	var senderId uint64
	var ids []interface{}
	senderId = HashValue(sender)
	if s.by_sender.Has(senderId) {
		ids = s.by_sender.Get(senderId).(mapset.Set).ToSlice()
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
func (s *Sima) Dispatch(sender *Symbol, context *context.Context) []interface{} {
	if s.receivers.Len() == 0 {
		return []interface{}{}
	}

	var result []interface{}
	for _, receiver := range s.GetReceiversFor(sender) {
		result = append(result, receiver(context))
	}

	return result
}

// Clean up remaining state
func (p *Sima) clearState() {
	if p != nil {
		for _, key := range p.by_receiver.Keys() {
			p.by_receiver.DeleteAndGet(key).(mapset.Set).Clear()
		}
		for _, key := range p.by_sender.Keys() {
			p.by_sender.DeleteAndGet(key).(mapset.Set).Clear()
		}
		for _, key := range p.receivers.Keys() {
			p.receivers.Delete(key)
		}
		// Set closed to true atomically
		atomic.StoreUint64(&(p.closed), uint64(1))
		fmt.Println("state cleared")
	}
}
