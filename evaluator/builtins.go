package evaluator

import (
	"github.com/senither/zen-lang/objects"
)

var builtins = map[string]*objects.Builtin{
	"len": {
		Fn: func(args ...objects.Object) objects.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got %d, want 1", len(args))
			}

			switch arg := args[0].(type) {
			case *objects.String:
				return &objects.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
