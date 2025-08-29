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
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got 2, want 1"},
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

			if errObj.Message != expected {
				t.Errorf("wrong error message. expected %q, got %q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayPushBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"array_push([], 1)", []any{1}},
		{"array_push([1], 2)", []any{1, 2}},
		{"array_push([1, 2], 3)", []any{1, 2, 3}},
		{"var x = []; array_push(x, 1);", []any{1}},
		{"var x = [1]; array_push(x, 2);", []any{1, 2}},
		{"var x = [1, 2]; array_push(x, 3);", []any{1, 2, 3}},
		{"var x = []; array_push(x, 1); x;", []any{1}},
		{"var x = [1]; array_push(x, 2); x;", []any{1, 2}},
		{"var x = [1, 2]; array_push(x, 3); x;", []any{1, 2, 3}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestArrayShiftBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"array_shift([])", nil},
		{"array_shift([1])", 1},
		{"array_shift([1, 2])", 1},
		{"var x = []; array_shift(x);", nil},
		{"var x = [1]; array_shift(x);", 1},
		{"var x = [1, 2]; array_shift(x);", 1},
		{"var x = []; array_shift(x); x;", []int64{}},
		{"var x = [1]; array_shift(x); x;", []int64{}},
		{"var x = [1, 2]; array_shift(x); x;", []int64{2}},
		{"var x = [1, 2, 3]; array_shift(x); x;", []int64{2, 3}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case []int64:
			converted := make([]any, len(expected))
			for i, v := range expected {
				converted[i] = v
			}
			testArrayObject(t, evaluated, converted, tt.input)
		}
	}
}

func TestArrayPopBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"array_pop([])", nil},
		{"array_pop([1])", 1},
		{"array_pop([1, 2])", 2},
		{"var x = []; array_pop(x);", nil},
		{"var x = [1]; array_pop(x);", 1},
		{"var x = [1, 2]; array_pop(x);", 2},
		{"var x = []; array_pop(x); x;", []int64{}},
		{"var x = [1]; array_pop(x); x;", []int64{}},
		{"var x = [1, 2]; array_pop(x); x;", []int64{1}},
		{"var x = [1, 2, 3]; array_pop(x); x;", []int64{1, 2}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case nil:
			testNullObject(t, evaluated)
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case []int64:
			converted := make([]any, len(expected))
			for i, v := range expected {
				converted[i] = v
			}
			testArrayObject(t, evaluated, converted, tt.input)
		}
	}
}
