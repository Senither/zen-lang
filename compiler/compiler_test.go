package compiler

import (
	"testing"

	"github.com/senither/zen-lang/code"
)

type compilerTestCase struct {
	name                 string
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestNullExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "null literal",
			input:             "null;",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpNull),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkNullExpression(b *testing.B) {
	runCompilationBenchmarks(b, []string{"null"})
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "integer literals",
			input:             "1; 2;",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "float literals",
			input:             "1.5; 3.14;",
			expectedConstants: []any{1.5, 3.14},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "addition",
			input:             "1 + 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "subtraction",
			input:             "1 - 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSub),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "multiplication",
			input:             "1 * 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMul),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "division",
			input:             "1 / 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDiv),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "exponentiation",
			input:             "2 ^ 3",
			expectedConstants: []any{2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPow),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "modulus",
			input:             "1 % 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMod),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "unary minus integer",
			input:             "-1",
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "unary minus float",
			input:             "-1.5",
			expectedConstants: []any{1.5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkIntegerArithmetic(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"1; 2;",
		"1.5; 3.14;",
		"1 + 2",
		"1 - 2",
		"1 * 2",
		"1 / 2",
		"2 ^ 3",
		"1 % 2",
		"-1",
		"-1.5",
	})
}

func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "boolean true literal",
			input:             "true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "boolean false literal",
			input:             "false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "greater than expression",
			input:             "1 > 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "less than expression",
			input:             "1 < 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "equal expression",
			input:             "1 == 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "not equal expression",
			input:             "1 != 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "greater than or equal expression",
			input:             "1 >= 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThanOrEqual),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "less than or equal expression",
			input:             "1 <= 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThanOrEqual),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "boolean equality expression",
			input:             "true == false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "boolean inequality expression",
			input:             "true != false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "boolean not expression",
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkBooleanExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"true",
		"false",
		"1 > 2",
		"1 < 2",
		"1 == 2",
		"1 != 2",
		"1 >= 2",
		"1 <= 2",
		"true == false",
		"true != false",
		"!true",
	})
}

func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "double quoted string literal",
			input:             `"hello world";`,
			expectedConstants: []any{"hello world"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "single quoted string literal",
			input:             `'hello world';`,
			expectedConstants: []any{"hello world"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "string concatenation",
			input:             `"hello" + 'world';`,
			expectedConstants: []any{"hello", "world"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkStringExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`"hello world";`,
		`'hello world';`,
		`"hello" + 'world';`,
	})
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "empty array literal",
			input:             "[]",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "array literal with integers",
			input:             "[1, 2, 3]",
			expectedConstants: []any{1, 2, 3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "array literal with expressions",
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkArrayLiterals(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"[]",
		"[1, 2, 3]",
		"[1 + 2, 3 - 4, 5 * 6]",
	})
}

func TestHashLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "empty hash literal",
			input:             "{}",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "hash literal with integer keys and values",
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "hash literal with expressions as values",
			input:             "{1: 2 + 3, 4: 5 * 6, 7: 8 - 9}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMul),
				code.Make(code.OpConstant, 6),
				code.Make(code.OpConstant, 7),
				code.Make(code.OpConstant, 8),
				code.Make(code.OpSub),
				code.Make(code.OpHash, 6),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "hash literal with string keys and values",
			input:             "{'key': 'value', 'another': 'pair'}",
			expectedConstants: []any{"another", "pair", "key", "value"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpHash, 4),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkHashLiterals(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"{}",
		"{1: 2, 3: 4, 5: 6}",
		"{1: 2 + 3, 4: 5 * 6, 7: 8 - 9}",
		"{'key': 'value', 'another': 'pair'}",
	})
}

func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "index expression with array",
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []any{1, 2, 3, 1, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "index expression with hash",
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []any{1, 2, 2, 1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSub),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkIndexExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"[1, 2, 3][1 + 1]",
		"{1: 2}[2 - 1]",
	})
}

func TestChainIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "chain index expression with single key",
			input:             "var test = {}; test.key",
			expectedConstants: []any{"key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "chain index expression with multiple keys",
			input:             "var test = {}; test.another.key",
			expectedConstants: []any{"another", "key"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpHash, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "chain index expression with array inside hash",
			input:             "var obj = {'key': [1, 2, 3]}; obj.key[0]",
			expectedConstants: []any{"key", 1, 2, 3, "key", 0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpArray, 3),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkChainIndexExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"var test = {}; test.key",
		"var test = {}; test.another.key",
		"var obj = {'key': [1, 2, 3]}; obj.key[0]",
	})
}

func TestChainIndexAssignments(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "index assignment with string key",
			input:             "var obj = {'key': 'value'}; obj['key'] = 42;",
			expectedConstants: []any{"key", "value", "key", 42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "index assignment with dot notation",
			input:             "var obj = {'key': 'value'}; obj.key = 42;",
			expectedConstants: []any{"key", "value", "key", 42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndexAssign),
			},
		},
		{
			name:              "chain index assignment with multiple keys using bracket notation",
			input:             "var obj = {'nested': {'key': 'value'}}; obj['nested']['key'] = 42;",
			expectedConstants: []any{"nested", "key", "value", "nested", "key", 42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "chain index assignment with multiple keys using dot notation",
			input:             "var obj = {'nested': {'key': 'value'}}; obj.nested.key = 42;",
			expectedConstants: []any{"nested", "key", "value", "nested", "key", 42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkChainIndexAssignments(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"var obj = {'key': 'value'}; obj['key'] = 42;",
		"var obj = {'key': 'value'}; obj.key = 42;",
		"var obj = {'nested': {'key': 'value'}}; obj['nested']['key'] = 42;",
		"var obj = {'nested': {'key': 'value'}}; obj.nested.key = 42;",
	})
}

func TestChainCallExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "chain call expression with single key",
			input: `
				var obj = { 'method': func() { 42 } };
				obj.method()
			`,
			expectedConstants: []any{
				"method",
				42,
				[]code.Instructions{
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
				"method",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndex),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "chain call expression with multiple keys",
			input: `
				var obj = {
					'a': {
						'b': {
							'c': func(a, b) { a + b }
						}
					}
				}

				obj.a.b.c(9, 42)
			`,
			expectedConstants: []any{
				"a",
				"b",
				"c",
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				"a",
				"b",
				"c",
				9,
				42,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpClosure, 3, 0),
				code.Make(code.OpHash, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 6),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 7),
				code.Make(code.OpConstant, 8),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
		},
		{
			name: "chain call expression with compiled function",
			input: `
				maps.keys({});
				var maps = {'keys': func (a) { a }};
				maps.keys({})
			`,
			expectedConstants: []any{
				"keys",
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				"keys",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpGetGlobalBuiltin, 512),
				code.Make(code.OpHash, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpIndex),
				code.Make(code.OpHash, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkChainCallExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var obj = { 'method': func() { 42 } };
			obj.method()
		`,
		`
			var obj = {
				'a': {
					'b': {
						'c': func(a, b) { a + b }
					}
				}
			}

			obj.a.b.c(9, 42)
		`,
	})
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:              "if statement with true condition",
			input:             "if (true) { 10 }; 5;",
			expectedConstants: []any{10, 5},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 11),
				// 0010
				code.Make(code.OpNull),
				// 0011
				code.Make(code.OpPop),
				// 0012
				code.Make(code.OpConstant, 1),
				// 0015
				code.Make(code.OpPop),
			},
		},
		{
			name:              "if-else statement",
			input:             "if (true) { 10 } else { 20 }; 5;",
			expectedConstants: []any{10, 20, 5},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpTrue),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 13),
				// 0010
				code.Make(code.OpConstant, 1),
				// 0013
				code.Make(code.OpPop),
				// 0014
				code.Make(code.OpConstant, 2),
				// 0017
				code.Make(code.OpPop),
			},
		},
		{
			name:              "if-else if-else statement",
			input:             "if (false) { 10 } else if (true) { 20 } else { 30 }; 5;",
			expectedConstants: []any{10, 20, 30, 5},
			expectedInstructions: []code.Instructions{
				// 0000
				code.Make(code.OpFalse),
				// 0001
				code.Make(code.OpJumpNotTruthy, 10),
				// 0004
				code.Make(code.OpConstant, 0),
				// 0007
				code.Make(code.OpJump, 23),
				// 0010
				code.Make(code.OpTrue),
				// 0011
				code.Make(code.OpJumpNotTruthy, 20),
				// 0014
				code.Make(code.OpConstant, 1),
				// 0017
				code.Make(code.OpJump, 23),
				// 0020
				code.Make(code.OpConstant, 2),
				// 0023
				code.Make(code.OpPop),
				// 0024
				code.Make(code.OpConstant, 3),
				// 0027
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkConditionals(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		"if (true) { 10 }; 5;",
		"if (true) { 10 } else { 20 }; 5;",
		"if (false) { 10 } else if (true) { 20 } else { 30 }; 5;",
	})
}

func TestGlobalVarStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "global variable declarations",
			input: `
				var one = 1;
				var two = 2;
			`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			name: "global mutable variable declarations",
			input: `
				var mut one = 1;
				var mut two = 2;
			`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			name: "global variable usage",
			input: `
				var one = 1;
				one;
			`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "global mutable variable usage",
			input: `
				var mut one = 1;
				var two = one;
				two;
			`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
		{
			name: "global variable arithmetic",
			input: `
				var one = 1;
				var two = 2;
				one + two;
			`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkGlobalVarStatements(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var one = 1;
			var two = 2;
		`,
		`
			var mut one = 1;
			var mut two = 2;
		`,
		`
			var one = 1;
			one;
		`,
		`
			var mut one = 1;
			var two = one;
			two;
		`,
		`
			var one = 1;
			var two = 2;
			one + two;
		`,
	})
}

func TestVarStatementScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "var statement in global scope",
			input: `
				var num = 55;
				func() { num }
			`,
			expectedConstants: []any{
				55,
				[]code.Instructions{
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "var statement in function scope",
			input: `
				func() {
					var num = 55;
					num
				}
			`,
			expectedConstants: []any{
				55,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "multiple var statements in function scope",
			input: `
				func() {
					var a = 1;
					var b = 2;
					a + b
				}
			`,
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkVarStatementScopes(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var num = 55;
			func() { num }
		`,
		`
			func() {
				var num = 55;
				num
			}
		`,
		`
			func() {
				var a = 1;
				var b = 2;
				a + b
			}
		`,
	})
}

func TestVarIncDec(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "global mutable variable increment",
			input: `
				var mut a = 1;
				a++;
			`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpIncGlobal),
				code.Make(code.OpPop),
			},
		},
		{
			name: "global mutable variable decrement",
			input: `
				var mut a = 1;
				a--;
			`,
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpDecGlobal),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkVarIncDec(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var mut a = 1;
			a++;
		`,
		`
			var mut a = 1;
			a--;
		`,
	})
}

func TestFuncLocalIncDec(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "func local mutable variable increment",
			input: `
				func() {
					var mut a = 1;
					a++;
				}
			`,
			expectedConstants: []any{
				1,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpIncLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "func local mutable variable decrement",
			input: `
				func() {
					var mut a = 1;
					a--;
				}
			`,
			expectedConstants: []any{
				1,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpDecLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkFuncLocalIncDec(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			func test() {
				var mut a = 1;
				a++;
			}
			test()
		`,
		`
			func test() {
				var mut a = 1;
				a--;
			}
			test()
		`,
	})
}

func TestCompilerScopes(t *testing.T) {
	compiler := New(nil)
	if compiler.scopeIndex != 0 {
		t.Fatalf("compiler at wrong initial scope index. got %d, want 0", compiler.scopeIndex)
	}

	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Fatalf("compiler did not enter scope correctly. got %d, want 1", compiler.scopeIndex)
	}

	compiler.emit(code.OpSub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Fatalf("instructions length wrong. got %d, want 1", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpSub {
		t.Fatalf("last instruction wrong. got %d, want %d", last.Opcode, code.OpSub)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Fatalf("compiler did not set outer symbol table correctly")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Fatalf("compiler did not leave scope correctly. got %d, want 0", compiler.scopeIndex)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Fatalf("compiler did not restore global symbol table")
	}

	if compiler.symbolTable.Outer != nil {
		t.Fatalf("compiler did not restore symbol table correctly")
	}

	compiler.emit(code.OpAdd)
	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Fatalf("instructions length wrong. got %d, want 2", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpAdd {
		t.Fatalf("last instruction wrong. got %d, want %d", last.Opcode, code.OpAdd)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	if previous.Opcode != code.OpMul {
		t.Fatalf("previous instruction wrong. got %d, want %d", previous.Opcode, code.OpMul)
	}
}

func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name:  "function with explicit return",
			input: "func() { return 5 + 10 }",
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:  "function with implicit return",
			input: "func() { 5 + 10 }",
			expectedConstants: []any{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:  "function with multiple statements",
			input: "func() { 1; 2 }",
			expectedConstants: []any{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:  "function with no return value",
			input: "func() { }",
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkFunctions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`func() { return 5 + 10 }`,
		`func() { 5 + 10 }`,
		`func() { 1; 2 }`,
		`func() { }`,
	})
}

func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "closure with one free variable",
			input: `
				func(a) {
					func(b) {
						a + b
					}
				}
			`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 0, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "closure with multiple free variables",
			input: `
				func(a) {
					func(b) {
						func(c) {
							a + b + c
						}
					}
				}
			`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetFree, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 0, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 1, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "nested closures with free variables",
			input: `
				var global = 55;

				func() {
					var a = 66;

					func() {
						var b = 77;

						func() {
							var c = 88;

							global + a + b + c;
						}
					}
				}
			`,
			expectedConstants: []any{
				55,
				66,
				77,
				88,
				[]code.Instructions{
					code.Make(code.OpConstant, 3),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpGetFree, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 2),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 4, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 5, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpClosure, 6, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkClosures(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			func(a) {
				func(b) {
					a + b
				}
			}
		`,
		`
			func(a) {
				func(b) {
					func(c) {
						a + b + c
					}
				}
			}
		`,
		`
			var global = 55;

			func() {
				var a = 66;

				func() {
					var b = 77;

					func() {
						var c = 88;

						global + a + b + c;
					}
				}
			}
		`,
	})
}

func TestNamedFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "named function with explicit return",
			input: `
				func example() {
					return 5;
				}
			`,
			expectedConstants: []any{
				5,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpSetGlobal, 0),
			},
		},
		{
			name: "named function with parameters",
			input: `
				func sum(a, b) {
					return a + b;
				}
			`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpSetGlobal, 0),
			},
		},
		{
			name: "nested named functions with free variables",
			input: `
				func alpha(a, b) {
					return func bravo(c, d) {
						return func charlie(e, f) {
							return a + b + c + d + e + f;
						}
					}
				}
			`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetFree, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetFree, 2),
					code.Make(code.OpAdd),
					code.Make(code.OpGetFree, 3),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetFree, 0),
					code.Make(code.OpGetFree, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpClosure, 0, 4),
					code.Make(code.OpSetLocal, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpClosure, 1, 2),
					code.Make(code.OpSetLocal, 2),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpSetGlobal, 0),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkNamedFunctions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			func example() {
				return 5;
			}
		`,
		`
			func sum(a, b) {
				return a + b;
			}
		`,
		`
			func alpha(a, b) {
				return func bravo(c, d) {
					return func charlie(e, f) {
						return a + b + c + d + e + f;
					}
				}
			}
		`,
	})
}

func TestRecursiveFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "recursive function in global scope",
			input: `
				var countDown = func(x) { countDown(x - 1) }
				countDown(1)
			`,
			expectedConstants: []any{
				1,
				[]code.Instructions{
					code.Make(code.OpCurrentClosure),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSub),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			name: "recursive function in local scope",
			input: `
				var wrapper = func() {
					var countDown = func(x) { countDown(x - 1) }
					countDown(1)
				}

				wrapper()
			`,
			expectedConstants: []any{
				1,
				[]code.Instructions{
					code.Make(code.OpCurrentClosure),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSub),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
				1,
				[]code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 3, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkRecursiveFunctions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var countDown = func(x) { countDown(x - 1) }
			countDown(1)
		`,
		`
			var wrapper = func() {
				var countDown = func(x) { countDown(x - 1) }
				countDown(1)
			}

			wrapper()
		`,
	})
}

func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "anonymous function call",
			input: `
				func() { 10 }()
			`,
			expectedConstants: []any{
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "function call with no arguments",
			input: `
				var noArgs = func() { 10 };
				noArgs()
			`,
			expectedConstants: []any{
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "function call with one argument",
			input: `
				var oneArg = func(a) { a };
				oneArg(55)
			`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				55,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			name: "function call with multiple arguments",
			input: `
				var manyArgs = func(a, b, c) { a; b; c};
				manyArgs(11, 22, 33)
			`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 2),
					code.Make(code.OpReturnValue),
				},
				11,
				22,
				33,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkFunctionCalls(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			func() { 10 }()
		`,
		`
			var noArgs = func() { 10 };
			noArgs()
		`,
		`
			var oneArg = func(a) { a };
			oneArg(55)
		`,
		`
			var manyArgs = func(a, b, c) { a; b; c};
			manyArgs(11, 22, 33)
		`,
	})
}

func TestBuiltins(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "builtins in global scope",
			input: `
				len([])
				print("Hello, World!", 42)
			`,
			expectedConstants: []any{
				"Hello, World!",
				42,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpGetBuiltin, 2),
				code.Make(code.OpArray, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
				code.Make(code.OpGetBuiltin, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
		},
		{
			name:  "builtins in function scope",
			input: `func() { len([]) }`,
			expectedConstants: []any{
				[]code.Instructions{
					code.Make(code.OpGetBuiltin, 2),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkBuiltins(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			len([])
			print("Hello, World!", 42)
		`,
		`func() { len([]) }`,
	})
}

func TestGlobalBuiltins(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "global builtins calls",
			input: `
				strings.contains("hello world", "world")
				arrays.push([1, 2, 3], 4)
			`,
			expectedConstants: []any{
				"hello world",
				"world",
				1,
				2,
				3,
				4,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpGetGlobalBuiltin, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
				code.Make(code.OpGetGlobalBuiltin, 256),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
		},
		{
			name:  "global builtins in function scope",
			input: `func() { strings.join([1, 2, 3], "-") }`,
			expectedConstants: []any{
				1,
				2,
				3,
				"-",
				[]code.Instructions{
					code.Make(code.OpGetGlobalBuiltin, 2),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpCall, 2),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 4, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkGlobalBuiltins(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			strings.contains("hello world", "world")
			arrays.push([1, 2, 3], 4)
		`,
		`func() { strings.join([1, 2, 3], "-") }`,
	})
}

func TestWhileLoop(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "while loop with a single statement",
			input: `
				while (true) { 10 }
			`,
			expectedConstants: []any{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 11),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpJump, 0),
				code.Make(code.OpLoopEnd),
			},
		},
		{
			name: "nested while loops",
			input: `
				while (true) {
					while (false) {
						10
					}
				}
			`,
			expectedConstants: []any{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 19),
				code.Make(code.OpFalse),
				code.Make(code.OpJumpNotTruthy, 15),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpJump, 4),
				code.Make(code.OpLoopEnd),
				code.Make(code.OpJump, 0),
				code.Make(code.OpLoopEnd),
			},
		},
		{
			name: "while loop with continue statement",
			input: `
				while (true) {
					if (true) {
						continue;
					}
				}
			`,
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 16),
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 11),
				code.Make(code.OpJump, 0),
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpJump, 0),
				code.Make(code.OpLoopEnd),
			},
		},
		{
			name: "while loop with break statement",
			input: `
				while (true) {
					if (true) {
						break;
					}
				}
			`,
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 16),
				code.Make(code.OpTrue),
				code.Make(code.OpJumpNotTruthy, 11),
				code.Make(code.OpJump, 16),
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpJump, 0),
				code.Make(code.OpLoopEnd),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkWhileLoop(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			while (true) { 10 }
		`,
		`
			while (true) {
				while (false) {
					10
				}
			}
		`,
		`
			while (true) {
				if (true) {
					continue;
				}
			}
		`,
		`
			while (true) {
				if (true) {
					break;
				}
			}
		`,
	})
}

func TestAssignmentExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "simple assignment expression",
			input: `
				var mut a = 1;
				a = 2;
			`,
			expectedConstants: []any{1, 2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name: "assignment expression with binary operation",
			input: `
				var mut a = 1;
				var b = 10;
				a = a + b;
			`,
			expectedConstants: []any{1, 10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkAssignmentExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var mut a = 1;
			a = 2;
		`,
		`
			var mut a = 1;
			var b = 10;
			a = a + b;
		`,
	})
}

func TestIndexAssignmentExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			name: "simple index assignment expression",
			input: `
				var mut arr = [1, 2, 3];
				arr[1] = 42;
			`,
			expectedConstants: []any{1, 2, 3, 1, 42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
		{
			name: "hash index assignment expression",
			input: `
				var mut hash = { "key": 1 };
				hash["key"] = 99;
			`,
			expectedConstants: []any{"key", 1, "key", 99},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
		{
			name: "nested object index assignment expression",
			input: `
				var mut obj = { "nested": { "key": 1 } };
				obj["nested"]["key"] = 123;
			`,
			expectedConstants: []any{"nested", "key", 1, "nested", "key", 123},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
		{
			name: "nested array index assignment expression",
			input: `
				var mut obj = { "nested": [1, 2, 3] };
				obj["nested"][2] = 42;
			`,
			expectedConstants: []any{"nested", 1, 2, 3, "nested", 2, 42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpArray, 3),
				code.Make(code.OpHash, 2),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpIndex),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpConstant, 6),
				code.Make(code.OpIndexAssign),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilationTests(t, tests)
}

func BenchmarkIndexAssignmentExpressions(b *testing.B) {
	runCompilationBenchmarks(b, []string{
		`
			var mut arr = [1, 2, 3];
			arr[1] = 42;
		`,
		`
			var mut hash = { "key": 1 };
			hash["key"] = 99;
		`,
		`
			var mut obj = { "nested": { "key": 1 } };
			obj["nested"]["key"] = 123;
		`,
		`
			var mut obj = { "nested": [1, 2, 3] };
			obj["nested"][2] = 42;
		`,
	})
}
