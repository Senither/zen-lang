package evaluator

import (
	"fmt"
	"testing"

	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN\n    at <unknown>:1:3",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN\n    at <unknown>:1:3",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN\n    at <unknown>:1:1",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN\n    at <unknown>:1:6",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN\n    at <unknown>:1:6",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN\n    at <unknown>:1:9",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN\n    at <unknown>:1:20",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}
		`,
			"unknown operator: BOOLEAN + BOOLEAN\n    at <unknown>:4:18",
		},
		{
			"foobar",
			"identifier not found: foobar\n    at <unknown>:1:1",
		},
		{
			`{"name": "value"}[func (x) { x }]`,
			"invalid type given as hash key: FUNCTION\n    at <unknown>:1:18",
		},
		{
			`
			func a(x) {
				return x();
			}
			func b(y) {
				return y(a);
			}
			func c(z) {
				return z(b);
			}

			println(c(func () {
				return true + false;
			}));
			`,
			"unknown operator: BOOLEAN + BOOLEAN\n    at <unknown>:13:17\n    at <unknown>:9:13\n    at <unknown>:12:13\n    at <unknown>:12:11",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*objects.Error)
		if !ok {
			t.Errorf("no error object returned. got %T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Inspect() != tt.expectedMessage {
			t.Errorf("wrong error message. expected %q, got %q", tt.expectedMessage, errObj.Inspect())
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func (x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*objects.Function)
	if !ok {
		t.Fatalf("object is not Function. got %T (%+v)\n%s", evaluated, evaluated, evaluated.Inspect())
	}

	if fn.Name != nil {
		t.Fatalf("function name is not nil. got %q", fn.Name.String())
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, fn.Body.String())
	}
}

func TestNamedFunctionObject(t *testing.T) {
	input := "func hello(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*objects.Function)
	if !ok {
		t.Fatalf("object is not Function. got %T (%+v)\n%s", evaluated, evaluated, evaluated.Inspect())
	}

	if fn.Name == nil || fn.Name.String() != "hello" {
		t.Fatalf("function name is not 'hello', got %q", fn.Name)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var identity = func(x) { x; }; identity(5);", 5},
		{"func identity(x) { x; }; identity(5);", 5},
		{"var identity = func(x) { return x; }; identity(5);", 5},
		{"func identity(x) { return x; }; identity(5);", 5},
		{"var double = func(x) { x * 2; }; double(5);", 10},
		{"func double(x) { x * 2; }; double(5);", 10},
		{"var add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"func add(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func add(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func testEval(input string) objects.Object {
	l := lexer.New(input)
	p := parser.New(l, nil)

	return Eval(p.ParseProgram(), objects.NewEnvironment(nil))
}

func testNullObject(t *testing.T, obj objects.Object) bool {
	if obj == objects.NULL {
		return true
	}

	t.Errorf("object is not NULL. got %T (%+v)\n%s", obj, obj, obj.Inspect())
	return false
}

func testStringObject(t *testing.T, obj objects.Object, expected string) bool {
	result, ok := obj.(*objects.String)
	if !ok {
		t.Errorf("object is not String. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %q, expected %q", result.Value, expected)
		return false
	}

	return true
}

func testIntegerObject(t *testing.T, obj objects.Object, expected int64) bool {
	result, ok := obj.(*objects.Integer)
	if !ok {
		t.Errorf("object is not Integer. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got %d, expected %d", result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj objects.Object, expected float64) bool {
	result, ok := obj.(*objects.Float)
	if !ok {
		t.Errorf("object is not Float. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if result.Inspect() != fmt.Sprintf("%f", expected) {
		t.Errorf("object has wrong value. got %v, expected %v", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj objects.Object, expected bool, input string) bool {
	boolean, ok := obj.(*objects.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if boolean.Value != expected {
		t.Errorf("object has wrong value. got %t, expected %t, input %q", boolean.Value, expected, input)
		return false
	}

	return true
}

func testArrayObject(t *testing.T, obj objects.Object, expected []any, input string) bool {
	array, ok := obj.(*objects.Array)
	if !ok {
		t.Errorf("object is not Array. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if len(array.Elements) != len(expected) {
		t.Errorf("array has wrong number of elements. got %d, expected %d", len(array.Elements), len(expected))
		return false
	}

	for i, elem := range expected {
		switch elem := elem.(type) {
		case string:
			testStringObject(t, array.Elements[i], elem)
		case int:
			testIntegerObject(t, array.Elements[i], int64(elem))
		case int64:
			testIntegerObject(t, array.Elements[i], elem)
		case bool:
			testBooleanObject(t, array.Elements[i], elem, input)
		case nil:
			testNullObject(t, array.Elements[i])
		default:
			t.Errorf("element type is not support for array testing objects. got %T (%+v)\n%s", elem, elem, obj.Inspect())
			return false
		}
	}

	return true
}

func testHashObject(t *testing.T, obj objects.Object, expected map[objects.HashKey]any, input string) bool {
	hash, ok := obj.(*objects.Hash)
	if !ok {
		t.Errorf("object is not Hash. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if len(hash.Pairs) != len(expected) {
		t.Errorf("hash has wrong number of pairs. got %d, expected %d", len(hash.Pairs), len(expected))
		return false
	}

	for key, expectedValue := range expected {
		pair, ok := hash.Pairs[key]
		if !ok {
			t.Errorf("hash is missing key %q", key)
			return false
		}

		switch expectedValue := expectedValue.(type) {
		case string:
			testStringObject(t, pair.Value, expectedValue)
		case int:
			testIntegerObject(t, pair.Value, int64(expectedValue))
		case int64:
			testIntegerObject(t, pair.Value, expectedValue)
		case bool:
			testBooleanObject(t, pair.Value, expectedValue, input)
		default:
			t.Errorf("element type is not support for hash testing objects. got %T (%+v)\n%s", expectedValue, expectedValue, obj.Inspect())
			return false
		}
	}

	return true
}

func testErrorObject(t *testing.T, obj objects.Object, expected string) bool {
	err, ok := obj.(*objects.Error)
	if !ok {
		t.Errorf("object is not Error. got %T (%+v)\n%s", obj, obj, obj.Inspect())
		return false
	}

	if err.Inspect() != expected {
		t.Errorf("object has wrong message.\ngot:  %q\nwant: %q", err.Inspect(), expected)
		return false
	}

	return true
}
