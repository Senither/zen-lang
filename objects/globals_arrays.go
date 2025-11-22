package objects

import "fmt"

func globalArraysPush(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, fmt.Errorf("argument to `push` must be an array, got %s", args[0].Type())
	}

	array.Elements = append(array.Elements, args[1])

	return array, nil
}

func globalArraysShift(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, fmt.Errorf("argument to `shift` must be an array, got %s", args[0].Type())
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
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, fmt.Errorf("argument to `pop` must be an array, got %s", args[0].Type())
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
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	array, ok := args[0].(*Array)
	if !ok {
		return nil, fmt.Errorf("argument to `filter` must be an array, got %s", args[0].Type())
	}

	callable, ok := args[1].(Callable)
	if !ok {
		return nil, fmt.Errorf("second argument to `filter` must be a function, got %s", args[1].Type())
	}

	if callable.ParametersCount() != 1 {
		return nil, fmt.Errorf("function passed to `filter` must take exactly one argument")
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
			return nil, fmt.Errorf("function passed to `filter` must return a boolean, got %s", rs.Type())
		}
	}

	return &Array{Elements: filtered}, nil
}

func globalArraysConcat(args ...Object) (Object, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want at least 2", len(args))
	}

	var elements []Object
	for _, arg := range args {
		array, ok := arg.(*Array)
		if !ok {
			return nil, fmt.Errorf("all arguments to `concat` must be arrays, got %s", arg.Type())
		}

		elements = append(elements, array.Elements...)
	}

	return &Array{Elements: elements}, nil
}
