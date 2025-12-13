package optimizer

import (
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

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
