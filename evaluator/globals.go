package evaluator

import (
	"github.com/senither/zen-lang/objects"
)

var globals map[string]*objects.ImmutableHash

func registerGlobals() {
	globals = make(map[string]*objects.ImmutableHash, len(objects.Globals))

	for _, global := range objects.Globals {
		builtins := make([]objects.HashPair, len(global.Builtins))

		for i, fn := range global.Builtins {
			builtins[i] = objects.WrapBuiltinFunctionInASTAwareMap(fn.Name, fn.Builtin)
		}

		globals[global.Name] = objects.BuildImmutableHash(builtins...)
	}
}
