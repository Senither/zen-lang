package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"string literal 'hello'", "\"hello\"", "hello"},
		{"string literal 'world'", "\"world\"", "world"},
		{"string literal 'Hello, world!'", "\"Hello, world!\"", "Hello, world!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	objects.AssertExpectedObject(t, "Hello World!", testEval(`"Hello" + " " + "World!"`))
}

func TestEvalStringCastConcatenation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"string and integer concatenation", `"The answer is: " + 42`, "The answer is: 42"},
		{"string and float concatenation", `"Pi is approximately " + 3.14`, "Pi is approximately 3.14"},
		{"string and true boolean concatenation", `"Value: " + true`, "Value: true"},
		{"string and false boolean concatenation", `"Value: " + false`, "Value: false"},
		{"string and expression concatenation", `"Number: " + (10 + 5)`, "Number: 15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
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
		{"2 * 3 ^ 4", 162},
		{"2 * 3 ^ 4 % 5", 2},
	}

	for _, tt := range tests {
		t.Run("integer expression: "+tt.input, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
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
		t.Run("float expression: "+tt.input, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
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
		{"true && true", true},
		{"true && false", false},
		{"false && true", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"false || false", false},
		{"true && (false || true)", true},
		{"(1 < 2) && (2 < 3)", true},
		{"(1 < 2) && (2 > 3)", false},
		{"(1 > 2) || (2 < 3)", true},
		{"(1 > 2) || (2 > 3)", false},
	}

	for _, tt := range tests {
		t.Run("boolean expression: "+tt.input, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestArrayLiterals(t *testing.T) {
	objects.AssertExpectedObject(t,
		[]int{1, 6, 9},
		testEval("[1, 2 * 3, 4 + 5]"),
	)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			"first element",
			"[1, 2, 3][0]",
			1,
		},
		{
			"second element",
			"[1, 2, 3][1]",
			2,
		},
		{
			"third element",
			"[1, 2, 3][2]",
			3,
		},
		{
			"negative index -3",
			"[1, 2, 3][-3]",
			1,
		},
		{
			"negative index -2",
			"[1, 2, 3][-2]",
			2,
		},
		{
			"negative index -1",
			"[1, 2, 3][-1]",
			3,
		},
		{
			"index with expression",
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"index stored in variable",
			"var i = 0; [1, 2][i];",
			1,
		},
		{
			"index stored in variable with expression",
			"var i = 0; [1, 2][i + 1];",
			2,
		},
		{
			"accessing array from variable",
			"var myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"sum of array elements",
			"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"chained array access",
			"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]",
			2,
		},
		{
			"out of bounds index",
			"[1, 2, 3][3]",
			&objects.Error{Message: "array index out of bounds: 3"},
		},
		{
			"negative out of bounds index",
			"[1, 2, 3][-4]",
			&objects.Error{Message: "array index out of bounds: -4"},
		},
	}

	for _, tt := range tests {
		t.Run("array index: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
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

		objects.AssertInteger(expectedValue, pair.Value)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			"hash string key found",
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			"hash string key not found",
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			"variable as hash key found",
			`var key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			"variable as hash key not found",
			`{}["foo"]`,
			nil,
		},
		{
			"hash integer key found",
			`{5: 5}[5]`,
			5,
		},
		{
			"hash integer key not found",
			`{true: 5}[true]`,
			5,
		},
		{
			"hash boolean key found",
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		t.Run("hash index: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestChainedHashIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			"string value",
			`var x = {"foo": 5}; x.foo`,
			5,
		},
		{
			"float value",
			`var x = {"foo": 3.14}; x.foo`,
			3.14,
		},
		{
			"boolean true value",
			`var x = {"foo": true}; x.foo`,
			true,
		},
		{
			"boolean false value",
			`var x = {"foo": false}; x.foo`,
			false,
		},
		{
			"array value",
			`var x = {"foo": [1, 2, 3]}; x.foo[1]`,
			2,
		},
		{
			"array value with expression index",
			`var x = {"foo": [1, 2, 3]}; x.foo[1 + 1]`,
			3,
		},
		{
			"hash value",
			`var x = {"foo": func (n) { return n + 1; }}; x.foo(5)`,
			6,
		},
		{
			"key not found",
			`var x = {"foo": 5}; x.bar`,
			&objects.Error{Message: "invalid chain expression for HASH, key not found: bar"},
		},
	}

	for _, tt := range tests {
		t.Run("chained hash index: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestChainedHashAssignmentExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"simple assignment",
			"var x = {'foo': 5}; x.foo = 10; x.foo;",
			10,
		},
		{
			"bracket notation assignment",
			"var x = {'foo': 5}; x['foo'] = 10; x.foo;",
			10,
		},
		{
			"nested assignment",
			"var x = {'foo': {'bar': 5}}; x.foo.bar = 10; x.foo.bar;",
			10,
		},
		{
			"nested bracket notation assignment",
			"var x = {'foo': {'bar': 5}}; x['foo']['bar'] = 10; x.foo.bar;",
			10,
		},
		{
			"assignment from another hash",
			"var x = {'foo': 5}; var y = {'bar': 10}; x.foo = y.bar; x.foo;",
			10,
		},
		{
			"new key assignment from another hash",
			"var x = {'foo': 5}; var y = {'bar': 10}; x.newKey = y.bar; x.newKey;",
			10,
		},
	}

	for _, tt := range tests {
		t.Run("chained hash assignment: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestChainedArrayIndexAssignmentExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
	}{
		{
			"simple assignment",
			"var x = {'arr': [1, 2, 3, 4]}; x.arr[2] = x.arr[2] + 40; x.arr;",
			[]any{1, 2, 43, 4},
		},
		{
			"assignment to index 1",
			"var x = {'arr': [1, 2, 3]}; x.arr[1] = 42; x.arr;",
			[]any{1, 42, 3},
		},
		{
			"assignment to last index",
			"var x = {'foo': {'bar': [5, 6]}}; x.foo.bar[0] = 99; x.foo.bar;",
			[]any{99, 6},
		},
	}

	for _, tt := range tests {
		t.Run("chained array index assignment: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestReassigningArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []any
	}{
		{
			"assignment to index 0",
			"var x = [1, 2, 3]; x[0] = 99; x;",
			[]any{99, 2, 3},
		},
		{
			"assignment to index 1",
			"var x = [1, 2, 3]; x[1] = 99; x;",
			[]any{1, 99, 3},
		},
		{
			"assignment to index 2",
			"var x = [1, 2, 3]; x[2] = 99; x;",
			[]any{1, 2, 99},
		},
		{
			"assignment to index 0 with string",
			"var x = [1, 2, 3]; x[0] = 'This is a test'; x;",
			[]any{"This is a test", 2, 3},
		},
		{
			"assignment to index 1 with float",
			"var x = [1, 2, 3]; x[1] = 3.14; x;",
			[]any{1, float64(3.14), 3},
		},
		{
			"assignment to index 2 with boolean",
			"var x = [1, 2, 3]; x[2] = true; x;",
			[]any{1, 2, true},
		},
	}

	for _, tt := range tests {
		t.Run("reassign array index: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestReassigningArrayIndexExpressionsErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *objects.Error
	}{
		{
			"index out of bounds (too high)",
			"var x = [1, 2, 3]; x[3] = 99;",
			&objects.Error{Message: "array index out of bounds: 3"},
		},
		{
			"index out of bounds (too low)",
			"var x = [1, 2, 3]; x[-4] = 99;",
			&objects.Error{Message: "array index out of bounds: -4"},
		},
	}

	for _, tt := range tests {
		t.Run("reassign array index: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"get value from truthy block", "if (true) { 10 }", 10},
		{"get nil from falsy block", "if (false) { 10 }", nil},
		{"get value from truthy block with integer condition", "if (1) { 10 }", 10},
		{"get value from truthy block with less than condition", "if (1 < 2) { 10 }", 10},
		{"get nil from falsy block with greater than condition", "if (1 > 2) { 10 }", nil},
		{"get value from else block with greater than condition", "if (1 > 2) { 10 } else { 20 }", 20},
		{"get value from truthy block with less than condition and else block", "if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestIfElseIfElseExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{
			"get value from truthy block",
			"if (true) { 10 } else if (false) { 20 } else { 30 }",
			10,
		},
		{
			"get value from else if block",
			"if (false) { 10 } else if (true) { 20 } else { 30 }",
			20,
		},
		{
			"get value from else block",
			"if (false) { 10 } else if (false) { 20 } else { 30 }",
			30,
		},
		{
			"get nil when no conditions are truthy and no else block",
			"if (false) { 10 } else if (false) { 20 }",
			nil,
		},
		{
			"get value from else if block when multiple else if blocks",
			"if (false) { 10 } else if (false) { 20 } else if (true) { 30 }",
			30,
		},
		{
			"get value from else block when multiple else if blocks",
			"if (false) { 10 } else if (false) { 20 } else if (false) { 30 } else { 40 }",
			40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}
func TestWhileExpressions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected any
	}{
		{"loop increments by 1", "var mut i = 0; while (i < 5) { i++; } i;", 5},
		{"loop never runs", "var mut i = 0; while (i > 5) { i++; } i;", 0},
		{"loop increments by 2", "var mut i = 0; while (i < 5) { i = i + 2; } i;", 6},
		{"loop with no return value", "var mut i = 0; while (i < 5) { i = i + 2; }", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}
