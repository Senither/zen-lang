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

// Deletes the instructions that are used to initialize an array or hash if the value is never used,
// this is done by calculating the stack delta of the instructions until it matches the number of
// elements in the array or hash, if a jump target is found during the process the optimization
// is aborted since it may be used somewhere else.
//
// Example:
//
//	OpConstant 0   (value 42)
//	OpConstant 1   (value "hello")
//	OpConstant 2   (value "world")
//	OpArray 3      (3 elements)
//
// -->
//
//	(nothing)
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

// Folds comparison and logical operations into OpTrue/OpFalse when both operands
// are known constants in advance.
//
// Example:
//
//	OpConstant 0   (value 5)
//	OpConstant 1   (value 10)
//	OpGreaterThan
//
// -->
//
//	OpFalse
func constantFoldComparisonLogicalOps(b *BytecodeOptimization) error {
	for i := range b.Infos {
		switch b.Infos[i].Op {
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual, code.OpAnd, code.OpOr:
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
			if leftConstIdx < 0 || leftConstIdx >= len(b.Constants) || rightConstIdx < 0 || rightConstIdx >= len(b.Constants) {
				continue
			}

			leftObj := b.Constants[leftConstIdx]
			rightObj := b.Constants[rightConstIdx]

			var result bool
			switch b.Infos[i].Op {
			case code.OpEqual:
				if objects.IsNumber(leftObj.Type()) && objects.IsNumber(rightObj.Type()) {
					result = objects.UnwrapNumberValue(leftObj) == objects.UnwrapNumberValue(rightObj)
				} else {
					result = objects.Equals(leftObj, rightObj) == objects.TRUE
				}
			case code.OpNotEqual:
				if objects.IsNumber(leftObj.Type()) && objects.IsNumber(rightObj.Type()) {
					result = objects.UnwrapNumberValue(leftObj) != objects.UnwrapNumberValue(rightObj)
				} else {
					result = objects.Equals(leftObj, rightObj) == objects.FALSE
				}
			case code.OpGreaterThan:
				if !objects.IsNumber(leftObj.Type()) || !objects.IsNumber(rightObj.Type()) {
					continue
				}

				result = objects.UnwrapNumberValue(leftObj) > objects.UnwrapNumberValue(rightObj)
			case code.OpGreaterThanOrEqual:
				if !objects.IsNumber(leftObj.Type()) || !objects.IsNumber(rightObj.Type()) {
					continue
				}

				result = objects.UnwrapNumberValue(leftObj) >= objects.UnwrapNumberValue(rightObj)
			case code.OpAnd:
				result = objects.IsTruthy(leftObj) && objects.IsTruthy(rightObj)
			case code.OpOr:
				result = objects.IsTruthy(leftObj) || objects.IsTruthy(rightObj)
			}

			rightInfo.Keep = false
			leftInfo.Keep = false

			if result {
				b.setInstructionInfoOpcode(i, code.OpTrue, nil)
			} else {
				b.setInstructionInfoOpcode(i, code.OpFalse, nil)
			}
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

// Replaces increments and decrements of a variable by 1 using OpConstant and OpAdd/Sub
// with a direct increment or decrement opcode.
//
// Example:
//
//	OpGetGlobal 0   (variable a)
//	OpConstant 0    (value 1)
//	OpAdd
//	OpSetGlobal 0   (variable a)
//
// -->
//
//	OpIncGlobal 0   (variable a)
func replaceIncrementsAndDecrementsWithDirectOperations(b *BytecodeOptimization) error {
	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		switch b.Infos[i].Op {
		case code.OpAdd, code.OpSub:
			if i+1 >= len(b.Infos) {
				continue
			}

			var expectedGetter code.Opcode
			switch b.Infos[i+1].Op {
			case code.OpSetLocal:
				expectedGetter = code.OpGetLocal
			case code.OpSetGlobal:
				expectedGetter = code.OpGetGlobal

			default:
				continue
			}

			if !b.Infos[i+1].Keep || len(b.Infos[i+1].Operands) == 0 {
				continue
			}

			infos, ok := b.getKeptInstructionsInfo(i, 2)
			if !ok {
				continue
			}

			rightInfo := infos[0]
			leftInfo := infos[1]
			if len(leftInfo.Operands) == 0 || len(rightInfo.Operands) == 0 {
				continue
			}

			var constInfo *InstructionInfo
			var getterInfo *InstructionInfo
			if leftInfo.Op == code.OpConstant && rightInfo.Op == expectedGetter {
				constInfo = leftInfo
				getterInfo = rightInfo
			} else if rightInfo.Op == code.OpConstant && leftInfo.Op == expectedGetter {
				constInfo = rightInfo
				getterInfo = leftInfo
			} else {
				continue
			}

			if b.Infos[i+1].Operands[0] != getterInfo.Operands[0] {
				continue
			}

			constIdx := constInfo.Operands[0]
			if constIdx < 0 || constIdx >= len(b.Constants) {
				continue
			}

			constObj := b.Constants[constIdx]
			if constObj.Type() != objects.INTEGER_OBJ {
				continue
			}

			if constObj.(*objects.Integer).Value != 1 {
				continue
			}

			var replacementOp code.Opcode
			switch b.Infos[i].Op {
			case code.OpAdd:
				switch expectedGetter {
				case code.OpGetLocal:
					replacementOp = code.OpIncLocal
				case code.OpGetGlobal:
					replacementOp = code.OpIncGlobal
				}
			case code.OpSub:
				switch expectedGetter {
				case code.OpGetLocal:
					replacementOp = code.OpDecLocal
				case code.OpGetGlobal:
					replacementOp = code.OpDecGlobal
				}
			}

			if replacementOp == 0 {
				continue
			}

			getterInfo.Keep = false
			constInfo.Keep = false
			b.setInstructionInfoOpcode(i+1, code.OpPop, nil)

			b.setInstructionInfoOpcode(i, replacementOp, []int{getterInfo.Operands[0]})
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

// Removes OpJumpNotTruthy instructions that jump to an instruction that is
// effectively a no-op, meaning it doesn't have any side effects and just
// continues to the next instruction, this includes jumps to OpNull or
// OpPop instructions that are not jump targets themselves.
//
// Example:
//
//	OpTrue
//	OpJumpNotTruthy X
//	OpJump Y
//	OpNull      (This would be target X)
//	OpPop       (This would be target Y)
//
// -->
//
// (nothing)
func removeRedundantJumpInstructions(b *BytecodeOptimization) error {
	for {
		changed := false

		offsetToIndex := map[int]int{}
		for i := range b.Infos {
			offsetToIndex[b.Infos[i].OldOffset] = i
		}

		incoming := map[int][]int{}
		for i := range b.Infos {
			info := &b.Infos[i]
			if !info.Keep || !info.IsJump || len(info.Operands) == 0 {
				continue
			}

			incoming[info.Operands[0]] = append(incoming[info.Operands[0]], i)
		}

		for i := range b.Infos {
			if !b.Infos[i].Keep || b.Infos[i].Op != code.OpJumpNotTruthy || len(b.Infos[i].Operands) == 0 {
				continue
			}

			condIdx := findPrevKeptInstructionIndex(b.Infos, i)
			if condIdx < 0 {
				continue
			}

			targetIdx, ok := resolveTargetInstructionIndex(b.Infos, offsetToIndex, b.Infos[i].Operands[0])
			if !ok {
				continue
			}

			if targetIdx <= i {
				continue
			}

			prevTargetIdx := findPrevKeptInstructionIndex(b.Infos, targetIdx)
			hasElseJump := false
			afterElseIdx := len(b.Infos)
			if prevTargetIdx > i && b.Infos[prevTargetIdx].Op == code.OpJump && len(b.Infos[prevTargetIdx].Operands) > 0 {
				afterElseIdx, ok = resolveTargetInstructionIndex(b.Infos, offsetToIndex, b.Infos[prevTargetIdx].Operands[0])
				if ok && b.Infos[prevTargetIdx].Operands[0] > b.Infos[i].Operands[0] {
					hasElseJump = true
				}
			}

			if isNoOpJump(b.Infos, i, targetIdx) {
				b.Infos[i].Keep = false
				changed = true
				break
			}

			truthy, known := evaluateKnownConditionTruthiness(b, condIdx)
			if !known {
				continue
			}

			toRemove := map[int]struct{}{
				condIdx: {},
				i:       {},
			}

			if isWhileJumpPattern(b.Infos, i, targetIdx, prevTargetIdx, hasElseJump) && !truthy {
				for idx := condIdx; idx <= targetIdx; idx++ {
					toRemove[idx] = struct{}{}
				}

				if !canRemoveIndexes(b.Infos, incoming, toRemove) {
					continue
				}

				for idx := range toRemove {
					b.Infos[idx].Keep = false
				}

				changed = true
				break
			}

			if hasElseJump {
				if afterElseIdx <= targetIdx {
					continue
				}

				if truthy {
					toRemove[prevTargetIdx] = struct{}{}
					for idx := targetIdx; idx < afterElseIdx; idx++ {
						toRemove[idx] = struct{}{}
					}
				} else {
					toRemove[prevTargetIdx] = struct{}{}
					for idx := i + 1; idx < prevTargetIdx; idx++ {
						toRemove[idx] = struct{}{}
					}
				}
			} else if !truthy {
				for idx := i + 1; idx < targetIdx; idx++ {
					toRemove[idx] = struct{}{}
				}
			}

			if !canRemoveIndexes(b.Infos, incoming, toRemove) {
				continue
			}

			for idx := range toRemove {
				b.Infos[idx].Keep = false
			}

			changed = true
			break
		}

		if !changed {
			break
		}
	}

	b.Targets = findJumpTargetsFromKeptInstructions(b.Infos)

	return nil
}

// Removes OpPop instructions at the beginning of the instructions since they
// are likely used to discard the result of a previous operation that is no
// longer needed, this is common after optimizations that remove
// instructions but leave the OpPop in place.
//
// Example:
//
//	OpPop
//	OpConstant 0   (value 42)
//	OpReturnValue
//
// -->
//
//	OpConstant 0   (value 42)
//	OpReturnValue
func removePopAtBeginningOfInstructions(b *BytecodeOptimization) error {
	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		if b.Infos[i].Op != code.OpPop {
			break
		}

		b.Infos[i].Keep = false
	}

	return nil
}

// Reorganizes constant references to remove unused constants and
// re-index the used and duplicated ones to a more compact range.
//
// Example:
//
//	OpConstant 7   (value 42)
//	OpConstant 19  (value "hello")
//	OpConstant 25  (value 42)
//
// -->
//
//	...
//	OpConstant 0   (value 42)
//	OpConstant 1   (value "hello")
//	OpConstant 0   (value 42)
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

	indexMap := make(map[int]int, len(used))
	newConstants := make([]objects.Object, 0, len(used))
	immutableConstantIndex := make(map[string]int, len(used))

	for oldIdx, c := range b.Constants {
		if _, ok := used[oldIdx]; !ok {
			continue
		}

		if key, immutable := immutableConstantReuseKey(c); immutable {
			if idx, exists := immutableConstantIndex[key]; exists {
				indexMap[oldIdx] = idx
				continue
			}
		}

		newIdx := len(newConstants)
		indexMap[oldIdx] = newIdx
		newConstants = append(newConstants, c)

		if key, immutable := immutableConstantReuseKey(c); immutable {
			immutableConstantIndex[key] = newIdx
		}
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
