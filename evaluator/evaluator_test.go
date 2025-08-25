package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

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
