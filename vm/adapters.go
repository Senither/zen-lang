package vm

import (
	"fmt"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type CompiledClosureAdapter struct {
	Closure *objects.Closure
	VM      *VM
}

func (ca *CompiledClosureAdapter) Type() objects.ObjectType {
	return "COMPILED_CLOSURE_ADAPTER"
}

func (ca *CompiledClosureAdapter) Inspect() string {
	return ca.Closure.Inspect()
}

func (ca *CompiledClosureAdapter) Call(args ...objects.Object) objects.Object {
	if len(args) != ca.Closure.Fn.NumParameters {
		return objects.NativeErrorToErrorObject(
			fmt.Errorf("wrong number of arguments: got %d, want %d", len(args), ca.Closure.Fn.NumParameters),
		)
	}

	funcVM := &VM{
		constants:   ca.VM.constants,
		stack:       make([]objects.Object, STACK_SIZE),
		sp:          0,
		globals:     ca.VM.globals,
		frames:      make([]*Frame, MAX_FRAMES),
		framesIndex: 0,
		settings:    ca.VM.settings,
	}

	frame := NewFrame(ca.Closure, 0)
	funcVM.pushFrame(frame)

	copy(funcVM.stack, args)
	funcVM.sp = ca.Closure.Fn.NumLocals

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

func WrapClosuresIfNeeded(vm *VM, obj objects.Object) objects.Object {
	if fn, ok := obj.(*objects.Closure); ok {
		return &CompiledClosureAdapter{Closure: fn, VM: vm}
	}

	return obj
}

func (ca *CompiledClosureAdapter) ParametersCount() int {
	return ca.Closure.Fn.NumParameters
}
