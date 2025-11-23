package objects

import (
	"math"
)

func globalMathMin(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("min", 2, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("min", GetNumberTypes(), 0, args)
	}

	if !IsNumber(args[1].Type()) {
		return nil, NewInvalidArgumentTypesError("min", GetNumberTypes(), 1, args)
	}

	return WrapNumberValue(math.Min(
		UnwrapNumberValue(args[0]),
		UnwrapNumberValue(args[1]),
	), args[0], args[1]), nil
}

func globalMathMax(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("max", 2, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("max", GetNumberTypes(), 0, args)
	}

	if !IsNumber(args[1].Type()) {
		return nil, NewInvalidArgumentTypesError("max", GetNumberTypes(), 1, args)
	}

	return WrapNumberValue(math.Max(
		UnwrapNumberValue(args[0]),
		UnwrapNumberValue(args[1]),
	), args[0], args[1]), nil
}

func globalMathCeil(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("ceil", 1, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("ceil", GetNumberTypes(), 0, args)
	}

	return &Float{Value: math.Ceil(UnwrapNumberValue(args[0]))}, nil
}

func globalMathFloor(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("floor", 1, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("floor", GetNumberTypes(), 0, args)
	}

	return &Float{Value: math.Floor(UnwrapNumberValue(args[0]))}, nil
}

func globalMathRound(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("round", 1, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("round", GetNumberTypes(), 0, args)
	}

	return &Float{Value: math.Round(UnwrapNumberValue(args[0]))}, nil
}

func globalMathLog(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("log", 1, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("log", GetNumberTypes(), 0, args)
	}

	return &Float{Value: math.Log10(UnwrapNumberValue(args[0]))}, nil
}

func globalMathSqrt(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("sqrt", 1, len(args))
	}

	if !IsNumber(args[0].Type()) {
		return nil, NewInvalidArgumentTypesError("sqrt", GetNumberTypes(), 0, args)
	}

	return &Float{Value: math.Sqrt(UnwrapNumberValue(args[0]))}, nil
}
