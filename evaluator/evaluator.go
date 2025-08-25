package evaluator

import (
	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/objects"
)

var (
	NULL  = &objects.Null{}
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
)

func Eval(node ast.Node) objects.Object {
	switch n := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.(*ast.Program).Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)

	// Expressions
	case *ast.StringLiteral:
		return &objects.String{Value: n.Value}
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)
	}

	return nil
}

func evalStatements(statements []ast.Statement) objects.Object {
	var result objects.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *objects.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
