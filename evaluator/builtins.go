package evaluator

import (
	"fmt"
	"math/big"
	"os"

	"github.com/senither/zen-lang/objects"
)

var builtins = map[string]*objects.Builtin{
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
				return &objects.Integer{Value: new(big.Int).SetInt64(int64(len(arg.Value)))}
			case *objects.Array:
				return &objects.Integer{Value: new(big.Int).SetInt64(int64(len(arg.Elements)))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
