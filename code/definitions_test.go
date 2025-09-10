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
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}

	expected := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
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
