package evaluator

import (
	"testing"

	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

func testEval(input string) objects.Object {
	l := lexer.New(input)
	p := parser.New(l, nil)

	return Eval(p.ParseProgram(), objects.NewEnvironment(nil))
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *objects.Error
	}{
		{
			"int+bool without semicolon",
			"5 + true;",
			&objects.Error{Message: "type mismatch: INTEGER + BOOLEAN"},
		},
		{
			"int+bool without semicolon and extra int",
			"5 + true; 5;",
			&objects.Error{Message: "type mismatch: INTEGER + BOOLEAN"},
		},
		{
			"minus boolean",
			"-true",
			&objects.Error{Message: "unknown operator: -BOOLEAN"},
		},
		{
			"boolean addition",
			"true + false;",
			&objects.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"},
		},
		{
			"multiple boolean additions",
			"true + false + true + false;",
			&objects.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"},
		},
		{
			"boolean addition with int",
			"5; true + false; 5",
			&objects.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"},
		},
		{
			"if statement with boolean addition",
			"if (10 > 1) { true + false; }",
			&objects.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"},
		},
		{
			"nested if statement with boolean addition",
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}

				return 1;
			}
		`,
			&objects.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"},
		},
		{
			"undefined identifier",
			"foobar",
			&objects.Error{Message: "identifier not found: foobar"},
		},
		{
			"invalid hash key type",
			`{"name": "value"}[func (x) { x }]`,
			&objects.Error{Message: "invalid type given as hash key: FUNCTION"},
		},
		{
			"deeply nested function calls with boolean addition",
			`
			func a(x) {
				return x();
			}
			func b(y) {
				return y(a);
			}
			func c(z) {
				return z(b);
			}

			println(c(func () {
				return true + false;
			}));
			`,
			&objects.Error{Message: "unknown operator: BOOLEAN + BOOLEAN"},
		},
	}

	for _, tt := range tests {
		t.Run("error handling: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func (x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*objects.Function)
	if !ok {
		t.Fatalf("object is not Function. got %T (%+v)\n%s", evaluated, evaluated, evaluated.Inspect())
	}

	if fn.Name != nil {
		t.Fatalf("function name is not nil. got %q", fn.Name.String())
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, fn.Body.String())
	}
}

func TestNamedFunctionObject(t *testing.T) {
	input := "func hello(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*objects.Function)
	if !ok {
		t.Fatalf("object is not Function. got %T (%+v)\n%s", evaluated, evaluated, evaluated.Inspect())
	}

	if fn.Name == nil || fn.Name.String() != "hello" {
		t.Fatalf("function name is not 'hello', got %q", fn.Name)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"var identity", "var identity = func(x) { x; }; identity(5);", 5},
		{"func identity", "func identity(x) { x; }; identity(5);", 5},
		{"var identity with return", "var identity = func(x) { return x; }; identity(5);", 5},
		{"func identity with return", "func identity(x) { return x; }; identity(5);", 5},
		{"var double", "var double = func(x) { x * 2; }; double(5);", 10},
		{"func double", "func double(x) { x * 2; }; double(5);", 10},
		{"var add", "var add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"func add", "func add(x, y) { x + y; }; add(5, 5);", 10},
		{"var add with nested calls", "var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func add with nested calls", "func add(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"anonymous function", "func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		t.Run("function application: "+tt.name, func(t *testing.T) {
			objects.AssertExpectedObject(t, tt.expected, testEval(tt.input))
		})
	}
}
