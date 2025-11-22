package objects

import (
	"fmt"
	"math"
)

func globalMathMin(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument 1 to `min` must be a number, got %s", args[0].Type())
	}

	if !IsNumber(args[1].Type()) {
		return nil, fmt.Errorf("argument 2 to `min` must be a number, got %s", args[1].Type())
	}

	return WrapNumberValue(math.Min(
		UnwrapNumberValue(args[0]),
		UnwrapNumberValue(args[1]),
	), args[0], args[1]), nil
}

func globalMathMax(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument 1 to `max` must be a number, got %s", args[0].Type())
	}

	if !IsNumber(args[1].Type()) {
		return nil, fmt.Errorf("argument 2 to `max` must be a number, got %s", args[1].Type())
	}

	return WrapNumberValue(math.Max(
		UnwrapNumberValue(args[0]),
		UnwrapNumberValue(args[1]),
	), args[0], args[1]), nil
}

func globalMathCeil(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument to `ceil` must be a number, got %s", args[0].Type())
	}

	return &Float{Value: math.Ceil(UnwrapNumberValue(args[0]))}, nil
}

func globalMathFloor(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument to `floor` must be a number, got %s", args[0].Type())
	}

	return &Float{Value: math.Floor(UnwrapNumberValue(args[0]))}, nil
}

func globalMathRound(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument to `round` must be a number, got %s", args[0].Type())
	}

	return &Float{Value: math.Round(UnwrapNumberValue(args[0]))}, nil
}

func globalMathLog(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument to `log` must be a number, got %s", args[0].Type())
	}

	return &Float{Value: math.Log10(UnwrapNumberValue(args[0]))}, nil
}

func globalMathSqrt(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, fmt.Errorf("argument to `sqrt` must be a number, got %s", args[0].Type())
	}

	return &Float{Value: math.Sqrt(UnwrapNumberValue(args[0]))}, nil
}
