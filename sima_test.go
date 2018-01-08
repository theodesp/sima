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

func (s *SimaSuite)TestMultipleSubscriptions(c *C) {
	tf := NewTopicFactory()
	hello := NewSima(tf)
	var i int

	f := func(context context.Context, sender *Topic) interface{} {
		i += 1
		return sender
	}

	hello.Connect(f, tf.GetNamed("hello"))
	hello.Connect(f, tf.GetNamed("hello"))

	response := hello.Dispatch(context.Background(), tf.GetNamed("hello"))

	// Called only once for same function and sender
	c.Assert(i, Equals, 1)
	c.Assert(len(response), Equals, 1)


	hello.Connect(func(context context.Context, sender *Topic) interface{} {
		i += 1
		return sender
	}, tf.GetNamed("hello"))

	response = hello.Dispatch(context.Background(), tf.GetNamed("hello"))

	// Called every-time for same function signature and sender
	c.Assert(i, Equals, 3)
	c.Assert(len(response), Equals, 2)
}