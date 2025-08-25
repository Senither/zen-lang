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
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)

	// Expression types
	case *ast.StringLiteral:
		return &objects.String{Value: n.Value}
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	// Expression operators
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return evalPrefixExpression(n.Operator, right)
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

func evalPrefixExpression(operator string, right objects.Object) objects.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(right objects.Object) objects.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right objects.Object) objects.Object {
	if right.Type() != objects.INTEGER_OBJ {
		return NULL
	}

	return &objects.Integer{Value: -right.(*objects.Integer).Value}
}
