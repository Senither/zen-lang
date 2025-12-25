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

	objects := []struct {
		name  string
		value Object
	}{
		{"true", TRUE},
		{"false", FALSE},
		{"null", NULL},
		{"string", &String{Value: "test"}},
		{"integer", &Integer{Value: 10}},
		{"float", &Float{Value: 3.14}},
		{"array", &Array{Elements: []Object{}}},
		{"hash", &Hash{Pairs: map[HashKey]HashPair{}}},
		{"nil", nil},
	}

	for _, obj := range objects {
		t.Run("checking is error for "+obj.name, func(t *testing.T) {
			if IsError(obj.value) {
				t.Errorf("expected IsError to return false for object of type %s", obj.value.Type())
			}
		})
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected bool
	}{
		{"true", TRUE, true},
		{"false", FALSE, false},
		{"null", NULL, false},
		{"integer positive", &Integer{Value: 10}, true},
		{"integer zero", &Integer{Value: 0}, true},
		{"string non-empty", &String{Value: "hello"}, true},
		{"string empty", &String{Value: ""}, true},
		{"array empty", &Array{Elements: []Object{}}, true},
		{"hash empty", &Hash{Pairs: map[HashKey]HashPair{}}, true},
	}

	for _, tt := range tests {
		t.Run("is truthy: "+tt.name, func(t *testing.T) {
			result := IsTruthy(tt.input)

			if result != tt.expected {
				t.Errorf("IsTruthy(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    ObjectType
		expected bool
	}{
		{"integer", INTEGER_OBJ, true},
		{"float", FLOAT_OBJ, true},
		{"string", STRING_OBJ, false},
		{"boolean", BOOLEAN_OBJ, false},
		{"array", ARRAY_OBJ, false},
		{"hash", HASH_OBJ, false},
		{"null", NULL_OBJ, false},
	}

	for _, tt := range tests {
		t.Run("is number: "+tt.name, func(t *testing.T) {
			result := IsNumber(tt.input)

			if result != tt.expected {
				t.Errorf("IsNumber(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestWrapNumberValue(t *testing.T) {
	tests := []struct {
		name   string
		value  float64
		left   Object
		right  Object
		result Object
	}{
		{"from integer to integer", 10.0, &Integer{Value: 5}, &Integer{Value: 5}, &Integer{Value: 10}},
		{"from integer to float", 10.2, &Integer{Value: 5}, &Integer{Value: 5}, &Float{Value: 10.2}},
		{"from float to float", 10.5, &Float{Value: 5.0}, &Integer{Value: 5}, &Float{Value: 10.5}},
		{"from integer to float", 10.0, &Integer{Value: 5}, &Float{Value: 5.0}, &Float{Value: 10.0}},
		{"to float with all floats", 10.7, &Float{Value: 5.0}, &Float{Value: 5.7}, &Float{Value: 10.7}},
	}

	for _, tt := range tests {
		t.Run("wrap number: "+tt.name, func(t *testing.T) {
			result := WrapNumberValue(tt.value, tt.left, tt.right)

			if !reflect.DeepEqual(result, tt.result) {
				t.Errorf(
					"WrapNumberValue(%v, %v, %v) = %v[%s], want %v[%s]",
					tt.value, tt.left.Inspect(), tt.right.Inspect(),
					result.Inspect(), result.Type(), tt.result.Inspect(), tt.result.Type(),
				)
			}
		})
	}
}

func TestUnwrapNumberValue(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected float64
	}{
		{"integer", &Integer{Value: 10}, 10.0},
		{"float", &Float{Value: 3.14}, 3.14},
		{"string", &String{Value: "hello"}, 0.0},
		{"true", TRUE, 0.0},
		{"null", NULL, 0.0},
	}

	for _, tt := range tests {
		t.Run("unwrap number: "+tt.name, func(t *testing.T) {
			result := UnwrapNumberValue(tt.input)

			if result != tt.expected {
				t.Errorf(
					"UnwrapNumberValue(%v) = %v, want %v",
					tt.input.Inspect(), result, tt.expected,
				)
			}
		})
	}
}

func TestIsStringable(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected bool
	}{
		{"integer", &Integer{Value: 10}, true},
		{"float", &Float{Value: 3.14}, true},
		{"true", TRUE, true},
		{"false", FALSE, true},
		{"string", &String{Value: "hello"}, false},
		{"array", &Array{Elements: []Object{}}, false},
		{"hash", &Hash{Pairs: map[HashKey]HashPair{}}, false},
		{"null", NULL, false},
	}

	for _, tt := range tests {
		t.Run("is stringable: "+tt.name, func(t *testing.T) {
			result := IsStringable(tt.input)

			if result != tt.expected {
				t.Errorf(
					"IsStringable(%v) = %v, want %v",
					tt.input.Inspect(), result, tt.expected,
				)
			}
		})
	}
}

func TestStringifyObject(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected string
	}{
		{"string", &String{Value: "hello"}, "hello"},
		{"integer", &Integer{Value: 10}, "10"},
		{"float", &Float{Value: 3.14}, "3.14"},
		{"float", &Float{Value: 42.198765}, "42.198765"},
		{"true", TRUE, "true"},
		{"false", FALSE, "false"},
		{"array", &Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}}}, "[1, 2]"},
		{"hash", &Hash{Pairs: map[HashKey]HashPair{
			(&String{Value: "key"}).HashKey(): {
				Key:   &String{Value: "key"},
				Value: &String{Value: "value"},
			}}},
			"{key: value}",
		},
		{"null", NULL, "null"},
	}

	for _, tt := range tests {
		t.Run("stringify object: "+tt.name, func(t *testing.T) {
			result := StringifyObject(tt.input)

			if result != tt.expected {
				t.Errorf(
					"StringifyObject(%v) = %v, want %v",
					tt.input.Inspect(), result, tt.expected,
				)
			}
		})
	}
}

func TestUnwrapReturnValue(t *testing.T) {
	tests := []struct {
		name     string
		input    Object
		expected Object
	}{
		{"return value with integer", &ReturnValue{Value: &Integer{Value: 10}}, &Integer{Value: 10}},
		{"integer", &Integer{Value: 20}, &Integer{Value: 20}},
		{"return value with string", &ReturnValue{Value: &String{Value: "hello"}}, &String{Value: "hello"}},
		{"string", &String{Value: "world"}, &String{Value: "world"}},
		{"null", NULL, NULL},
	}

	for _, tt := range tests {
		t.Run("unwrap return value: "+tt.name, func(t *testing.T) {
			result := UnwrapReturnValue(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf(
					"UnwrapReturnValue(%v) = %v, want %v",
					tt.input.Inspect(), result.Inspect(), tt.expected.Inspect(),
				)
			}
		})
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
		name     string
		left     Object
		right    Object
		expected *Boolean
	}{
		{"integer", &Integer{Value: 10}, &Integer{Value: 10}, TRUE},
		{"float", &Float{Value: 3.14}, &Float{Value: 3.14}, TRUE},
		{"string", &String{Value: "hello"}, &String{Value: "hello"}, TRUE},
		{
			"array",
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			TRUE,
		},
		{
			"hash",
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
			"nested hash",
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
			"hash with array",
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
		{"true", TRUE, TRUE, TRUE},
		{"false", FALSE, FALSE, TRUE},
		{"null", NULL, NULL, TRUE},

		{"integer", &Integer{Value: 10}, &Integer{Value: 20}, FALSE},
		{"integer vs float", &Integer{Value: 10}, &Float{Value: 10.0}, FALSE},
		{"float", &Float{Value: 3.14}, &Float{Value: 3.13}, FALSE},
		{"string", &String{Value: "hello"}, &String{Value: "world"}, FALSE},
		{"string vs integer", &String{Value: "10"}, &Integer{Value: 10}, FALSE},
		{
			"array values differ",
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 2}, &Integer{Value: 3}}},
			&Array{Elements: []Object{&Integer{Value: 1}, &Integer{Value: 0}, &Integer{Value: 3}}},
			FALSE,
		},
		{
			"hash values differ",
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
			"nested hash values differ",
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
			"hash with array values differ",
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
		{"boolean differ", TRUE, FALSE, FALSE},
		{"string 'null' vs null", NULL, &String{Value: "null"}, FALSE},
	}

	for _, tt := range tests {
		t.Run("equals: "+tt.name, func(t *testing.T) {
			result := Equals(tt.left, tt.right)

			if result != tt.expected {
				t.Errorf(
					"Equals(%v, %v) = %v, want %v",
					tt.left.Inspect(), tt.right.Inspect(), result.Inspect(), tt.expected.Inspect(),
				)
			}
		})
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

func TestCopyNatives(t *testing.T) {
	tests := []struct {
		input Object
	}{
		{&Integer{Value: 10}},
		{&Float{Value: 3.14}},
		{&String{Value: "Hello, Zen"}},
		{&String{Value: ""}},
	}

	for _, tt := range tests {
		result := Copy(tt.input)

		if !reflect.DeepEqual(result, tt.input) {
			t.Errorf(
				"Copy(%v) = %v, want a deep copy of input",
				tt.input.Inspect(), result.Inspect(),
			)
		}

		if fmt.Sprintf("%p", result) == fmt.Sprintf("%p", tt.input) {
			t.Errorf(
				"Copy(%v) returned the same reference, expected a different one",
				tt.input.Inspect(),
			)
		}
	}
}

func TestCopyStaticNatives(t *testing.T) {
	tests := []struct {
		input Object
	}{
		{TRUE},
		{FALSE},
		{NULL},
	}

	for _, tt := range tests {
		result := Copy(tt.input)

		if result != tt.input {
			t.Errorf(
				"Copy(%v) = %v, want the same reference as input",
				tt.input.Inspect(), result.Inspect(),
			)
		}
	}
}

func TestCopyComplexObjects(t *testing.T) {
	arrayInput := &Array{
		Elements: []Object{
			&Integer{Value: 1},
			&String{Value: "two"},
			&Float{Value: 3.0},
		},
	}

	result := Copy(arrayInput)

	if !reflect.DeepEqual(result, arrayInput) {
		t.Errorf(
			"Copy(%v) = %v, want a deep copy of input",
			arrayInput.Inspect(), result.Inspect(),
		)
	}

	if fmt.Sprintf("%p", result) == fmt.Sprintf("%p", arrayInput) {
		t.Errorf(
			"Copy(%v) returned the same reference, expected a different one",
			arrayInput.Inspect(),
		)
	}

	arrayResult := result.(*Array)
	for i, elem := range arrayInput.Elements {
		if fmt.Sprintf("%p", arrayResult.Elements[i]) == fmt.Sprintf("%p", elem) {
			t.Errorf(
				"Copy(%v) element at index %d returned the same reference, expected a different one",
				arrayInput.Inspect(), i,
			)
		}
	}

	hashInput := &Hash{
		Pairs: map[HashKey]HashPair{
			(&String{Value: "key1"}).HashKey(): {
				Key:   &String{Value: "key1"},
				Value: &Integer{Value: 1},
			},
			(&String{Value: "key2"}).HashKey(): {
				Key:   &String{Value: "key2"},
				Value: &String{Value: "value2"},
			},
		},
	}

	result = Copy(hashInput)

	if !reflect.DeepEqual(result, hashInput) {
		t.Errorf(
			"Copy(%v) = %v, want a deep copy of input",
			hashInput.Inspect(), result.Inspect(),
		)
	}

	if fmt.Sprintf("%p", result) == fmt.Sprintf("%p", hashInput) {
		t.Errorf(
			"Copy(%v) returned the same reference, expected a different one",
			hashInput.Inspect(),
		)
	}

	hashResult := result.(*Hash)
	for key, pair := range hashInput.Pairs {
		resultPair, ok := hashResult.Pairs[key]
		if !ok {
			t.Errorf(
				"Copy(%v) missing key %v in result",
				hashInput.Inspect(), pair.Key.Inspect(),
			)
			continue
		}
		if fmt.Sprintf("%p", resultPair.Key) == fmt.Sprintf("%p", pair.Key) {
			t.Errorf(
				"Copy(%v) key %v returned the same reference, expected a different one",
				hashInput.Inspect(), pair.Key.Inspect(),
			)
		}

		if fmt.Sprintf("%p", resultPair.Value) == fmt.Sprintf("%p", pair.Value) {
			t.Errorf(
				"Copy(%v) value for key %v returned the same reference, expected a different one",
				hashInput.Inspect(), pair.Key.Inspect(),
			)
		}
	}
}
