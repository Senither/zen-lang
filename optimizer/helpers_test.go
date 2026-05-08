package optimizer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func runOptimizerTests(t *testing.T, tests []optimizerTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parse(tt.input)
			if err != nil {
				t.Fatalf("parser error: %s", err)
			}

			compiler := compiler.New(nil)
			if err := compiler.Compile(program); err != nil {
				t.Fatalf("compiler error: %s", err)
			}

			bytecode, err := Optimize(compiler.Bytecode())
			if err != nil {
				t.Fatalf("optimizer error: %s", err)
			}

			if err := testOptimizerInstructions(tt.expectedInstructions, bytecode.Instructions); err != nil {
				t.Fatalf("instruction test failed: %s", err)
			}

			if err := testOptimizerConstants(t, tt.expectedConstants, bytecode.Constants); err != nil {
				t.Fatalf("constants test failed: %s", err)
			}
		})
	}
}

func runOptimizerBenchmarks(b *testing.B, inputs []string) {
	b.Helper()

	compiledBytecode := make([]*compiler.Bytecode, len(inputs))
	for i, input := range inputs {
		program, err := parse(input)
		if err != nil {
			b.Fatalf("parser error: %s", err)
		}

		compiler := compiler.New(nil)
		if err := compiler.Compile(program); err != nil {
			b.Fatalf("compiler error: %s", err)
		}

		compiledBytecode[i] = compiler.Bytecode()
	}

	for i, bytecode := range compiledBytecode {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			for b.Loop() {
				if _, err := Optimize(bytecode); err != nil {
					b.Fatalf("optimizer error: %s", err)
				}
			}
		})
	}
}

func parse(input string) (*ast.Program, error) {
	l := lexer.New(input)
	p := parser.New(l, nil)

	program := p.ParseProgram()
	if len(p.Errors()) == 0 {
		return program, nil
	}

	var buf strings.Builder
	for _, msg := range p.Errors() {
		fmt.Fprintf(&buf, "%s\n", msg.String())
	}

	return nil, fmt.Errorf("parser errors encountered\n%s", buf.String())
}

func testOptimizerInstructions(expected []code.Instructions, actual code.Instructions) error {
	combined := flattenOptimizerInstructions(expected)
	if len(actual) != len(combined) {
		return fmt.Errorf("wrong instructions length.\n\twant:\n%s\n\tgot:\n%s", combined, actual)
	}

	for i, ins := range combined {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\n\tinstruction: %d\n\twant:\n%s\n\tgot:\n%s", i, ins, combined, actual)
		}
	}

	return nil
}

func flattenOptimizerInstructions(parts []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, part := range parts {
		out = append(out, part...)
	}

	return out
}

func testOptimizerConstants(t *testing.T, expected []any, actual []objects.Object) error {
	t.Helper()

	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got %d, want %d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			if err := objects.AssertInteger(int64(constant), actual[i]); err != nil {
				return fmt.Errorf("constant %d - integer assertion failed: %s", i, err)
			}
		case float64:
			if err := objects.AssertFloat(constant, actual[i]); err != nil {
				return fmt.Errorf("constant %d - float assertion failed: %s", i, err)
			}
		case string:
			if err := objects.AssertString(constant, actual[i]); err != nil {
				return fmt.Errorf("constant %d - string assertion failed: %s", i, err)
			}
		case []code.Instructions:
			if err := testOptimizerCodeInstructions(constant, actual[i]); err != nil {
				return fmt.Errorf("constant %d - code instructions assertion failed: %s", i, err)
			}

		default:
			return fmt.Errorf("unknown constant type %T", constant)
		}
	}

	return nil
}

func testOptimizerCodeInstructions(expected []code.Instructions, actual objects.Object) error {
	fn, ok := actual.(*objects.CompiledFunction)
	if !ok {
		return fmt.Errorf("object is not CompiledFunction. got %T (%+v)", actual, actual)
	}

	if err := testOptimizerInstructions(expected, fn.Instructions()); err != nil {
		return fmt.Errorf("instructions do not match the CompiledFunction instructions: %s", err)
	}

	return nil
}

func concatInstructions(parts ...[]byte) code.Instructions {
	out := code.Instructions{}

	for _, p := range parts {
		out = append(out, p...)
	}

	return out
}

func assertIntSet(t *testing.T, got map[int]struct{}, want []int) {
	t.Helper()

	if got == nil {
		got = map[int]struct{}{}
	}

	wantSet := map[int]struct{}{}
	for _, v := range want {
		wantSet[v] = struct{}{}
	}

	if !reflect.DeepEqual(got, wantSet) {
		t.Fatalf("unexpected set:\ngot:\n\t%v\nwant\n\t%v", got, wantSet)
	}
}

func TestEqualConstants(t *testing.T) {
	t.Run("len mismatch", func(t *testing.T) {
		a := []objects.Object{&objects.Integer{Value: 1}}
		b := []objects.Object{&objects.Integer{Value: 1}, &objects.Integer{Value: 2}}

		if equalConstants(a, b) {
			t.Fatalf("expected constants to not be equal")
		}
	})

	t.Run("different inspect values", func(t *testing.T) {
		a := []objects.Object{&objects.Integer{Value: 1}, &objects.String{Value: "hello"}}
		b := []objects.Object{&objects.Integer{Value: 2}, &objects.String{Value: "hello"}}

		if equalConstants(a, b) {
			t.Fatalf("expected constants to not be equal")
		}
	})

	t.Run("same inspect values", func(t *testing.T) {
		a := []objects.Object{&objects.Integer{Value: 1}, &objects.String{Value: "hello"}}
		b := []objects.Object{&objects.Integer{Value: 1}, &objects.String{Value: "hello"}}

		if !equalConstants(a, b) {
			t.Fatalf("expected constants to be equal")
		}
	})
}

func TestFindJumpTargets(t *testing.T) {
	infos := []InstructionInfo{
		{Op: code.OpJump, IsJump: true, Operands: []int{10}},
		{Op: code.OpJumpNotTruthy, IsJump: true, Operands: []int{42}},
		{Op: code.OpJump, IsJump: true, Operands: nil},
		{Op: code.OpAdd, IsJump: false, Operands: []int{999}},
	}

	assertIntSet(t, findJumpTargets(infos), []int{10, 42})
}

func TestFindUsedGlobals(t *testing.T) {
	got := findUsedGlobals([]InstructionInfo{
		{Op: code.OpGetGlobal, Operands: []int{1}},
	}, []objects.Object{
		&objects.CompiledFunction{OpcodeInstructions: concatInstructions(
			code.Make(code.OpGetGlobal, 2),
			code.Make(code.OpReturn),
		)},
		// Non-function constant should be ignored
		&objects.CompiledFunction{OpcodeInstructions: code.Instructions{255}},
		&objects.Integer{Value: 123},
	})

	assertIntSet(t, got, []int{1, 2})
}

func TestFindUsedGlobalsInInstructions(t *testing.T) {
	infos := []InstructionInfo{
		{Op: code.OpGetGlobal, Operands: []int{1}},
		{Op: code.OpIncGlobal, Operands: []int{2}},
		{Op: code.OpDecGlobal, Operands: []int{3}},
		{Op: code.OpSetGlobal, Operands: []int{99}}, // ignored
		{Op: code.OpGetGlobal, Operands: nil},       // ignored
	}

	assertIntSet(t, findUsedGlobalsInInstructions(infos), []int{1, 2, 3})
}

func TestFindChangedGlobals(t *testing.T) {
	// Top-level changes:
	// - global 5 set twice => changed
	// - global 7 set once, but nested fn also mutates once => changed
	// - global 9 set once => not changed
	topLevelInfos := []InstructionInfo{
		{Op: code.OpSetGlobal, Operands: []int{5}},
		{Op: code.OpSetGlobal, Operands: []int{5}},
		{Op: code.OpSetGlobal, Operands: []int{7}},
		{Op: code.OpSetGlobal, Operands: []int{9}},
	}

	nestedFn := &objects.CompiledFunction{OpcodeInstructions: concatInstructions(
		code.Make(code.OpIncGlobal, 7),
		code.Make(code.OpReturn),
	)}

	got := findChangedGlobals(topLevelInfos, []objects.Object{nestedFn})
	assertIntSet(t, got, []int{5, 7})
}

func TestFindChangedGlobalsInInstructions(t *testing.T) {
	infos := []InstructionInfo{
		{Op: code.OpSetGlobal, Operands: []int{1}},
		{Op: code.OpIncGlobal, Operands: []int{1}},
		{Op: code.OpDecGlobal, Operands: []int{2}},
		{Op: code.OpGetGlobal, Operands: []int{3}}, // ignored
		{Op: code.OpSetGlobal, Operands: nil},      // ignored
	}

	got := findChangedGlobalsInInstructions(infos)
	if got[1] != 2 {
		t.Fatalf("expected global 1 count 2, got %d", got[1])
	}

	if got[2] != 1 {
		t.Fatalf("expected global 2 count 1, got %d", got[2])
	}

	if _, ok := got[3]; ok {
		t.Fatalf("did not expect global 3")
	}
}

func TestFindPrevKeptInstructionIndex(t *testing.T) {
	infos := []InstructionInfo{
		{Keep: false},
		{Keep: true},
		{Keep: false},
		{Keep: false},
		{Keep: true},
	}

	for i, want := range []int{-1, -1, 1, 1, 1, 4} {
		t.Run("index "+strconv.Itoa(i), func(t *testing.T) {
			index := findPrevKeptInstructionIndex(infos, i)
			if index != want {
				t.Fatalf("unexpected value for index of %d: got %d want %d", i, index, want)
			}
		})
	}
}

func TestFindNextKeptInstructionIndex(t *testing.T) {
	infos := []InstructionInfo{
		{Keep: false},
		{Keep: true},
		{Keep: false},
		{Keep: false},
		{Keep: true},
	}

	for i, want := range []int{1, 4, 4, 4, -1, -1} {
		t.Run("index "+strconv.Itoa(i), func(t *testing.T) {
			index := findNextKeptInstructionIndex(infos, i)
			if index != want {
				t.Fatalf("unexpected value for index of %d: got %d want %d", i, index, want)
			}
		})
	}
}

func TestFindJumpTargetsFromKeptInstructions(t *testing.T) {
	infos := []InstructionInfo{
		{Keep: false, IsJump: true, Operands: []int{9}},
		{Keep: true, IsJump: true, Operands: []int{10}},
		{Keep: false, IsJump: true, Operands: []int{42}},
		{Keep: true, IsJump: false, Operands: []int{123}},
		{Keep: true, IsJump: true, Operands: []int{999}},
	}

	kept := findJumpTargetsFromKeptInstructions(infos)

	if len(kept) != 2 {
		t.Fatalf("unexpected number of kept targets: got %d want 2", len(kept))
	}

	if _, ok := kept[10]; !ok {
		t.Fatalf("expected target 10 to be kept")
	}

	if _, ok := kept[999]; !ok {
		t.Fatalf("expected target 999 to be kept")
	}
}

func TestIsNoOpJump(t *testing.T) {
	infos := []InstructionInfo{
		{Keep: true, IsJump: true, Operands: []int{3}},
		{Keep: false},
		{Keep: false},
		{Keep: true},
		{Keep: false},
		{Keep: true},
	}

	results := map[int]bool{
		0: false,
		1: false,
		2: false,
		3: true,
		4: false,
		5: false,
	}

	for idx, want := range results {
		t.Run("index "+strconv.Itoa(idx), func(t *testing.T) {
			if isNoOpJump(infos, 0, idx) != want {
				t.Fatalf("unexpected value for index %d: got %v want %v", idx, !want, want)
			}
		})
	}
}

func TestIsWhileJumpPattern(t *testing.T) {
	t.Run("returns false when hasElseJump is true", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJump, IsJump: true, Operands: []int{0}, OldOffset: 0},
		}

		if isWhileJumpPattern(infos, 0, 0, 0, true) {
			t.Fatalf("expected false when hasElseJump is true")
		}
	})

	t.Run("returns false when prevTargetIdx <= jumpNotTruthyIdx", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJump, IsJump: true, Operands: []int{0}, OldOffset: 0},
		}

		// prevTargetIdx=0, jumpNotTruthyIdx=0 (equal)
		if isWhileJumpPattern(infos, 0, 0, 0, false) {
			t.Fatalf("expected false when prevTargetIdx == jumpNotTruthyIdx")
		}

		// prevTargetIdx=0, jumpNotTruthyIdx=1 (less than)
		if isWhileJumpPattern(infos, 1, 1, 0, false) {
			t.Fatalf("expected false when prevTargetIdx < jumpNotTruthyIdx")
		}
	})

	t.Run("returns false when prevTarget is not OpJump", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJumpNotTruthy, IsJump: true, Operands: nil, OldOffset: 20},
			{Op: code.OpAdd, IsJump: false, Operands: []int{5}, OldOffset: 10},
		}

		// prevTargetIdx=1 > jumpNotTruthyIdx=0, but prevTarget at index 1 is OpAdd, not OpJump
		if isWhileJumpPattern(infos, 0, 100, 1, false) {
			t.Fatalf("expected false when prevTarget is not OpJump")
		}
	})

	t.Run("returns false when prevTarget has no operands", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJumpNotTruthy, IsJump: true, Operands: []int{50}, OldOffset: 10},
			{Op: code.OpJump, IsJump: true, Operands: nil, OldOffset: 20},
		}

		// prevTargetIdx=1 > jumpNotTruthyIdx=0, prevTarget has OpJump but no operands
		if isWhileJumpPattern(infos, 0, 100, 1, false) {
			t.Fatalf("expected false when prevTarget has no operands")
		}
	})

	t.Run("returns false when prevTarget jump target > jumpNotTruthyIdx OldOffset", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJumpNotTruthy, IsJump: true, Operands: []int{100}, OldOffset: 10},
			{Op: code.OpJump, IsJump: true, Operands: []int{50}, OldOffset: 20},
		}

		// prevTargetIdx=1 > jumpNotTruthyIdx=0
		// prevTarget.Operands[0]=50 > jumpNotTruthyIdx.OldOffset=10
		if isWhileJumpPattern(infos, 0, 100, 1, false) {
			t.Fatalf("expected false when prevTarget jump target > jumpNotTruthyIdx OldOffset")
		}
	})

	t.Run("returns true for valid while loop pattern", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJumpNotTruthy, IsJump: true, Operands: []int{100}, OldOffset: 20},
			{Op: code.OpJump, IsJump: true, Operands: []int{10}, OldOffset: 5},
		}

		// prevTargetIdx=1 > jumpNotTruthyIdx=0
		// prevTarget.Op == OpJump
		// prevTarget has operands
		// prevTarget.Operands[0]=10 <= jumpNotTruthyIdx.OldOffset=20
		if !isWhileJumpPattern(infos, 0, 100, 1, false) {
			t.Fatalf("expected true for valid while loop pattern")
		}
	})

	t.Run("returns true when prevTarget jump target equals jumpNotTruthyIdx OldOffset", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpJumpNotTruthy, IsJump: true, Operands: []int{100}, OldOffset: 20},
			{Op: code.OpJump, IsJump: true, Operands: []int{20}, OldOffset: 5},
		}

		// prevTarget.Operands[0]=20 == jumpNotTruthyIdx.OldOffset=20
		if !isWhileJumpPattern(infos, 0, 100, 1, false) {
			t.Fatalf("expected true when jump target equals OldOffset")
		}
	})

	t.Run("returns true in complex multi-instruction scenario", func(t *testing.T) {
		infos := []InstructionInfo{
			{Op: code.OpConstant, Operands: []int{0}, OldOffset: 0},
			{Op: code.OpSetLocal, Operands: []int{0}, OldOffset: 3},
			{Op: code.OpGetLocal, Operands: []int{0}, OldOffset: 5},
			{Op: code.OpJumpNotTruthy, IsJump: true, Operands: []int{50}, OldOffset: 11},
			{Op: code.OpJump, IsJump: true, Operands: []int{5}, OldOffset: 8},
		}

		// prevTargetIdx=4 > jumpNotTruthyIdx=3
		// prevTarget.Op == OpJump
		// prevTarget.Operands[0]=5 <= jumpNotTruthyIdx.OldOffset=11
		if !isWhileJumpPattern(infos, 3, 100, 4, false) {
			t.Fatalf("expected true in complex scenario")
		}
	})
}

func TestComputeGlobalSwaps(t *testing.T) {
	t.Run("swaps only for non-reassigned OpSetGlobal preceded by constant", func(t *testing.T) {
		instructions := concatInstructions(
			code.Make(code.OpConstant, 0),
			code.Make(code.OpSetGlobal, 1),
			code.Make(code.OpConstant, 1),
			code.Make(code.OpSetGlobal, 2),
			code.Make(code.OpConstant, 2),
			code.Make(code.OpSetGlobal, 3),
			// reassigned global 3, so swaps should not include it at all
			code.Make(code.OpConstant, 3),
			code.Make(code.OpSetGlobal, 3),
		)

		swaps := computeGlobalSwaps(instructions, nil)
		if swaps == nil {
			t.Fatalf("expected swaps map")
		}

		s1, ok := swaps[1]
		if !ok || s1.Op != code.OpConstant || len(s1.Operands) != 1 || s1.Operands[0] != 0 {
			t.Fatalf("unexpected swap for global 1: %+v", s1)
		}

		s2, ok := swaps[2]
		if !ok || s2.Op != code.OpConstant || len(s2.Operands) != 1 || s2.Operands[0] != 1 {
			t.Fatalf("unexpected swap for global 2: %+v", s2)
		}

		if _, ok := swaps[3]; ok {
			t.Fatalf("did not expect swap for reassigned global 3")
		}
	})

	t.Run("no swap when OpSetGlobal is first instruction", func(t *testing.T) {
		instructions := concatInstructions(
			code.Make(code.OpSetGlobal, 1),
		)

		swaps := computeGlobalSwaps(instructions, nil)
		if len(swaps) != 0 {
			t.Fatalf("expected no swaps, got %v", swaps)
		}
	})

	t.Run("no swap when previous instruction has no operands (OpNull)", func(t *testing.T) {
		instructions := concatInstructions(
			code.Make(code.OpNull),
			code.Make(code.OpSetGlobal, 10),
		)

		swaps := computeGlobalSwaps(instructions, nil)
		if len(swaps) != 0 {
			t.Fatalf("expected no swaps, got %v", swaps)
		}
	})
}

func TestStackDelta(t *testing.T) {
	tests := []struct {
		name string
		info InstructionInfo
		want int
	}{
		{
			name: "constant",
			info: InstructionInfo{Op: code.OpConstant, Operands: []int{0}},
			want: 1,
		},
		{
			name: "null",
			info: InstructionInfo{Op: code.OpNull},
			want: 1,
		},
		{
			name: "true",
			info: InstructionInfo{Op: code.OpTrue},
			want: 1,
		},
		{
			name: "false",
			info: InstructionInfo{Op: code.OpFalse},
			want: 1,
		},
		{
			name: "array with operand",
			info: InstructionInfo{Op: code.OpArray, Operands: []int{3}},
			want: -2,
		},
		{
			name: "hash with operand",
			info: InstructionInfo{Op: code.OpHash, Operands: []int{4}},
			want: -3,
		},
		{
			name: "array missing operands",
			info: InstructionInfo{Op: code.OpArray},
			want: 0,
		},
		{
			name: "binary op",
			info: InstructionInfo{Op: code.OpAdd},
			want: -1,
		},
		{
			name: "compare op",
			info: InstructionInfo{Op: code.OpEqual},
			want: -1,
		},
		{
			name: "prefix op",
			info: InstructionInfo{Op: code.OpMinus},
			want: 0,
		},
		{
			name: "pop",
			info: InstructionInfo{Op: code.OpPop},
			want: -1,
		},
		{
			name: "set global",
			info: InstructionInfo{Op: code.OpSetGlobal, Operands: []int{0}},
			want: -1,
		},
		{
			name: "set local",
			info: InstructionInfo{Op: code.OpSetLocal, Operands: []int{0}},
			want: -1,
		},
		{
			name: "default",
			info: InstructionInfo{Op: code.OpReturn},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stackDelta(&tt.info)
			if got != tt.want {
				t.Fatalf("unexpected delta: got %d want %d", got, tt.want)
			}
		})
	}
}

func TestStackDeltaBinaryAndComparisonOp(t *testing.T) {
	tests := []code.Opcode{
		code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod, code.OpIndex,
		code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual,
	}

	for _, op := range tests {
		def, err := code.Lookup(op)
		if err != nil {
			t.Fatalf("unknown op: %v", op)
		}

		t.Run(def.Name, func(t *testing.T) {
			info := InstructionInfo{Op: op}
			delta := stackDelta(&info)
			if delta != -1 {
				t.Fatalf("unexpected delta for op %v: got %d want -1", op, delta)
			}
		})
	}
}

func TestResolveTargetInstructionIndex(t *testing.T) {
	t.Run("found in offsetToIndex map", func(t *testing.T) {
		infos := []InstructionInfo{
			{OldOffset: 0, Width: 3},
			{OldOffset: 3, Width: 3},
		}
		offsetToIndex := map[int]int{0: 0, 3: 1}

		idx, ok := resolveTargetInstructionIndex(infos, offsetToIndex, 3)
		if !ok {
			t.Fatalf("expected offset index to exist")
		}

		if idx != 1 {
			t.Fatalf("unexpected index: got %d want 1", idx)
		}
	})

	t.Run("not found in map but matches end offset", func(t *testing.T) {
		infos := []InstructionInfo{
			{OldOffset: 0, Width: 3},
			{OldOffset: 3, Width: 4},
		}
		offsetToIndex := map[int]int{0: 0, 3: 1}

		// end offset = 3 + 4 = 7
		idx, ok := resolveTargetInstructionIndex(infos, offsetToIndex, 7)
		if !ok {
			t.Fatalf("expected offset index to exist")
		}

		if idx != len(infos) {
			t.Fatalf("unexpected index: got %d want %d", idx, len(infos))
		}
	})

	t.Run("empty infos returns not found", func(t *testing.T) {
		idx, ok := resolveTargetInstructionIndex([]InstructionInfo{}, map[int]int{}, 10)
		if ok {
			t.Fatalf("expected offset to not be found")
		}

		if idx != -1 {
			t.Fatalf("unexpected index: got %d want -1", idx)
		}
	})

	t.Run("nil map and empty infos returns not found", func(t *testing.T) {
		idx, ok := resolveTargetInstructionIndex([]InstructionInfo{}, nil, 0)
		if ok {
			t.Fatalf("expected offset to not be found")
		}

		if idx != -1 {
			t.Fatalf("unexpected index: got %d want -1", idx)
		}
	})

	t.Run("offset not in map and not at end returns not found", func(t *testing.T) {
		infos := []InstructionInfo{
			{OldOffset: 0, Width: 3},
			{OldOffset: 3, Width: 3},
		}
		offsetToIndex := map[int]int{0: 0, 3: 1}

		idx, ok := resolveTargetInstructionIndex(infos, offsetToIndex, 99)
		if ok {
			t.Fatalf("expected offset to not be found")
		}

		if idx != -1 {
			t.Fatalf("unexpected index: got %d want -1", idx)
		}
	})

	t.Run("single instruction end offset resolves to len(infos)", func(t *testing.T) {
		infos := []InstructionInfo{
			{OldOffset: 5, Width: 2},
		}
		offsetToIndex := map[int]int{5: 0}

		// end offset = 5 + 2 = 7
		idx, ok := resolveTargetInstructionIndex(infos, offsetToIndex, 7)
		if !ok {
			t.Fatalf("expected offset index to exist")
		}

		if idx != 1 {
			t.Fatalf("unexpected index: got %d want 1", idx)
		}
	})
}

func TestEvaluateKnownConditionTruthiness(t *testing.T) {
	t.Run("OpTrue is truthy", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpTrue},
			},
		}

		truthy, ok := evaluateKnownConditionTruthiness(b, 0)
		if !ok {
			t.Fatalf("expected condition to be known")
		}

		if !truthy {
			t.Fatalf("expected OpTrue to be truthy")
		}
	})

	t.Run("OpFalse is not truthy", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpFalse},
			},
		}

		truthy, ok := evaluateKnownConditionTruthiness(b, 0)
		if !ok {
			t.Fatalf("expected condition to be known")
		}

		if truthy {
			t.Fatalf("expected OpFalse to not be truthy")
		}
	})

	t.Run("OpNull is not truthy", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpNull},
			},
		}

		truthy, ok := evaluateKnownConditionTruthiness(b, 0)
		if !ok {
			t.Fatalf("expected condition to be known")
		}

		if truthy {
			t.Fatalf("expected OpNull to not be truthy")
		}
	})

	t.Run("OpConstant with integer is truthy", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpConstant, Operands: []int{0}},
			},
			Constants: []objects.Object{
				&objects.Integer{Value: 42},
			},
		}

		truthy, ok := evaluateKnownConditionTruthiness(b, 0)
		if !ok {
			t.Fatalf("expected condition to be known")
		}

		if !truthy {
			t.Fatalf("expected non-zero integer to be truthy")
		}
	})

	t.Run("OpConstant with no operands is unknown", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpConstant},
			},
			Constants: []objects.Object{},
		}

		_, ok := evaluateKnownConditionTruthiness(b, 0)
		if ok {
			t.Fatalf("expected condition to be unknown")
		}
	})

	t.Run("OpConstant with out-of-bounds index is unknown", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpConstant, Operands: []int{5}},
			},
			Constants: []objects.Object{
				&objects.Integer{Value: 1},
			},
		}

		_, ok := evaluateKnownConditionTruthiness(b, 0)
		if ok {
			t.Fatalf("expected condition to be unknown")
		}
	})

	t.Run("unknown opcode is unknown", func(t *testing.T) {
		b := &BytecodeOptimization{
			Infos: []InstructionInfo{
				{Op: code.OpAdd},
			},
		}

		_, ok := evaluateKnownConditionTruthiness(b, 0)
		if ok {
			t.Fatalf("expected condition to be unknown for OpAdd")
		}
	})
}

func TestCanRemoveIndexes(t *testing.T) {
	t.Run("empty toRemove returns true", func(t *testing.T) {
		infos := []InstructionInfo{
			{Keep: true, OldOffset: 0},
		}

		if !canRemoveIndexes(infos, map[int][]int{}, map[int]struct{}{}) {
			t.Fatalf("expected true for empty toRemove")
		}
	})

	t.Run("toRemove index not kept is skipped", func(t *testing.T) {
		infos := []InstructionInfo{
			{Keep: false, OldOffset: 0},
		}

		if !canRemoveIndexes(infos, map[int][]int{}, map[int]struct{}{0: {}}) {
			t.Fatalf("expected true when removed instruction is not kept")
		}
	})

	t.Run("kept instruction with no incoming sources can be removed", func(t *testing.T) {
		infos := []InstructionInfo{
			{Keep: true, OldOffset: 10},
		}

		if !canRemoveIndexes(infos, map[int][]int{}, map[int]struct{}{0: {}}) {
			t.Fatalf("expected true when no incoming sources")
		}
	})

	t.Run("returns true when all sources are also being removed", func(t *testing.T) {
		infos := []InstructionInfo{
			{Keep: true, OldOffset: 0},
			{Keep: true, OldOffset: 5},
		}

		// instruction at index 1 (offset 5) has source at index 0
		incoming := map[int][]int{5: {0}}
		toRemove := map[int]struct{}{0: {}, 1: {}}

		if !canRemoveIndexes(infos, incoming, toRemove) {
			t.Fatalf("expected true when all sources are also being removed")
		}
	})

	t.Run("returns false when kept source is not being removed", func(t *testing.T) {
		infos := []InstructionInfo{
			{Keep: true, OldOffset: 0},
			{Keep: true, OldOffset: 5},
		}

		// instruction at index 1 (offset 5) has source at index 0, but index 0 is not in toRemove
		incoming := map[int][]int{5: {0}}
		toRemove := map[int]struct{}{1: {}}

		if canRemoveIndexes(infos, incoming, toRemove) {
			t.Fatalf("expected false when source is kept and not being removed")
		}
	})

	t.Run("returns true when source is not kept", func(t *testing.T) {
		infos := []InstructionInfo{
			{Keep: false, OldOffset: 0},
			{Keep: true, OldOffset: 5},
		}

		// source at index 0 is not kept, so it is skipped
		incoming := map[int][]int{5: {0}}
		toRemove := map[int]struct{}{1: {}}

		if !canRemoveIndexes(infos, incoming, toRemove) {
			t.Fatalf("expected true when source is not kept")
		}
	})
}

func TestImmutableConstantReuseKey(t *testing.T) {
	tests := []struct {
		name string
		obj  objects.Object
		want string
		ok   bool
	}{
		{
			name: "null",
			obj:  &objects.Null{},
			want: "null",
			ok:   true,
		},
		{
			name: "boolean true",
			obj:  &objects.Boolean{Value: true},
			want: "bool:1",
			ok:   true,
		},
		{
			name: "boolean false",
			obj:  &objects.Boolean{Value: false},
			want: "bool:0",
			ok:   true,
		},
		{
			name: "integer",
			obj:  &objects.Integer{Value: 99},
			want: "int:99",
			ok:   true,
		},
		{
			name: "negative integer",
			obj:  &objects.Integer{Value: -7},
			want: "int:-7",
			ok:   true,
		},
		{
			name: "float",
			obj:  &objects.Float{Value: 1.5},
			want: "float:3ff8000000000000",
			ok:   true,
		},
		{
			name: "string",
			obj:  &objects.String{Value: "hello"},
			want: "string:5:hello",
			ok:   true,
		},
		{
			name: "empty string",
			obj:  &objects.String{Value: ""},
			want: "string:0:",
			ok:   true,
		},
		{
			name: "unsupported type returns empty key",
			obj:  &objects.Array{},
			want: "",
			ok:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := immutableConstantReuseKey(tt.obj)
			if ok != tt.ok {
				t.Fatalf("unexpected ok: got %v want %v", ok, tt.ok)
			}

			if got != tt.want {
				t.Fatalf("unexpected key: got %q want %q", got, tt.want)
			}
		})
	}
}
