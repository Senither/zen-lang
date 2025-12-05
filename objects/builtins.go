package objects

import (
	"fmt"
	"os"

	"github.com/senither/zen-lang/ast"
)

var Builtins = []BuiltinDefinition{
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
		Schema: BuiltinSchema{
			NewRequiredArgument(STRING_OBJ, ARRAY_OBJ, NULL_OBJ),
		},
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, NewWrongNumberOfArgumentsError("len", 1, len(args))
			}

			switch arg := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(arg.Value))}, nil
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}, nil
			case *Null:
				return &Integer{Value: 0}, nil

			default:
				return nil, NewInvalidArgumentTypesError("len", []ObjectType{STRING_OBJ, ARRAY_OBJ, NULL_OBJ}, 0, args)
			}
		}},
	},
	{
		Name:   "string",
		Schema: BuiltinSchema{NewRequiredArgument()},
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, NewWrongNumberOfArgumentsError("string", 1, len(args))
			}

			return &String{Value: args[0].Inspect()}, nil
		}},
	},
	{
		Name: "int",
		Schema: BuiltinSchema{
			NewRequiredArgument(INTEGER_OBJ, FLOAT_OBJ, STRING_OBJ, BOOLEAN_OBJ, NULL_OBJ),
		},
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, NewWrongNumberOfArgumentsError("int", 1, len(args))
			}

			switch arg := args[0].(type) {
			case *Integer:
				return arg, nil
			case *Float:
				return &Integer{Value: int64(arg.Value)}, nil
			case *Boolean:
				if arg.Value {
					return &Integer{Value: 1}, nil
				} else {
					return &Integer{Value: 0}, nil
				}
			case *Null:
				return &Integer{Value: 0}, nil
			case *String:
				var intValue int64
				_, err := fmt.Sscan(arg.Value, &intValue)
				if err != nil {
					return nil, NewErrorf("int", "failed to convert `%s` to %s", arg.Value, INTEGER_OBJ)
				}

				return &Integer{Value: intValue}, nil

			default:
				return nil, NewInvalidArgumentTypesError("int", []ObjectType{
					INTEGER_OBJ, FLOAT_OBJ, STRING_OBJ, BOOLEAN_OBJ, NULL_OBJ,
				}, 0, args)
			}
		}},
	},
	{
		Name: "float",
		Schema: BuiltinSchema{
			NewRequiredArgument(INTEGER_OBJ, FLOAT_OBJ, STRING_OBJ, BOOLEAN_OBJ, NULL_OBJ),
		},
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, NewWrongNumberOfArgumentsError("float", 1, len(args))
			}

			switch arg := args[0].(type) {
			case *Float:
				return arg, nil
			case *Integer:
				return &Float{Value: float64(arg.Value)}, nil
			case *Boolean:
				if arg.Value {
					return &Float{Value: 1}, nil
				} else {
					return &Float{Value: 0}, nil
				}
			case *Null:
				return &Float{Value: 0}, nil
			case *String:
				var intValue float64
				_, err := fmt.Sscan(arg.Value, &intValue)
				if err != nil {
					return nil, NewErrorf("float", "failed to convert `%s` to %s", arg.Value, FLOAT_OBJ)
				}

				return &Float{Value: float64(intValue)}, nil

			default:
				return nil, NewInvalidArgumentTypesError("float", []ObjectType{
					INTEGER_OBJ, FLOAT_OBJ, STRING_OBJ, BOOLEAN_OBJ, NULL_OBJ,
				}, 0, args)
			}
		}},
	},
	{
		Name:   "type",
		Schema: BuiltinSchema{NewRequiredArgument()},
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, NewWrongNumberOfArgumentsError("type", 1, len(args))
			}

			switch args[0].Type() {
			case FUNCTION_OBJ, BUILTIN_OBJ, COMPILED_FUNCTION_OBJ, CLOSURE_OBJ:
				return &String{Value: "FUNCTION"}, nil

			default:
				_, ok := args[0].(Callable)
				if ok {
					return &String{Value: "FUNCTION"}, nil
				}

				return &String{Value: string(args[0].Type())}, nil
			}
		}},
	},
	{
		Name:   "isNaN",
		Schema: BuiltinSchema{NewRequiredArgument()},
		Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, NewWrongNumberOfArgumentsError("isNaN", 1, len(args))
			}

			val, ok := args[0].(*Float)
			if !ok {
				return FALSE, nil
			}

			if val.Value != val.Value {
				return TRUE, nil
			}

			return FALSE, nil
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
				return NewError(node.Token, env.GetFileDescriptorContext(), "%s", err.Error())
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
