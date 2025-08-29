package evaluator

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello\"", "hello"},
		{"\"world\"", "world"},
		{"\"Hello, world!\"", "Hello, world!"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*objects.String)
	if !ok {
		t.Fatalf("object is not String. got %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got %q", str.Value)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"5 + 5 - 10", 0},
		{"5 - 10", -5},
		{"5 + 5 + 5 - 10", 5},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"(5 + 10) * 2", 30},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, new(big.Int).SetInt64(int64(tt.expected)))
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.5", 5.5},
		{"10.5", 10.5},
		{"-5.5", -5.5},
		{"-10.0", -10.0},
		{"1.2 + 3.4 + 5.6 + 7.8 - 10.0", 8.0},
		{"5.5 + 5.45 - 10", 0.95},
		{"5.123 - 10.321", -5.198},
		{"5.5 + 5.5 + 5.5 - 10", 6.5},
		{"2 / 5", 0.4},
		{"5 / 2", 2.5},
		{"5.0 / 2.0", 2.5},
		{"5.5 / 2.5", 2.2},
		{"5.5 / -2.5", -2.2},
		{"-5.5 / 2.5", -2.2},
		{"14 / (10 - 3.75)", 2.24},
		{"5.5 * 2.5", 13.75},
		{"10.5 * 9.15", 96.075},
		{"2.5 * (2.5 + 25)", 68.75},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, objects.NewFloat(tt.expected))
	}
}

func TestEvalFloatAtHighPrecision(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"12.3456789", 12.3456789},
		{"12.3456789 + 0.0000001", 12.3456790},
		{"12.3456789 - 0.0000001", 12.3456788},
		{"12.3456789 * 2", 24.6913578},
		{"12.3456789 / 2", 6.17283945},
		{"0.123456789123456789123456789", 0.123456789123456789123456789},
		{"0.123456789123456789123456789 + 0.000000000000000000000000001", 0.123456789123456789123456790},
		{"0.123456789123456789123456789 + 0.000001000002000003000004000", 0.123457789125456789123456789},
		{"0.123456789123456789123456789 - 0.000000000000000000000000001", 0.123456789123456789123456788},
		{"0.123456789123456789123456789 - 0.000001000002000003000004000", 0.123455789121456789123456789},
		{"0.123456789123456789123456789 * 1.234567891234567891234567890", 0.152415787800000000000000000},
		{"0.123456789123456789123456789 / 2", 0.061728394561728394561728394},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		fmt.Printf("Input: %q, Output: %q\n", tt.input, evaluated)
		testFloatObject(t, evaluated, objects.NewFloat(tt.expected))
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"1 >= 0", true},
		{"1 >= 1", true},
		{"1 >= 2", false},
		{"1 <= 0", false},
		{"1 <= 1", true},
		{"1 <= 2", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 3, 4 + 5]"

	evaluated := testEval(input)
	result, ok := evaluated.(*objects.Array)
	if !ok {
		t.Fatalf("object is not Array. got %T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Errorf("array has wrong number of elements. got %d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], new(big.Int).SetInt64(1))
	testIntegerObject(t, result.Elements[1], new(big.Int).SetInt64(6))
	testIntegerObject(t, result.Elements[2], new(big.Int).SetInt64(9))
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"[1, 2, 3][-3]",
			1,
		},
		{
			"[1, 2, 3][-2]",
			2,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"var i = 0; [1, 2][i];",
			1,
		},
		{
			"var i = 0; [1, 2][i + 1];",
			2,
		},
		{
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			"array index out of bounds: 3",
		},
		{
			"[1, 2, 3][-4]",
			"array index out of bounds: -4",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, new(big.Int).SetInt64(int64(expected)))
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}

}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, new(big.Int).SetInt64(int64(integer)))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestIfElseIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 } else if (false) { 20 } else { 30 }", 10},
		{"if (false) { 10 } else if (true) { 20 } else { 30 }", 20},
		{"if (false) { 10 } else if (false) { 20 } else { 30 }", 30},
		{"if (false) { 10 } else if (false) { 20 }", nil},
		{"if (false) { 10 } else if (false) { 20 } else if (true) { 30 }", 30},
		{"if (false) { 10 } else if (false) { 20 } else if (false) { 30 } else { 40 }", 40},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, new(big.Int).SetInt64(int64(integer)))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
