package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpPop, []int{}, []byte{byte(OpPop)}},
		// Jumps
		{OpJump, []int{1024}, []byte{byte(OpJump), 4, 0}},
		{OpJumpNotTruthy, []int{1024}, []byte{byte(OpJumpNotTruthy), 4, 0}},
		// Globals
		{OpSetGlobal, []int{255}, []byte{byte(OpSetGlobal), 0, 255}},
		{OpGetGlobal, []int{255}, []byte{byte(OpGetGlobal), 0, 255}},
		// Locals
		{OpSetLocal, []int{255}, []byte{byte(OpSetLocal), 255}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		// Arithmetic
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMul, []int{}, []byte{byte(OpMul)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpPow, []int{}, []byte{byte(OpPow)}},
		{OpMod, []int{}, []byte{byte(OpMod)}},
		// Booleans
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
		// Comparisons
		{OpEqual, []int{}, []byte{byte(OpEqual)}},
		{OpNotEqual, []int{}, []byte{byte(OpNotEqual)}},
		{OpGreaterThan, []int{}, []byte{byte(OpGreaterThan)}},
		{OpGreaterThanOrEqual, []int{}, []byte{byte(OpGreaterThanOrEqual)}},
		// Prefixes
		{OpMinus, []int{}, []byte{byte(OpMinus)}},
		{OpBang, []int{}, []byte{byte(OpBang)}},
		// Suffixes
		{OpIndex, []int{}, []byte{byte(OpIndex)}},
		{OpIndexAssign, []int{}, []byte{byte(OpIndexAssign)}},
		// Objects
		{OpArray, []int{255}, []byte{byte(OpArray), 0, 255}},
		{OpHash, []int{255}, []byte{byte(OpHash), 0, 255}},
		// Loop control
		{OpLoopEnd, []int{}, []byte{byte(OpLoopEnd)}},
		// Functions
		{OpCall, []int{255}, []byte{byte(OpCall), 255}},
		{OpReturnValue, []int{}, []byte{byte(OpReturnValue)}},
		{OpReturn, []int{}, []byte{byte(OpReturn)}},
		// Internal Functions
		{OpGetBuiltin, []int{255}, []byte{byte(OpGetBuiltin), 255}},
		{OpGetGlobalBuiltin, []int{65535}, []byte{byte(OpGetGlobalBuiltin), 255, 255}},
		// Closures
		{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
		{OpGetFree, []int{255}, []byte{byte(OpGetFree), 255}},
		{OpCurrentClosure, []int{}, []byte{byte(OpCurrentClosure)}},
		// Import/Export
		{OpImport, []int{65534}, []byte{byte(OpImport), 255, 254}},
		{OpExport, []int{65534}, []byte{byte(OpExport), 255, 254}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Fatalf("instruction has wrong length. got %d, want %d", len(instruction), len(tt.expected))
		}

		for i := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Fatalf("wrong byte at position %d. got %d, want %d", i, instruction[i], tt.expected[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	instructions := []Instructions{
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpClosure, 65535, 255),
	}

	expected := `0000 OpAdd
0001 OpGetLocal 1
0003 OpConstant 2
0006 OpConstant 65535
0009 OpClosure 65535 255
`
	combined := Instructions{}
	for _, ins := range instructions {
		combined = append(combined, ins...)
	}

	if combined.String() != expected {
		t.Errorf("instructions wrongly formatted.\n\tgot:\n%s\n\twant:\n%s", combined.String(), expected)
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpGetLocal, []int{255}, 1},
		{OpConstant, []int{65535}, 2},
		{OpClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		def, err := Lookup(tt.op)
		if err != nil {
			t.Fatalf("definition not found: %q\n", err)
		}

		operandsRead, n := ReadOperands(def, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. got %d, want %d", n, tt.bytesRead)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Fatalf("operand %d wrong. got %d, want %d", i, operandsRead[i], want)
			}
		}
	}
}
