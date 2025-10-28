package evaluator

import (
	"github.com/senither/zen-lang/objects"
)

var globals = map[string]*objects.ImmutableHash{}

func registerGlobals() {
	globals = map[string]*objects.ImmutableHash{
		"arrays": objects.BuildImmutableHash(
			objects.WrapBuiltinFunctionInASTAwareMap("push", objects.GetGlobalBuiltinByName("arrays", "push")),
			objects.WrapBuiltinFunctionInASTAwareMap("shift", objects.GetGlobalBuiltinByName("arrays", "shift")),
			objects.WrapBuiltinFunctionInASTAwareMap("pop", objects.GetGlobalBuiltinByName("arrays", "pop")),
			objects.WrapBuiltinFunctionInASTAwareMap("filter", objects.GetGlobalBuiltinByName("arrays", "filter")),
		),

		"strings": objects.BuildImmutableHash(
			objects.WrapBuiltinFunctionInASTAwareMap("contains", objects.GetGlobalBuiltinByName("strings", "contains")),
			objects.WrapBuiltinFunctionInASTAwareMap("split", objects.GetGlobalBuiltinByName("strings", "split")),
			objects.WrapBuiltinFunctionInASTAwareMap("join", objects.GetGlobalBuiltinByName("strings", "join")),
			objects.WrapBuiltinFunctionInASTAwareMap("format", objects.GetGlobalBuiltinByName("strings", "format")),
		),
	}
}
