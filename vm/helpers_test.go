package vm

import (
	"testing"

	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		Stdout.Clear()

		compiler, err := compile(tt.input)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(compiler.Bytecode())
		vm.EnableStdoutCapture()

		stackElem := Stdout.Mute(func() objects.Object {
			err = vm.Run()
			if err != nil {
				t.Fatalf("VM run error: %s", err)
			}

			return vm.LastPoppedStackElem()
		})

		objects.AssertExpectedObject(t, tt.expected, stackElem)
	}
}

func runVmBenchmark(b *testing.B, input string) {
	compiler, err := compile(input)
	if err != nil {
		b.Fatalf("compiler error: %s", err)
	}

	for b.Loop() {
		New(compiler.Bytecode()).Run()
	}
}

func compile(input string) (*compiler.Compiler, error) {
	l := lexer.New(input)
	p := parser.New(l, nil)
	c := compiler.New(nil)

	return c, c.Compile(p.ParseProgram())
}
