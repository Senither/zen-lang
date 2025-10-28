package vm

import (
	"fmt"

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
	// This still needs to be implemented, however we now have
	// access to the VM instance as well as the function
	// itself, outside of the evaluator package.
	return objects.NativeErrorToErrorObject(
		fmt.Errorf("unsupported function call for the Virtual Machine"),
	)
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
