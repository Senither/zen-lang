package objects

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}

	BREAK    = &Break{}
	CONTINUE = &Continue{}
)

func NativeErrorToErrorObject(err error) *Error {
	return &Error{Message: err.Error()}
}

func NewError(token tokens.Token, env *Environment, format string, a ...interface{}) *Error {
	err := Error{
		Message: fmt.Sprintf(format, a...),
		Line:    token.Line,
		Column:  token.Column,
	}

	if env.GetFile() != nil {
		err.File = env.GetFile().Name
		err.Path = env.GetFile().Path
	}

	return &err
}

func NewEmptyErrorWithParent(parent *Error, token tokens.Token, env *Environment) *Error {
	err := NewError(token, env, "")
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

func WrapBuiltinFunctionInMap(
	name string,
	fn func(node *ast.CallExpression, env *Environment, args ...Object) Object,
) HashPair {
	return HashPair{
		Key:   &String{Value: name},
		Value: &ASTAwareBuiltin{Fn: fn},
	}
}
