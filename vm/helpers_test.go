package vm

import (
	"fmt"
	"testing"

	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		compiler, err := compile(tt.input)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(compiler.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("VM run error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()
		testExpectedObject(t, tt.expected, stackElem)
	}
}

func compile(input string) (*compiler.Compiler, error) {
	l := lexer.New(input)
	p := parser.New(l, nil)
	c := compiler.New()

	return c, c.Compile(p.ParseProgram())
}

func testExpectedObject(t *testing.T, expected interface{}, actual objects.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case float64:
		err := testFloatObject(expected, actual)
		if err != nil {
			t.Errorf("testFloatObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}

	case nil:
		if actual != NULL {
			t.Errorf("object is not NULL. got %T (%+v)", actual, actual)
		}

	default:
		t.Errorf("unsupported type %T", expected)
	}
}

func testIntegerObject(expected int64, actual objects.Object) error {
	result, ok := actual.(*objects.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %d, want %d", result.Value, expected)
	}

	return nil
}

func testFloatObject(expected float64, actual objects.Object) error {
	result, ok := actual.(*objects.Float)
	if !ok {
		return fmt.Errorf("object is not Float. got %T (%+v)", actual, actual)
	}

	if result.Inspect() != fmt.Sprintf("%f", expected) {
		return fmt.Errorf("object has wrong value. got %f, expected %f", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual objects.Object) error {
	result, ok := actual.(*objects.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %t, expected %t", result.Value, expected)
	}

	return nil
}
