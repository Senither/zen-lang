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
	return fmt.Errorf(
		"argument %d to `%s` has invalid type: got %s, want %s",
		index+1,
		name,
		args[index].Type(),
		StringifyObjectTypes(expected),
	)
}

func NewInvalidArgumentTypesErrorWithQualifiers(
	name string,
	expected []ObjectType,
	qualifier []ObjectType,
	index int,
	args []Object,
) error {
	return fmt.Errorf(
		"argument %d to `%s` has invalid type for %s: got %s, want %s",
		index+1,
		name,
		StringifyObjectTypes(qualifier),
		args[index].Type(),
		StringifyObjectTypes(expected),
	)
}

func StringifyObjectTypes(types []ObjectType) string {
	s := make([]string, len(types))

	for i, t := range types {
		s[i] = string(t)
	}

	return strings.Join(s, "|")
}
