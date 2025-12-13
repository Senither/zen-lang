package optimizer

import (
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
			swap, ok := swaps[b.Infos[i].Operands[0]]
			if !ok {
				continue
			}

			b.Infos[i].Op = swap.Op
			b.Infos[i].Operands = swap.Operands
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
