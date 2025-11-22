package objects

import (
	"fmt"
	"strings"
)

func globalStringsContains(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, fmt.Errorf("argument to `contains` must be a string, got %s", args[0].Type())
	}

	substr, ok := args[1].(*String)
	if !ok {
		return nil, fmt.Errorf("second argument to `contains` must be a string, got %s", args[1].Type())
	}

	if strings.Contains(str.Value, substr.Value) {
		return TRUE, nil
	}

	return FALSE, nil
}

func globalStringsSplit(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, fmt.Errorf("argument to `contains` must be a string, got %s", args[0].Type())
	}

	substr, ok := args[1].(*String)
	if !ok {
		return nil, fmt.Errorf("second argument to `contains` must be a string, got %s", args[1].Type())
	}

	arr := strings.Split(str.Value, substr.Value)

	var elements []Object
	for _, s := range arr {
		elements = append(elements, &String{Value: s})
	}

	return &Array{Elements: elements}, nil
}

func globalStringsJoin(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
	}

	arr, ok := args[0].(*Array)
	if !ok {
		return nil, fmt.Errorf("argument to `join` must be an array, got %s", args[0].Type())
	}

	sep, ok := args[1].(*String)
	if !ok {
		return nil, fmt.Errorf("second argument to `join` must be a string, got %s", args[1].Type())
	}

	var elements []string
	for _, elem := range arr.Elements {
		switch elem := elem.(type) {
		case *Float:
			elements = append(elements, fmt.Sprintf("%v", elem.Value))

		default:
			elements = append(elements, elem.Inspect())
		}
	}

	return &String{Value: strings.Join(elements, sep.Value)}, nil
}

func globalStringsFormat(args ...Object) (Object, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("wrong number of arguments. got %d, want at least 2", len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, fmt.Errorf("argument to `contains` must be a string, got %s", args[0].Type())
	}

	var values []any
	for _, arg := range args[1:] {
		switch arg := arg.(type) {
		case *String:
			values = append(values, arg.Value)
		case *Integer:
			values = append(values, arg.Value)
		case *Float:
			values = append(values, arg.Value)
		case *Boolean:
			values = append(values, arg.Value)
		case *Null:
			values = append(values, nil)

		default:
			values = append(values, arg.Inspect())
		}
	}

	return &String{Value: fmt.Sprintf(str.Value, values...)}, nil
}
