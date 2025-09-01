package evaluator

import (
	"fmt"
	"strings"

	"github.com/senither/zen-lang/objects"
)

var globals = map[string]*objects.ImmutableHash{}

func registerGlobals() {
	globals = map[string]*objects.ImmutableHash{
		"arrays": objects.BuildImmutableHash(
			objects.WrapBuiltinFunctionInMap("push", func(args ...objects.Object) objects.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d, want 2", len(args))
				}

				array, ok := args[0].(*objects.Array)
				if !ok {
					return newError("argument to `push` must be an array, got %s", args[0].Type())
				}

				array.Elements = append(array.Elements, args[1])

				return array
			}),
			objects.WrapBuiltinFunctionInMap("shift", func(args ...objects.Object) objects.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				array, ok := args[0].(*objects.Array)
				if !ok {
					return newError("argument to `shift` must be an array, got %s", args[0].Type())
				}

				if len(array.Elements) == 0 {
					return NULL
				}

				first := array.Elements[0]
				array.Elements = array.Elements[1:]

				return first
			}),
			objects.WrapBuiltinFunctionInMap("pop", func(args ...objects.Object) objects.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				array, ok := args[0].(*objects.Array)
				if !ok {
					return newError("argument to `pop` must be an array, got %s", args[0].Type())
				}

				if len(array.Elements) == 0 {
					return NULL
				}

				last := array.Elements[len(array.Elements)-1]
				array.Elements = array.Elements[:len(array.Elements)-1]

				return last
			}),
			objects.WrapBuiltinFunctionInMap("filter", func(args ...objects.Object) objects.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d, want 2", len(args))
				}

				array, ok := args[0].(*objects.Array)
				if !ok {
					return newError("argument to `filter` must be an array, got %s", args[0].Type())
				}

				callback, ok := args[1].(*objects.Function)
				if !ok {
					return newError("second argument to `filter` must be a function, got %s", args[1].Type())
				}

				if callback.Parameters == nil || len(callback.Parameters) != 1 {
					return newError("function passed to `filter` must take exactly one argument")
				}

				filtered := make([]objects.Object, 0)
				for _, elem := range array.Elements {
					env := objects.NewEnclosedEnvironment(callback.Env)
					env.SetImmutableForcefully(callback.Parameters[0].Value, elem)

					rs := unwrapReturnValue(Eval(callback.Body, env))

					switch rs := rs.(type) {
					case *objects.Boolean:
						if rs == TRUE {
							filtered = append(filtered, elem)
						}
					default:
						return newError("function passed to `filter` must return a boolean, got %s", rs.Type())
					}
				}

				return &objects.Array{Elements: filtered}
			}),
		),

		"strings": objects.BuildImmutableHash(
			objects.WrapBuiltinFunctionInMap("contains", func(args ...objects.Object) objects.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d, want 2", len(args))
				}

				str, ok := args[0].(*objects.String)
				if !ok {
					return newError("argument to `contains` must be a string, got %s", args[0].Type())
				}

				substr, ok := args[1].(*objects.String)
				if !ok {
					return newError("second argument to `contains` must be a string, got %s", args[1].Type())
				}

				if strings.Contains(str.Value, substr.Value) {
					return TRUE
				}

				return FALSE
			}),
			objects.WrapBuiltinFunctionInMap("split", func(args ...objects.Object) objects.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d, want 2", len(args))
				}

				str, ok := args[0].(*objects.String)
				if !ok {
					return newError("argument to `contains` must be a string, got %s", args[0].Type())
				}

				substr, ok := args[1].(*objects.String)
				if !ok {
					return newError("second argument to `contains` must be a string, got %s", args[1].Type())
				}

				arr := strings.Split(str.Value, substr.Value)

				var elements []objects.Object
				for _, s := range arr {
					elements = append(elements, &objects.String{Value: s})
				}

				return &objects.Array{Elements: elements}
			}),
			objects.WrapBuiltinFunctionInMap("join", func(args ...objects.Object) objects.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d, want 2", len(args))
				}

				arr, ok := args[0].(*objects.Array)
				if !ok {
					return newError("argument to `join` must be an array, got %s", args[0].Type())
				}

				sep, ok := args[1].(*objects.String)
				if !ok {
					return newError("second argument to `join` must be a string, got %s", args[1].Type())
				}

				var elements []string
				for _, elem := range arr.Elements {
					switch elem := elem.(type) {
					case *objects.Float:
						elements = append(elements, fmt.Sprintf("%v", elem.Value))
					default:
						elements = append(elements, elem.Inspect())
					}
				}

				return &objects.String{Value: strings.Join(elements, sep.Value)}
			}),
			objects.WrapBuiltinFunctionInMap("format", func(args ...objects.Object) objects.Object {
				if len(args) < 2 {
					return newError("wrong number of arguments. got %d, want at least 2", len(args))
				}

				str, ok := args[0].(*objects.String)
				if !ok {
					return newError("argument to `contains` must be a string, got %s", args[0].Type())
				}

				var values []any
				for _, arg := range args[1:] {
					switch arg := arg.(type) {
					case *objects.String:
						values = append(values, arg.Value)
					case *objects.Integer:
						values = append(values, arg.Value)
					case *objects.Float:
						values = append(values, arg.Value)
					case *objects.Boolean:
						values = append(values, arg.Value)
					default:
						values = append(values, arg.Inspect())
					}
				}

				return &objects.String{Value: fmt.Sprintf(str.Value, values...)}
			}),
		),
	}
}
