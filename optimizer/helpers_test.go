package optimizer

import (
	"reflect"
	"testing"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

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
