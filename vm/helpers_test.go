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

		stackElem := vm.StackTop()
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
