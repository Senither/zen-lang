package evaluator

import (
	"bytes"
	"fmt"
	"math/big"
	"os"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/objects"
)

var (
	NULL  = &objects.Null{}
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
)

func Eval(node ast.Node, env *objects.Environment) objects.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}

		return &objects.ReturnValue{Value: val}
	case *ast.VariableStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}

		return env.Set(node.Name.Value, val, node.Mutable)

	// Expression types
	case *ast.StringLiteral:
		return &objects.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &objects.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &objects.Array{Elements: elements}

	// Expression operators
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.AssignmentExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalAssignmentExpression(node.Left, right, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)

	// Functions & Builtins
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body

		function := &objects.Function{
			Name:       node.Name,
			Parameters: params,
			Env:        env,
			Body:       body,
		}

		if function.Name != nil {
			rs := env.Set(function.Name.Value, function, false)
			if isError(rs) {
				return rs
			}
		}

		return function
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
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

func evalProgram(statements []ast.Statement, env *objects.Environment) objects.Object {
	var result objects.Object

	for _, stmt := range statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *objects.ReturnValue:
			return result.Value
		case *objects.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *objects.Environment) objects.Object {
	var result objects.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == objects.RETURN_VALUE_OBJ || rt == objects.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIfExpression(ie *ast.IfExpression, env *objects.Environment) objects.Object {
	condition := Eval(ie.Condition, env)

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Intermediary != nil {
		return Eval(ie.Intermediary, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func evalIdentifier(node *ast.Identifier, env *objects.Environment) objects.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("%s: %s", "identifier not found", node.Value)
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

func isNumber(obj objects.ObjectType) bool {
	switch obj {
	case objects.INTEGER_OBJ, objects.FLOAT_OBJ:
		return true
	default:
		return false
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
	switch right := right.(type) {
	case *objects.Integer:
		return &objects.Integer{Value: right.Value.Neg(right.Value)}
	case *objects.Float:
		return &objects.Float{Value: right.Value.Neg(right.Value)}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right objects.Object) objects.Object {
	switch {
	case isNumber(left.Type()) && isNumber(right.Type()):
		return evalNumberInfixExpression(operator, left, right)
	case left.Type() == objects.STRING_OBJ && right.Type() == objects.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
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

func evalIndexExpression(left, index objects.Object) objects.Object {
	switch {
	case left.Type() == objects.ARRAY_OBJ && index.Type() == objects.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left, index objects.Object) objects.Object {
	arrObj := left.(*objects.Array)
	idxObj := index.(*objects.Integer)

	length := int64(len(arrObj.Elements))

	var idx int64
	if idxObj.Value.Cmp(big.NewInt(0)) < 0 {
		idx = length + idxObj.Value.Int64()
	} else {
		idx = idxObj.Value.Int64()
	}

	if idx < 0 || idx >= length {
		return newError("array index out of bounds: %d", idxObj.Value)
	}

	return arrObj.Elements[idx]
}

func evalAssignmentExpression(left ast.Expression, right objects.Object, env *objects.Environment) objects.Object {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return newError("left hand side of assignment is not an identifier: %s", left)
	}

	return env.Set(ident.Value, right, false)
}

func evalNumberInfixExpression(operator string, left, right objects.Object) objects.Object {
	leftVal := unwrapNumberValue(left)
	rightVal := unwrapNumberValue(right)

	switch operator {
	case "+":
		return wrapNumberValue(leftVal.Add(leftVal, rightVal), left, right)
	case "-":
		return wrapNumberValue(leftVal.Sub(leftVal, rightVal), left, right)
	case "*":
		return wrapNumberValue(leftVal.Mul(leftVal, rightVal), left, right)
	case "/":
		return wrapNumberValue(leftVal.Quo(leftVal, rightVal), left, right)
	case "<":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) < 0)
	case ">":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) > 0)
	case "==":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) == 0)
	case "!=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) != 0)
	case "<=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) <= 0)
	case ">=":
		return nativeBoolToBooleanObject(leftVal.Cmp(rightVal) >= 0)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right objects.Object) objects.Object {
	leftVal := left.(*objects.String).Value
	rightVal := right.(*objects.String).Value

	switch operator {
	case "+":
		return &objects.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalExpressions(exps []ast.Expression, env *objects.Environment) []objects.Object {
	var result []objects.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []objects.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn objects.Object, args []objects.Object) objects.Object {
	switch fn := fn.(type) {
	case *objects.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *objects.Builtin:
		return captureStdoutForBuiltin(fn, args)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *objects.Function, args []objects.Object) *objects.Environment {
	env := objects.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx], false)
	}

	return env
}

func unwrapReturnValue(obj objects.Object) objects.Object {
	if returnValue, ok := obj.(*objects.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func wrapNumberValue(value *big.Float, left, right objects.Object) objects.Object {
	if left.Type() == objects.FLOAT_OBJ || right.Type() == objects.FLOAT_OBJ {
		return &objects.Float{Value: value}
	}

	if value.IsInt() {
		valueInt, _ := value.Int(nil)
		return &objects.Integer{Value: valueInt}
	}

	return &objects.Float{Value: value}
}

func unwrapNumberValue(obj objects.Object) *big.Float {
	switch n := obj.(type) {
	case *objects.Integer:
		return objects.NewFloatFromBigInt(n.Value).Value
	case *objects.Float:
		return n.Value
	default:
		return objects.NewFloatFromInt64(int64(0)).Value
	}
}

func captureStdoutForBuiltin(fn *objects.Builtin, args []objects.Object) objects.Object {
	var buf bytes.Buffer

	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rs := fn.Fn(args...)

	w.Close()
	buf.ReadFrom(r)
	os.Stdout = originalStdout

	output := buf.String()
	if output != "" && output != "\n" {
		Stdout.Write(output)
	}

	return rs
}
