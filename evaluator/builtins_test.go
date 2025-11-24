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
		{`print(null)`, "null"},
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
		{`println(null)`, "null\n"},
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
		{`len(null)`, 0},
		{`len(1)`, "argument 1 to `len` has invalid type: got INTEGER, want STRING|ARRAY|NULL\n    at <unknown>:1:4"},
		{`len("one", "two")`, "wrong number of arguments to `len`: got 2, want 1\n    at <unknown>:1:4"},
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
				t.Errorf("wrong error message.\nexpected:\n\t%q\ngot:\n\t%q", expected, errObj.Inspect())
			}
		}
	}
}

func TestStringBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`string(123)`, "123"},
		{`string(45.67)`, "45.670000"},
		{`string(true)`, "true"},
		{`string(false)`, "false"},
		{`string("hello")`, "hello"},
		{`string("world")`, "world"},
		{`string(1 + 2)`, "3"},
		{`string(3.14 * 2)`, "6.280000"},
		{`string(!true)`, "false"},
		{`string("foo" + "bar")`, "foobar"},
		{`string("foo" + " " + "bar")`, "foo bar"},
		{`string("foo" + " " + "bar" + "!")`, "foo bar!"},
		{`string(null)`, "null"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestIntBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`int(123)`, int64(123)},
		{`int(45.67)`, int64(45)},
		{`int("789")`, int64(789)},
		{`int(0)`, int64(0)},
		{`int(-42)`, int64(-42)},
		{`int(-3.99)`, int64(-3)},
		{`int("0")`, int64(0)},
		{`int("   456   ")`, int64(456)},
		{`int(true)`, int64(1)},
		{`int(false)`, int64(0)},
		{`int(null)`, int64(0)},
		{`int("")`, "error in `int`: failed to convert `` to INTEGER\n    at <unknown>:1:4"},
		{`int("not a number")`, "error in `int`: failed to convert `not a number` to INTEGER\n    at <unknown>:1:4"},
		{`int(1, 2)`, "wrong number of arguments to `int`: got 2, want 1\n    at <unknown>:1:4"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int64:
			testIntegerObject(t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*objects.Error)
			if !ok {
				t.Errorf("object is not Error. got %T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Inspect() != expected {
				t.Errorf("wrong error message.\ngot:  %q\nwant: %q", errObj.Inspect(), expected)
			}
		}
	}
}

func TestFloatBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`float(123)`, float64(123)},
		{`float(45.67)`, float64(45.67)},
		{`float("789.01")`, float64(789.01)},
		{`float(0)`, float64(0)},
		{`float(-42)`, float64(-42)},
		{`float(-3.99)`, float64(-3.99)},
		{`float("0")`, float64(0)},
		{`float("   456.78   ")`, float64(456.78)},
		{`float(true)`, float64(1)},
		{`float(false)`, float64(0)},
		{`float(null)`, float64(0)},
		{`float("")`, "error in `float`: failed to convert `` to FLOAT\n    at <unknown>:1:6"},
		{`float("not a number")`, "error in `float`: failed to convert `not a number` to FLOAT\n    at <unknown>:1:6"},
		{`float(1, 2)`, "wrong number of arguments to `float`: got 2, want 1\n    at <unknown>:1:6"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case float64:
			testFloatObject(t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*objects.Error)
			if !ok {
				t.Errorf("object is not Error. got %T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Inspect() != expected {
				t.Errorf("wrong error message.\ngot:  %q\nwant: %q", errObj.Inspect(), expected)
			}
		}
	}
}

func TestTypeBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`type(123)`, "INTEGER"},
		{`type(45.67)`, "FLOAT"},
		{`type(true)`, "BOOLEAN"},
		{`type(false)`, "BOOLEAN"},
		{`type("hello")`, "STRING"},
		{`type(null)`, "NULL"},
		{`type([1, 2, 3])`, "ARRAY"},
		{`type({"key": "value"})`, "HASH"},
		{`type(func(x) { return x; })`, "FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestIsNaNBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`isNaN(0/0)`, true},
		{`isNaN(1/0)`, false},
		{`isNaN(-1/0)`, false},
		{`isNaN(5)`, false},
		{`isNaN(3.14)`, false},
		{`isNaN("not a number")`, false},
		{`isNaN(true)`, false},
		{`isNaN(null)`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}
