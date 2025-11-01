package vm

import (
	"fmt"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type CompiledFunctionAdapter struct {
	Fn *objects.CompiledFunction
	VM *VM
}

func (fa *CompiledFunctionAdapter) Type() objects.ObjectType {
	return "COMPILED_FUNCTION_ADAPTER"
}

func (fa *CompiledFunctionAdapter) Inspect() string {
	return fa.Fn.Inspect()
}

func (fa *CompiledFunctionAdapter) Call(args ...objects.Object) objects.Object {
	if len(args) != fa.Fn.NumParameters {
		return objects.NativeErrorToErrorObject(
			fmt.Errorf("wrong number of arguments: got %d, want %d", len(args), fa.Fn.NumParameters),
		)
	}

	funcVM := &VM{
		constants:   fa.VM.constants,
		stack:       make([]objects.Object, STACK_SIZE),
		sp:          0,
		globals:     fa.VM.globals,
		frames:      make([]*Frame, MAX_FRAMES),
		framesIndex: 0,
		settings:    fa.VM.settings,
	}

	frame := NewFrame(fa.Fn, 0)
	funcVM.pushFrame(frame)

	copy(funcVM.stack, args)
	funcVM.sp = fa.Fn.NumLocals

	for funcVM.currentFrame().ip < len(funcVM.currentFrame().Instructions())-1 {
		funcVM.currentFrame().ip++

		ip := funcVM.currentFrame().ip
		ins := funcVM.currentFrame().Instructions()

		if ip >= len(ins) {
			break
		}

		op := code.Opcode(ins[ip])

		switch op {
		case code.OpReturnValue:
			return funcVM.pop()
		case code.OpReturn:
			return objects.NULL
		}

		err := funcVM.executeInstructions(op, ins, ip)
		if err != nil {
			return objects.NativeErrorToErrorObject(err)
		}
	}

	return objects.NULL
}

func WrapFunctionIfNeeded(vm *VM, obj objects.Object) objects.Object {
	if fn, ok := obj.(*objects.CompiledFunction); ok {
		return &CompiledFunctionAdapter{Fn: fn, VM: vm}
	}

	return obj
}

func (fa *CompiledFunctionAdapter) ParametersCount() int {
	return fa.Fn.NumParameters
}
