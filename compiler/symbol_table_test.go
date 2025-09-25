package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := []Symbol{
		{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		{Name: "b", Mutable: true, Scope: GlobalScope, Index: 1},
		{Name: "c", Mutable: false, Scope: GlobalScope, Index: 2},
	}

	global := NewSymbolTable()

	for _, sym := range expected {
		defined := global.Define(sym.Name, sym.Mutable)

		if defined != sym {
			t.Errorf("expected %+v got %+v", sym, defined)
		}
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()

	global.Define("a", false)
	global.Define("b", true)
	global.Define("c", false)

	expected := []Symbol{
		{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		{Name: "b", Mutable: true, Scope: GlobalScope, Index: 1},
		{Name: "c", Mutable: false, Scope: GlobalScope, Index: 2},
	}

	for _, sym := range expected {
		result, ok := global.Resolve(sym.Name)
		if !ok {
			t.Fatalf("expected to resolve %q", sym.Name)
		}

		if result != sym {
			t.Errorf("expected %+v, got %+v", sym, result)
		}
	}
}
