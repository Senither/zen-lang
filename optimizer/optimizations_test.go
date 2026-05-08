package optimizer

import (
	"testing"

	"github.com/senither/zen-lang/code"
)

type optimizerTestCase struct {
	name                 string
	input                string
	expectedConstants    []any
	expectedInstructions []code.Instructions
}

func TestOptimizerUnfoldNonReassignedVariables(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "replaces non-reassigned global reads with constant",
			input:             "var a = 42; a;",
			expectedConstants: []any{42},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "replaces non-reassigned global reads with boolean",
			input:             "var flag = true; flag;",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkUnfoldNonReassignedVariables(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"var a = 42; a;",
		"var flag = true; flag;",
	})
}

func TestOptimizerRemoveUnusedVariableInitializations(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "removes unused scalar global initialization",
			input:             "var a = 42; 99;",
			expectedConstants: []any{99},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "removes unused array global initialization",
			input:             "var arr = [1, 2, 3]; 5;",
			expectedConstants: []any{5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkRemoveUnusedVariableInitializations(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"var a = 42; 99;",
		"var arr = [1, 2, 3]; 5;",
	})
}

func TestOptimizerDeleteArrayOrHashInitializerIndirect(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "drops unused hash initializer chain via parent pass",
			input:             "var payload = {\"a\": 1, \"b\": 2}; 7;",
			expectedConstants: []any{7},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkDeleteArrayOrHashInitializerIndirect(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"var payload = {\"a\": 1, \"b\": 2}; 7;",
	})
}

func TestOptimizerPreCalculateNumberConstants(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "folds integer addition",
			input:             "1 + 2",
			expectedConstants: []any{3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds integer subtraction",
			input:             "5 - 2",
			expectedConstants: []any{3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds integer multiplication",
			input:             "3 * 4",
			expectedConstants: []any{12},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds integer division",
			input:             "12 / 4",
			expectedConstants: []any{3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds integer power",
			input:             "2 ^ 3",
			expectedConstants: []any{8},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds integer modulo",
			input:             "12 % 4",
			expectedConstants: []any{0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds floating-point addition",
			input:             "1.5 + 2.5",
			expectedConstants: []any{4.0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds floating-point subtraction",
			input:             "5.5 - 2.5",
			expectedConstants: []any{3.0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds floating-point multiplication",
			input:             "2.0 * 3.0",
			expectedConstants: []any{6.0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds floating-point division",
			input:             "6.0 / 2.0",
			expectedConstants: []any{3.0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds floating-point power",
			input:             "2.0 ^ 3.0",
			expectedConstants: []any{8.0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds floating-point modulo",
			input:             "12.0 % 4.0",
			expectedConstants: []any{0.0},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds nested expressions",
			input:             "(1 + 2) * (3 - 4)",
			expectedConstants: []any{-3},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds mixed integer and floating-point expressions",
			input:             "1 + 2.5 * 3",
			expectedConstants: []any{8.5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkPreCalculateNumberConstants(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"1 + 2",
		"5 - 2",
		"3 * 4",
		"12 / 4",
		"2 ^ 3",
		"12 % 4",
	})
}

func TestOptimizerConstantFoldComparisonLogicalOps(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "folds greater-than comparison",
			input:             "5 > 10",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "folds mixed numeric equality",
			input:             "5 == 5.0",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "keeps short-circuit logical and form",
			input:             "true && false",
			expectedConstants: []any{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpAnd),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkConstantFoldComparisonLogicalOps(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"5 > 10",
		"5 == 5.0",
		"true && false",
	})
}

func TestOptimizerConcatenateStringableConstants(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "concatenates string and integer",
			input:             "\"Value=\" + 42",
			expectedConstants: []any{"Value=42"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "concatenates integer and string",
			input:             "42 + \" items\"",
			expectedConstants: []any{"42 items"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "concatenates string and float",
			input:             "\"Pi=\" + 3.14",
			expectedConstants: []any{"Pi=3.14"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "concatenates float and string",
			input:             "3.14 + \" is Pi\"",
			expectedConstants: []any{"3.14 is Pi"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "concatenates string and string",
			input:             "\"Hello\" + \" World\"",
			expectedConstants: []any{"Hello World"},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkConcatenateStringableConstants(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"\"Value=\" + 42",
		"42 + \" items\"",
		"\"Pi=\" + 3.14",
		"3.14 + \" is Pi\"",
		"\"Hello\" + \" World\"",
	})
}

func TestOptimizerRemoveUnusedGettersAfterAssignments(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "keeps assignment result getter/pop when at instruction tail",
			input:             "var mut a = 1; a = 2;",
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
	}

	runOptimizerTests(t, tests)
}

func BenchmarkRemoveUnusedGettersAfterAssignments(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"var mut a = 1; a = 2;",
	})
}

func TestOptimizerReplaceIncrementsAndDecrementsWithDirectOperations(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "replaces increment assignment with OpIncGlobal",
			input:             "var mut a = 1; a = a + 1;",
			expectedConstants: []any{1},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpIncGlobal, 0),
				code.Make(code.OpPop),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "replaces decrement assignment with OpDecGlobal",
			input:             "var mut a = 10; a = a - 1;",
			expectedConstants: []any{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpDecGlobal, 0),
				code.Make(code.OpPop),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkReplaceIncrementsAndDecrementsWithDirectOperations(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"var mut a = 1; a = a + 1;",
		"var mut a = 10; a = a - 1;",
	})
}

func TestOptimizerCallBuiltinsWithKnownConstantParameters(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "folds builtin len call with constant string",
			input:             "len(\"hello\")",
			expectedConstants: []any{5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkCallBuiltinsWithKnownConstantParameters(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"len(\"hello\")",
	})
}

func TestOptimizerRemoveInstructionsAfterReturn(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:  "removes unreachable instructions after return in function",
			input: "func() { return 42; 99; };",
			expectedConstants: []any{
				42,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkRemoveInstructionsAfterReturn(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"func() { return 42; 99; };",
	})
}

func TestOptimizerRemoveRedundantJumpInstructions(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "removes jumps for constant true condition",
			input:             "if (true) { 5 } else { 10 }",
			expectedConstants: []any{5},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			name:              "removes jumps for constant false condition",
			input:             "if (false) { 5 } else { 10 }",
			expectedConstants: []any{10},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkRemoveRedundantJumpInstructions(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"if (true) { 5 } else { 10 }",
		"if (false) { 5 } else { 10 }",
	})
}

func TestOptimizerRemovePopAtBeginningOfInstructions(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "retains leading null/pop when jump target is preserved",
			input:             "if (false) { 1 }; 2;",
			expectedConstants: []any{2},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpNull),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkRemovePopAtBeginningOfInstructions(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"if (false) { 1 }; 2;",
	})
}

func TestOptimizerReorganizeConstantReferences(t *testing.T) {
	tests := []optimizerTestCase{
		{
			name:              "deduplicates immutable constants and rewrites references",
			input:             "42; 42; 99; 42;",
			expectedConstants: []any{42, 99},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runOptimizerTests(t, tests)
}

func BenchmarkReorganizeConstantReferences(b *testing.B) {
	runOptimizerBenchmarks(b, []string{
		"42; 42; 99; 42;",
	})
}
