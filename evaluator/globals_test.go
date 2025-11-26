package evaluator

import (
	"testing"
	"time"

	"github.com/senither/zen-lang/objects/timer"
)

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

func TestArraysFilterGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"arrays.filter([1, 2, 3, 4], func(x) { x % 2 == 0 })", []any{2, 4}},
		{"arrays.filter([1, 2, 3, 4, 5], func(x) { x > 3 })", []any{4, 5}},
		{"arrays.filter([], func(x) { x > 0 })", []any{}},
		{"var x = [1, 2, 3, 4]; arrays.filter(x, func(y) { y < 3 });", []any{1, 2}},
		{"var x = [10, 15, 20, 25]; arrays.filter(x, func(x) { x >= 20 });", []any{20, 25}},
		{"arrays.filter([null, 1, null, 2], func(a) { a != null })", []any{1, 2}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayObject(t, evaluated, tt.expected, tt.input)
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

func TestArraysFirstGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"arrays.first([100, 200, 300], func (x) { x >= 100 })", 100},
		{"arrays.first([100, 200, 300], func (x) { x > 100 })", 200},
		{"arrays.first([100, 200, 300], func (x, i) { i == 2 })", 300},
		{"arrays.first([100, 200, 300], func (x) { x > 500 })", nil},
		{"arrays.first(5, func (a) { })", "argument 1 to `first` has invalid type: got INTEGER, want ARRAY\n    at <unknown>:1:13"},
		{"arrays.first([100, 200, 300], func () { })", "error in `first`: function passed to `first` must take at least one argument\n    at <unknown>:1:13"},
		{"arrays.first([100, 200, 300], func (a, b, c) { })", "error in `first`: function passed to `first` must take at most two arguments\n    at <unknown>:1:13"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			testErrorObject(t, evaluated, expected)
		case nil:
			testNullObject(t, evaluated)
		}
	}
}

func TestArraysSortGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{"arrays.sort([5, 3, 1, 4, 2])", []any{1, 2, 3, 4, 5}},
		{"arrays.sort([-1, 0, 3, -5, 2, 1])", []any{-5, -1, 0, 1, 2, 3}},
		{"arrays.sort([])", []any{}},
		{"var x = [3, 2, 1]; arrays.sort(x);", []any{1, 2, 3}},
		{"var x = [10, 5, -5, -10, 0]; arrays.sort(x);", []any{-10, -5, 0, 5, 10}},
		{"arrays.sort([1.5, 2.2, 0.3, -1.1])", []any{-1.1, 0.3, 1.5, 2.2}},
		{"arrays.sort(['banana', 'apple', 'cherry'])", []any{"apple", "banana", "cherry"}},
		{"var x = ['zen', 'lang', 'is', 'awesome']; arrays.sort(x);", []any{"awesome", "is", "lang", "zen"}},
		{"arrays.sort([true, false, true, false])", []any{false, false, true, true}},
		{"arrays.sort([5, 3, 1, 4, 2], func (a, b) { a < b })", []any{1, 2, 3, 4, 5}},
		{"arrays.sort([5, 3, 1, 4, 2], func (a, b) { a > b })", []any{5, 4, 3, 2, 1}},
		{"arrays.sort(['bb', 'a', 'ccc'], func (a, b) { len(a) < len(b) })", []any{"a", "bb", "ccc"}},
		{"arrays.sort(['bb', 'a', 'ccc'], func (a, b) { len(a) > len(b) })", []any{"ccc", "bb", "a"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testArrayObject(t, evaluated, tt.expected, tt.input)
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

func TestStringsStartsWithGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`strings.startsWith("hello", "he")`, true},
		{`strings.startsWith("hello", "hello")`, true},
		{`strings.startsWith("hello", "")`, true},
		{`strings.startsWith("hello", "world")`, false},
		{`strings.startsWith("hello", ["Hello", "World"])`, false},
		{`strings.startsWith("hello", ["hello", "world"])`, true},
		{`var x = "foobar"; strings.startsWith(x, "foo");`, true},
		{`var x = "foobar"; strings.startsWith(x, "bar");`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, expected, tt.input)
		}
	}
}

func TestStringsEndsWithGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`strings.endsWith("hello", "lo")`, true},
		{`strings.endsWith("hello", "hello")`, true},
		{`strings.endsWith("hello", "")`, true},
		{`strings.endsWith("hello", "world")`, false},
		{`strings.endsWith("hello", ["Hello", "World"])`, false},
		{`strings.endsWith("hello", ["lo", "world"])`, true},
		{`var x = "foobar"; strings.endsWith(x, "bar");`, true},
		{`var x = "foobar"; strings.endsWith(x, "foo");`, false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case bool:
			testBooleanObject(t, evaluated, expected, tt.input)
		}
	}
}

func TestStringsToUpperGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`strings.toUpper("hello")`, "HELLO"},
		{`strings.toUpper("Hello World!")`, "HELLO WORLD!"},
		{`var x = "Test"; strings.toUpper(x);`, "TEST"},
		{`strings.toUpper("")`, ""},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringsToLowerGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`strings.toLower("HELLO")`, "hello"},
		{`strings.toLower("Hello World!")`, "hello world!"},
		{`var x = "TeSt"; strings.toLower(x);`, "test"},
		{`strings.toLower("")`, ""},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestStringsTrimGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`strings.trim("  hello  ")`, "hello"},
		{"strings.trim(\"\n\thello world \t\n\")", "hello world"},
		{`strings.trim("hello")`, "hello"},
		{`var x = "  spaced  "; strings.trim(x);`, "spaced"},
		{`strings.trim("")`, ""},
		{`strings.trim("  hello  ", " ")`, "hello"},
		{`strings.trim("xxhelloxx", "x")`, "hello"},
		{`strings.trim("!!wow!!", "!")`, "wow"},
		{`strings.trim("!!!!!wow!!!", "!")`, "wow"},
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

func TestTimeNowUnfrozenGlobalFunction(t *testing.T) {
	evaluated := testEval("time.now()")
	testIntegerObject(t, evaluated, time.Now().UnixMilli())
}

func TestTimeNowFrozenGlobalFunction(t *testing.T) {
	timer.Freeze(1767606155000)

	evaluated := testEval("time.now()")
	testIntegerObject(t, evaluated, 1767606155000)

	timer.Unfreeze()
}

func TestTimeSleepUnfrozenGlobalFunction(t *testing.T) {
	start := time.Now().UnixMilli()
	testEval("time.sleep(100)")
	end := time.Now().UnixMilli()

	if end-start < 100 {
		t.Fatalf("time.sleep did not sleep long enough: expected at least 100ms, got %dms", end-start)
	}
}

func TestTimeSleepFrozenGlobalFunction(t *testing.T) {
	timer.Freeze(1767606155000)

	testEval("time.sleep(10_000)")

	now := timer.Now()
	timer.Unfreeze()

	if now != 1767606165000 {
		t.Fatalf("time.sleep did not advance frozen time: expected 1767606165000, got %d", now)
	}
}

func TestTimeParseGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`time.parse("02-01-2026 16:23:48", "%d-%m-%Y %h:%i:%s")`, 1767371028000},
		{`time.parse("2026/03/25 08:45:25 PM", "%Y/%m/%d %H:%i:%s %A")`, 1774471525000},
		{`time.parse("03-25-2026 08:45:25 am", "%m-%d-%Y %H:%i:%s %a")`, 1774428325000},
		{`time.parse("Fri, Jan, 26", "%D, %M, %y")`, 1767225600000},
		{`time.parse("Fri, Jan, 26", "%D, %M, %y")`, 1767225600000},
		{`time.parse("27 February 1993", "%d %F %Y")`, 730771200000},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestTimeFormatGlobalFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`(time.format(1767371028000, "%Y-%m-%d %H:%i:%s"))`, "2026-01-02 04:23:48"},
		{`(time.format(1774471525000, "%Y/%m/%d %H-%i-%s"))`, "2026/03/25 08-45-25"},
		{`(time.format(1774428325000, "%d-%m-%Y %s:%i:%H"))`, "25-03-2026 25:45:08"},
		{`(time.format(1767225600000, "%Y-%m-%d"))`, "2026-01-01"},
		{`(time.format(1767225600000, "%Y/%m/%d"))`, "2026/01/01"},
		{`(time.format(730771200000, "%Y"))`, "1993"},
		{`(time.format(730771200000, "%D %d, %F, %y"))`, "Sat 27, February, 93"},
		{`(time.format(1767606155000, "%H:%i:%s %a"))`, "09:42:35 am"},
		{`(time.format(1767606155000, "%H:%i:%s %A"))`, "09:42:35 AM"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
