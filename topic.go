package sima

import (
	"github.com/OneOfOne/cmap"
)

// This initial signal name is imported by default
const ANY = "ANY"
const NONE = "NONE"

// A constant topic symbol
type Topic struct {
	name string
}

type TopicFactory struct {
	topics *cmap.CMap
}

func NewTopic(name string) *Topic {
	return &Topic{name}
}

// Repeated calls of TopicFactory('name') will all return the same instance.
func NewTopicFactory() *TopicFactory {
	factory := &TopicFactory{
		cmap.New(),
	}
	factory.topics.Set(ANY, NewTopic(ANY))

	return factory
}

func (f *TopicFactory) GetNamed(name string) *Topic {
	if f.topics.Has(name) {
		return f.topics.Get(name).(*Topic)
	} else {
		topic := NewTopic(name)
		f.topics.Set(name, topic)
		return topic
	}
}

func (f *TopicFactory) GetNames() []interface{} {
	return f.topics.Keys()
}
