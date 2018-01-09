package sima

import (
	. "gopkg.in/check.v1"
	"testing"
	"context"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type SimaSuite struct{}

var _ = Suite(&SimaSuite{})


func (s *SimaSuite)TestSingleSubscriptions(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	},nil)

	response := hello.Dispatch(context.Background(), nil)

	c.Assert(len(response), Equals, 0)
}

func (s *SimaSuite)TestHasReceivers(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)

	c.Check(hello.HasReceiversFor(nil), Equals, false)
	c.Check(hello.HasReceiversFor(tf.GetByName(ANY)), Equals, false)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	},nil)

	c.Check(hello.HasReceiversFor(tf.GetByName(ANY)), Equals, true)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	},tf.GetByName("Hello"))

	c.Check(hello.HasReceiversFor(tf.GetByName("Hello")), Equals, true)
}

func (s *SimaSuite)TestMultipleSubscriptions(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)
	var i int

	f := func(context context.Context, sender *Topic) interface{} {
		i += 1
		return sender
	}

	hello.Connect(f, tf.GetByName("hello"))
	hello.Connect(f, tf.GetByName("hello"))

	response := hello.Dispatch(context.Background(), tf.GetByName("hello"))

	// Called only once for same function and sender
	c.Assert(i, Equals, 1)
	c.Assert(len(response), Equals, 1)


	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		i += 1
		return sender
	}, tf.GetByName("hello"))

	response = hello.Dispatch(context.Background(), tf.GetByName("hello"))

	// Called every-time for same function signature and sender
	c.Assert(i, Equals, 3)
	c.Assert(len(response), Equals, 2)
}