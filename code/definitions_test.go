package code

import (
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		name     string
		op       Opcode
		operands []int
		expected []byte
	}{
		{"OpConstant", OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{"OpPop", OpPop, []int{}, []byte{byte(OpPop)}},
		// Jumps
		{"OpJump", OpJump, []int{1024}, []byte{byte(OpJump), 4, 0}},
		{"OpJumpNotTruthy", OpJumpNotTruthy, []int{1024}, []byte{byte(OpJumpNotTruthy), 4, 0}},
		// Globals
		{"OpSetGlobal", OpSetGlobal, []int{255}, []byte{byte(OpSetGlobal), 0, 255}},
		{"OpGetGlobal", OpGetGlobal, []int{255}, []byte{byte(OpGetGlobal), 0, 255}},
		// Locals
		{"OpSetLocal", OpSetLocal, []int{255}, []byte{byte(OpSetLocal), 255}},
		{"OpGetLocal", OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		// Arithmetic
		{"OpAdd", OpAdd, []int{}, []byte{byte(OpAdd)}},
		{"OpSub", OpSub, []int{}, []byte{byte(OpSub)}},
		{"OpMul", OpMul, []int{}, []byte{byte(OpMul)}},
		{"OpDiv", OpDiv, []int{}, []byte{byte(OpDiv)}},
		{"OpPow", OpPow, []int{}, []byte{byte(OpPow)}},
		{"OpMod", OpMod, []int{}, []byte{byte(OpMod)}},
		// Increment/Decrement
		{"OpIncGlobal", OpIncGlobal, []int{65535}, []byte{byte(OpIncGlobal), 255, 255}},
		{"OpDecGlobal", OpDecGlobal, []int{65535}, []byte{byte(OpDecGlobal), 255, 255}},
		{"OpIncLocal", OpIncLocal, []int{255}, []byte{byte(OpIncLocal), 255}},
		{"OpDecLocal", OpDecLocal, []int{255}, []byte{byte(OpDecLocal), 255}},
		// Booleans
		{"OpTrue", OpTrue, []int{}, []byte{byte(OpTrue)}},
		{"OpFalse", OpFalse, []int{}, []byte{byte(OpFalse)}},
		// Comparisons
		{"OpEqual", OpEqual, []int{}, []byte{byte(OpEqual)}},
		{"OpNotEqual", OpNotEqual, []int{}, []byte{byte(OpNotEqual)}},
		{"OpGreaterThan", OpGreaterThan, []int{}, []byte{byte(OpGreaterThan)}},
		{"OpGreaterThanOrEqual", OpGreaterThanOrEqual, []int{}, []byte{byte(OpGreaterThanOrEqual)}},
		{"OpAnd", OpAnd, []int{}, []byte{byte(OpAnd)}},
		{"OpOr", OpOr, []int{}, []byte{byte(OpOr)}},
		// Prefixes
		{"OpMinus", OpMinus, []int{}, []byte{byte(OpMinus)}},
		{"OpBang", OpBang, []int{}, []byte{byte(OpBang)}},
		// Suffixes
		{"OpIndex", OpIndex, []int{}, []byte{byte(OpIndex)}},
		{"OpIndexAssign", OpIndexAssign, []int{}, []byte{byte(OpIndexAssign)}},
		// Objects
		{"OpArray", OpArray, []int{255}, []byte{byte(OpArray), 0, 255}},
		{"OpHash", OpHash, []int{255}, []byte{byte(OpHash), 0, 255}},
		// Loop control
		{"OpLoopEnd", OpLoopEnd, []int{}, []byte{byte(OpLoopEnd)}},
		// Functions
		{"OpCall", OpCall, []int{255}, []byte{byte(OpCall), 255}},
		{"OpReturnValue", OpReturnValue, []int{}, []byte{byte(OpReturnValue)}},
		{"OpReturn", OpReturn, []int{}, []byte{byte(OpReturn)}},
		// Internal Functions
		{"OpGetBuiltin", OpGetBuiltin, []int{255}, []byte{byte(OpGetBuiltin), 255}},
		{"OpGetGlobalBuiltin", OpGetGlobalBuiltin, []int{65535}, []byte{byte(OpGetGlobalBuiltin), 255, 255}},
		// Closures
		{"OpClosure", OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
		{"OpGetFree", OpGetFree, []int{255}, []byte{byte(OpGetFree), 255}},
		{"OpCurrentClosure", OpCurrentClosure, []int{}, []byte{byte(OpCurrentClosure)}},
		// Import/Export
		{"OpImport", OpImport, []int{65534}, []byte{byte(OpImport), 255, 254}},
		{"OpExport", OpExport, []int{}, []byte{byte(OpExport)}},
	}

	for _, tt := range tests {
		t.Run("making "+tt.name, func(t *testing.T) {

			instruction := Make(tt.op, tt.operands...)

			if len(instruction) != len(tt.expected) {
				t.Fatalf("instruction has wrong length. got %d, want %d", len(instruction), len(tt.expected))
			}

			for i := range tt.expected {
				if instruction[i] != tt.expected[i] {
					t.Fatalf("wrong byte at position %d. got %d, want %d", i, instruction[i], tt.expected[i])
				}
			}
		})
	}
}

func BenchmarkMakeNoOperandInstruction(b *testing.B) {
	for b.Loop() {
		Make(OpTrue)
	}
}

func BenchmarkMakeOneOperandInstruction(b *testing.B) {
	for b.Loop() {
		Make(OpConstant, 1)
	}
}

func BenchmarkMakeTwoOperandsInstruction(b *testing.B) {
	for b.Loop() {
		Make(OpClosure, 1, 2)
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
		name      string
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{"OpGetLocal", OpGetLocal, []int{255}, 1},
		{"OpConstant", OpConstant, []int{65535}, 2},
		{"OpClosure", OpClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		t.Run("reading operands for "+tt.name, func(t *testing.T) {
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
		})
	}
}

func BenchmarkLookup(b *testing.B) {
	for b.Loop() {
		Lookup(OpAdd)
	}
}

func BenchmarkReadNoOperands(b *testing.B) {
	instruction := Make(OpAdd)
	def, _ := Lookup(OpAdd)

	for b.Loop() {
		ReadOperands(def, instruction[1:])
	}
}

func BenchmarkReadOneOperand(b *testing.B) {
	instruction := Make(OpGetGlobal, 255)
	def, _ := Lookup(OpGetGlobal)

	for b.Loop() {
		ReadOperands(def, instruction[1:])
	}
}

func BenchmarkReadTwoOperands(b *testing.B) {
	instruction := Make(OpClosure, 65535, 255)
	def, _ := Lookup(OpClosure)

	for b.Loop() {
		ReadOperands(def, instruction[1:])
	}
}
