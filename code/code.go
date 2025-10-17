package code

import (
	"bytes"
	"fmt"
)

type Opcode byte

const (
	OpConstant Opcode = iota
	OpPop

	OpNull

	// Jumps
	OpJump
	OpJumpNotTruthy

	// Globals
	OpSetGlobal
	OpGetGlobal

	// Arithmetic
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPow
	OpMod

	// Booleans
	OpTrue
	OpFalse

	// Comparisons
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpGreaterThanOrEqual

	// Prefixes
	OpMinus
	OpBang

	// Suffixes
	OpIndex

	// Objects
	OpArray
	OpHash

	// Functions
	OpCall
	OpReturnValue
	OpReturn
)

type Instructions []byte

func (in Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(in) {
		def, err := Lookup(Opcode(in[i]))
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, in[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, in.FormatInstruction(def, operands))

		i += read + 1
	}

	return out.String()
}

func (in Instructions) FormatInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])

	default:
		return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
	}
}
