package vm

import (
	"reflect"
	"testing"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type MockCompiledInstructionsObject struct {
	instructions code.Instructions
}

func (m *MockCompiledInstructionsObject) Type() objects.ObjectType {
	return "CompiledInstructionsObject"
}
func (m *MockCompiledInstructionsObject) Inspect() string {
	return "MockCompiledInstructionsObject"
}
func (m *MockCompiledInstructionsObject) Instructions() code.Instructions {
	return m.instructions
}

func TestNewFrame(t *testing.T) {
	mockInstructions := &MockCompiledInstructionsObject{instructions: code.Instructions{0x01, 0x02, 0x03}}
	closure := &objects.Closure{Fn: &objects.CompiledFunction{OpcodeInstructions: mockInstructions.instructions}}

	frame := NewFrame(closure, 0)

	if frame.ip != -1 {
		t.Errorf("Expected ip to be -1, got %d", frame.ip)
	}

	if frame.obj != closure {
		t.Errorf("Expected obj to be the provided Closure, got %v", frame.obj)
	}

	instructions := frame.Instructions()
	expectedInstructions := code.Instructions{0x01, 0x02, 0x03}

	if !reflect.DeepEqual(instructions, expectedInstructions) {
		t.Errorf("Expected instructions to be %v, got %v", expectedInstructions, instructions)
	}
}
