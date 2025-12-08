package objects

import "testing"

func TestNewErrorf(t *testing.T) {
	err := NewErrorf("myFunc", "an error occurred: %s", "details")
	expected := "error in `myFunc`: an error occurred: details"

	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewWrongNumberOfArgumentsError(t *testing.T) {
	err := NewWrongNumberOfArgumentsError("myFunc", 2, 3)
	expected := "wrong number of arguments to `myFunc`: got 3, want 2"

	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewWrongNumberOfArgumentsWantAtLeastError(t *testing.T) {
	err := NewWrongNumberOfArgumentsWantAtLeastError("myFunc", 2, 1)
	expected := "wrong number of arguments to `myFunc`: got 1, want at least 2"

	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewInvalidArgumentTypeError(t *testing.T) {
	args := []Object{&Integer{Value: 42}}
	err := NewInvalidArgumentTypeError("myFunc", STRING_OBJ, 0, args)
	expected := "argument 1 to `myFunc` has invalid type: got INTEGER, want STRING"

	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewInvalidArgumentTypesError(t *testing.T) {
	args := []Object{&Integer{Value: 42}}
	err := NewInvalidArgumentTypesError("myFunc", []ObjectType{STRING_OBJ, BOOLEAN_OBJ}, 0, args)
	expected := "argument 1 to `myFunc` has invalid type: got INTEGER, want STRING|BOOLEAN"

	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestNewInvalidArgumentTypesErrorWithQualifiers(t *testing.T) {
	args := []Object{&Integer{Value: 42}}
	err := NewInvalidArgumentTypesErrorWithQualifiers(
		"myFunc",
		[]ObjectType{STRING_OBJ, BOOLEAN_OBJ},
		[]ObjectType{ARRAY_OBJ, HASH_OBJ},
		0,
		args,
	)
	expected := "argument 1 to `myFunc` has invalid type for ARRAY|HASH: got INTEGER, want STRING|BOOLEAN"

	if err.Error() != expected {
		t.Errorf("expected '%s', got '%s'", expected, err.Error())
	}
}

func TestStringifyObjectTypes(t *testing.T) {
	types := []ObjectType{STRING_OBJ, INTEGER_OBJ, BOOLEAN_OBJ}
	result := StringifyObjectTypes(types)
	expected := "STRING|INTEGER|BOOLEAN"

	if result != expected {
		t.Errorf("expected '%s', got '%s'", expected, result)
	}
}
