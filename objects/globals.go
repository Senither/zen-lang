package objects

import (
	"fmt"
	"math"
	"strings"
)

var Globals = []struct {
	Name     string
	Builtins []*BuiltinDefinition
}{
	{
		Name: "arrays",
		Builtins: []*BuiltinDefinition{
			{
				Name: "push",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 2 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
					}

					array, ok := args[0].(*Array)
					if !ok {
						return nil, fmt.Errorf("argument to `push` must be an array, got %s", args[0].Type())
					}

					array.Elements = append(array.Elements, args[1])

					return array, nil
				}},
			},
			{
				Name: "shift",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 1 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
					}

					array, ok := args[0].(*Array)
					if !ok {
						return nil, fmt.Errorf("argument to `shift` must be an array, got %s", args[0].Type())
					}

					if len(array.Elements) == 0 {
						return NULL, nil
					}

					first := array.Elements[0]
					array.Elements = array.Elements[1:]

					return first, nil
				}},
			},
			{
				Name: "pop",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 1 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
					}

					array, ok := args[0].(*Array)
					if !ok {
						return nil, fmt.Errorf("argument to `pop` must be an array, got %s", args[0].Type())
					}

					if len(array.Elements) == 0 {
						return NULL, nil
					}

					last := array.Elements[len(array.Elements)-1]
					array.Elements = array.Elements[:len(array.Elements)-1]

					return last, nil
				}},
			},
			{
				Name: "filter",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 2 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
					}

					array, ok := args[0].(*Array)
					if !ok {
						return nil, fmt.Errorf("argument to `filter` must be an array, got %s", args[0].Type())
					}

					callable, ok := args[1].(Callable)
					if !ok {
						return nil, fmt.Errorf("second argument to `filter` must be a function, got %s", args[1].Type())
					}

					if callable.ParametersCount() != 1 {
						return nil, fmt.Errorf("function passed to `filter` must take exactly one argument")
					}

					filtered := make([]Object, 0)
					for _, elem := range array.Elements {
						rs := callable.Call(elem)

						switch rs := rs.(type) {
						case *Boolean:
							if rs == TRUE {
								filtered = append(filtered, elem)
							}
						case *Error:
							return rs, nil

						default:
							return nil, fmt.Errorf("function passed to `filter` must return a boolean, got %s", rs.Type())
						}
					}

					return &Array{Elements: filtered}, nil
				}},
			},
		},
	},

	{
		Name: "strings",
		Builtins: []*BuiltinDefinition{
			{
				Name: "contains",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 2 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
					}

					str, ok := args[0].(*String)
					if !ok {
						return nil, fmt.Errorf("argument to `contains` must be a string, got %s", args[0].Type())
					}

					substr, ok := args[1].(*String)
					if !ok {
						return nil, fmt.Errorf("second argument to `contains` must be a string, got %s", args[1].Type())
					}

					if strings.Contains(str.Value, substr.Value) {
						return TRUE, nil
					}

					return FALSE, nil
				}},
			},
			{
				Name: "split",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 2 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
					}

					str, ok := args[0].(*String)
					if !ok {
						return nil, fmt.Errorf("argument to `contains` must be a string, got %s", args[0].Type())
					}

					substr, ok := args[1].(*String)
					if !ok {
						return nil, fmt.Errorf("second argument to `contains` must be a string, got %s", args[1].Type())
					}

					arr := strings.Split(str.Value, substr.Value)

					var elements []Object
					for _, s := range arr {
						elements = append(elements, &String{Value: s})
					}

					return &Array{Elements: elements}, nil
				}},
			},
			{
				Name: "join",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) != 2 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
					}

					arr, ok := args[0].(*Array)
					if !ok {
						return nil, fmt.Errorf("argument to `join` must be an array, got %s", args[0].Type())
					}

					sep, ok := args[1].(*String)
					if !ok {
						return nil, fmt.Errorf("second argument to `join` must be a string, got %s", args[1].Type())
					}

					var elements []string
					for _, elem := range arr.Elements {
						switch elem := elem.(type) {
						case *Float:
							elements = append(elements, fmt.Sprintf("%v", elem.Value))

						default:
							elements = append(elements, elem.Inspect())
						}
					}

					return &String{Value: strings.Join(elements, sep.Value)}, nil
				}},
			},
			{
				Name: "format",
				Builtin: &Builtin{Fn: func(args ...Object) (Object, error) {
					if len(args) < 2 {
						return nil, fmt.Errorf("wrong number of arguments. got %d, want at least 2", len(args))
					}

					str, ok := args[0].(*String)
					if !ok {
						return nil, fmt.Errorf("argument to `contains` must be a string, got %s", args[0].Type())
					}

					var values []any
					for _, arg := range args[1:] {
						switch arg := arg.(type) {
						case *String:
							values = append(values, arg.Value)
						case *Integer:
							values = append(values, arg.Value)
						case *Float:
							values = append(values, arg.Value)
						case *Boolean:
							values = append(values, arg.Value)
						case *Null:
							values = append(values, nil)

						default:
							values = append(values, arg.Inspect())
						}
					}

					return &String{Value: fmt.Sprintf(str.Value, values...)}, nil
				}},
			},
		},
	},

	{
		Name: "math",
		Builtins: []*BuiltinDefinition{
			{
				Name: "min",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 2 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument 1 to `min` must be a number, got %s", args[0].Type())
						}

						if !IsNumber(args[1].Type()) {
							return nil, fmt.Errorf("argument 2 to `min` must be a number, got %s", args[1].Type())
						}

						return WrapNumberValue(math.Min(
							UnwrapNumberValue(args[0]),
							UnwrapNumberValue(args[1]),
						), args[0], args[1]), nil
					},
				},
			},
			{
				Name: "max",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 2 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 2", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument 1 to `max` must be a number, got %s", args[0].Type())
						}

						if !IsNumber(args[1].Type()) {
							return nil, fmt.Errorf("argument 2 to `max` must be a number, got %s", args[1].Type())
						}

						return WrapNumberValue(math.Max(
							UnwrapNumberValue(args[0]),
							UnwrapNumberValue(args[1]),
						), args[0], args[1]), nil
					},
				},
			},
			{
				Name: "ceil",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 1 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument to `ceil` must be a number, got %s", args[0].Type())
						}

						return &Float{Value: math.Ceil(UnwrapNumberValue(args[0]))}, nil
					},
				},
			},
			{
				Name: "floor",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 1 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument to `floor` must be a number, got %s", args[0].Type())
						}

						return &Float{Value: math.Floor(UnwrapNumberValue(args[0]))}, nil
					},
				},
			},
			{
				Name: "round",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 1 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument to `round` must be a number, got %s", args[0].Type())
						}

						return &Float{Value: math.Round(UnwrapNumberValue(args[0]))}, nil
					},
				},
			},
			{
				Name: "log",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 1 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument to `log` must be a number, got %s", args[0].Type())
						}

						return &Float{Value: math.Log10(UnwrapNumberValue(args[0]))}, nil
					},
				},
			},
			{
				Name: "sqrt",
				Builtin: &Builtin{
					Fn: func(args ...Object) (Object, error) {
						if len(args) != 1 {
							return nil, fmt.Errorf("wrong number of arguments. got %d, want 1", len(args))
						}

						if !IsNumber(args[0].Type()) {
							return nil, fmt.Errorf("argument to `sqrt` must be a number, got %s", args[0].Type())
						}

						return &Float{Value: math.Sqrt(UnwrapNumberValue(args[0]))}, nil
					},
				},
			},
		},
	},
}

func GetGlobalBuiltinByName(scope, name string) *Builtin {
	for _, grp := range Globals {
		if grp.Name == scope {
			for _, def := range grp.Builtins {
				if def.Name == name {
					return def.Builtin
				}
			}
		}
	}

	return nil
}
