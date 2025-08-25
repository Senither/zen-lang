package evaluator

import (
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

func testEval(input string) objects.Object {
	l := lexer.New(input)
	p := parser.New(l)

	return Eval(p.ParseProgram())
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
