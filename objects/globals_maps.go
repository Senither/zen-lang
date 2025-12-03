package objects

import (
	"maps"
)

func globalMapsKeys(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("keys", 1, len(args))
	}

	obj, ok := args[0].(*Hash)
	if !ok {
		return nil, NewInvalidArgumentTypeError("keys", HASH_OBJ, 0, args)
	}

	keys := make([]Object, 0, len(obj.Pairs))
	for k := range obj.Pairs {
		keys = append(keys, obj.Pairs[k].Key)
	}

	return &Array{Elements: keys}, nil
}

func globalMapsValues(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("values", 1, len(args))
	}

	obj, ok := args[0].(*Hash)
	if !ok {
		return nil, NewInvalidArgumentTypeError("values", HASH_OBJ, 0, args)
	}

	values := make([]Object, 0, len(obj.Pairs))
	for _, pair := range obj.Pairs {
		values = append(values, pair.Value)
	}

	return &Array{Elements: values}, nil
}

func globalMapsHas(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("has", 2, len(args))
	}

	obj, ok := args[0].(*Hash)
	if !ok {
		return nil, NewInvalidArgumentTypeError("has", HASH_OBJ, 0, args)
	}

	key, ok := args[1].(Hashable)
	if !ok {
		return nil, NewInvalidArgumentTypesError("has", []ObjectType{
			STRING_OBJ, INTEGER_OBJ, FLOAT_OBJ, BOOLEAN_OBJ,
		}, 1, args)
	}

	_, exists := obj.Pairs[key.HashKey()]

	if !exists {
		return FALSE, nil
	}

	return TRUE, nil
}

func globalMapsEach(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("each", 2, len(args))
	}

	obj, ok := args[0].(*Hash)
	if !ok {
		return nil, NewInvalidArgumentTypeError("each", HASH_OBJ, 0, args)
	}

	fn, ok := args[1].(Callable)
	if !ok {
		return nil, NewInvalidArgumentTypeError("each", FUNCTION_OBJ, 1, args)
	}

	if fn.ParametersCount() != 2 {
		return nil, NewErrorf("each", "function must take exactly 2 parameters")
	}

	for _, pair := range obj.Pairs {
		rs := fn.Call(pair.Key, pair.Value)
		if IsError(rs) {
			return rs, nil
		}
	}

	return NULL, nil
}

func globalMapsMerge(args ...Object) (Object, error) {
	if len(args) < 2 {
		return nil, NewWrongNumberOfArgumentsWantAtLeastError("merge", 2, len(args))
	}

	result := &Hash{Pairs: make(map[HashKey]HashPair)}

	for i := range args {
		otherHash, ok := args[i].(*Hash)
		if !ok {
			return nil, NewInvalidArgumentTypeError("merge", HASH_OBJ, i, args)
		}

		result = deepMerge(result, otherHash)
	}

	return result, nil
}

func deepMerge(h1, h2 *Hash) *Hash {
	result := &Hash{Pairs: make(map[HashKey]HashPair)}
	maps.Copy(result.Pairs, h1.Pairs)

	for key, pair2 := range h2.Pairs {
		if _, exists := result.Pairs[key]; !exists {
			result.Pairs[key] = pair2
			continue
		}

		switch v1 := result.Pairs[key].Value.(type) {
		case *Null:
			result.Pairs[key] = pair2
		case *Hash:
			if v2, ok := pair2.Value.(*Hash); ok {
				merged := deepMerge(v1, v2)
				result.Pairs[key] = HashPair{Key: pair2.Key, Value: merged}
			}
		}
	}

	return result
}
