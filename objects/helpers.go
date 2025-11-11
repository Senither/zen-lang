package objects

import (
	"fmt"
	"path/filepath"

	"github.com/senither/zen-lang/tokens"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}

	BREAK    = &Break{}
	CONTINUE = &Continue{}
)

type FileDescriptorContext struct {
	Name string
	Path string
}

func NewFileDescriptorContext(path string) *FileDescriptorContext {
	name := filepath.Base(path)

	return &FileDescriptorContext{
		Name: name,
		Path: path[:len(path)-len(name)-1],
	}
}

func NativeErrorToErrorObject(err error) *Error {
	return &Error{Message: err.Error()}
}

func NewError(token tokens.Token, fileCtx *FileDescriptorContext, format string, a ...interface{}) *Error {
	err := Error{
		Message: fmt.Sprintf(format, a...),
		Line:    token.Line,
		Column:  token.Column,
	}

	if fileCtx != nil {
		err.File = fileCtx.Name
		err.Path = fileCtx.Path
	}

	return &err
}

func NewEmptyErrorWithParent(parent *Error, token tokens.Token, fileCtx *FileDescriptorContext) *Error {
	err := NewError(token, fileCtx, "")
	err.Parent = parent

	return err
}

func IsError(obj Object) bool {
	return obj != nil && obj.Type() == ERROR_OBJ
}

func IsTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case *Boolean:
		return obj.Value
	case *Null:
		return false

	default:
		return true
	}
}

func IsNumber(t ObjectType) bool {
	return t == INTEGER_OBJ || t == FLOAT_OBJ
}

func WrapNumberValue(value float64, left, right Object) Object {
	if left.Type() == FLOAT_OBJ || right.Type() == FLOAT_OBJ {
		return &Float{Value: value}
	}

	if float64(int64(value)) == value {
		return &Integer{Value: int64(value)}
	}

	return &Float{Value: value}
}

func UnwrapNumberValue(obj Object) float64 {
	switch n := obj.(type) {
	case *Integer:
		return float64(n.Value)
	case *Float:
		return n.Value

	default:
		return 0
	}
}

func UnwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func NativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func CreateImmutableHashFromEnvExports(env *Environment) *ImmutableHash {
	hashPairs := []HashPair{}

	for key, val := range env.GetExports() {
		hashPairs = append(hashPairs, HashPair{
			Key:   &String{Value: key},
			Value: val,
		})
	}

	hash := BuildImmutableHash(hashPairs...)

	return hash
}

func BuildImmutableHash(args ...HashPair) *ImmutableHash {
	pairs := make(map[HashKey]HashPair)

	for _, arg := range args {
		hash := arg.Key.(Hashable)
		pairs[hash.HashKey()] = arg
	}

	return &ImmutableHash{Value: Hash{Pairs: pairs}}
}

func WrapBuiltinFunctionInASTAwareMap(
	name string,
	fn *Builtin,
) HashPair {
	return HashPair{
		Key:   &String{Value: name},
		Value: BuiltinToASTAwareBuiltin(fn),
	}
}
