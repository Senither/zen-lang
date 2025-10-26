package evaluator

import (
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
		{"2 ^ 0", 1},
		{"2 ^ 3", 8},
		{"2 + 5 ^ 2", 27},
		{"(2 + 5) ^ 2 * 2", 98},
		{"20 % 2", 0},
		{"20 % 3", 2},
		{"100 % 17", 15},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
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
		{"5.5 ^ 0", 1.0},
		{"5.5 ^ 2", 30.25},
		{"2.5 + 5.5 ^ 2", 32.75},
		{"(2.5 + 5.5) ^ 2 * 2", 128.0},
		{"10.75 % 3", 1.75},
		{"12.34 % 5", 2.34},
		{"3.14 % 2", 1.14},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
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

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 6)
	testIntegerObject(t, result.Elements[2], 9)
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
			"array index out of bounds: 3\n    at <unknown>:1:10",
		},
		{
			"[1, 2, 3][-4]",
			"array index out of bounds: -4\n    at <unknown>:1:10",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testErrorObject(t, evaluated, expected)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
		var two = "two";
		{
			"one": 1,
			two: 2,
			"thr" + "ee": 3,
			4: 4,
			true: 5,
			false: 6
		};
	`

	evaluated := testEval(input)
	result, ok := evaluated.(*objects.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got %T (%+v)", evaluated, evaluated)
	}

	expected := map[objects.HashKey]int64{
		(&objects.String{Value: "one"}).HashKey():   1,
		(&objects.String{Value: "two"}).HashKey():   2,
		(&objects.String{Value: "three"}).HashKey(): 3,
		(&objects.Integer{Value: 4}).HashKey():      4,
		objects.TRUE.HashKey():                      5,
		objects.FALSE.HashKey():                     6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong number of pairs. got %d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("Hash is missing key. got %v", result.Pairs)
			continue
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestChainedHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`var x = {"foo": 5}; x.foo`,
			5,
		},
		{
			`var x = {"foo": 3.14}; x.foo`,
			3.14,
		},
		{
			`var x = {"foo": true}; x.foo`,
			true,
		},
		{
			`var x = {"foo": false}; x.foo`,
			false,
		},
		{
			`var x = {"foo": [1, 2, 3]}; x.foo[1]`,
			2,
		},
		{
			`var x = {"foo": [1, 2, 3]}; x.foo[1 + 1]`,
			3,
		},
		{
			`var x = {"foo": func (n) { return n + 1; }}; x.foo(5)`,
			6,
		},
		{
			`var x = {"foo": 5}; x.bar`,
			"invalid chain expression for HASH, key not found: bar\n    at <unknown>:1:22",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case float64:
			testFloatObject(t, evaluated, expected)
		case bool:
			testBooleanObject(t, evaluated, expected, tt.input)
		case string:
			testErrorObject(t, evaluated, expected)
		default:
			t.Errorf("Unexpected type of object. got %T (%+v)", evaluated, evaluated)
		}
	}
}

func TestReassigningArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{
			"var x = [1, 2, 3]; x[0] = 99; x;",
			[]any{99, 2, 3},
		},
		{
			"var x = [1, 2, 3]; x[1] = 99; x;",
			[]any{1, 99, 3},
		},
		{
			"var x = [1, 2, 3]; x[2] = 99; x;",
			[]any{1, 2, 99},
		},
		{
			"var x = [1, 2, 3]; x[0] = 'This is a test'; x;",
			[]any{"This is a test", 2, 3},
		},
		{
			"var x = [1, 2, 3]; x[1] = 3.14; x;",
			[]any{1, float64(3.14), 3},
		},
		{
			"var x = [1, 2, 3]; x[2] = true; x;",
			[]any{1, 2, true},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		arr, ok := evaluated.(*objects.Array)
		if !ok {
			t.Fatalf("object is not Array. got %T (%+v)", evaluated, evaluated)
		}

		for i, expected := range tt.expected {
			switch expected := expected.(type) {
			case int:
				testIntegerObject(t, arr.Elements[i], int64(expected))
			case float64:
				testFloatObject(t, arr.Elements[i], expected)
			case bool:
				testBooleanObject(t, arr.Elements[i], expected, tt.input)
			case string:
				testStringObject(t, arr.Elements[i], expected)
			}
		}
	}
}

func TestReassigningArrayIndexExpressionsErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"var x = [1, 2, 3]; x[3] = 99;",
			"array index out of bounds: 3\n    at <unknown>:1:25",
		},
		{
			"var x = [1, 2, 3]; x[-4] = 99;",
			"array index out of bounds: -4\n    at <unknown>:1:26",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testErrorObject(t, evaluated, tt.expected)
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
			testIntegerObject(t, evaluated, int64(integer))
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
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"var mut i = 0; while (i < 5) { i++; } i;", 5},
		{"var mut i = 0; while (i > 5) { i++; } i;", 0},
		{"var mut i = 0; while (i < 5) { i = i + 2; } i;", 6},
		{"var mut i = 0; while (i < 5) { i = i + 2; }", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
