package compiler

import (
	"fmt"
	"testing"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func runCompilationTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)
		if err != nil {
			t.Fatalf("instruction test failed: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)
		if err != nil {
			t.Fatalf("constants test failed: %s", err)
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l, nil)

	return p.ParseProgram()
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	combined := concatInstructions(expected)
	if len(actual) != len(combined) {
		return fmt.Errorf("wrong instructions length.\n\twant:\n%s\n\tgot:\n%s", combined, actual)
	}

	for i, ins := range combined {
		if actual[i] != ins {
			return fmt.Errorf("wrong instruction at %d.\n\tinstruction: %d\n\twant:\n%s\n\tgot:\n%s", i, ins, combined, actual)
		}
	}

	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, ins := range s {
		out = append(out, ins...)
	}

	return out
}

func testConstants(t *testing.T, expected []interface{}, actual []objects.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants. got %d, want %d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testIntegerObject failed: %s", i, err)
			}
		case float64:
			err := testFloatObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testFloatObject failed: %s", i, err)
			}
		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - testStringObject failed: %s", i, err)
			}

		default:
			return fmt.Errorf("unknown constant type %T", constant)
		}
	}

	return nil
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
