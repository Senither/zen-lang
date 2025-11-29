package compiler

import (
	"fmt"

	"github.com/senither/zen-lang/objects"
)

type SymbolScope string

const (
	GlobalScope        SymbolScope = "GLOBAL"
	LocalScope         SymbolScope = "LOCAL"
	FreeScope          SymbolScope = "FREE"
	BuiltinScope       SymbolScope = "BUILTIN"
	GlobalBuiltinScope SymbolScope = "GLOBAL_BUILTIN"
	FunctionScope      SymbolScope = "FUNCTION"
)

type SymbolKind string

const (
	NativeKind SymbolKind = "NATIVE"
	ArrayKind  SymbolKind = "ARRAY"
	HashKind   SymbolKind = "HASH"
)

type Symbol struct {
	Name    string
	Mutable bool
	Scope   SymbolScope
	Index   int
	Kind    SymbolKind
}

func (s *Symbol) isEligibleForFreeing() bool {
	return s.Scope != GlobalScope &&
		s.Scope != BuiltinScope &&
		s.Scope != GlobalBuiltinScope
}

type SymbolTable struct {
	Outer *SymbolTable

	store          map[string]Symbol
	numDefinitions int

	FreeSymbols []Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
		FreeSymbols:    []Symbol{},
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
	symbol := Symbol{Name: name, Mutable: false, Scope: BuiltinScope, Index: index, Kind: NativeKind}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineGlobalBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Mutable: false, Scope: GlobalBuiltinScope, Index: index, Kind: NativeKind}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string, mutable bool) Symbol {
	symbol := Symbol{Name: name, Mutable: mutable, Index: 0, Scope: FunctionScope, Kind: NativeKind}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) Define(name string, mutable bool) Symbol {
	symbol := Symbol{Name: name, Mutable: mutable, Index: s.numDefinitions, Kind: NativeKind}

	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) defineFree(original Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, original)

	symbol := Symbol{
		Name:    original.Name,
		Mutable: original.Mutable,
		Scope:   FreeScope,
		Index:   len(s.FreeSymbols) - 1,
		Kind:    original.Kind,
	}

	s.store[original.Name] = symbol

	return symbol
}

func (s *SymbolTable) UpdateKind(name string, kind SymbolKind) error {
	symbol, ok := s.store[name]
	if !ok {
		return fmt.Errorf("symbol %s not found", name)
	}

	symbol.Kind = kind
	s.store[name] = symbol

	return nil
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := s.store[name]
	if !ok && s.Outer != nil {
		symbol, ok = s.Outer.Resolve(name)
		if !ok {
			return symbol, ok
		}

		if !symbol.isEligibleForFreeing() {
			return symbol, ok
		}

		free := s.defineFree(symbol)
		return free, true
	}

	return symbol, ok
}
