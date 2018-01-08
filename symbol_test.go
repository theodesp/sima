package sima

import (
	. "gopkg.in/check.v1"
)

type SignalSuite struct{}

var _ = Suite(&SignalSuite{})

func (s *SimaSuite) TestNewSymbolFactory(c *C) {
	f := NewSymbolFactory()

	c.Assert(f.GetNames(), DeepEquals, []interface{}{ANY})

	hello := f.GetNamed("hello")
	f.GetNamed("world")
	f.GetNamed("hello")

	c.Assert(f.GetNames(), DeepEquals, []interface{}{"world", "hello"})
	// Signal reusage
	c.Assert(f.GetNamed("hello"), Equals, hello)
}
