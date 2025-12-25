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
		name     string
		input    string
		expected int64
	}{
		{"return int", "return 10;", 10},
		{"return int with extra statement", "return 10; 9;", 10},
		{"return expression with extra statement", "return 2 * 5; 9;", 10},
		{"expression before return with extra statement", "9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestNestedReturnStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"nested return",
			"if (10 > 1) { if (10 > 1) { return 10; } return 1; }",
			10,
		},
		{
			"nested return in function",
			`
				func test() {
					if (10 > 1) {
						if (10 > 1) {
							return 10;
						}
						return 1;
					}
				}

				test()
			`,
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"variable declaration", "var x = 5; x;", 5},
		{"variable declaration with expression", "var x = 5 * 5; x;", 25},
		{"multiple variable declarations", "var x = 5; var y = 10; x + y;", 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestCompoundAssignments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"addition assignment", "var mut x = 5; x += 5;", 10},
		{"subtraction assignment", "var mut x = 5; x -= 5;", 0},
		{"multiplication assignment", "var mut x = 5; x *= 5;", 25},
		{"division assignment", "var mut x = 5; x /= 5;", 1},
		{"modulus assignment", "var mut x = 5; x %= 5;", 0},
		{"exponentiation assignment", "var mut x = 5; x ^= 5;", 3125},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestWhileBreakStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"break in while loop with increment of 1",
			"var mut i = 0; while (true) { if (i > 10) { break; } i = i + 1; } i;",
			11,
		},
		{
			"break in while loop with increment of 2",
			"var mut i = 0; while (true) { if (i > 10) { break; } i = i + 2; } i;",
			12,
		},
		{
			"break in while loop with increment of 2 and greater or equal condition",
			"var mut i = 0; while (true) { if (i >= 10) { break; } i = i + 2; } i;",
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestExportStatement(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			"export single function",
			`export func name(x) {x};`,
			[]string{"name"},
		},
		{
			"export after function declaration",
			`
				func name(x) {x};
				export name;
			`,
			[]string{"name"},
		},
		{
			"export multiple functions",
			`
				export func functionOne() {}
				func functionTwo(x) {x};
				export func functionThree(x) {x};
			`,
			[]string{"functionOne", "functionThree"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l, nil)

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
		})
	}
}
