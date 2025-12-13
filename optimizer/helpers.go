package optimizer

import "github.com/senither/zen-lang/code"

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

func findUsedGlobals(infos []InstructionInfo) map[int]struct{} {
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
