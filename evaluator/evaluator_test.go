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
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
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
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*objects.Error)
		if !ok {
			t.Errorf("no error object returned. got %T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected %q, got %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func (x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*objects.Function)
	if !ok {
		t.Fatalf("object is not Function. got %T (%+v)", evaluated, evaluated)
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
		t.Fatalf("object is not Function. got %T (%+v)", evaluated, evaluated)
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
	p := parser.New(l)

	return Eval(p.ParseProgram(), objects.NewEnvironment())
}

func testNullObject(t *testing.T, obj objects.Object) bool {
	if obj == NULL {
		return true
	}

	t.Errorf("object is not NULL. got %T (%+v)", obj, obj)
	return false
}

func testStringObject(t *testing.T, obj objects.Object, expected string) bool {
	result, ok := obj.(*objects.String)
	if !ok {
		t.Errorf("object is not String. got %T (%+v)", obj, obj)
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
		t.Errorf("object is not Integer. got %T (%+v)", obj, obj)
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
		t.Errorf("object is not Float. got %T (%+v)", obj, obj)
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
		t.Errorf("object is not Boolean. got %T (%+v)", obj, obj)
		return false
	}

	if boolean.Value != expected {
		t.Errorf("object has wrong value. got %t, expected %t, input %q", boolean.Value, expected, input)
		return false
	}

	return true
}
