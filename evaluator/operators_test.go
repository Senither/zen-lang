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
		{`var x = 5; x = 10; x;`, "ERROR: Cannot modify immutable variable 'x'"},
		{`var x = "hello"; x = "world"; x;`, "ERROR: Cannot modify immutable variable 'x'"},
		{`var name = "Senither"; name = "test"; name;`, "ERROR: Cannot modify immutable variable 'name'"},
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		if !isError(evaluated) {
			t.Fatalf("expected error, got: %q", evaluated)
		}

		if evaluated.Inspect() != tt.expected {
			t.Fatalf("expected: %q, got: %q", tt.expected, evaluated.Inspect())
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
	}

	for _, tt := range input {
		evaluated := testEval(tt.input)

		testHashObject(t, evaluated, tt.expected, tt.input)
	}
}
