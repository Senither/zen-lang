package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!!!true", false},
		{"!!!false", true},
		{"!5", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestVarReassignmentStatements(t *testing.T) {
	input := []struct {
		input    string
		expected interface{}
	}{
		{`var mut x = 5; x = 10; x;`, 10},
		{`var mut x = "hello"; x = "world"; x;`, "world"},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testStringObject(t, evaluated, expected)
		}
	}
}

func TestVarReassignmentFailure(t *testing.T) {
	input := []struct {
		input    string
		expected string
	}{
		{`var x = 5; x = 10; x;`, "cannot modify immutable variable: x\n    at <unknown>:1:14"},
		{`var x = 3.14; x = 4.14; x;`, "cannot modify immutable variable: x\n    at <unknown>:1:17"},
		{`var x = "hello"; x = "world"; x;`, "cannot modify immutable variable: x\n    at <unknown>:1:20"},
		{`var name = "Senither"; name = "test"; name;`, "cannot modify immutable variable: name\n    at <unknown>:1:29"},
		{`var arr = [1,2,3]; arr = ['another', 'array']; arr;`, "cannot modify immutable variable: arr\n    at <unknown>:1:24"},
		{`var obj = {"key": "value"}; obj = {"key": "new value"}; obj;`, "cannot modify immutable variable: obj\n    at <unknown>:1:33"},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		if !objects.IsError(evaluated) {
			t.Fatalf("expected error, got: %q", evaluated)
		}

		if evaluated.Inspect() != tt.expected {
			t.Fatalf("expected: %q, got: %q", tt.expected, evaluated.Inspect())
		}
	}
}

func TestVarIncrementingStatements(t *testing.T) {
	input := []struct {
		input    string
		expected interface{}
	}{
		{`var x = 5; x++; x;`, 6},
		{`var mut x = 5; x++; x;`, 6},
		{`var x = 3.14; x++; x;`, float64(4.14)},
		{`var mut x = 3.14; x++; x;`, float64(4.14)},
		{`var x = "hello"; x++; x;`, "unknown operator: ++STRING\n    at <unknown>:1:18"},
		{`var mut x = "hello"; x++; x;`, "unknown operator: ++STRING\n    at <unknown>:1:22"},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case float64:
			testFloatObject(t, evaluated, expected)
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}
}

func TestVarDecrementingStatements(t *testing.T) {
	input := []struct {
		input    string
		expected interface{}
	}{
		{`var x = 5; x--; x;`, 4},
		{`var mut x = 5; x--; x;`, 4},
		{`var x = 3.14; x--; x;`, float64(2.14)},
		{`var mut x = 3.14; x--; x;`, float64(2.14)},
		{`var x = "hello"; x--; x;`, "unknown operator: --STRING\n    at <unknown>:1:18"},
		{`var mut x = "hello"; x--; x;`, "unknown operator: --STRING\n    at <unknown>:1:22"},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case float64:
			testFloatObject(t, evaluated, expected)
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}
}

func TestArrayReassignmentStatements(t *testing.T) {
	input := []struct {
		input    string
		expected []any
	}{
		{`var arr = [1, 2, 3]; arr[0] = 10; arr;`, []any{10, 2, 3}},
		{`var mut arr = [1, 2, 3]; arr[0] = 10; arr;`, []any{10, 2, 3}},
		{`var arr = [1, 2, 3]; arr[1] = 10; arr;`, []any{1, 10, 3}},
		{`var mut arr = [1, 2, 3]; arr[1] = 10; arr;`, []any{1, 10, 3}},
		{`var arr = ["foo", "bar"]; arr[1] = "baz"; arr;`, []any{"foo", "baz"}},
		{`var mut arr = ["foo", "bar"]; arr[1] = "baz"; arr;`, []any{"foo", "baz"}},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		testArrayObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestHashReassignmentStatements(t *testing.T) {
	input := []struct {
		input    string
		expected map[objects.HashKey]any
	}{
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash["foo"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash["foo"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash["baz"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash["baz"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash.foo = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash.foo = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash.baz = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash.baz = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		testHashObject(t, evaluated, tt.expected, tt.input)
	}
}
