package objects

import (
	"fmt"
	"strings"
)

func NewErrorf(name string, format string, a ...any) error {
	return fmt.Errorf("error in `%s`: %s", name, fmt.Sprintf(format, a...))
}

func NewWrongNumberOfArgumentsError(name string, expected int, got int) error {
	return fmt.Errorf("wrong number of arguments to `%s`: got %d, want %d", name, got, expected)
}

func NewWrongNumberOfArgumentsWantAtLeastError(name string, expected int, got int) error {
	return fmt.Errorf("wrong number of arguments to `%s`: got %d, want at least %d", name, got, expected)
}

func NewInvalidArgumentTypeError(name string, expected ObjectType, index int, args []Object) error {
	return fmt.Errorf(
		"argument %d to `%s` has invalid type: got %s, want %s",
		index+1, name, args[index].Type(), expected,
	)
}

func NewInvalidArgumentTypesError(name string, expected []ObjectType, index int, args []Object) error {
	expectedStr := make([]string, len(expected))
	for i, e := range expected {
		expectedStr[i] = string(e)
	}

	return fmt.Errorf(
		"argument %d to `%s` has invalid type: got %s, want %s",
		index+1, name, args[index].Type(), strings.Join(expectedStr, "|"),
	)
}
