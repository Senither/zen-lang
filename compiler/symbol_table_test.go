package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Mutable: false, Scope: GlobalScope, Index: 1},
		"c": {Name: "c", Mutable: false, Scope: GlobalScope, Index: 2},
	}

	global := NewSymbolTable()

	for name := range expected {
		sym := global.Define(expected[name].Name, expected[name].Mutable)

		if sym != expected[name] {
			t.Errorf("expected %+v got %+v", expected[name], sym)
		}
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()

	global.Define("a", false)
	global.Define("b", true)

	expected := map[string]Symbol{
		"a": {Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Mutable: true, Scope: GlobalScope, Index: 1},
	}

	for name, expectedSymbol := range expected {
		result, ok := global.Resolve(name)
		if !ok {
			t.Fatalf("expected to resolve %q", name)
		}

		if result != expectedSymbol {
			t.Errorf("expected %+v, got %+v", expectedSymbol, result)
		}
	}
}
