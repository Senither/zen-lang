package objects

import (
	"fmt"
	"testing"
)

func AssertExpectedObject(t *testing.T, expected interface{}, actual Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := AssertInteger(int64(expected), actual)
		if err != nil {
			t.Errorf("integer assertion failed: %s", err)
		}
	case int64:
		err := AssertInteger(expected, actual)
		if err != nil {
			t.Errorf("integer assertion failed: %s", err)
		}
	case float64:
		err := AssertFloat(expected, actual)
		if err != nil {
			t.Errorf("float assertion failed: %s", err)
		}
	case bool:
		err := AssertBoolean(expected, actual)
		if err != nil {
			t.Errorf("boolean assertion failed: %s", err)
		}
	case string:
		err := AssertString(expected, actual)
		if err != nil {
			t.Errorf("string assertion failed: %s", err)
		}
	case []int:
		arr := make([]any, len(expected))
		for i, v := range expected {
			arr[i] = v
		}

		err := AssertArray(arr, actual)
		if err != nil {
			t.Errorf("integer array assertion failed: %s", err)
		}
	case []float64:
		arr := make([]any, len(expected))
		for i, v := range expected {
			arr[i] = v
		}

		err := AssertArray(arr, actual)
		if err != nil {
			t.Errorf("float array assertion failed: %s", err)
		}
	case []bool:
		arr := make([]any, len(expected))
		for i, v := range expected {
			arr[i] = v
		}

		err := AssertArray(arr, actual)
		if err != nil {
			t.Errorf("boolean array assertion failed: %s", err)
		}
	case []string:
		arr := make([]any, len(expected))
		for i, v := range expected {
			arr[i] = v
		}

		err := AssertArray(arr, actual)
		if err != nil {
			t.Errorf("string array assertion failed: %s", err)
		}
	case []any:
		err := AssertArray(expected, actual)
		if err != nil {
			t.Errorf("any array assertion failed: %s", err)
		}
	case map[string]int:
		hash := make(map[string]any)
		for key, value := range expected {
			hash[key] = value
		}

		err := AssertMapObject(hash, actual)
		if err != nil {
			t.Errorf("string map assertion failed: %s", err)
		}
	case map[HashKey]int64:
		hash := make(map[HashKey]any)
		for key, value := range expected {
			hash[key] = value
		}

		err := AssertHashKeyMapObject(hash, actual)
		if err != nil {
			t.Errorf("hash key map assertion failed: %s", err)
		}
	case map[HashKey]any:
		err := AssertHashKeyMapObject(expected, actual)
		if err != nil {
			t.Errorf("hash key map assertion failed: %s", err)
		}
	case *Error:
		actualError, ok := actual.(*Error)
		if !ok {
			t.Errorf("object is not Error. got %T (%+v)", actual, actual)
			return
		}

		unwrappedErrorMessage := unwrapErrorMessage(actualError)
		if unwrappedErrorMessage != expected.Message {
			t.Errorf("wrong error message.\ngot:\n%s\nwant:\n%s", unwrappedErrorMessage, expected.Message)
		}

	case nil:
		if actual != NULL {
			t.Errorf("object is not NULL. got %T (%+v)", actual, actual)
		}

	case Error:
		t.Errorf("unsupported assertion type objects.Error, errors must be passed as pointers")

	default:
		t.Errorf("unsupported assertion type %T", expected)
	}
}

func AssertInteger(expected int64, actual Object) error {
	result, ok := actual.(*Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %d, want %d", result.Value, expected)
	}

	return nil
}

func AssertFloat(expected float64, actual Object) error {
	result, ok := actual.(*Float)
	if !ok {
		return fmt.Errorf("object is not Float. got %T (%+v)", actual, actual)
	}

	if result.Inspect() != fmt.Sprintf("%f", expected) {
		return fmt.Errorf("object has wrong value. got %f, expected %f", result.Value, expected)
	}

	return nil
}

func AssertBoolean(expected bool, actual Object) error {
	result, ok := actual.(*Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %t, expected %t", result.Value, expected)
	}

	return nil
}

func AssertString(expected string, actual Object) error {
	result, ok := actual.(*String)
	if !ok {
		return fmt.Errorf("object is not String. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %q, expected %q", result.Value, expected)
	}

	return nil
}

func AssertArray(expected []any, actual Object) error {
	array, ok := actual.(*Array)
	if !ok {
		return fmt.Errorf("object is not Array. got %T (%+v)", actual, actual)
	}

	if len(array.Elements) != len(expected) {
		return fmt.Errorf("array has wrong length. got %d, want %d", len(array.Elements), len(expected))
	}

	for i, expectedElem := range expected {
		err := assertInterfaceMatchesActual(expectedElem, array.Elements[i])
		if err != nil {
			return fmt.Errorf("array[%d] - %s", i, err)
		}
	}

	return nil
}

func AssertMapObject(expected map[string]any, actual Object) error {
	hash, ok := actual.(*Hash)
	if !ok {
		return fmt.Errorf("object is not Hash. got %T (%+v)", actual, actual)
	}

	if len(hash.Pairs) != len(expected) {
		return fmt.Errorf("hash has wrong number of pairs. got %d, want %d", len(hash.Pairs), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		keyObj := &String{Value: expectedKey}
		hashKey := keyObj.HashKey()

		pair, ok := hash.Pairs[hashKey]
		if !ok {
			return fmt.Errorf("no pair found for given key in Pairs: %q", expectedKey)
		}

		err := assertInterfaceMatchesActual(expectedValue, pair.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func AssertHashKeyMapObject(expected map[HashKey]any, actual Object) error {
	hash, ok := actual.(*Hash)
	if !ok {
		return fmt.Errorf("object is not Hash. got %T (%+v)", actual, actual)
	}

	if len(hash.Pairs) != len(expected) {
		return fmt.Errorf("hash has wrong number of pairs. got %d, want %d", len(hash.Pairs), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := hash.Pairs[expectedKey]
		if !ok {
			return fmt.Errorf("no pair found for given key in Pairs: %d", expectedKey.Value)
		}

		err := assertInterfaceMatchesActual(expectedValue, pair.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

func assertInterfaceMatchesActual(expectedValue any, actual Object) error {
	switch expectedValue := expectedValue.(type) {
	case int:
		err := AssertInteger(int64(expectedValue), actual)
		if err != nil {
			return fmt.Errorf("integer assertion failed: %s", err)
		}
	case int64:
		err := AssertInteger(expectedValue, actual)
		if err != nil {
			return fmt.Errorf("integer assertion failed: %s", err)
		}
	case bool:
		err := AssertBoolean(expectedValue, actual)
		if err != nil {
			return fmt.Errorf("boolean assertion failed: %s", err)
		}
	case float64:
		err := AssertFloat(expectedValue, actual)
		if err != nil {
			return fmt.Errorf("float assertion failed: %s", err)
		}
	case string:
		err := AssertString(expectedValue, actual)
		if err != nil {
			return fmt.Errorf("string assertion failed: %s", err)
		}
	case nil:
		if actual != NULL {
			return fmt.Errorf("object is not NULL. got %T (%+v)", actual, actual)
		}

	default:
		return fmt.Errorf("unsupported assertion type %T", expectedValue)
	}

	return nil
}

func unwrapErrorMessage(err *Error) string {
	if err.Parent != nil {
		return unwrapErrorMessage(err.Parent)
	}

	return err.Message
}
