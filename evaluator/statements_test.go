package evaluator

import (
	"slices"
	"testing"

	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestNestedReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var x = 5; x;", 5},
		{"var x = 5 * 5; x;", 25},
		{"var x = 5; var y = 10; x + y;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestExportStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			`export func name(x) {x};`,
			[]string{"name"},
		},
		{
			`
				func name(x) {x};
				export name;
			`,
			[]string{"name"},
		},
		{
			`
				export func functionOne() {}
				func functionTwo(x) {x};
				export func functionThree(x) {x};
			`,
			[]string{"functionOne", "functionThree"},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		env := objects.NewEnvironment(nil)
		evaluated := Eval(p.ParseProgram(), env)
		if evaluated == nil {
			t.Errorf("Failed to evaluate program, evaluation returned nil")
		}

		exports := env.GetExports()
		exportKeys := []string{}
		for k := range exports {
			exportKeys = append(exportKeys, k)
		}

		if len(exportKeys) != len(tt.expected) {
			t.Errorf("Expected %d exported keys, got %d", len(tt.expected), len(exportKeys))
		}

		for _, v := range tt.expected {
			if !slices.Contains(exportKeys, v) {
				t.Errorf("Expected exported key %q to be in %v", v, exportKeys)
			}
		}
	}
}
