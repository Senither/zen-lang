package evaluator

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/objects"
)

var (
	NULL  = &objects.Null{}
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
)

func Eval(node ast.Node) objects.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}

		return &objects.ReturnValue{Value: val}

	// Expression types
	case *ast.StringLiteral:
		return &objects.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	// Expression operators
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if isError(left) {
			return left
		}

		right := Eval(node.Right)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node)
	}

	return nil
}

func newError(format string, a ...interface{}) *objects.Error {
	return &objects.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj objects.Object) bool {
	if obj == nil {
		return false
	}

	return obj.Type() == objects.ERROR_OBJ
}

func evalProgram(statements []ast.Statement) objects.Object {
	var result objects.Object

	for _, stmt := range statements {
		result = Eval(stmt)

		switch result := result.(type) {
		case *objects.ReturnValue:
			return result.Value
		case *objects.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement) objects.Object {
	var result objects.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil {
			rt := result.Type()
			if rt == objects.RETURN_VALUE_OBJ || rt == objects.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression) objects.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	}

	if ie.Intermediary != nil {
		return Eval(ie.Intermediary)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative)
	}

	return NULL
}

func nativeBoolToBooleanObject(input bool) *objects.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj objects.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

func evalPrefixExpression(operator string, right objects.Object) objects.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: -%s", right.Type())
	}

	return &objects.Integer{Value: -right.(*objects.Integer).Value}
}

func evalInfixExpression(operator string, left, right objects.Object) objects.Object {
	switch {
	case left.Type() == objects.INTEGER_OBJ && right.Type() == objects.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)

	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right objects.Object) objects.Object {
	leftVal := left.(*objects.Integer).Value
	rightVal := right.(*objects.Integer).Value

	switch operator {
	case "+":
		return &objects.Integer{Value: leftVal + rightVal}
	case "-":
		return &objects.Integer{Value: leftVal - rightVal}
	case "*":
		return &objects.Integer{Value: leftVal * rightVal}
	case "/":
		return &objects.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
