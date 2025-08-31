package evaluator

import "testing"

func TestArraysPushGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"arrays.push([], 1)", []any{1}},
		{"arrays.push([1], 2)", []any{1, 2}},
		{"arrays.push([1, 2], 3)", []any{1, 2, 3}},
		{"var x = []; arrays.push(x, 1);", []any{1}},
		{"var x = [1]; arrays.push(x, 2);", []any{1, 2}},
		{"var x = [1, 2]; arrays.push(x, 3);", []any{1, 2, 3}},
		{"var x = []; arrays.push(x, 1); x;", []any{1}},
		{"var x = [1]; arrays.push(x, 2); x;", []any{1, 2}},
		{"var x = [1, 2]; arrays.push(x, 3); x;", []any{1, 2, 3}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestArraysShiftGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"arrays.shift([])", nil},
		{"arrays.shift([1])", 1},
		{"arrays.shift([1, 2])", 1},
		{"var x = []; arrays.shift(x);", nil},
		{"var x = [1]; arrays.shift(x);", 1},
		{"var x = [1, 2]; arrays.shift(x);", 1},
		{"var x = []; arrays.shift(x); x;", []int64{}},
		{"var x = [1]; arrays.shift(x); x;", []int64{}},
		{"var x = [1, 2]; arrays.shift(x); x;", []int64{2}},
		{"var x = [1, 2, 3]; arrays.shift(x); x;", []int64{2, 3}},
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

func TestArraysPopGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"arrays.pop([])", nil},
		{"arrays.pop([1])", 1},
		{"arrays.pop([1, 2])", 2},
		{"var x = []; arrays.pop(x);", nil},
		{"var x = [1]; arrays.pop(x);", 1},
		{"var x = [1, 2]; arrays.pop(x);", 2},
		{"var x = []; arrays.pop(x); x;", []int64{}},
		{"var x = [1]; arrays.pop(x); x;", []int64{}},
		{"var x = [1, 2]; arrays.pop(x); x;", []int64{1}},
		{"var x = [1, 2, 3]; arrays.pop(x); x;", []int64{1, 2}},
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
