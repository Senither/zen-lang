package optimizer

import (
	"math"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type instructionSwap struct {
	Operands []int
	Op       code.Opcode
}

func unfoldNonReassignedVariables(b *BytecodeOptimization) error {
	swaps := map[int]instructionSwap{}

	for i := range b.Infos {
		if len(b.Infos[i].Operands) == 0 {
			continue
		}

		switch b.Infos[i].Op {
		case code.OpSetGlobal:
			globalIdx := b.Infos[i].Operands[0]
			_, reassigned := b.ChangedGlobals[globalIdx]
			if reassigned {
				continue
			}

			prev := &b.Infos[i-1]
			if !prev.Keep {
				continue
			}

			switch prev.Op {
			case code.OpConstant, code.OpNull, code.OpTrue, code.OpFalse:
				swaps[globalIdx] = instructionSwap{
					Operands: prev.Operands,
					Op:       prev.Op,
				}

				prev.Keep = false
				b.Infos[i].Keep = false
			}

		case code.OpGetGlobal:
			globalIdx := b.Infos[i].Operands[0]

			if b.GlobalSwaps != nil {
				if swap, ok := b.GlobalSwaps[globalIdx]; ok {
					b.setInstructionInfoOpcode(i, swap.Op, swap.Operands)
					continue
				}
			}

			if swap, ok := swaps[globalIdx]; ok {
				b.setInstructionInfoOpcode(i, swap.Op, swap.Operands)
			}
		}
	}

	return nil
}

func removeUnusedVariableInitializations(b *BytecodeOptimization) error {
	for i := range b.Infos {
		if b.Infos[i].Op != code.OpSetGlobal || len(b.Infos[i].Operands) == 0 {
			continue
		}

		globalIdx := b.Infos[i].Operands[0]
		if _, used := b.UsedGlobals[globalIdx]; used {
			continue
		}

		if b.isJumpTarget(b.Infos[i].OldOffset) {
			continue
		}

		b.Infos[i].Keep = false

		if i > 0 {
			prev := &b.Infos[i-1]
			if b.isJumpTarget(prev.OldOffset) {
				continue
			}

			if prev.Keep {
				switch prev.Op {
				case code.OpConstant, code.OpNull, code.OpTrue, code.OpFalse:
					prev.Keep = false

				case code.OpArray, code.OpHash:
					deleteArrayOrHashInitializer(b, i-1)

				case code.OpClosure:
					if len(prev.Operands) == 0 {
						continue
					}

					b.Constants[prev.Operands[0]] = objects.NULL
					prev.Keep = false
				}
			}
		}
	}

	return nil
}

func deleteArrayOrHashInitializer(b *BytecodeOptimization, idx int) {
	info := &b.Infos[idx]
	if len(info.Operands) == 0 {
		return
	}

	targetDelta := info.Operands[0]
	currentDelta := 0

	toDelete := map[int]struct{}{idx: {}}
	for i := idx - 1; i >= 0 && currentDelta < targetDelta; i-- {
		inst := &b.Infos[i]
		if !inst.Keep {
			continue
		}

		if b.isJumpTarget(inst.OldOffset) {
			return
		}

		currentDelta += stackDelta(inst)
		toDelete[i] = struct{}{}
	}

	if currentDelta != targetDelta {
		return
	}

	for i := range toDelete {
		b.Infos[i].Keep = false
	}
}

func preCalculateNumberConstants(b *BytecodeOptimization) error {
	for i := range b.Infos {
		switch b.Infos[i].Op {
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpPow, code.OpMod:
			if !b.Infos[i].Keep {
				continue
			}

			infos, ok := b.getKeptInstructionsInfo(i, 2)
			if !ok {
				continue
			}

			rightInfo := infos[0]
			leftInfo := infos[1]

			if leftInfo.Op != code.OpConstant || rightInfo.Op != code.OpConstant {
				continue
			}

			leftConstIdx := leftInfo.Operands[0]
			rightConstIdx := rightInfo.Operands[0]

			rightObj, leftObj := b.Constants[rightConstIdx], b.Constants[leftConstIdx]
			if !objects.IsNumber(leftObj.Type()) || !objects.IsNumber(rightObj.Type()) {
				continue
			}

			leftVal := objects.UnwrapNumberValue(leftObj)
			rightVal := objects.UnwrapNumberValue(rightObj)

			var result float64

			switch b.Infos[i].Op {
			case code.OpAdd:
				result = leftVal + rightVal
			case code.OpSub:
				result = leftVal - rightVal
			case code.OpMul:
				result = leftVal * rightVal
			case code.OpDiv:
				result = leftVal / rightVal
			case code.OpPow:
				result = math.Pow(leftVal, rightVal)
			case code.OpMod:
				result = math.Mod(leftVal, rightVal)
			}

			newConst := objects.WrapNumberValue(result, leftObj, rightObj)

			newConstIdx := len(b.Constants)
			b.Constants = append(b.Constants, newConst)

			rightInfo.Keep = false
			leftInfo.Keep = false

			b.setInstructionInfoOpcode(i, code.OpConstant, []int{newConstIdx})
		}
	}

	return nil
}
