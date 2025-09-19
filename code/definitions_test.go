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
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
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
		{OpConstant, []int{65535}, 2},
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
