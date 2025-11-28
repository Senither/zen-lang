package objects

import (
	"encoding/json"
	"strings"
)

func globalJSONParse(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("parse", 1, len(args))
	}

	str, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("parse", STRING_OBJ, 0, args)
	}

	jsonStr := strings.TrimSpace(str.Value)
	if jsonStr == "" {
		return NULL, nil
	}

	if strings.HasPrefix(jsonStr, "{") {
		var value map[string]any
		err := json.Unmarshal([]byte(jsonStr), &value)
		if err != nil {
			return nil, NewErrorf("parse", "%s", err.Error())
		}

		return mapToObject(value), nil
	} else if strings.HasPrefix(jsonStr, "[") {
		var value []any
		err := json.Unmarshal([]byte(jsonStr), &value)
		if err != nil {
			return nil, NewErrorf("parse", "%s", err.Error())
		}

		return arrayToObject(value), nil
	}

	return NULL, NewErrorf("parse", "failed to parse `%s` as JSON", str.Value)
}

func mapToObject(m map[string]any) Object {
	pairs := make(map[HashKey]HashPair)

	for k, v := range m {
		key := &String{Value: k}

		pairs[key.HashKey()] = HashPair{
			Key:   key,
			Value: nativeValueToObject(v),
		}
	}

	return &Hash{Pairs: pairs}
}

func arrayToObject(arr []any) Object {
	elements := make([]Object, len(arr))

	for i, v := range arr {
		elements[i] = nativeValueToObject(v)
	}

	return &Array{Elements: elements}
}

func nativeValueToObject(value any) Object {
	switch val := value.(type) {
	case string:
		return &String{Value: val}
	case int:
		return &Integer{Value: int64(val)}
	case float64:
		if val == float64(int64(val)) {
			return &Integer{Value: int64(val)}
		} else {
			return &Float{Value: val}
		}
	case bool:
		return &Boolean{Value: val}
	case map[string]any:
		return mapToObject(val)
	case []any:
		return arrayToObject(val)

	default:
		return NULL
	}
}

func globalJSONStringify(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("stringify", 1, len(args))
	}

	var data []byte
	var err error

	switch v := args[0].(type) {
	case *Hash:
		nativeMap := mapObjectToNative(v)
		data, err = json.Marshal(nativeMap)
	case *ImmutableHash:
		nativeMap := mapObjectToNative(&v.Value)
		data, err = json.Marshal(nativeMap)
	case *Array:
		nativeArray := arrayObjectToNative(v)
		data, err = json.Marshal(nativeArray)

	default:
		return nil, NewInvalidArgumentTypesError("stringify", []ObjectType{HASH_OBJ, ARRAY_OBJ}, 0, args)
	}

	if err != nil {
		return nil, NewErrorf("stringify", "%s", err.Error())
	}

	return &String{Value: string(data)}, nil
}

func mapObjectToNative(m *Hash) map[string]any {
	nativeMap := make(map[string]any)

	for _, pair := range m.Pairs {
		keyStr, ok := pair.Key.(*String)
		if !ok {
			continue
		}

		nativeMap[keyStr.Value] = objectToNativeValue(pair.Value)
	}

	return nativeMap
}

func arrayObjectToNative(a *Array) []any {
	nativeArray := make([]any, len(a.Elements))

	for i, elem := range a.Elements {
		nativeArray[i] = objectToNativeValue(elem)
	}

	return nativeArray
}

func objectToNativeValue(obj Object) any {
	switch val := obj.(type) {
	case *String:
		return val.Value
	case *Integer:
		return int64(val.Value)
	case *Float:
		return val.Value
	case *Boolean:
		return val.Value
	case *Hash:
		return mapObjectToNative(val)
	case *Array:
		return arrayObjectToNative(val)

	default:
		return nil
	}
}
