package objects

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

func TestNewFileDescriptorContext(t *testing.T) {
	path := "/path/to/file.zen"
	fileCtx := NewFileDescriptorContext(path)

	if fileCtx.Name != "file.zen" {
		t.Errorf("expected file name to be 'file.zen', got '%s'", fileCtx.Name)
	}

	expectedPath := "/path/to"
	if fileCtx.Path != expectedPath {
		t.Errorf(
			"expected file path to be '%s', got '%s'",
			expectedPath, fileCtx.Path,
		)
	}
}

func TestNativeErrorToErrorObject(t *testing.T) {
	nativeErr := fmt.Errorf("this is a native error")
	errObj := NativeErrorToErrorObject(nativeErr)

	if errObj.Message != nativeErr.Error() {
		t.Errorf(
			"expected error message to be '%s', got '%s'",
			nativeErr.Error(), errObj.Message,
		)
	}
}

func TestNewErrorWithFileContext(t *testing.T) {
	token := tokens.Token{Line: 10, Column: 5}
	fileCtx := &FileDescriptorContext{
		Name: "file.zen",
		Path: "/path/to",
	}

	errObj := NewError(token, fileCtx, "an error occurred: %s", "details")

	if errObj.Message != "an error occurred: details" {
		t.Errorf(
			"expected error message to be 'an error occurred: details', got '%s'",
			errObj.Message,
		)
	}

	if errObj.Line != token.Line {
		t.Errorf(
			"expected error line to be %d, got %d",
			token.Line, errObj.Line,
		)
	}

	if errObj.Column != token.Column {
		t.Errorf(
			"expected error column to be %d, got %d",
			token.Column, errObj.Column,
		)
	}

	if errObj.File != fileCtx.Name {
		t.Errorf(
			"expected error file to be '%s', got '%s'",
			fileCtx.Name, errObj.File,
		)
	}

	if errObj.Path != fileCtx.Path {
		t.Errorf(
			"expected error path to be '%s', got '%s'",
			fileCtx.Path, errObj.Path,
		)
	}
}

func TestNewErrorWithoutFileContext(t *testing.T) {
	token := tokens.Token{Line: 15, Column: 8}

	errObj := NewError(token, nil, "an error occurred: %s", "details")

	if errObj.Message != "an error occurred: details" {
		t.Errorf(
			"expected error message to be 'an error occurred: details', got '%s'",
			errObj.Message,
		)
	}

	if errObj.Line != token.Line {
		t.Errorf(
			"expected error line to be %d, got %d",
			token.Line, errObj.Line,
		)
	}

	if errObj.Column != token.Column {
		t.Errorf(
			"expected error column to be %d, got %d",
			token.Column, errObj.Column,
		)
	}

	if errObj.File != "" {
		t.Errorf(
			"expected error file to be empty, got '%s'",
			errObj.File,
		)
	}
}

func TestNewEmptyErrorWithParent(t *testing.T) {
	parentErr := &Error{Message: "parent error"}
	token := tokens.Token{Line: 20, Column: 12}
	fileCtx := &FileDescriptorContext{
		Name: "file.zen",
		Path: "/path/to",
	}

	errObj := NewEmptyErrorWithParent(parentErr, token, fileCtx)

	if errObj.Message != "" {
		t.Errorf(
			"expected error message to be empty, got '%s'",
			errObj.Message,
		)
	}

	if errObj.Line != token.Line {
		t.Errorf(
			"expected error line to be %d, got %d",
			token.Line, errObj.Line,
		)
	}

	if errObj.Column != token.Column {
		t.Errorf(
			"expected error column to be %d, got %d",
			token.Column, errObj.Column,
		)
	}

	if errObj.File != fileCtx.Name {
		t.Errorf(
			"expected error file to be '%s', got '%s'",
			fileCtx.Name, errObj.File,
		)
	}

	if errObj.Path != fileCtx.Path {
		t.Errorf(
			"expected error path to be '%s', got '%s'",
			fileCtx.Path, errObj.Path,
		)
	}

	if errObj.Parent != parentErr {
		t.Errorf(
			"expected error parent to be '%v', got '%v'",
			parentErr, errObj.Parent,
		)
	}
}

func TestIsError(t *testing.T) {
	if !IsError(&Error{}) {
		t.Errorf("expected IsError to return true for Error object")
	}

	objects := []Object{
		TRUE, FALSE, NULL,
		&String{Value: "test"},
		&Integer{Value: 10},
		&Float{Value: 3.14},
		&Array{Elements: []Object{}},
		&Hash{Pairs: map[HashKey]HashPair{}},
		nil,
	}

	for _, obj := range objects {
		if IsError(obj) {
			t.Errorf("expected IsError to return false for object of type %s", obj.Type())
		}
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		input    Object
		expected bool
	}{
		{TRUE, true},
		{FALSE, false},
		{NULL, false},
		{&Integer{Value: 10}, true},
		{&Integer{Value: 0}, true},
		{&String{Value: "hello"}, true},
		{&String{Value: ""}, true},
		{&Array{Elements: []Object{}}, true},
		{&Hash{Pairs: map[HashKey]HashPair{}}, true},
	}

	for _, tt := range tests {
		result := IsTruthy(tt.input)

		if result != tt.expected {
			t.Errorf("IsTruthy(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		input    ObjectType
		expected bool
	}{
		{INTEGER_OBJ, true},
		{FLOAT_OBJ, true},
		{STRING_OBJ, false},
		{BOOLEAN_OBJ, false},
		{ARRAY_OBJ, false},
		{HASH_OBJ, false},
		{NULL_OBJ, false},
	}

	for _, tt := range tests {
		result := IsNumber(tt.input)

		if result != tt.expected {
			t.Errorf("IsNumber(%v) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestWrapNumberValue(t *testing.T) {
	tests := []struct {
		value  float64
		left   Object
		right  Object
		result Object
	}{
		{10.0, &Integer{Value: 5}, &Integer{Value: 5}, &Integer{Value: 10}},
		{10.2, &Integer{Value: 5}, &Integer{Value: 5}, &Float{Value: 10.2}},
		{10.5, &Float{Value: 5.0}, &Integer{Value: 5}, &Float{Value: 10.5}},
		{10.0, &Integer{Value: 5}, &Float{Value: 5.0}, &Float{Value: 10.0}},
		{10.7, &Float{Value: 5.0}, &Float{Value: 5.7}, &Float{Value: 10.7}},
	}

	for _, tt := range tests {
		result := WrapNumberValue(tt.value, tt.left, tt.right)

		if !reflect.DeepEqual(result, tt.result) {
			t.Errorf(
				"WrapNumberValue(%v, %v, %v) = %v[%s], want %v[%s]",
				tt.value, tt.left.Inspect(), tt.right.Inspect(),
				result.Inspect(), result.Type(), tt.result.Inspect(), tt.result.Type(),
			)
		}
	}
}

func TestUnwrapNumberValue(t *testing.T) {
	tests := []struct {
		input    Object
		expected float64
	}{
		{&Integer{Value: 10}, 10.0},
		{&Float{Value: 3.14}, 3.14},
		{&String{Value: "hello"}, 0.0},
		{TRUE, 0.0},
		{NULL, 0.0},
	}

	for _, tt := range tests {
		result := UnwrapNumberValue(tt.input)

		if result != tt.expected {
			t.Errorf(
				"UnwrapNumberValue(%v) = %v, want %v",
				tt.input.Inspect(), result, tt.expected,
			)
		}
	}
}

func TestIsStringable(t *testing.T) {
	tests := []struct {
		input    Object
		expected bool
	}{
		{&Integer{Value: 10}, true},
		{&Float{Value: 3.14}, true},
		{TRUE, true},
		{FALSE, true},
		{&String{Value: "hello"}, false},
		{&Array{Elements: []Object{}}, false},
		{&Hash{Pairs: map[HashKey]HashPair{}}, false},
		{NULL, false},
	}

	for _, tt := range tests {
		result := IsStringable(tt.input)

		if result != tt.expected {
			t.Errorf(
				"IsStringable(%v) = %v, want %v",
				tt.input.Inspect(), result, tt.expected,
			)
		}
	}
}

func TestStringifyObject(t *testing.T) {
	tests := []struct {
		input    Object
		expected string
	}{
		{&String{Value: "hello"}, "hello"},
		{&Integer{Value: 10}, "10"},
		{&Float{Value: 3.14}, "3.14"},
		{&Float{Value: 42.198765}, "42.198765"},
		{TRUE, "true"},
		{FALSE, "false"},
		{&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}}, "[1, 2]"},
		{&Hash{Pairs: map[HashKey]HashPair{
			(&String{Value: "key"}).HashKey(): {
				Key:   &String{Value: "key"},
				Value: &String{Value: "value"},
			}}},
			"{key: value}",
		},
		{NULL, "null"},
	}

	for _, tt := range tests {
		result := StringifyObject(tt.input)

		if result != tt.expected {
			t.Errorf(
				"StringifyObject(%v) = %v, want %v",
				tt.input.Inspect(), result, tt.expected,
			)
		}
	}
}

func TestUnwrapReturnValue(t *testing.T) {
	tests := []struct {
		input    Object
		expected Object
	}{
		{&ReturnValue{Value: &Integer{Value: 10}}, &Integer{Value: 10}},
		{&Integer{Value: 20}, &Integer{Value: 20}},
		{&ReturnValue{Value: &String{Value: "hello"}}, &String{Value: "hello"}},
		{&String{Value: "world"}, &String{Value: "world"}},
		{NULL, NULL},
	}

	for _, tt := range tests {
		result := UnwrapReturnValue(tt.input)

		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf(
				"UnwrapReturnValue(%v) = %v, want %v",
				tt.input.Inspect(), result.Inspect(), tt.expected.Inspect(),
			)
		}
	}
}

func TestNativeBoolToBooleanObject(t *testing.T) {
	if NativeBoolToBooleanObject(true) != TRUE {
		t.Errorf(
			"NativeBoolToBooleanObject(true) = %v, want TRUE",
			NativeBoolToBooleanObject(true),
		)
	}

	if NativeBoolToBooleanObject(false) != FALSE {
		t.Errorf(
			"NativeBoolToBooleanObject(false) = %v, want FALSE",
			NativeBoolToBooleanObject(false),
		)
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		left     Object
		right    Object
		expected *Boolean
	}{
		{&Integer{Value: 10}, &Integer{Value: 10}, TRUE},
		{&Float{Value: 3.14}, &Float{Value: 3.14}, TRUE},
		{&String{Value: "hello"}, &String{Value: "hello"}, TRUE},
		{
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			TRUE,
		},
		{
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "key1"}).HashKey(): {
					Key:   &String{Value: "key1"},
					Value: &Integer{Value: 1},
				},
				(&String{Value: "key2"}).HashKey(): {
					Key:   &String{Value: "key2"},
					Value: &Integer{Value: 2},
				},
			}},
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "key1"}).HashKey(): {
					Key:   &String{Value: "key1"},
					Value: &Integer{Value: 1},
				},
				(&String{Value: "key2"}).HashKey(): {
					Key:   &String{Value: "key2"},
					Value: &Integer{Value: 2},
				},
			}},
			TRUE,
		},
		{
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Hash{Pairs: map[HashKey]HashPair{
						(&String{Value: "nestedKey"}).HashKey(): {
							Key:   &String{Value: "nestedKey"},
							Value: &String{Value: "nestedValue"},
						},
					}},
				},
			}},
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Hash{Pairs: map[HashKey]HashPair{
						(&String{Value: "nestedKey"}).HashKey(): {
							Key:   &String{Value: "nestedKey"},
							Value: &String{Value: "nestedValue"},
						},
					}},
				},
			}},
			TRUE,
		},
		{
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Array{Elements: []Object{
						&Integer{Value: 1},
						&String{Value: "two"},
						&Float{Value: 3.0},
					}},
				},
			}},
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Array{Elements: []Object{
						&Integer{Value: 1},
						&String{Value: "two"},
						&Float{Value: 3.0},
					}},
				},
			}},
			TRUE,
		},
		{TRUE, TRUE, TRUE},
		{FALSE, FALSE, TRUE},
		{NULL, NULL, TRUE},

		{&Integer{Value: 10}, &Integer{Value: 20}, FALSE},
		{&Integer{Value: 10}, &Float{Value: 10.0}, FALSE},
		{&Float{Value: 3.14}, &Float{Value: 3.13}, FALSE},
		{&String{Value: "hello"}, &String{Value: "world"}, FALSE},
		{&String{Value: "10"}, &Integer{Value: 10}, FALSE},
		{
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 0}, &Integer{Value: 3}}},
			FALSE,
		},
		{
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "key1"}).HashKey(): {
					Key:   &String{Value: "key1"},
					Value: &Integer{Value: 1},
				},
				(&String{Value: "key2"}).HashKey(): {
					Key:   &String{Value: "key2"},
					Value: &Integer{Value: 2},
				},
			}},
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "key1"}).HashKey(): {
					Key:   &String{Value: "key1"},
					Value: &Integer{Value: 1},
				},
				(&String{Value: "key2"}).HashKey(): {
					Key:   &String{Value: "key2"},
					Value: &Integer{Value: 20},
				},
			}},
			FALSE,
		},
		{
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Hash{Pairs: map[HashKey]HashPair{
						(&String{Value: "nestedKey"}).HashKey(): {
							Key:   &String{Value: "nestedKey"},
							Value: &String{Value: "nestedValue"},
						},
					}},
				},
			}},
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Hash{Pairs: map[HashKey]HashPair{
						(&String{Value: "nestedKey"}).HashKey(): {
							Key:   &String{Value: "nestedKey"},
							Value: &String{Value: "new-value"},
						},
					}},
				},
			}},
			FALSE,
		},
		{
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Array{Elements: []Object{
						&Integer{Value: 1},
						&String{Value: "two"},
						&Float{Value: 3.0},
					}},
				},
			}},
			&Hash{Pairs: map[HashKey]HashPair{
				(&String{Value: "obj"}).HashKey(): {
					Key: &String{Value: "obj"},
					Value: &Array{Elements: []Object{
						&Integer{Value: 1},
						&String{Value: "two"},
						&Float{Value: 3.01},
					}},
				},
			}},
			FALSE,
		},
		{TRUE, FALSE, FALSE},
		{NULL, &String{Value: "null"}, FALSE},
	}

	for _, tt := range tests {
		result := Equals(tt.left, tt.right)

		if result != tt.expected {
			t.Errorf(
				"Equals(%v, %v) = %v, want %v",
				tt.left.Inspect(), tt.right.Inspect(), result.Inspect(), tt.expected.Inspect(),
			)
		}
	}
}

func TestCreateImmutableHashFromEnvExports(t *testing.T) {
	env := &Environment{
		exports: map[string]Object{
			"var1": &Integer{Value: 10},
			"var2": &String{Value: "hello"},
			"var3": &Float{Value: 3.14},
		},
	}

	immutableHash := CreateImmutableHashFromEnvExports(env)

	expectedPairs := map[HashKey]HashPair{
		(&String{Value: "var1"}).HashKey(): {
			Key:   &String{Value: "var1"},
			Value: &Integer{Value: 10},
		},
		(&String{Value: "var2"}).HashKey(): {
			Key:   &String{Value: "var2"},
			Value: &String{Value: "hello"},
		},
		(&String{Value: "var3"}).HashKey(): {
			Key:   &String{Value: "var3"},
			Value: &Float{Value: 3.14},
		},
	}

	if len(immutableHash.Value.Pairs) != len(expectedPairs) {
		t.Errorf(
			"expected %d pairs, got %d pairs",
			len(expectedPairs), len(immutableHash.Value.Pairs),
		)
	}

	for key, expectedPair := range expectedPairs {
		actualPair, ok := immutableHash.Value.Pairs[key]
		if !ok {
			t.Errorf("expected key '%s' to be present in immutable hash", expectedPair.Key.Inspect())
			continue
		}

		if !Equals(actualPair.Key, expectedPair.Key).Value {
			t.Errorf(
				"expected key %v, got %v",
				expectedPair.Key.Inspect(), actualPair.Key.Inspect(),
			)
		}

		if !Equals(actualPair.Value, expectedPair.Value).Value {
			t.Errorf(
				"expected value %v, got %v",
				expectedPair.Value.Inspect(), actualPair.Value.Inspect(),
			)
		}
	}
}

func TestBuiltImmutableHash(t *testing.T) {
	pairs := []HashPair{
		{
			Key:   &String{Value: "key1"},
			Value: &Integer{Value: 1},
		},
		{
			Key:   &String{Value: "key2"},
			Value: &String{Value: "value2"},
		},
		{
			Key:   &String{Value: "key3"},
			Value: &Float{Value: 3.14},
		},
	}

	immutableHash := BuildImmutableHash(pairs...)

	if immutableHash.Type() != IMMUTABLE_HASH_OBJ {
		t.Errorf(
			"expected type %s, got %s",
			IMMUTABLE_HASH_OBJ, immutableHash.Type(),
		)
	}

	if len(immutableHash.Value.Pairs) != len(pairs) {
		t.Errorf(
			"expected %d pairs, got %d pairs",
			len(pairs), len(immutableHash.Value.Pairs),
		)
	}

	for _, expectedPair := range pairs {
		actualPair, ok := immutableHash.Value.Pairs[expectedPair.Key.(Hashable).HashKey()]
		if !ok {
			t.Errorf("expected key '%s' to be present in immutable hash", expectedPair.Key.Inspect())
			continue
		}

		if !Equals(actualPair.Key, expectedPair.Key).Value {
			t.Errorf(
				"expected key %v, got %v",
				expectedPair.Key.Inspect(), actualPair.Key.Inspect(),
			)
		}

		if !Equals(actualPair.Value, expectedPair.Value).Value {
			t.Errorf(
				"expected value %v, got %v",
				expectedPair.Value.Inspect(), actualPair.Value.Inspect(),
			)
		}
	}
}

func TestWrapBuiltinFunctionInASTAwareMap(t *testing.T) {
	builtinFn := &Builtin{Fn: func(args ...Object) (Object, error) {
		return &String{Value: "hello from builtin"}, nil
	}}

	astAwareMap := WrapBuiltinFunctionInASTAwareMap("testFunc", builtinFn)

	if astAwareMap.Key.(*String).Value != "testFunc" {
		t.Errorf(
			"expected key to be 'testFunc', got '%s'",
			astAwareMap.Key.(*String).Value,
		)
	}

	astAwareBuiltin, ok := astAwareMap.Value.(*ASTAwareBuiltin)
	if !ok {
		t.Fatalf(
			"expected value to be of type ASTAwareBuiltin, got %s",
			astAwareMap.Value.Type(),
		)
	}

	result := astAwareBuiltin.Fn(&ast.CallExpression{}, nil)
	if result.Inspect() != "hello from builtin" {
		t.Errorf(
			"expected builtin function to return 'hello from builtin', got '%s'",
			result.Inspect(),
		)
	}
}
