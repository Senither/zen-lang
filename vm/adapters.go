package vm

import (
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
		fnName := "<anonymous>"
		if ca.Closure.Fn.Name != "" {
			fnName = ca.Closure.Fn.Name
		}

		return objects.NativeErrorToErrorObject(
			objects.NewWrongNumberOfArgumentsError(fnName, ca.Closure.Fn.NumParameters, len(args)),
		)
	}

	funcVM := ca.VM.Copy()

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
