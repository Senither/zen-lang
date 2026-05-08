package optimizer

import (
	"fmt"
	"math"
	"strconv"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

func equalConstants(a, b []objects.Object) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if fmt.Sprintf("%s", a[i]) != fmt.Sprintf("%s", b[i]) {
			return false
		}
	}

	return true
}

func findJumpTargets(infos []InstructionInfo) map[int]struct{} {
	targets := map[int]struct{}{}

	for _, info := range infos {
		if !info.IsJump || len(info.Operands) == 0 {
			continue
		}

		jumpTarget := info.Operands[0]
		targets[jumpTarget] = struct{}{}
	}

	return targets
}

func findUsedGlobals(infos []InstructionInfo, constants []objects.Object) map[int]struct{} {
	usedGlobals := findUsedGlobalsInInstructions(infos)

	for _, constant := range constants {
		fn, ok := constant.(*objects.CompiledFunction)
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

		nestedGlobals := findUsedGlobalsInInstructions(nestedInfos)
		for idx := range nestedGlobals {
			usedGlobals[idx] = struct{}{}
		}
	}

	return usedGlobals
}

func findUsedGlobalsInInstructions(infos []InstructionInfo) map[int]struct{} {
	usedGlobals := map[int]struct{}{}

	for _, info := range infos {
		if len(info.Operands) == 0 {
			continue
		}

		switch info.Op {
		case code.OpGetGlobal, code.OpIncGlobal, code.OpDecGlobal:
			usedGlobals[info.Operands[0]] = struct{}{}
		}
	}

	return usedGlobals
}

func findChangedGlobals(infos []InstructionInfo, constants []objects.Object) map[int]struct{} {
	globalUpdateCounters := findChangedGlobalsInInstructions(infos)

	for _, constant := range constants {
		fn, ok := constant.(*objects.CompiledFunction)
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

		nestedGlobals := findChangedGlobalsInInstructions(nestedInfos)
		for idx := range nestedGlobals {
			if _, exists := globalUpdateCounters[idx]; !exists {
				globalUpdateCounters[idx] = 0
			}

			globalUpdateCounters[idx] += nestedGlobals[idx]
		}
	}

	var globalUpdates map[int]struct{} = make(map[int]struct{})
	for idx, count := range globalUpdateCounters {
		if count > 1 {
			globalUpdates[idx] = struct{}{}
		}
	}

	return globalUpdates
}

func findChangedGlobalsInInstructions(infos []InstructionInfo) map[int]int {
	globals := map[int]int{}

	for _, info := range infos {
		if len(info.Operands) == 0 {
			continue
		}

		switch info.Op {
		case code.OpSetGlobal, code.OpIncGlobal, code.OpDecGlobal:
			if _, exists := globals[info.Operands[0]]; !exists {
				globals[info.Operands[0]] = 0
			}

			globals[info.Operands[0]]++
		}
	}

	return globals
}

func findPrevKeptInstructionIndex(infos []InstructionInfo, from int) int {
	for i := from - 1; i >= 0; i-- {
		if infos[i].Keep {
			return i
		}
	}

	return -1
}

func findNextKeptInstructionIndex(infos []InstructionInfo, from int) int {
	for i := from + 1; i < len(infos); i++ {
		if infos[i].Keep {
			return i
		}
	}

	return -1
}

func findJumpTargetsFromKeptInstructions(infos []InstructionInfo) map[int]struct{} {
	targets := map[int]struct{}{}

	for _, info := range infos {
		if !info.Keep || !info.IsJump || len(info.Operands) == 0 {
			continue
		}

		targets[info.Operands[0]] = struct{}{}
	}

	return targets
}

func isNoOpJump(infos []InstructionInfo, jumpIdx, targetIdx int) bool {
	nextIdx := findNextKeptInstructionIndex(infos, jumpIdx)
	if nextIdx < 0 {
		return false
	}

	return nextIdx == targetIdx
}

func isWhileJumpPattern(infos []InstructionInfo, jumpNotTruthyIdx, targetIdx, prevTargetIdx int, hasElseJump bool) bool {
	if hasElseJump || prevTargetIdx <= jumpNotTruthyIdx {
		return false
	}

	prevTarget := infos[prevTargetIdx]
	if prevTarget.Op != code.OpJump || len(prevTarget.Operands) == 0 {
		return false
	}

	return prevTarget.Operands[0] <= infos[jumpNotTruthyIdx].OldOffset
}

func computeGlobalSwaps(instructions code.Instructions, constants []objects.Object) map[int]instructionSwap {
	infos, err := decodeInstructions(instructions)
	if err != nil {
		return nil
	}

	changedGlobals := findChangedGlobals(infos, constants)
	swaps := map[int]instructionSwap{}

	for i := range infos {
		if len(infos[i].Operands) == 0 {
			continue
		}

		switch infos[i].Op {
		case code.OpSetGlobal:
			globalIdx := infos[i].Operands[0]
			if _, reassigned := changedGlobals[globalIdx]; reassigned {
				continue
			}

			if i == 0 {
				continue
			}

			prev := &infos[i-1]
			if len(prev.Operands) == 0 {
				continue
			}

			switch prev.Op {
			case code.OpConstant, code.OpNull, code.OpTrue, code.OpFalse:
				swaps[globalIdx] = instructionSwap{
					Operands: prev.Operands,
					Op:       prev.Op,
				}
			}
		}
	}

	return swaps
}

func stackDelta(info *InstructionInfo) int {
	switch info.Op {
	case code.OpConstant, code.OpNull, code.OpTrue, code.OpFalse:
		return 1
	case code.OpArray, code.OpHash:
		if len(info.Operands) == 0 {
			return 0
		}

		return 1 - info.Operands[0]
	case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod, code.OpIndex,
		code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
		return -1
	case code.OpMinus, code.OpBang:
		return 0
	case code.OpPop, code.OpSetGlobal, code.OpSetLocal:
		return -1

	default:
		return 0
	}
}

func resolveTargetInstructionIndex(infos []InstructionInfo, offsetToIndex map[int]int, targetOffset int) (int, bool) {
	if idx, ok := offsetToIndex[targetOffset]; ok {
		return idx, true
	}

	if len(infos) == 0 {
		return -1, false
	}

	endOffset := infos[len(infos)-1].OldOffset + infos[len(infos)-1].Width
	if targetOffset == endOffset {
		return len(infos), true
	}

	return -1, false
}

func evaluateKnownConditionTruthiness(b *BytecodeOptimization, condIdx int) (bool, bool) {
	condInfo := b.Infos[condIdx]

	switch condInfo.Op {
	case code.OpTrue:
		return objects.IsTruthy(objects.TRUE), true
	case code.OpFalse:
		return objects.IsTruthy(objects.FALSE), true
	case code.OpNull:
		return objects.IsTruthy(objects.NULL), true
	case code.OpConstant:
		if len(condInfo.Operands) == 0 {
			return false, false
		}

		constIdx := condInfo.Operands[0]
		if constIdx < 0 || constIdx >= len(b.Constants) {
			return false, false
		}

		return objects.IsTruthy(b.Constants[constIdx]), true

	default:
		return false, false
	}
}

func canRemoveIndexes(
	infos []InstructionInfo,
	incoming map[int][]int,
	toRemove map[int]struct{},
) bool {
	for idx := range toRemove {
		if idx < 0 || idx >= len(infos) || !infos[idx].Keep {
			continue
		}

		sources := incoming[infos[idx].OldOffset]
		for _, srcIdx := range sources {
			if !infos[srcIdx].Keep {
				continue
			}

			if _, deletingSource := toRemove[srcIdx]; !deletingSource {
				return false
			}
		}
	}

	return true
}

func immutableConstantReuseKey(obj objects.Object) (string, bool) {
	switch value := obj.(type) {
	case *objects.Null:
		return "null", true
	case *objects.Boolean:
		if value.Value {
			return "bool:1", true
		}

		return "bool:0", true
	case *objects.Integer:
		return "int:" + strconv.FormatInt(value.Value, 10), true
	case *objects.Float:
		return "float:" + strconv.FormatUint(math.Float64bits(value.Value), 16), true
	case *objects.String:
		return "string:" + strconv.Itoa(len(value.Value)) + ":" + value.Value, true

	default:
		return "", false
	}
}
