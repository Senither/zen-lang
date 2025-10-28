package objects

import (
	"fmt"
	"os"

	"github.com/senither/zen-lang/ast"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		Name: "print",
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			for _, arg := range args {
				fmt.Fprint(os.Stdout, arg.Inspect())
			}

			return NULL, nil
		}},
	},
	{
		Name: "println",
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			for _, arg := range args {
				fmt.Fprint(os.Stdout, arg.Inspect(), "\n")
			}

			return NULL, nil
		}},
	},
	{
		Name: "len",
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}, nil
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}, nil
			case *Null:
				return &Integer{Value: 0}, nil

			default:
				return nil, fmt.Errorf("argument to `len` not supported, got %s", args[0].Type())
			}
		}},
	},
	{
		Name: "string",
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
			}

			return &String{Value: args[0].Inspect()}, nil
		}},
	},
}

func BuiltinToASTAwareBuiltin(builtin *Builtin) *ASTAwareBuiltin {
	return &ASTAwareBuiltin{
		Fn: func(
			node *ast.CallExpression,
			env *Environment,
			args ...Object,
		) Object {
			result, err := builtin.Fn(args...)
			if err != nil {
				return NewError(node.Token, env, "%s", err.Error())
			}

			return result
		},
	}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}
