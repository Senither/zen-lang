package objects

func globalArraysPush(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("push", 2, len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("push", ARRAY_OBJ, 0, args)
	}

	array.Elements = append(array.Elements, args[1])

	return array, nil
}

func globalArraysShift(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("shift", 1, len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("shift", ARRAY_OBJ, 0, args)
	}

	if len(array.Elements) == 0 {
		return NULL, nil
	}

	first := array.Elements[0]
	array.Elements = array.Elements[1:]

	return first, nil
}

func globalArraysPop(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("pop", 1, len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("pop", ARRAY_OBJ, 0, args)
	}

	if len(array.Elements) == 0 {
		return NULL, nil
	}

	last := array.Elements[len(array.Elements)-1]
	array.Elements = array.Elements[:len(array.Elements)-1]

	return last, nil
}

func globalArraysFilter(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("filter", 2, len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("filter", ARRAY_OBJ, 0, args)
	}

	callable, ok := args[1].(Callable)
	if !ok {
		return nil, NewInvalidArgumentTypesError("filter", []ObjectType{FUNCTION_OBJ, CLOSURE_OBJ}, 1, args)
	}

	if callable.ParametersCount() != 1 {
		return nil, NewErrorf("filter", "function passed to `filter` must take exactly one argument")
	}

	filtered := make([]Object, 0)
	for _, elem := range array.Elements {
		rs := callable.Call(elem)

		switch rs := rs.(type) {
		case *Boolean:
			if rs == TRUE {
				filtered = append(filtered, elem)
			}
		case *Error:
			return rs, nil

		default:
			return nil, NewErrorf("filter", "function passed to `filter` must return a %s, got %s", BOOLEAN_OBJ, rs.Type())
		}
	}

	return &Array{Elements: filtered}, nil
}

func globalArraysConcat(args ...Object) (Object, error) {
	if len(args) < 2 {
		return nil, NewWrongNumberOfArgumentsError("concat", 2, len(args))
	}

	var elements []Object
	for i := 0; i < len(args); i++ {
		array, ok := args[i].(*Array)
		if !ok {
			return nil, NewInvalidArgumentTypeError("concat", ARRAY_OBJ, i, args)
		}

		elements = append(elements, array.Elements...)
	}

	return &Array{Elements: elements}, nil
}
