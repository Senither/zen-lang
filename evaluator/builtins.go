package evaluator

import (
	"github.com/senither/zen-lang/objects"
)

var builtins map[string]*objects.ASTAwareBuiltin

func registerBuiltins() {
	builtins = make(map[string]*objects.ASTAwareBuiltin, len(objects.Builtins))

	for _, fn := range objects.Builtins {
		builtins[fn.Name] = objects.BuiltinToASTAwareBuiltin(fn.Builtin)
	}
}
