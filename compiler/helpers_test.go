package compiler

import (
	"fmt"
	"strings"
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

		compiler := New(nil)
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

func runCompilationBenchmarks(b *testing.B, benchmarks []string) {
	programs := make([]*ast.Program, len(benchmarks))

	for i, input := range benchmarks {
		programs[i] = parse(input)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, program := range programs {
			err := New(nil).Compile(program)
			if err != nil {
				b.Fatalf("compiler error: %s", err)
			}
		}
	}
}

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l, nil)

	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		var buf strings.Builder
		for _, msg := range p.Errors() {
			fmt.Fprintf(&buf, "%s\n", msg.String())
		}

		panic("parser errors encountered\n" + buf.String())
	}

	return program
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
			err := objects.AssertInteger(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - integer assertion failed: %s", i, err)
			}
		case float64:
			err := objects.AssertFloat(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - float assertion failed: %s", i, err)
			}
		case string:
			err := objects.AssertString(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - string assertion failed: %s", i, err)
			}
		case []code.Instructions:
			err := testCodeInstructions(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant %d - code instructions assertion failed: %s", i, err)
			}

		default:
			return fmt.Errorf("unknown constant type %T", constant)
		}
	}

	return nil
}

func testCodeInstructions(expected []code.Instructions, actual objects.Object) error {
	fn, ok := actual.(*objects.CompiledFunction)
	if !ok {
		return fmt.Errorf("object is not CompiledFunction. got %T (%+v)", actual, actual)
	}

	err := testInstructions(expected, fn.Instructions())
	if err != nil {
		return fmt.Errorf("instructions do not match the CompiledFunction instructions: %s", err)
	}

	return nil
}
