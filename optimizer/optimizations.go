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

// Removes initial OpSetGlobal instructions that are never reassigned a new value,
// and replaces all OpGetGlobal references with the OpConstant equivalent.
//
// Example:
//
//	OpConstant 0   (value 42)
//	OpSetGlobal 0  (variable a)
//	...
//	OpGetGlobal 0  (variable a)
//
// -->
//
//	...
//	OpConstant 0   (value 42)
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

// Removes OpSetGlobal instructions if the global variable is never referenced anywhere, including
// the initialization of the value being stored in the global (array or hash construction).
//
// Example:
//
//	OpConstant 0   (value 42)
//	OpConstant 1   (value "hello")
//	OpConstant 2   (value "world")
//	OpArray 3      (3 elements)
//	OpSetGlobal 0  (variable a)
//
// -->
//
//	(nothing)
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

// Pre-calculates operations that only involve constant values that are numbers, such as addition,
// subtraction, multiplication, division, etc. The result is stored as a new constant.
//
// Example:
//
//	OpConstant 0   (value 9)
//	OpConstant 1   (value 10)
//	OpConstant 2   (value 42)
//	OpMul          (multiplies top two constants)
//	OpAdd          (adds top two constants)
//
// -->
//
//	OpConstant 0   (value 429)
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

// Concatenates stringable constants by using the objects.StringifyObject function,
// and then storing the result as a new constant, at least one of the two
// constants must be a string object to perform the optimization.
//
// Example:
//
//	OpConstant 0   (value "Value=")
//	OpConstant 1   (value 42)
//	OpAdd          (concatenates top two constants)
//
// -->
//
//	OpConstant 0   (value "Value=42")
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

// Removes OpGetGlobal and OpGetLocal instructions that immediately follow
// an OpSetGlobal or OpSetLocal instruction, and are themselves followed
// by an OpPop instruction so that the value is never used.
//
// This will only remove instructions that exists within the instruction set, if the
// OpPop opcode is at the end of the instructions it will be kept as is since it
// may be used to validate the VMs last popped element in tests or outputs.
//
// Example:
//
//	OpConstant 0   (value 42)
//	OpSetLocal 0   (variable a)
//	OpGetLocal 0   (variable a)
//	OpPop
//
// -->
//
//	OpConstant 0   (value 42)
//	OpSetLocal 0   (variable a)
func removeUnusedGettersAfterAssignments(b *BytecodeOptimization) error {
	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		switch b.Infos[i].Op {
		case code.OpSetLocal, code.OpSetGlobal:
			if i+3 >= len(b.Infos) {
				continue
			}

			getterInfo := &b.Infos[i+1]
			popInfo := &b.Infos[i+2]
			if !getterInfo.Keep || !popInfo.Keep {
				continue
			}

			var expectedGetter code.Opcode
			switch b.Infos[i].Op {
			case code.OpSetLocal:
				expectedGetter = code.OpGetLocal
			case code.OpSetGlobal:
				expectedGetter = code.OpGetGlobal
			}

			if getterInfo.Op != expectedGetter || getterInfo.Operands[0] != b.Infos[i].Operands[0] {
				continue
			}

			if popInfo.Op != code.OpPop {
				continue
			}

			getterInfo.Keep = false
			popInfo.Keep = false
		}
	}

	return nil
}

// Calls built-in functions if all the parameters are known constants, and stores the result as a new constant.
// Some builtins are skipped because they may have side effects or are non-deterministic.
//
// Example:
//
//	OpGetBuiltin 2  (builtin "len")
//	OpConstant 0    (value "hello")
//	OpCall 1        (1 argument)
//
// -->
//
//	OpConstant 0    (value 5)
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

// Removes instructions that are unreachable because they are after a return statement.
//
// Example:
//
//	...
//	OpConstant 0   (value 42)
//	OpReturnValue
//	OpConstant 1   (value "unreachable")
//
// -->
//
//	...
//	OpConstant 0   (value 42)
//	OpReturnValue
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

// Reorganizes constant references to remove unused constants and
// re-index the used ones to a more compact range.
//
// Example:
//
//	OpConstant 7   (value 42)
//	OpConstant 19  (value "hello")
//
// -->
//
//	...
//	OpConstant 0   (value 42)
//	OpConstant 1   (value "hello")
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
