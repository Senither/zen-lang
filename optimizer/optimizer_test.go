package optimizer

import (
	"testing"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
)

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func TestOptimizeRoundsZero(t *testing.T) {
	b := &compiler.Bytecode{
		Instructions: code.Instructions{byte(code.OpNull)},
		Constants:    []objects.Object{&objects.String{Value: "test"}},
	}

	result, err := OptimizeRounds(b, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != b {
		t.Error("expected same bytecode returned when rounds <= 0")
	}
}

func TestOptimizeEmptyBytecode(t *testing.T) {
	b := &compiler.Bytecode{
		Instructions: code.Instructions{},
		Constants:    []objects.Object{},
	}

	result, err := Optimize(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestOptimizeDoesNotMutateInput(t *testing.T) {
	original := code.Instructions{byte(code.OpNull), byte(code.OpPop)}
	b := &compiler.Bytecode{
		Instructions: make(code.Instructions, len(original)),
		Constants:    []objects.Object{},
	}
	copy(b.Instructions, original)

	_, err := Optimize(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, v := range original {
		if b.Instructions[i] != v {
			t.Errorf("input instructions mutated at index %d", i)
		}
	}
}

func TestDecodeInstructionsEmpty(t *testing.T) {
	infos, err := decodeInstructions(code.Instructions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(infos) != 0 {
		t.Errorf("expected 0 infos, got %d", len(infos))
	}
}

func TestDecodeInstructionsSingleOpNull(t *testing.T) {
	instructions := code.Make(code.OpNull)
	infos, err := decodeInstructions(instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(infos) != 1 {
		t.Fatalf("expected 1 info, got %d", len(infos))
	}

	if infos[0].Op != code.OpNull {
		t.Errorf("expected OpNull, got %v", infos[0].Op)
	}

	if !infos[0].Keep {
		t.Error("expected Keep to be true")
	}

	if infos[0].OldOffset != 0 {
		t.Errorf("expected OldOffset 0, got %d", infos[0].OldOffset)
	}
}

func TestDecodeInstructionsJumpMarked(t *testing.T) {
	instructions := code.Make(code.OpJump, 5)
	infos, err := decodeInstructions(instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(infos) != 1 {
		t.Fatalf("expected 1 info, got %d", len(infos))
	}

	if !infos[0].IsJump {
		t.Error("expected IsJump to be true for OpJump")
	}
}

func TestDecodeInstructionsJumpNotTruthyMarked(t *testing.T) {
	instructions := code.Make(code.OpJumpNotTruthy, 5)
	infos, err := decodeInstructions(instructions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !infos[0].IsJump {
		t.Error("expected IsJump to be true for OpJumpNotTruthy")
	}
}

func TestRunOptimizationPassesError(t *testing.T) {
	b := &BytecodeOptimization{}

	errPass := func(b *BytecodeOptimization) error {
		return &testError{"pass failed"}
	}

	err := b.runOptimizationPasses(errPass)
	if err == nil {
		t.Error("expected error from failing pass")
	}
}

func TestRunOptimizationPassesStopsOnFirstError(t *testing.T) {
	b := &BytecodeOptimization{}
	called := 0

	pass1 := func(b *BytecodeOptimization) error {
		called++
		return &testError{"fail"}
	}
	pass2 := func(b *BytecodeOptimization) error {
		called++
		return nil
	}

	_ = b.runOptimizationPasses(pass1, pass2)

	if called != 1 {
		t.Errorf("expected only 1 pass to be called, got %d", called)
	}
}

func TestAssembleInstructionsSkipsNotKept(t *testing.T) {
	b := &BytecodeOptimization{
		Infos: []InstructionInfo{
			{Op: code.OpNull, Width: 1, Keep: true},
			{Op: code.OpPop, Width: 1, Keep: false},
			{Op: code.OpNull, Width: 1, Keep: true},
		},
	}

	out := b.assembleInstructions()

	expected := append(code.Make(code.OpNull), code.Make(code.OpNull)...)
	if len(out) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(out))
	}

	for i, v := range expected {
		if out[i] != v {
			t.Errorf("expected byte %d at index %d, got %d", v, i, out[i])
		}
	}
}

func TestIsJumpTarget(t *testing.T) {
	b := &BytecodeOptimization{
		Targets: map[int]struct{}{
			5: {},
		},
	}

	if !b.isJumpTarget(5) {
		t.Error("expected offset 5 to be a jump target")
	}

	if b.isJumpTarget(3) {
		t.Error("expected offset 3 to not be a jump target")
	}
}

func TestGetKeptInstructionsInfo(t *testing.T) {
	b := &BytecodeOptimization{
		Infos: []InstructionInfo{
			{Op: code.OpNull, Keep: true},
			{Op: code.OpPop, Keep: false},
			{Op: code.OpNull, Keep: true},
			{Op: code.OpPop, Keep: true},
		},
	}

	infos, ok := b.getKeptInstructionsInfo(4, 2)
	if !ok {
		t.Fatal("expected to get 2 kept instructions")
	}

	if len(infos) != 2 {
		t.Errorf("expected 2 infos, got %d", len(infos))
	}
}

func TestGetKeptInstructionsInfoNotEnough(t *testing.T) {
	b := &BytecodeOptimization{
		Infos: []InstructionInfo{
			{Op: code.OpNull, Keep: true},
		},
	}

	_, ok := b.getKeptInstructionsInfo(1, 3)
	if ok {
		t.Error("expected false when not enough kept instructions")
	}
}

func TestSetInstructionInfoOpcode(t *testing.T) {
	b := &BytecodeOptimization{
		Infos: []InstructionInfo{
			{Op: code.OpNull, Width: 1, Keep: true},
		},
	}

	b.setInstructionInfoOpcode(0, code.OpPop, []int{})

	if b.Infos[0].Op != code.OpPop {
		t.Errorf("expected OpPop, got %v", b.Infos[0].Op)
	}
}

func TestOptimizeConstantsCopied(t *testing.T) {
	str := &objects.String{Value: "hello"}
	b := &compiler.Bytecode{
		Instructions: append(
			code.Make(code.OpConstant, 0),
			code.Make(code.OpPop)...,
		),
		Constants: []objects.Object{str},
	}

	result, err := Optimize(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Constants) == 0 {
		t.Fatal("expected constants to be copied")
	}
}
