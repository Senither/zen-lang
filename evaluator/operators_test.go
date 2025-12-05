package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/objects"
)

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
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
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
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
	}
}

func TestVarReassignmentFailure(t *testing.T) {
	input := []struct {
		input    string
		expected *objects.Error
	}{
		{
			`var x = 5; x = 10; x;`,
			&objects.Error{Message: "cannot modify immutable variable: x"},
		},
		{
			`var x = 3.14; x = 4.14; x;`,
			&objects.Error{Message: "cannot modify immutable variable: x"},
		},
		{
			`var x = "hello"; x = "world"; x;`,
			&objects.Error{Message: "cannot modify immutable variable: x"},
		},
		{
			`var name = "Senither"; name = "test"; name;`,
			&objects.Error{Message: "cannot modify immutable variable: name"},
		},
		{
			`var arr = [1,2,3]; arr = ['another', 'array']; arr;`,
			&objects.Error{Message: "cannot modify immutable variable: arr"},
		},
		{
			`var obj = {"key": "value"}; obj = {"key": "new value"}; obj;`,
			&objects.Error{Message: "cannot modify immutable variable: obj"},
		},
	}

	for _, tt := range input {
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
	}
}

func TestVarIncrementingStatements(t *testing.T) {
	input := []struct {
		input    string
		expected interface{}
	}{
		{`var x = 5; x++; x;`, 6},
		{`var mut x = 5; x++; x;`, 6},
		{`var x = 3.14; x++; x;`, float64(4.14)},
		{`var mut x = 3.14; x++; x;`, float64(4.14)},
		{`var x = "hello"; x++; x;`, &objects.Error{Message: "unknown operator: ++STRING"}},
		{`var mut x = "hello"; x++; x;`, &objects.Error{Message: "unknown operator: ++STRING"}},
	}

	for _, tt := range input {
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
	}
}

func TestVarDecrementingStatements(t *testing.T) {
	input := []struct {
		input    string
		expected interface{}
	}{
		{`var x = 5; x--; x;`, 4},
		{`var mut x = 5; x--; x;`, 4},
		{`var x = 3.14; x--; x;`, float64(2.14)},
		{`var mut x = 3.14; x--; x;`, float64(2.14)},
		{`var x = "hello"; x--; x;`, &objects.Error{Message: "unknown operator: --STRING"}},
		{`var mut x = "hello"; x--; x;`, &objects.Error{Message: "unknown operator: --STRING"}},
	}

	for _, tt := range input {
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
	}
}

func TestArrayReassignmentStatements(t *testing.T) {
	input := []struct {
		input    string
		expected []any
	}{
		{`var arr = [1, 2, 3]; arr[0] = 10; arr;`, []any{10, 2, 3}},
		{`var mut arr = [1, 2, 3]; arr[0] = 10; arr;`, []any{10, 2, 3}},
		{`var arr = [1, 2, 3]; arr[1] = 10; arr;`, []any{1, 10, 3}},
		{`var mut arr = [1, 2, 3]; arr[1] = 10; arr;`, []any{1, 10, 3}},
		{`var arr = ["foo", "bar"]; arr[1] = "baz"; arr;`, []any{"foo", "baz"}},
		{`var mut arr = ["foo", "bar"]; arr[1] = "baz"; arr;`, []any{"foo", "baz"}},
	}

	for _, tt := range input {
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
	}
}

func TestHashReassignmentStatements(t *testing.T) {
	input := []struct {
		input    string
		expected map[objects.HashKey]any
	}{
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash["foo"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash["foo"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash["baz"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash["baz"] = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash.foo = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash.foo = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 10,
				(&objects.String{Value: "bar"}).HashKey(): 2,
			},
		},
		{
			`var mut hash = {"foo": 1, "bar": 2}; hash.baz = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
		{
			`var hash = {"foo": 1, "bar": 2}; hash.baz = 10; hash;`,
			map[objects.HashKey]any{
				(&objects.String{Value: "foo"}).HashKey(): 1,
				(&objects.String{Value: "bar"}).HashKey(): 2,
				(&objects.String{Value: "baz"}).HashKey(): 10,
			},
		},
	}

	for _, tt := range input {
		objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
	}
}
