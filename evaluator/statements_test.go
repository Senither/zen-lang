package evaluator

import (
	"math/big"
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, new(big.Int).SetInt64(int64(tt.expected)))
	}
}

func TestNestedReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, new(big.Int).SetInt64(int64(tt.expected)))
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var x = 5; x;", 5},
		{"var x = 5 * 5; x;", 25},
		{"var x = 5; var y = 10; x + y;", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, objects.NewInteger(int64(tt.expected)).Value)
	}
}
