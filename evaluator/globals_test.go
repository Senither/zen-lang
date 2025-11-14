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
		{"arrays.push([1], null)", []any{1, nil}},
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
		{"arrays.shift([null])", nil},
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
		{"arrays.pop([null])", nil},
		{"arrays.pop([1])", 1},
		{"arrays.pop([1, 2])", 2},
		{"arrays.pop([1, null])", nil},
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

func TestArraysConcatGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"arrays.concat([], [])", []int64{}},
		{"arrays.concat([1], [2])", []int64{1, 2}},
		{"arrays.concat([1, 2], [3, 4])", []int64{1, 2, 3, 4}},
		{"arrays.concat([1, 2], [3, 4], [5, 6])", []int64{1, 2, 3, 4, 5, 6}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case []int64:
			converted := make([]any, len(expected))
			for i, v := range expected {
				converted[i] = v
			}
			testArrayObject(t, evaluated, converted, tt.input)
		}
	}
}

func TestStringsContainsGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`strings.contains("hello", "he")`, true},
		{`strings.contains("hello", "lo")`, true},
		{`strings.contains("hello", "world")`, false},
		{`var x = "hello"; strings.contains(x, "he");`, true},
		{`var x = "hello"; strings.contains(x, "lo");`, true},
		{`var x = "hello"; strings.contains(x, "world");`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, expected, tt.input)
		}
	}
}

func TestStringsSplitGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{`strings.split("a,b,c", ",")`, []any{"a", "b", "c"}},
		{`strings.split("a b c", " ")`, []any{"a", "b", "c"}},
		{`strings.split("a b c", "-")`, []any{"a b c"}},
		{`var x = "a,b,c"; strings.split(x, ",");`, []any{"a", "b", "c"}},
		{`var x = "a b c"; strings.split(x, " ");`, []any{"a", "b", "c"}},
		{`var x = "a b c"; strings.split(x, "-");`, []any{"a b c"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayObject(t, evaluated, tt.expected, tt.input)
	}
}

func TestStringsJoinGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`strings.join([1, 2, 3, 4, 5], "")`, "12345"},
		{`strings.join([1, 2, 3, 4, 5], ", ")`, "1, 2, 3, 4, 5"},
		{`strings.join([], ", ")`, ""},
		{`strings.join([1, 2.22, true, false, [5,6,7], {"key": "value"}], ", ")`, "1, 2.22, true, false, [5, 6, 7], {key: value}"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringsFormatGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`strings.format("%s", "hello")`, "hello"},
		{`strings.format("%d", 5)`, "5"},
		{`strings.format("%f", 3.14)`, "3.140000"},
		{`strings.format("%v", 3.14)`, "3.14"},
		{`strings.format("%t", true)`, "true"},
		{`strings.format("%t", false)`, "false"},
		{`strings.format("%+v", null)`, "<nil>"},
		{`strings.format("%v", [1, 2, 3])`, "[1, 2, 3]"},
		{`strings.format("%v", {"key": "value"})`, "{key: value}"},
		{`strings.format("%#T", "test")`, "string"},
		{`strings.format("%s %d %f %v", "hello", 5, 3.14, true)`, "hello 5 3.140000 true"},
		{`strings.format("%s %s", "test")`, "test %!s(MISSING)"},
		{`strings.format("%s", "test", 5)`, "test%!(EXTRA int64=5)"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestMathMinGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"math.min(1, 2)", 1},
		{"math.min(2, 1)", 1},
		{"math.min(-1, 1)", -1},
		{"math.min(0, 0)", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMathMaxGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"math.max(1, 2)", 2},
		{"math.max(2, 1)", 2},
		{"math.max(-1, 1)", 1},
		{"math.max(0, 0)", 0},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestMathCeilGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"math.ceil(1.1)", 2},
		{"math.ceil(1.9)", 2},
		{"math.ceil(-1.1)", -1},
		{"math.ceil(-1.9)", -1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestMathFloorGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"math.floor(1.1)", 1},
		{"math.floor(1.9)", 1},
		{"math.floor(-1.1)", -2},
		{"math.floor(-1.9)", -2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestMathRoundGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"math.round(1.1)", 1},
		{"math.round(1.5)", 2},
		{"math.round(1.9)", 2},
		{"math.round(-1.1)", -1},
		{"math.round(-1.5)", -2},
		{"math.round(-1.9)", -2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestMathLogGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"math.log(1)", 0},
		{"math.log(10)", 1},
		{"math.log(100)", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestMathSqrtGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"math.sqrt(0)", 0},
		{"math.sqrt(1)", 1},
		{"math.sqrt(4)", 2},
		{"math.sqrt(9)", 3},
		{"math.sqrt(1764)", 42},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}
