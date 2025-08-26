package evaluator

import "testing"

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
