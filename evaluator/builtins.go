package evaluator

import (
	"github.com/senither/zen-lang/objects"
)

var builtins = map[string]*objects.ASTAwareBuiltin{}

func registerBuiltins() {
	builtins = map[string]*objects.ASTAwareBuiltin{
		"print":   objects.BuiltinToASTAwareBuiltin(objects.GetBuiltinByName("print")),
		"println": objects.BuiltinToASTAwareBuiltin(objects.GetBuiltinByName("println")),
		"len":     objects.BuiltinToASTAwareBuiltin(objects.GetBuiltinByName("len")),
		"string":  objects.BuiltinToASTAwareBuiltin(objects.GetBuiltinByName("string")),
	}
}
