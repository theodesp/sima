package sima

import (
	"github.com/OneOfOne/cmap"
)

// This initial signal name is imported by default
const ANY = "ANY"

// A constant symbol
type Symbol struct {
	name string
}

type SymbolFactory struct {
	symbols *cmap.CMap
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name}
}

// Repeated calls of SymbolFactory('name') will all return the same instance.
func NewSymbolFactory() *SymbolFactory {
	factory := &SymbolFactory{
		cmap.New(),
	}
	factory.symbols.Set(ANY, NewSymbol(ANY))

	return factory
}

func (f *SymbolFactory) GetNamed(name string) *Symbol {
	if f.symbols.Has(name) {
		return f.symbols.Get(name).(*Symbol)
	} else {
		symbol := NewSymbol(name)
		f.symbols.Set(name, symbol)
		return symbol
	}
}

func (f *SymbolFactory) GetNames() []interface{} {
	return f.symbols.Keys()
}
