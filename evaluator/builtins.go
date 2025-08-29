package evaluator

import (
	"fmt"
	"os"

	"github.com/senither/zen-lang/objects"
)

var builtins = map[string]*objects.Builtin{}

func registerBuiltins() {
	builtins = map[string]*objects.Builtin{
		"print": {
			Fn: func(args ...objects.Object) objects.Object {
				for _, arg := range args {
					fmt.Fprint(os.Stdout, arg.Inspect())
				}

				return NULL
			},
		},
		"println": {
			Fn: func(args ...objects.Object) objects.Object {
				for _, arg := range args {
					fmt.Fprint(os.Stdout, arg.Inspect(), "\n")
				}

				return NULL
			},
		},
		"len": {
			Fn: func(args ...objects.Object) objects.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got %d, want 1", len(args))
				}

				switch arg := args[0].(type) {
				case *objects.String:
					return &objects.Integer{Value: int64(len(arg.Value))}
				case *objects.Array:
					return &objects.Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
		"array_push": {
			Fn: func(args ...objects.Object) objects.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got %d, want 2", len(args))
				}

				array, ok := args[0].(*objects.Array)
				if !ok {
					return newError("argument to `push` must be an array, got %s", args[0].Type())
				}

				array.Elements = append(array.Elements, args[1])

				return array
			},
		},
		"array_shift": {
			Fn: func(args ...objects.Object) objects.Object {
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
			},
		},
		"array_pop": {
			Fn: func(args ...objects.Object) objects.Object {
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
			},
		},
		"array_filter": {
			Fn: func(args ...objects.Object) objects.Object {
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
			},
		},
	}
}
