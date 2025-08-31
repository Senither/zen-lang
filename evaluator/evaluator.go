package evaluator

import (
	"bytes"
	"fmt"
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
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

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
	case *ast.SuffixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		return evalSuffixExpression(node.Operator, left)
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
	case *ast.ChainExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		return evalChainExpression(left, node.Right, env)
	case *ast.AssignmentExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalAssignmentExpression(node.Left, right, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
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

		return evalCallExpression(node, function, env)
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
	registerBuiltins()

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

func evalWhileExpression(we *ast.WhileExpression, env *objects.Environment) objects.Object {
	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		body := Eval(we.Body, env)
		if isError(body) {
			return body
		}
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

func evalHashLiteral(node *ast.HashLiteral, env *objects.Environment) objects.Object {
	pairs := make(map[objects.HashKey]objects.HashPair)

	for key, value := range node.Pairs {
		keyObj := Eval(key, env)
		if isError(keyObj) {
			return keyObj
		}

		hashKey, ok := keyObj.(objects.Hashable)
		if !ok {
			return newError("key is not hashable: %s", keyObj.Type())
		}

		valueObj := Eval(value, env)
		if isError(valueObj) {
			return valueObj
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = objects.HashPair{Key: keyObj, Value: valueObj}
	}

	return &objects.Hash{Pairs: pairs}
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
		return &objects.Integer{Value: -right.Value}
	case *objects.Float:
		return &objects.Float{Value: -right.Value}
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

func evalSuffixExpression(operator string, left objects.Object) objects.Object {
	switch operator {
	case "++":
		return evalIncrementExpression(left)
	case "--":
		return evalDecrementExpression(left)
	default:
		return newError("unknown operator: %s%s", operator, left.Type())
	}
}

func evalIncrementExpression(left objects.Object) objects.Object {
	switch left := left.(type) {
	case *objects.Integer:
		left.Value++
		return left
	case *objects.Float:
		left.Value++
		return left
	default:
		return newError("unknown operator: %s%s", "++", left.Type())
	}
}

func evalDecrementExpression(left objects.Object) objects.Object {
	switch left := left.(type) {
	case *objects.Integer:
		left.Value--
		return left
	case *objects.Float:
		left.Value--
		return left
	default:
		return newError("unknown operator: %s%s", "--", left.Type())
	}
}

func evalIndexExpression(left, index objects.Object) objects.Object {
	switch {
	case left.Type() == objects.ARRAY_OBJ && index.Type() == objects.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == objects.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(left, index objects.Object) objects.Object {
	arrObj := left.(*objects.Array)
	idxObj := index.(*objects.Integer)

	length := int64(len(arrObj.Elements))

	var idx int64
	if idxObj.Value < 0 {
		idx = length + idxObj.Value
	} else {
		idx = idxObj.Value
	}

	if idx < 0 || idx >= length {
		return newError("array index out of bounds: %d", idxObj.Value)
	}

	return arrObj.Elements[idx]
}

func evalHashIndexExpression(left, index objects.Object) objects.Object {
	hashObj := left.(*objects.Hash)

	key, ok := index.(objects.Hashable)
	if !ok {
		return newError("invalid type given as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalChainExpression(left objects.Object, right ast.Expression, env *objects.Environment) objects.Object {
	switch left := left.(type) {
	case *objects.Hash:
		return evalHashChainExpression(left, right, env)
	default:
		return newError("invalid chain expression for %s", left.Type())
	}
}

func evalHashChainExpression(hash *objects.Hash, right ast.Expression, env *objects.Environment) objects.Object {
	switch right := right.(type) {
	case *ast.Identifier:
		pair, ok := hash.Pairs[(&objects.String{Value: right.Value}).HashKey()]
		if !ok {
			return newError("invalid chain expression for %s, key not found: %s", hash.Type(), right.Value)
		}

		return pair.Value
	case *ast.CallExpression:
		name, ok := right.Function.(*ast.Identifier)
		if !ok {
			return newError("invalid chain expression for %s, expected identifier, got %s", hash.Type(), right.Function.TokenLiteral())
		}

		pair, ok := hash.Pairs[(&objects.String{Value: name.Value}).HashKey()]
		if !ok {
			return newError("invalid chain expression for %s, key not found: %s", hash.Type(), name.Value)
		}

		return evalCallExpression(right, pair.Value, env)
	case *ast.ChainExpression:
		leftInner, ok := right.Left.(*ast.Identifier)
		if !ok {
			return newError("invalid chain expression for %s, expected identifier, got %s", hash.Type(), right.Left.TokenLiteral())
		}

		pair, ok := hash.Pairs[(&objects.String{Value: leftInner.Value}).HashKey()]
		if !ok {
			return newError("invalid chain expression for %s, key not found: %s", hash.Type(), leftInner.Value)
		}

		return evalChainExpression(pair.Value, right.Right, env)

	default:
		return newError("invalid chain expression for %s, got %s", hash.Type(), right.TokenLiteral())
	}
}

func evalAssignmentExpression(left ast.Expression, right objects.Object, env *objects.Environment) objects.Object {
	switch left := left.(type) {
	case *ast.Identifier:
		return env.Set(left.Value, right, false)
	case *ast.IndexExpression:
		leftObj := Eval(left.Left, env)
		if isError(leftObj) {
			return leftObj
		}

		switch leftObj := leftObj.(type) {
		case *objects.Array:
			return evalArrayAssignmentExpression(leftObj, left.Index, right, env)
		case *objects.Hash:
			return evalHashAssignmentExpression(leftObj, left.Index, right, env)
		default:
			return newError("left hand side of index assignment is not a valid indexable type: %s (%T)", leftObj, leftObj)
		}

	default:
		return newError("left hand side of assignment is not a valid expression: %s (%T)", left, left)
	}
}

func evalArrayAssignmentExpression(
	arr *objects.Array,
	index ast.Expression,
	value objects.Object,
	env *objects.Environment,
) objects.Object {
	idx := Eval(index, env)
	if isError(idx) {
		return idx
	}

	switch idx := idx.(type) {
	case *objects.Integer:
		if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
			return newError("array index out of bounds: %d", idx.Value)
		}

		arr.Elements[idx.Value] = value
	default:
		return newError("index operator not supported: %s", idx.Type())
	}

	return value
}

func evalHashAssignmentExpression(
	hash *objects.Hash,
	index ast.Expression,
	value objects.Object,
	env *objects.Environment,
) objects.Object {
	idx := Eval(index, env)
	if isError(idx) {
		return idx
	}

	key, ok := idx.(objects.Hashable)
	if !ok {
		return newError("invalid type given as hash key: %s", idx.Type())
	}

	hash.Pairs[key.HashKey()] = objects.HashPair{Key: idx, Value: value}

	return value
}

func evalNumberInfixExpression(operator string, left, right objects.Object) objects.Object {
	leftVal := unwrapNumberValue(left)
	rightVal := unwrapNumberValue(right)

	switch operator {
	case "+":
		return wrapNumberValue(leftVal+rightVal, left, right)
	case "-":
		return wrapNumberValue(leftVal-rightVal, left, right)
	case "*":
		return wrapNumberValue(leftVal*rightVal, left, right)
	case "/":
		return wrapNumberValue(leftVal/rightVal, left, right)
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
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

func evalCallExpression(node *ast.CallExpression, function objects.Object, env *objects.Environment) objects.Object {
	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	return applyFunction(function, args)
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

func wrapNumberValue(value float64, left, right objects.Object) objects.Object {
	if left.Type() == objects.FLOAT_OBJ || right.Type() == objects.FLOAT_OBJ {
		return &objects.Float{Value: value}
	}

	if float64(int64(value)) == value {
		return &objects.Integer{Value: int64(value)}
	}

	return &objects.Float{Value: value}
}

func unwrapNumberValue(obj objects.Object) float64 {
	switch n := obj.(type) {
	case *objects.Integer:
		return float64(n.Value)
	case *objects.Float:
		return n.Value
	default:
		return 0
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
