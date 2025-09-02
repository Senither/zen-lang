package objects

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

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

func IsError(obj Object) bool {
	if obj == nil {
		return false
	}

	return obj.Type() == ERROR_OBJ
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
		Value: &Builtin{Fn: fn},
	}
}
