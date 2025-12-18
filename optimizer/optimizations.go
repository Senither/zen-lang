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

func concatenateStringableConstants(b *BytecodeOptimization) error {
	isStringableMatch := func(a, b objects.Object) bool {
		if a.Type() == objects.STRING_OBJ && objects.IsStringable(b) {
			return true
		}

		if objects.IsStringable(a) && b.Type() == objects.STRING_OBJ {
			return true
		}

		return false
	}

	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		if b.Infos[i].Op != code.OpAdd {
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

		var newConst *objects.String = nil
		if leftObj.Type() == objects.STRING_OBJ && rightObj.Type() == objects.STRING_OBJ {
			leftStr := leftObj.(*objects.String).Value
			rightStr := rightObj.(*objects.String).Value

			newConst = &objects.String{Value: leftStr + rightStr}
		} else if isStringableMatch(leftObj, rightObj) {
			leftStr := objects.StringifyObject(leftObj)
			rightStr := objects.StringifyObject(rightObj)

			newConst = &objects.String{Value: leftStr + rightStr}
		}

		if newConst == nil {
			continue
		}

		newConstIdx := len(b.Constants)
		b.Constants = append(b.Constants, newConst)

		rightInfo.Keep = false
		leftInfo.Keep = false

		b.setInstructionInfoOpcode(i, code.OpConstant, []int{newConstIdx})
	}

	return nil
}

func callBuiltinsWithKnownConstantParameters(b *BytecodeOptimization) error {
MAIN_LOOP:
	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		if b.Infos[i].Op != code.OpCall {
			continue
		}

		infos, ok := b.getKeptInstructionsInfo(i, b.Infos[i].Operands[0]+1)
		if !ok {
			continue
		}

		for j := 0; j < len(infos)-1; j++ {
			if infos[j].Op != code.OpConstant {
				continue MAIN_LOOP
			}
		}

		builtinInfo := infos[len(infos)-1]
		if builtinInfo.Op != code.OpGetBuiltin && builtinInfo.Op != code.OpGetGlobalBuiltin {
			continue
		}

		var definition *objects.BuiltinDefinition
		builtinIdx := builtinInfo.Operands[0]

		switch builtinInfo.Op {
		case code.OpGetBuiltin:
			definition = &objects.Builtins[builtinIdx]
		case code.OpGetGlobalBuiltin:
			scopeIdx := uint8(builtinIdx >> 8)
			builtIdx := uint8(builtinIdx & 0xFF)

			definition = objects.Globals[scopeIdx].Builtins[builtIdx]
		}

		if definition == nil || definition.OmitOptimization {
			continue
		}

		args := make([]objects.Object, len(infos)-1)
		for j := 0; j < len(infos)-1; j++ {
			constIdx := infos[j].Operands[0]
			args[len(infos)-2-j] = b.Constants[constIdx]
		}

		result, err := definition.Builtin.Fn(args...)
		if err != nil {
			continue
		}

		newConstIdx := len(b.Constants)
		b.Constants = append(b.Constants, result)

		for j := range infos {
			infos[j].Keep = false
		}

		b.setInstructionInfoOpcode(i, code.OpConstant, []int{newConstIdx})
	}

	return nil
}

func removeInstructionsAfterReturn(b *BytecodeOptimization) error {
	foundReturn := false

	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		if foundReturn {
			if b.isJumpTarget(b.Infos[i].OldOffset) {
				foundReturn = false
				continue
			}

			b.Infos[i].Keep = false
			continue
		}

		if b.Infos[i].Op == code.OpReturnValue || b.Infos[i].Op == code.OpReturn {
			foundReturn = true
		}
	}

	return nil
}

func reorganizeConstantReferences(b *BytecodeOptimization) error {
	used := map[int]struct{}{}

	markUsedFromInfos := func(infos []InstructionInfo) {
		for _, info := range infos {
			switch info.Op {
			case code.OpConstant, code.OpClosure, code.OpImport:
				used[info.Operands[0]] = struct{}{}
			}
		}
	}

	markUsedFromInfos(b.Infos)

	for _, c := range b.Constants {
		fn, ok := c.(*objects.CompiledFunction)
		if !ok {
			continue
		}

		ins := fn.Instructions()
		if len(ins) == 0 {
			continue
		}

		nestedInfos, err := decodeInstructions(ins)
		if err != nil {
			continue
		}

		markUsedFromInfos(nestedInfos)
	}

	if len(used) == len(b.Constants) {
		return nil
	}

	indexMap := make(map[int]int, len(used))
	newConstants := make([]objects.Object, 0, len(used))

	for oldIdx, c := range b.Constants {
		if _, ok := used[oldIdx]; !ok {
			continue
		}

		newIdx := len(newConstants)
		indexMap[oldIdx] = newIdx
		newConstants = append(newConstants, c)
	}

	for i := range b.Infos {
		switch b.Infos[i].Op {
		case code.OpConstant, code.OpClosure, code.OpImport:
			oldIdx := b.Infos[i].Operands[0]

			if newIdx, ok := indexMap[oldIdx]; ok {
				b.Infos[i].Operands[0] = newIdx
			}
		}
	}

	for _, c := range b.Constants {
		fn, ok := c.(*objects.CompiledFunction)
		if !ok {
			continue
		}

		ins := fn.Instructions()
		if len(ins) == 0 {
			continue
		}

		nestedInfos, err := decodeInstructions(ins)
		if err != nil {
			continue
		}

		changed := false
		for i := range nestedInfos {
			switch nestedInfos[i].Op {
			case code.OpConstant, code.OpClosure, code.OpImport:
				oldIdx := nestedInfos[i].Operands[0]

				if newIdx, ok := indexMap[oldIdx]; ok {
					nestedInfos[i].Operands[0] = newIdx
					changed = true
				}
			}
		}

		if !changed {
			continue
		}

		var newIns code.Instructions
		for _, info := range nestedInfos {
			newIns = append(newIns, code.Make(info.Op, info.Operands...)...)
		}

		fn.OpcodeInstructions = newIns
	}

	b.Constants = newConstants

	return nil
}
