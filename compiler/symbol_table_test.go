package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": {Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		"b": {Name: "b", Mutable: true, Scope: GlobalScope, Index: 1},
		"c": {Name: "c", Mutable: false, Scope: GlobalScope, Index: 2},
		"d": {Name: "d", Mutable: true, Scope: LocalScope, Index: 0},
		"e": {Name: "e", Mutable: false, Scope: LocalScope, Index: 1},
		"f": {Name: "f", Mutable: false, Scope: LocalScope, Index: 0},
		"g": {Name: "g", Mutable: false, Scope: LocalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a", false)
	if a != expected["a"] {
		t.Errorf("expected %+v, got %+v", expected["a"], a)
	}

	b := global.Define("b", true)
	if b != expected["b"] {
		t.Errorf("expected %+v, got %+v", expected["b"], b)
	}

	c := global.Define("c", false)
	if c != expected["c"] {
		t.Errorf("expected %+v, got %+v", expected["c"], c)
	}

	firstLocal := NewEnclosedSymbolTable(global)

	d := firstLocal.Define("d", true)
	if d != expected["d"] {
		t.Errorf("expected %+v, got %+v", expected["d"], d)
	}

	e := firstLocal.Define("e", false)
	if e != expected["e"] {
		t.Errorf("expected %+v, got %+v", expected["e"], e)
	}

	secondLocal := NewEnclosedSymbolTable(firstLocal)

	f := secondLocal.Define("f", false)
	if f != expected["f"] {
		t.Errorf("expected %+v, got %+v", expected["f"], f)
	}

	g := secondLocal.Define("g", false)
	if g != expected["g"] {
		t.Errorf("expected %+v, got %+v", expected["g"], g)
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

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a", false)
	global.Define("b", true)

	local := NewEnclosedSymbolTable(global)
	local.Define("c", false)
	local.Define("d", true)

	expected := []Symbol{
		{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		{Name: "b", Mutable: true, Scope: GlobalScope, Index: 1},
		{Name: "c", Mutable: false, Scope: LocalScope, Index: 0},
		{Name: "d", Mutable: true, Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := local.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a", true)
	global.Define("b", false)

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c", true)
	firstLocal.Define("d", false)

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", true)
	secondLocal.Define("f", false)

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				Symbol{Name: "a", Mutable: true, Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Mutable: false, Scope: GlobalScope, Index: 1},
				Symbol{Name: "c", Mutable: true, Scope: LocalScope, Index: 0},
				Symbol{Name: "d", Mutable: false, Scope: LocalScope, Index: 1},
			},
		},
		{
			secondLocal,
			[]Symbol{
				Symbol{Name: "a", Mutable: true, Scope: GlobalScope, Index: 0},
				Symbol{Name: "b", Mutable: false, Scope: GlobalScope, Index: 1},
				Symbol{Name: "e", Mutable: true, Scope: LocalScope, Index: 0},
				Symbol{Name: "f", Mutable: false, Scope: LocalScope, Index: 1},
			},
		},
	}
	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		{Name: "a", Mutable: false, Scope: BuiltinScope, Index: 0},
		{Name: "b", Mutable: true, Scope: BuiltinScope, Index: 1},
		{Name: "c", Mutable: false, Scope: BuiltinScope, Index: 2},
	}

	for i, v := range expected {
		global.DefineBuiltin(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Fatalf("expected to resolve %q", sym.Name)
			}

			// Builtins are always immutable
			sym.Mutable = false

			if result != sym {
				t.Errorf("expected %+v, got %+v", sym, result)
			}
		}
	}
}

func TestResolveFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a", false)
	global.Define("b", false)

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c", false)
	firstLocal.Define("d", false)

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)

	tests := []struct {
		table               *SymbolTable
		expectedSymbols     []Symbol
		expectedFreeSymbols []Symbol
	}{
		{
			firstLocal,
			[]Symbol{
				{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
				{Name: "b", Mutable: false, Scope: GlobalScope, Index: 1},
				{Name: "c", Mutable: false, Scope: LocalScope, Index: 0},
				{Name: "d", Mutable: false, Scope: LocalScope, Index: 1},
			},
			[]Symbol{},
		},
		{
			secondLocal,
			[]Symbol{
				{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
				{Name: "b", Mutable: false, Scope: GlobalScope, Index: 1},
				{Name: "c", Mutable: false, Scope: FreeScope, Index: 0},
				{Name: "d", Mutable: false, Scope: FreeScope, Index: 1},
				{Name: "e", Mutable: false, Scope: LocalScope, Index: 0},
				{Name: "f", Mutable: false, Scope: LocalScope, Index: 1},
			},
			[]Symbol{
				{Name: "c", Mutable: false, Scope: LocalScope, Index: 0},
				{Name: "d", Mutable: false, Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}

			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
			}
		}

		if len(tt.table.FreeSymbols) != len(tt.expectedFreeSymbols) {
			t.Errorf("wrong number of free symbols. got %d, want %d", len(tt.table.FreeSymbols), len(tt.expectedFreeSymbols))
			continue
		}

		for i, sym := range tt.expectedFreeSymbols {
			result := tt.table.FreeSymbols[i]
			if result != sym {
				t.Errorf("wrong free symbol. got %+v, want %+v", result, sym)
			}
		}
	}
}

func TestResolveUnresolvableFree(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a", false)

	firstLocal := NewEnclosedSymbolTable(global)
	firstLocal.Define("c", false)

	secondLocal := NewEnclosedSymbolTable(firstLocal)
	secondLocal.Define("e", false)
	secondLocal.Define("f", false)

	expected := []Symbol{
		{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0},
		{Name: "c", Mutable: false, Scope: FreeScope, Index: 0},
		{Name: "e", Mutable: false, Scope: LocalScope, Index: 0},
		{Name: "f", Mutable: false, Scope: LocalScope, Index: 1},
	}

	for _, sym := range expected {
		result, ok := secondLocal.Resolve(sym.Name)
		if !ok {
			t.Errorf("name %s not resolvable", sym.Name)
			continue
		}

		if result != sym {
			t.Errorf("expected %s to resolve to %+v, got %+v", sym.Name, sym, result)
		}
	}

	for _, name := range []string{"b", "d"} {
		_, ok := secondLocal.Resolve(name)
		if ok {
			t.Errorf("name %s resolved, but was expected not to", name)
		}
	}
}

func TestDefineAndResolveFunctionName(t *testing.T) {
	global := NewSymbolTable()
	global.DefineFunctionName("a", false)

	expected := Symbol{Name: "a", Mutable: false, Scope: FunctionScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got %+v", expected.Name, expected, result)
	}
}

func TestShadowingFunctionName(t *testing.T) {
	global := NewSymbolTable()
	global.DefineFunctionName("a", false)
	global.Define("a", false)

	expected := Symbol{Name: "a", Mutable: false, Scope: GlobalScope, Index: 0}

	result, ok := global.Resolve(expected.Name)
	if !ok {
		t.Fatalf("function name %s not resolvable", expected.Name)
	}

	if result != expected {
		t.Errorf("expected %s to resolve to %+v, got %+v", expected.Name, expected, result)
	}
}
