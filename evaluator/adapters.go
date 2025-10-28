package evaluator

import "github.com/senither/zen-lang/objects"

type FunctionAdapter struct {
	Fn *objects.Function
}

func (fa *FunctionAdapter) Type() objects.ObjectType {
	return "FUNCTION_ADAPTER"
}

func (fa *FunctionAdapter) Inspect() string {
	return fa.Fn.Inspect()
}

func (fa *FunctionAdapter) Call(args ...objects.Object) objects.Object {
	env := objects.NewEnclosedEnvironment(fa.Fn.Env)

	for paramIdx, param := range fa.Fn.Parameters {
		env.SetImmutableForcefully(param.Value, args[paramIdx])
	}

	return objects.UnwrapReturnValue(Eval(fa.Fn.Body, env))
}

func WrapFunctionIfNeeded(obj objects.Object) objects.Object {
	if fn, ok := obj.(*objects.Function); ok {
		return &FunctionAdapter{Fn: fn}
	}

	return obj
}

func (fa *FunctionAdapter) ParametersCount() int {
	return len(fa.Fn.Parameters)
}
