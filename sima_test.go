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


func (s *SimaSuite)TestSingleSubscriptionsWithNoSenders(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	},"")

	response := hello.Dispatch(context.Background(), "")

	c.Assert(len(response), Equals, 1)
	c.Assert(response[0], DeepEquals, tf.GetByName(ANY))
}

func (s *SimaSuite)TestSingleSubscriptionsWithSenders(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	},"hello")

	response1 := hello.Dispatch(context.Background(), "")
	c.Assert(len(response1), Equals, 0)

	response2 := hello.Dispatch(context.Background(), "world")
	c.Assert(len(response2), Equals, 0)

	response3 := hello.Dispatch(context.Background(), "hello")
	c.Assert(len(response3), Equals, 1)
	c.Assert(response3[0], DeepEquals, tf.GetByName("hello"))
}

func (s *SimaSuite)TestHasReceiversFor(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)

	c.Check(hello.HasReceiversFor(""), Equals, false)
	c.Check(hello.HasReceiversFor(""), Equals, false)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	},"")

	c.Check(hello.HasReceiversFor(""), Equals, true)

	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		return sender
	}, "Hello")

	c.Check(hello.HasReceiversFor("Hello"), Equals, true)
}

func (s *SimaSuite)TestMultipleSubscriptions(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)
	var i int

	f := func(context context.Context, sender *Topic) interface{} {
		i += 1
		return sender
	}

	hello.Connect(f, "hello")
	hello.Connect(f, "hello")

	response := hello.Dispatch(context.Background(), "hello")

	// Called only once for same function and sender
	c.Assert(i, Equals, 1)
	c.Assert(len(response), Equals, 1)


	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		i += 1
		return sender
	}, "hello")

	response = hello.Dispatch(context.Background(), "hello")

	// Called every-time for same function signature and sender
	c.Assert(i, Equals, 3)
	c.Assert(len(response), Equals, 2)
}