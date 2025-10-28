package compiler

import (
	"fmt"

	"github.com/senither/zen-lang/objects"
)

type SymbolScope string

const (
	GlobalScope        SymbolScope = "GLOBAL"
	LocalScope         SymbolScope = "LOCAL"
	BuiltinScope       SymbolScope = "BUILTIN"
	GlobalBuiltinScope SymbolScope = "GLOBAL_BUILTIN"
)

type Symbol struct {
	Name    string
	Mutable bool
	Scope   SymbolScope
	Index   int
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer

	return s
}

func WriteBuiltinSymbols(table *SymbolTable) {
	for i, v := range objects.Builtins {
		table.DefineBuiltin(i, v.Name)
	}

	for sIdx, s := range objects.Globals {
		for bIdx, b := range s.Builtins {
			name := fmt.Sprintf("%s.%s", s.Name, b.Name)
			idx := (uint16(sIdx) << 8) | uint16(bIdx)

			table.DefineGlobalBuiltin(int(idx), name)
		}
	}
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Mutable: false, Scope: BuiltinScope, Index: index}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineGlobalBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Mutable: false, Scope: GlobalBuiltinScope, Index: index}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) Define(name string, mutable bool) Symbol {
	symbol := Symbol{Name: name, Mutable: mutable, Index: s.numDefinitions}

	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.Outer != nil {
		return s.Outer.Resolve(name)
	}

	return symbol, ok
}
