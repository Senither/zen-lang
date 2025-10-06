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
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject failed: %s", err)
		}
	case []int:
		err := testIntegerArrayObject(expected, actual)
		if err != nil {
			t.Errorf("testIntegerArrayObject failed: %s", err)
		}
	case map[objects.HashKey]int64:
		err := testHashMapIntegerObject(expected, actual)
		if err != nil {
			t.Errorf("testHashMapObject failed: %s", err)
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

func testStringObject(expected string, actual objects.Object) error {
	result, ok := actual.(*objects.String)
	if !ok {
		return fmt.Errorf("object is not String. got %T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got %q, expected %q", result.Value, expected)
	}

	return nil
}

func testIntegerArrayObject(expected []int, actual objects.Object) error {
	array, ok := actual.(*objects.Array)
	if !ok {
		return fmt.Errorf("object is not Array. got %T (%+v)", actual, actual)
	}

	if len(array.Elements) != len(expected) {
		return fmt.Errorf("array has wrong length. got %d, want %d", len(array.Elements), len(expected))
	}

	for i, expectedElem := range expected {
		err := testIntegerObject(int64(expectedElem), array.Elements[i])
		if err != nil {
			return fmt.Errorf("array[%d] - %s", i, err)
		}
	}

	return nil
}

func testHashMapIntegerObject(expected map[objects.HashKey]int64, actual objects.Object) error {
	hash, ok := actual.(*objects.Hash)
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

		err := testIntegerObject(expectedValue, pair.Value)
		if err != nil {
			return fmt.Errorf("testIntegerObject failed: %s", err)
		}
	}

	return nil
}
