package optimizer

import (
	"bytes"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
)

const DEFAULT_OPTIMIZATION_ROUNDS = 10

type OptimizerPass func(b *BytecodeOptimization) error

type BytecodeOptimization struct {
	Infos          []InstructionInfo
	Constants      []objects.Object
	Targets        map[int]struct{}
	UsedGlobals    map[int]struct{}
	ChangedGlobals map[int]struct{}
	GlobalSwaps    map[int]instructionSwap
}

type InstructionInfo struct {
	Op        code.Opcode
	Operands  []int
	Width     int
	OldOffset int
	NewOffset int
	Keep      bool
	IsJump    bool
}

func Optimize(b *compiler.Bytecode) (*compiler.Bytecode, error) {
	return OptimizeRounds(b, DEFAULT_OPTIMIZATION_ROUNDS)
}

func OptimizeRounds(b *compiler.Bytecode, rounds int) (*compiler.Bytecode, error) {
	if rounds <= 0 {
		return b, nil
	}

	out := &compiler.Bytecode{
		Constants:    make([]objects.Object, len(b.Constants)),
		Instructions: make(code.Instructions, len(b.Instructions)),
	}

	copy(out.Instructions, b.Instructions)
	for i, constant := range b.Constants {
		out.Constants[i] = objects.Copy(constant)
	}

	for range rounds {
		globalSwaps := computeGlobalSwaps(out.Instructions, out.Constants)

		for i := 0; i < len(out.Constants); i++ {
			switch obj := out.Constants[i].(type) {
			case *objects.CompiledFunction:
				optimized, constants, err := optimizeInstructions(obj.OpcodeInstructions, out.Constants, globalSwaps)
				if err != nil {
					return nil, err
				}

				obj.OpcodeInstructions = optimized
				out.Constants = constants
			case *objects.CompiledZenFileImport:
				b, err := OptimizeRounds(&compiler.Bytecode{
					Instructions: obj.OpcodeInstructions,
					Constants:    obj.Constants,
				}, 1)

				if err != nil {
					return nil, err
				}

				obj.OpcodeInstructions = b.Instructions
				obj.Constants = b.Constants
			}
		}

		optimized, constants, err := optimizeInstructions(out.Instructions, out.Constants, globalSwaps)
		if err != nil {
			return nil, err
		}

		if len(optimized) == len(out.Instructions) &&
			bytes.Equal(optimized, out.Instructions) &&
			len(constants) == len(out.Constants) &&
			equalConstants(constants, out.Constants) {
			break
		}

		out.Instructions = optimized
		out.Constants = constants
	}

	optimized, constants, err := optimizeConstantsReferences(out.Instructions, out.Constants)
	if err != nil {
		return nil, err
	}

	out.Instructions = optimized
	out.Constants = constants

	return out, nil
}

func optimizeInstructions(
	instructions code.Instructions,
	constants []objects.Object,
	globalSwaps map[int]instructionSwap,
) (code.Instructions, []objects.Object, error) {
	infos, err := decodeInstructions(instructions)
	if err != nil {
		return nil, nil, err
	}

	b := &BytecodeOptimization{
		Infos:          infos,
		Constants:      constants,
		Targets:        findJumpTargets(infos),
		UsedGlobals:    findUsedGlobals(infos, constants),
		ChangedGlobals: findChangedGlobals(infos, constants),
		GlobalSwaps:    globalSwaps,
	}

	err = b.runOptimizationPasses(
		unfoldNonReassignedVariables,
		removeUnusedVariableInitializations,
		preCalculateNumberConstants,
		concatenateStringableConstants,
		removeInstructionsAfterReturn,
	)

	if err != nil {
		return nil, nil, err
	}

	return b.reassembleBytecodeParameters()
}

func optimizeConstantsReferences(
	instructions code.Instructions,
	constants []objects.Object,
) (code.Instructions, []objects.Object, error) {
	infos, err := decodeInstructions(instructions)
	if err != nil {
		return nil, nil, err
	}

	opt := &BytecodeOptimization{
		Infos:     infos,
		Constants: constants,
	}

	if err := opt.runOptimizationPasses(reorganizeConstantReferences); err != nil {
		return nil, nil, err
	}

	return opt.reassembleBytecodeParameters()
}

func decodeInstructions(instructions code.Instructions) ([]InstructionInfo, error) {
	var infos []InstructionInfo

	for offset := 0; offset < len(instructions); {
		op := code.Opcode(instructions[offset])

		def, err := code.Lookup(op)
		if err != nil {
			return nil, err
		}

		operands, read := code.ReadOperands(def, instructions[offset+1:])
		width := 1 + read

		infos = append(infos, InstructionInfo{
			Op:        op,
			Operands:  operands,
			Width:     width,
			OldOffset: offset,
			IsJump:    op == code.OpJump || op == code.OpJumpNotTruthy,
			Keep:      true,
		})

		offset += width
	}

	return infos, nil
}

func (b *BytecodeOptimization) runOptimizationPasses(passes ...OptimizerPass) error {
	for _, pass := range passes {
		if err := pass(b); err != nil {
			return err
		}
	}

	return nil
}

func (b *BytecodeOptimization) reassembleBytecodeParameters() (code.Instructions, []objects.Object, error) {
	for i := range b.Infos {
		if b.isJumpTarget(b.Infos[i].OldOffset) {
			b.Infos[i].Keep = true
		}
	}

	b.rewriteJumpOperands()

	return b.assembleInstructions(), b.Constants, nil
}

func (b *BytecodeOptimization) rewriteJumpOperands() {
	newOffset := 0
	oldToNew := make(map[int]int, len(b.Infos))

	for i := range b.Infos {
		if !b.Infos[i].Keep {
			continue
		}

		b.Infos[i].NewOffset = newOffset
		oldToNew[b.Infos[i].OldOffset] = newOffset
		newOffset += b.Infos[i].Width
	}

	for i := range b.Infos {
		if !b.Infos[i].Keep || !b.Infos[i].IsJump || len(b.Infos[i].Operands) == 0 {
			continue
		}

		newTarget, ok := oldToNew[b.Infos[i].Operands[0]]
		if !ok {
			continue
		}

		b.Infos[i].Operands[0] = newTarget
	}
}

func (b *BytecodeOptimization) assembleInstructions() code.Instructions {
	var out code.Instructions
	for _, info := range b.Infos {
		if !info.Keep {
			continue
		}

		out = append(out, code.Make(info.Op, info.Operands...)...)
	}

	return out
}

func (b *BytecodeOptimization) isJumpTarget(offset int) bool {
	_, exists := b.Targets[offset]
	return exists
}

func (b *BytecodeOptimization) setInstructionInfoOpcode(idx int, op code.Opcode, operands []int) {
	b.Infos[idx].Op = op
	b.Infos[idx].Operands = operands

	def, err := code.Lookup(op)
	if err == nil {
		width := 1
		for _, w := range def.OperandWidths {
			width += w
		}

		b.Infos[idx].Width = width
	}
}

func (b *BytecodeOptimization) getKeptInstructionsInfo(startIdx, count int) ([]*InstructionInfo, bool) {
	var keptInfos []*InstructionInfo

	for i := startIdx - 1; i >= 0 && len(keptInfos) < count; i-- {
		if b.Infos[i].Keep {
			keptInfos = append(keptInfos, &b.Infos[i])
		}
	}

	if len(keptInfos) < count {
		return nil, false
	}

	return keptInfos, true
}
