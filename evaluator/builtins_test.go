package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestPrintBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"print string", `print("hello world")`, "hello world"},
		{"print int", `print(5)`, "5"},
		{"print float", `print(5.5)`, "5.500000"},
		{"print true", `print(true)`, "true"},
		{"print false", `print(false)`, "false"},
		{"print multi string", `print("hello", "world")`, "helloworld"},
		{"print null", `print(null)`, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

func TestPrintlnBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"println string", `println("hello world")`, "hello world\n"},
		{"println int", `println(5)`, "5\n"},
		{"println float", `println(5.5)`, "5.500000\n"},
		{"println true", `println(true)`, "true\n"},
		{"println false", `println(false)`, "false\n"},
		{"println multi string", `println("hello", "world")`, "hello\nworld\n"},
		{"println null", `println(null)`, "null\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		})
	}
}

func TestLenBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"len empty string", `len("")`, 0},
		{"len four", `len("four")`, 4},
		{"len hello world", `len("hello world")`, 11},
		{"len null", `len(null)`, 0},
		{"len int", `len(1)`, &objects.Error{Message: "argument 1 to `len` has invalid type: got INTEGER, want STRING|ARRAY|NULL"}},
		{"len too many arguments", `len("one", "two")`, &objects.Error{Message: "wrong number of arguments to `len`: got 2, want 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestStringBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"string int", `string(123)`, "123"},
		{"string float", `string(45.67)`, "45.670000"},
		{"string true", `string(true)`, "true"},
		{"string false", `string(false)`, "false"},
		{"string hello", `string("hello")`, "hello"},
		{"string world", `string("world")`, "world"},
		{"string expression 1 + 2", `string(1 + 2)`, "3"},
		{"string expression 3.14 * 2", `string(3.14 * 2)`, "6.280000"},
		{"string expression !true", `string(!true)`, "false"},
		{"string expression \"foo\" + \"bar\"", `string("foo" + "bar")`, "foobar"},
		{"string expression \"foo\" + \" \" + \"bar\"", `string("foo" + " " + "bar")`, "foo bar"},
		{"string expression \"foo\" + \" \" + \"bar\" + \"!\"", `string("foo" + " " + "bar" + "!")`, "foo bar!"},
		{"string null", `string(null)`, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestIntBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"int 123", `int(123)`, int64(123)},
		{"int 45.67", `int(45.67)`, int64(45)},
		{"int '789'", `int("789")`, int64(789)},
		{"int 0", `int(0)`, int64(0)},
		{"int -42", `int(-42)`, int64(-42)},
		{"int -3.99", `int(-3.99)`, int64(-3)},
		{"int 0", `int("0")`, int64(0)},
		{"int 456", `int("   456   ")`, int64(456)},
		{"int true", `int(true)`, int64(1)},
		{"int false", `int(false)`, int64(0)},
		{"int null", `int(null)`, int64(0)},
		{"int empty string", `int("")`, &objects.Error{Message: "error in `int`: failed to convert `` to INTEGER"}},
		{"int non-number", `int("not a number")`, &objects.Error{Message: "error in `int`: failed to convert `not a number` to INTEGER"}},
		{"int multiple parameters", `int(1, 2)`, &objects.Error{Message: "wrong number of arguments to `int`: got 2, want 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestFloatBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"float 123", `float(123)`, float64(123)},
		{"float 45.67", `float(45.67)`, float64(45.67)},
		{"float '789.01'", `float("789.01")`, float64(789.01)},
		{"float 0", `float(0)`, float64(0)},
		{"float -42", `float(-42)`, float64(-42)},
		{"float -3.99", `float(-3.99)`, float64(-3.99)},
		{"float '0'", `float("0")`, float64(0)},
		{"float '456.78' with padding", `float("   456.78   ")`, float64(456.78)},
		{"float true", `float(true)`, float64(1)},
		{"float false", `float(false)`, float64(0)},
		{"float null", `float(null)`, float64(0)},
		{"float empty string", `float("")`, &objects.Error{Message: "error in `float`: failed to convert `` to FLOAT"}},
		{"float non-number", `float("not a number")`, &objects.Error{Message: "error in `float`: failed to convert `not a number` to FLOAT"}},
		{"float multiple parameters", `float(1, 2)`, &objects.Error{Message: "wrong number of arguments to `float`: got 2, want 1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestTypeBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"type integer", `type(123)`, "INTEGER"},
		{"type float", `type(45.67)`, "FLOAT"},
		{"type boolean true", `type(true)`, "BOOLEAN"},
		{"type boolean false", `type(false)`, "BOOLEAN"},
		{"type string", `type("hello")`, "STRING"},
		{"type null", `type(null)`, "NULL"},
		{"type array", `type([1, 2, 3])`, "ARRAY"},
		{"type hash", `type({"key": "value"})`, "HASH"},
		{"type function", `type(func(x) { return x; })`, "FUNCTION"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestIsNaNBuiltinFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"isNaN zero divided by zero", `isNaN(0/0)`, true},
		{"isNaN one divided by zero", `isNaN(1/0)`, false},
		{"isNaN negative one divided by zero", `isNaN(-1/0)`, false},
		{"isNaN integer", `isNaN(5)`, false},
		{"isNaN float", `isNaN(3.14)`, false},
		{"isNaN string not a number", `isNaN("not a number")`, false},
		{"isNaN boolean true", `isNaN(true)`, false},
		{"isNaN boolean false", `isNaN(false)`, false},
		{"isNaN null", `isNaN(null)`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}
