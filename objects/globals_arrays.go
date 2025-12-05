package objects

import "sort"

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
		return nil, NewWrongNumberOfArgumentsWantAtLeastError("concat", 2, len(args))
	}

	var elements []Object
	for i := range args {
		array, ok := args[i].(*Array)
		if !ok {
			return nil, NewInvalidArgumentTypeError("concat", ARRAY_OBJ, i, args)
		}

		elements = append(elements, array.Elements...)
	}

	return &Array{Elements: elements}, nil
}

func globalArraysFlatten(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("flatten", 1, len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("flatten", ARRAY_OBJ, 0, args)
	}

	return &Array{Elements: flattenArray(array)}, nil
}

func flattenArray(arr *Array) []Object {
	var result []Object

	for _, elem := range arr.Elements {
		if nestedArr, ok := elem.(*Array); ok {
			result = append(result, flattenArray(nestedArr)...)
		} else {
			result = append(result, elem)
		}
	}

	return result
}

func globalArraysFirst(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("first", 2, len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("first", ARRAY_OBJ, 0, args)
	}

	callable, ok := args[1].(Callable)
	if !ok {
		return nil, NewInvalidArgumentTypesError("first", []ObjectType{FUNCTION_OBJ, CLOSURE_OBJ}, 1, args)
	}

	if callable.ParametersCount() == 0 {
		return nil, NewErrorf("first", "function passed to `first` must take at least one argument")
	}

	if callable.ParametersCount() > 2 {
		return nil, NewErrorf("first", "function passed to `first` must take at most two arguments")
	}

	for i, elem := range array.Elements {
		var rs Object

		if callable.ParametersCount() == 1 {
			rs = callable.Call(elem)
		} else {
			rs = callable.Call(elem, &Integer{Value: int64(i)})
		}

		switch rs := rs.(type) {
		case *Boolean:
			if rs == TRUE {
				return elem, nil
			}
		case *Error:
			return rs, nil

		default:
			return nil, NewErrorf("first", "function passed to `first` must return a %s, got %s", BOOLEAN_OBJ, rs.Type())
		}
	}

	return NULL, nil
}

func globalArraysSort(args ...Object) (Object, error) {
	if len(args) == 0 {
		return nil, NewWrongNumberOfArgumentsWantAtLeastError("sort", 1, len(args))
	}

	arr, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("sort", ARRAY_OBJ, 0, args)
	}

	sortedArr := &Array{Elements: make([]Object, len(arr.Elements))}
	copy(sortedArr.Elements, arr.Elements)
	arr = sortedArr

	if len(args) == 1 {
		sort.Sort(sortedArr)
		return sortedArr, nil
	}

	callable, ok := args[1].(Callable)
	if !ok {
		return nil, NewInvalidArgumentTypesError("sort", []ObjectType{FUNCTION_OBJ, CLOSURE_OBJ}, 1, args)
	}

	if callable.ParametersCount() != 2 {
		return nil, NewErrorf("sort", "function passed to `sort` must take exactly two arguments")
	}

	sort.SliceStable(sortedArr.Elements, func(i, j int) bool {
		rs := callable.Call(sortedArr.Elements[i], sortedArr.Elements[j])

		switch rs := rs.(type) {
		case *Boolean:
			return rs == TRUE

		default:
			return false
		}
	})

	return sortedArr, nil
}
