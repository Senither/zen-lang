package objects

import (
	"fmt"
	"strings"
)

func globalStringsContains(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("contains", 2, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("contains", STRING_OBJ, 0, args)
	}

	substr, ok := args[1].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("contains", STRING_OBJ, 1, args)
	}

	if strings.Contains(str.Value, substr.Value) {
		return TRUE, nil
	}

	return FALSE, nil
}

func globalStringsSplit(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("split", 2, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("split", STRING_OBJ, 0, args)
	}

	substr, ok := args[1].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("split", STRING_OBJ, 1, args)
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
		return nil, NewWrongNumberOfArgumentsError("join", 2, len(args))
	}

	arr, ok := args[0].(*Array)
	if !ok {
		return nil, NewInvalidArgumentTypeError("join", ARRAY_OBJ, 0, args)
	}

	sep, ok := args[1].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("join", STRING_OBJ, 1, args)
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
		return nil, NewWrongNumberOfArgumentsError("format", 2, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("format", STRING_OBJ, 0, args)
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

func globalStringsStartsWith(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("startsWith", 2, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("startsWith", STRING_OBJ, 0, args)
	}

	switch prefix := args[1].(type) {
	case *String:
		if strings.HasPrefix(str.Value, prefix.Value) {
			return TRUE, nil
		}
		return FALSE, nil
	case *Array:
		for _, elem := range prefix.Elements {
			prefixStr, ok := elem.(*String)
			if !ok {
				return nil, NewErrorf("startsWith", "elements of the prefix array must be %s, got %s", STRING_OBJ, elem.Type())
			}

			if strings.HasPrefix(str.Value, prefixStr.Value) {
				return TRUE, nil
			}
		}
		return FALSE, nil

	default:
		return nil, NewInvalidArgumentTypesError("startsWith", []ObjectType{STRING_OBJ, ARRAY_OBJ}, 1, args)
	}
}

func globalStringsEndsWith(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("endsWith", 2, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("endsWith", STRING_OBJ, 0, args)
	}

	switch suffix := args[1].(type) {
	case *String:
		if strings.HasSuffix(str.Value, suffix.Value) {
			return TRUE, nil
		}
		return FALSE, nil
	case *Array:
		for _, elem := range suffix.Elements {
			prefixStr, ok := elem.(*String)
			if !ok {
				return nil, NewErrorf("endsWith", "elements of the suffix array must be %s, got %s", STRING_OBJ, elem.Type())
			}

			if strings.HasSuffix(str.Value, prefixStr.Value) {
				return TRUE, nil
			}
		}
		return FALSE, nil

	default:
		return nil, NewInvalidArgumentTypesError("endsWith", []ObjectType{STRING_OBJ, ARRAY_OBJ}, 1, args)
	}
}

func globalStringsToUpper(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("toUpper", 1, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("toUpper", STRING_OBJ, 0, args)
	}

	return &String{Value: strings.ToUpper(str.Value)}, nil
}

func globalStringsToLower(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("toLower", 1, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("toLower", STRING_OBJ, 0, args)
	}

	return &String{Value: strings.ToLower(str.Value)}, nil
}

func globalStringsTrim(args ...Object) (Object, error) {
	if len(args) == 0 {
		return nil, NewWrongNumberOfArgumentsError("trim", 1, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("trim", STRING_OBJ, 0, args)
	}

	if len(args) == 1 {
		return &String{Value: strings.TrimSpace(str.Value)}, nil
	}

	chars, ok := args[1].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("trim", STRING_OBJ, 1, args)
	}

	return &String{Value: strings.Trim(str.Value, chars.Value)}, nil
}
