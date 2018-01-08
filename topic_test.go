package sima

import (
	. "gopkg.in/check.v1"
)

type SignalSuite struct{}

var _ = Suite(&SignalSuite{})

func (s *SimaSuite) TestNewSymbolFactory(c *C) {
	f := NewTopicFactory()

	c.Assert(f.Names(), DeepEquals, []interface{}{ANY})

	hello := f.GetByName("hello")
	f.GetByName("world")
	f.GetByName("hello")

	c.Assert(f.Names(), DeepEquals, []interface{}{ANY, "world", "hello"})
	// Signal re-usage
	c.Assert(f.GetByName("hello"), Equals, hello)
}
