package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestBangOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"bang operator on true", "!true", false},
		{"bang operator on false", "!false", true},
		{"double bang operator on true", "!!true", true},
		{"double bang operator on false", "!!false", false},
		{"triple bang operator on true", "!!!true", false},
		{"triple bang operator on false", "!!!false", true},
		{"bang operator on integer", "!5", false},
		{"double bang operator on integer", "!!5", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestVarReassignmentStatements(t *testing.T) {
	input := []struct {
		name     string
		input    string
		expected any
	}{
		{"var reassignment int with mut", `var mut x = 5; x = 10; x;`, 10},
		{"var reassignment string with mut", `var mut x = "hello"; x = "world"; x;`, "world"},
	}

	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestVarReassignmentFailure(t *testing.T) {
	input := []struct {
		name     string
		input    string
		expected *objects.Error
	}{
		{
			"var reassignment int without mut",
			`var x = 5; x = 10; x;`,
			&objects.Error{Message: "cannot modify immutable variable: x"},
		},
		{
			"var reassignment float without mut",
			`var x = 3.14; x = 4.14; x;`,
			&objects.Error{Message: "cannot modify immutable variable: x"},
		},
		{
			"var reassignment string without mut",
			`var x = "hello"; x = "world"; x;`,
			&objects.Error{Message: "cannot modify immutable variable: x"},
		},
		{
			"var reassignment string without mut",
			`var name = "Senither"; name = "test"; name;`,
			&objects.Error{Message: "cannot modify immutable variable: name"},
		},
		{
			"var reassignment array without mut",
			`var arr = [1,2,3]; arr = ['another', 'array']; arr;`,
			&objects.Error{Message: "cannot modify immutable variable: arr"},
		},
		{
			"var reassignment hash without mut",
			`var obj = {"key": "value"}; obj = {"key": "new value"}; obj;`,
			&objects.Error{Message: "cannot modify immutable variable: obj"},
		},
	}

	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestVarIncrementingStatements(t *testing.T) {
	input := []struct {
		name     string
		input    string
		expected any
	}{
		{"var incrementing int without mut", `var x = 5; x++; x;`, 6},
		{"var incrementing int with mut", `var mut x = 5; x++; x;`, 6},
		{"var incrementing float without mut", `var x = 3.14; x++; x;`, float64(4.14)},
		{"var incrementing float with mut", `var mut x = 3.14; x++; x;`, float64(4.14)},
		{"var incrementing string without mut", `var x = "hello"; x++; x;`, &objects.Error{Message: "unknown operator: ++STRING"}},
		{"var incrementing string with mut", `var mut x = "hello"; x++; x;`, &objects.Error{Message: "unknown operator: ++STRING"}},
	}

	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestVarDecrementingStatements(t *testing.T) {
	input := []struct {
		name     string
		input    string
		expected any
	}{
		{"var decrementing int without mut", `var x = 5; x--; x;`, 4},
		{"var decrementing int with mut", `var mut x = 5; x--; x;`, 4},
		{"var decrementing float without mut", `var x = 3.14; x--; x;`, float64(2.14)},
		{"var decrementing float with mut", `var mut x = 3.14; x--; x;`, float64(2.14)},
		{"var decrementing string without mut", `var x = "hello"; x--; x;`, &objects.Error{Message: "unknown operator: --STRING"}},
		{"var decrementing string with mut", `var mut x = "hello"; x--; x;`, &objects.Error{Message: "unknown operator: --STRING"}},
	}

	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestArrayReassignmentStatements(t *testing.T) {
	input := []struct {
		name     string
		input    string
		expected []any
	}{
		{"var reassignment array without mut", `var arr = [1, 2, 3]; arr[0] = 10; arr;`, []any{10, 2, 3}},
		{"var reassignment array with mut", `var mut arr = [1, 2, 3]; arr[0] = 10; arr;`, []any{10, 2, 3}},
		{"var reassignment array without mut", `var arr = [1, 2, 3]; arr[1] = 10; arr;`, []any{1, 10, 3}},
		{"var reassignment array with mut", `var mut arr = [1, 2, 3]; arr[1] = 10; arr;`, []any{1, 10, 3}},
		{"var reassignment array without mut", `var arr = ["foo", "bar"]; arr[1] = "baz"; arr;`, []any{"foo", "baz"}},
		{"var reassignment array with mut", `var mut arr = ["foo", "bar"]; arr[1] = "baz"; arr;`, []any{"foo", "baz"}},
	}

	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestHashReassignmentStatements(t *testing.T) {
	input := []struct {
		name     string
		input    string
		expected map[objects.HashKey]any
	}{
		{
			"var reassignment hash with mut using index",
			`var mut hash = {"foo": 1, "bar": 2}; hash["foo"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			"var reassignment hash without mut using index",
			`var hash = {"foo": 1, "bar": 2}; hash["foo"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			"var reassignment hash with mut using index to new key",
			`var mut hash = {"foo": 1, "bar": 2}; hash["baz"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			"var reassignment hash without mut using index to new key",
			`var hash = {"foo": 1, "bar": 2}; hash["baz"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			"var reassignment hash with mut using dot notation",
			`var mut hash = {"foo": 1, "bar": 2}; hash.foo = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			"var reassignment hash without mut using dot notation",
			`var hash = {"foo": 1, "bar": 2}; hash.foo = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			"var reassignment hash with mut using dot notation to new key",
			`var mut hash = {"foo": 1, "bar": 2}; hash.baz = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			"var reassignment hash without mut using dot notation to new key",
			`var hash = {"foo": 1, "bar": 2}; hash.baz = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
	}

	for _, tt := range input {
		t.Run(tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}
