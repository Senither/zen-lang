package optimizer

import (
	"fmt"

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
