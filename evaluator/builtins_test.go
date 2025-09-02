package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestPrintBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`print("hello world")`, "hello world"},
		{`print(5)`, "5"},
		{`print(5.5)`, "5.500000"},
		{`print(true)`, "true"},
		{`print(false)`, "false"},
		{`print("hello", "world")`, "helloworld"},
	}

	for _, tt := range tests {
		Stdout.Clear()
		Stdout.Mute(func() objects.Object {
			return testEval(tt.input)
		})

		output := Stdout.ReadAll()
		if len(output) != 1 {
			t.Errorf("expected 1 lines of output, got %d for %q\n\tOutput: %q", len(output), tt.input, output)
			return
		}

		if output[0] != tt.expected {
			t.Errorf("expected output to be %q, got %q", tt.expected, output[0])
		}
	}
}

func TestPrintlnBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`println("hello world")`, "hello world\n"},
		{`println(5)`, "5\n"},
		{`println(5.5)`, "5.500000\n"},
		{`println(true)`, "true\n"},
		{`println(false)`, "false\n"},
		{`println("hello", "world")`, "hello\nworld\n"},
	}

	for _, tt := range tests {
		Stdout.Clear()
		Stdout.Mute(func() objects.Object {
			return testEval(tt.input)
		})

		output := Stdout.ReadAll()
		if len(output) != 1 {
			t.Errorf("expected 1 lines of output, got %d for %q\n\tOutput: %q", len(output), tt.input, output)
			return
		}

		if output[0] != tt.expected {
			t.Errorf("expected output to be %q, got %q", tt.expected, output[0])
		}
	}
}

func TestLenBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER\n    at <unknown>:1:4"},
		{`len("one", "two")`, "wrong number of arguments. got 2, want 1\n    at <unknown>:1:4"},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*objects.Error)
			if !ok {
				t.Errorf("object is not Error. got %T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Inspect() != expected {
				t.Errorf("wrong error message. expected %q, got %q", expected, errObj.Inspect())
			}
		}
	}
}
