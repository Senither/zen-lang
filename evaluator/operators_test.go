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
		{`var x = 5; x = 10; x;`, 10},
		{`var x = "hello"; x = "world"; x;`, "world"},
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
