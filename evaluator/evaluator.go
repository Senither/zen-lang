package evaluator

import (
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
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
		if objects.IsError(val) {
			return val
		}

		return &objects.ReturnValue{Value: val}
	case *ast.VariableStatement:
		val := Eval(node.Value, env)
		if objects.IsError(val) {
			return val
		}

		return env.Set(node, node.Name.Value, val, node.Mutable)

	// Loop controls
	case *ast.BreakStatement:
		return &objects.ReturnValue{Value: objects.BREAK}
	case *ast.ContinueStatement:
		return &objects.ReturnValue{Value: objects.CONTINUE}

	// Expression types
	case *ast.NullLiteral:
		return objects.NULL
	case *ast.StringLiteral:
		return &objects.String{Value: node.Value}
	case *ast.IntegerLiteral:
		return &objects.Integer{Value: node.Value}
	case *ast.FloatLiteral:
		return &objects.Float{Value: node.Value}
	case *ast.BooleanLiteral:
		return objects.NativeBoolToBooleanObject(node.Value)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && objects.IsError(elements[0]) {
			return elements[0]
		}

		return &objects.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	// Expression operators
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if objects.IsError(right) {
			return right
		}

		return evalPrefixExpression(node, right, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if objects.IsError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if objects.IsError(right) {
			return right
		}

		return evalInfixExpression(node, left, right, env)
	case *ast.SuffixExpression:
		left := Eval(node.Left, env)
		if objects.IsError(left) {
			return left
		}

		return evalSuffixExpression(node, left, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if objects.IsError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if objects.IsError(index) {
			return index
		}

		return evalIndexExpression(node, left, index, env)
	case *ast.ChainExpression:
		left := Eval(node.Left, env)
		if objects.IsError(left) {
			return left
		}

		return evalChainExpression(node, left, node.Right, env)
	case *ast.AssignmentExpression:
		right := Eval(node.Right, env)
		if objects.IsError(right) {
			return right
		}

		return evalAssignmentExpression(node, node.Left, right, env)
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
			rs := env.Set(node, function.Name.Value, function, false)
			if objects.IsError(rs) {
				return rs
			}
		}

		return function
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if objects.IsError(function) {
			return function
		}

		return evalCallExpression(node, function, env)

	// Import & Export statements
	case *ast.ImportStatement:
		return evalImportStatement(node, env)
	case *ast.ExportStatement:
		return evalExportStatement(node, env)
	}

	return nil
}

func evalProgram(statements []ast.Statement, env *objects.Environment) objects.Object {
	registerGlobals()
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

	if objects.IsTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Intermediary != nil {
		return Eval(ie.Intermediary, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return objects.NULL
}

func evalWhileExpression(we *ast.WhileExpression, env *objects.Environment) objects.Object {
	for {
		condition := Eval(we.Condition, env)
		if objects.IsError(condition) {
			return condition
		}

		if !objects.IsTruthy(condition) {
			break
		}

		body := objects.UnwrapReturnValue(Eval(we.Body, env))
		if objects.IsError(body) {
			return body
		}

		if body == objects.BREAK {
			break
		}
	}

	return objects.NULL
}

func evalIdentifier(node *ast.Identifier, env *objects.Environment) objects.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	if global, ok := globals[node.Value]; ok {
		return global
	}

	return objects.NewError(
		node.Token, env.GetFileDescriptorContext(),
		"%s: %s",
		"identifier not found", node.Value,
	)
}

func evalHashLiteral(node *ast.HashLiteral, env *objects.Environment) objects.Object {
	pairs := make(map[objects.HashKey]objects.HashPair)

	for key, value := range node.Pairs {
		keyObj := Eval(key, env)
		if objects.IsError(keyObj) {
			return keyObj
		}

		hashKey, ok := keyObj.(objects.Hashable)
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"key is not hashable: %s",
				keyObj.Type(),
			)
		}

		valueObj := Eval(value, env)
		if objects.IsError(valueObj) {
			return valueObj
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = objects.HashPair{Key: keyObj, Value: valueObj}
	}

	return &objects.Hash{Pairs: pairs}
}

func evalPrefixExpression(
	node *ast.PrefixExpression,
	right objects.Object,
	env *objects.Environment,
) objects.Object {
	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right, env)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s%s",
			node.Operator, right.Type(),
		)
	}
}

func evalBangOperatorExpression(right objects.Object) objects.Object {
	switch right {
	case objects.TRUE:
		return objects.FALSE
	case objects.FALSE:
		return objects.TRUE
	case objects.NULL:
		return objects.TRUE

	default:
		return objects.FALSE
	}
}

func evalMinusPrefixOperatorExpression(
	node *ast.PrefixExpression,
	right objects.Object,
	env *objects.Environment,
) objects.Object {
	switch right := right.(type) {
	case *objects.Integer:
		return &objects.Integer{Value: -right.Value}
	case *objects.Float:
		return &objects.Float{Value: -right.Value}

	default:
		return objects.NewError(node.Token, env.GetFileDescriptorContext(), "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(node *ast.InfixExpression, left, right objects.Object, env *objects.Environment) objects.Object {
	switch {
	case objects.IsNumber(left.Type()) && objects.IsNumber(right.Type()):
		return evalNumberInfixExpression(node, left, right, env)
	case left.Type() == objects.STRING_OBJ && right.Type() == objects.STRING_OBJ:
		return evalStringInfixExpression(node, left, right, env)
	case node.Operator == "==":
		return objects.NativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return objects.NativeBoolToBooleanObject(left != right)
	case left.Type() == objects.STRING_OBJ && objects.IsStringable(right), right.Type() == objects.STRING_OBJ && objects.IsStringable(left):
		return evalStringableInfixExpression(node, left, right, env)

	case left.Type() != right.Type():
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"type mismatch: %s %s %s",
			left.Type(), node.Operator, right.Type(),
		)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s %s %s",
			left.Type(), node.Operator, right.Type(),
		)
	}
}

func evalSuffixExpression(node *ast.SuffixExpression, left objects.Object, env *objects.Environment) objects.Object {
	switch node.Operator {
	case "++":
		return evalIncrementExpression(node, left, env)
	case "--":
		return evalDecrementExpression(node, left, env)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s%s",
			node.Operator, left.Type(),
		)
	}
}

func evalIncrementExpression(node *ast.SuffixExpression, left objects.Object, env *objects.Environment) objects.Object {
	switch left := left.(type) {
	case *objects.Integer:
		left.Value++
		return left
	case *objects.Float:
		left.Value++
		return left

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s%s",
			node.Operator, left.Type(),
		)
	}
}

func evalDecrementExpression(node *ast.SuffixExpression, left objects.Object, env *objects.Environment) objects.Object {
	switch left := left.(type) {
	case *objects.Integer:
		left.Value--
		return left
	case *objects.Float:
		left.Value--
		return left

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s%s",
			node.Operator, left.Type(),
		)
	}
}

func evalIndexExpression(
	node *ast.IndexExpression,
	left, index objects.Object,
	env *objects.Environment,
) objects.Object {
	switch {
	case left.Type() == objects.ARRAY_OBJ && index.Type() == objects.INTEGER_OBJ:
		return evalArrayIndexExpression(node, left, index, env)
	case left.Type() == objects.HASH_OBJ:
		return evalHashIndexExpression(node, left, index, env)
	case left.Type() == objects.IMMUTABLE_HASH_OBJ:
		iHash := left.(*objects.ImmutableHash)
		return evalHashIndexExpression(node, &iHash.Value, index, env)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"index operator not supported: %s",
			left.Type(),
		)
	}
}

func evalArrayIndexExpression(
	node *ast.IndexExpression,
	left, index objects.Object,
	env *objects.Environment,
) objects.Object {
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
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"array index out of bounds: %d",
			idxObj.Value,
		)
	}

	return arrObj.Elements[idx]
}

func evalHashIndexExpression(
	node *ast.IndexExpression,
	left, index objects.Object,
	env *objects.Environment,
) objects.Object {
	hashObj := left.(*objects.Hash)

	key, ok := index.(objects.Hashable)
	if !ok {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"invalid type given as hash key: %s",
			index.Type(),
		)
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return objects.NULL
	}

	return pair.Value
}

func evalChainExpression(
	node *ast.ChainExpression,
	left objects.Object,
	right ast.Expression,
	env *objects.Environment,
) objects.Object {
	switch left := left.(type) {
	case *objects.Hash:
		return evalHashChainExpression(node, left, right, env)
	case *objects.ImmutableHash:
		switch right := right.(type) {
		case *ast.AssignmentExpression:
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"cannot assign to immutable hash keys",
			)

		default:
			return evalHashChainExpression(node, &left.Value, right, env)
		}

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"invalid chain expression for %s",
			left.Type(),
		)
	}
}

func evalHashChainExpression(
	node *ast.ChainExpression,
	hash *objects.Hash,
	right ast.Expression,
	env *objects.Environment,
) objects.Object {
	switch right := right.(type) {
	case *ast.Identifier:
		pair, ok := hash.Pairs[(&objects.String{Value: right.Value}).HashKey()]
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, key not found: %s",
				hash.Type(), right.Value,
			)
		}

		return pair.Value
	case *ast.CallExpression:
		name, ok := right.Function.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, expected identifier, got %s",
				hash.Type(), right.Function.TokenLiteral(),
			)
		}

		pair, ok := hash.Pairs[(&objects.String{Value: name.Value}).HashKey()]
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, key not found: %s",
				hash.Type(), name.Value,
			)
		}

		return evalCallExpression(right, pair.Value, env)
	case *ast.IndexExpression:
		leftInner, ok := right.Left.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, expected identifier, got %s",
				hash.Type(), right.Left.TokenLiteral(),
			)
		}

		pair, ok := hash.Pairs[(&objects.String{Value: leftInner.Value}).HashKey()]
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, key not found: %s",
				hash.Type(), leftInner.Value,
			)
		}

		index := Eval(right.Index, env)
		if objects.IsError(index) {
			return index
		}

		return evalIndexExpression(right, pair.Value, index, env)
	case *ast.ChainExpression:
		leftInner, ok := right.Left.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, expected identifier, got %s",
				hash.Type(), right.Left.TokenLiteral(),
			)
		}

		pair, ok := hash.Pairs[(&objects.String{Value: leftInner.Value}).HashKey()]
		if !ok {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"invalid chain expression for %s, key not found: %s",
				hash.Type(), leftInner.Value,
			)
		}

		return evalChainExpression(node, pair.Value, right.Right, env)
	case *ast.AssignmentExpression:
		assign := right.Right.(*ast.AssignmentExpression)
		leftKey, ok := assign.Left.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				assign.Token, env.GetFileDescriptorContext(),
				"invalid assignment expression for %s, expected identifier, got %s",
				hash.Type(), assign.Left.TokenLiteral(),
			)
		}

		obj := Eval(assign.Right, env)
		if objects.IsError(obj) {
			return obj
		}

		hashKey := (&objects.String{Value: leftKey.Value}).HashKey()
		hash.Pairs[hashKey] = objects.HashPair{Key: &objects.String{Value: leftKey.Value}, Value: obj}

		return obj

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"invalid chain expression for %s, got %s",
			hash.Type(), right.TokenLiteral(),
		)
	}
}

func evalAssignmentExpression(
	node *ast.AssignmentExpression,
	left ast.Expression,
	right objects.Object,
	env *objects.Environment,
) objects.Object {
	switch left := left.(type) {
	case *ast.Identifier:
		if !env.Has(left.Value) {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"assignment to undeclared variable: %s",
				left.Value,
			)
		}

		return env.Set(node, left.Value, right, false)
	case *ast.IndexExpression:
		leftObj := Eval(left.Left, env)
		if objects.IsError(leftObj) {
			return leftObj
		}

		switch leftObj := leftObj.(type) {
		case *objects.Array:
			return evalArrayAssignmentExpression(node, leftObj, left.Index, right, env)
		case *objects.Hash:
			return evalHashAssignmentExpression(node, leftObj, left.Index, right, env)
		case *objects.ImmutableHash:
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"cannot assign to immutable hash keys",
			)

		default:
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"left hand side of index assignment is not a valid indexable type: %s (%T)",
				leftObj, leftObj,
			)
		}

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"left hand side of assignment is not a valid expression: %s (%T)",
			left, left,
		)
	}
}

func evalArrayAssignmentExpression(
	node *ast.AssignmentExpression,
	arr *objects.Array,
	index ast.Expression,
	value objects.Object,
	env *objects.Environment,
) objects.Object {
	idx := Eval(index, env)
	if objects.IsError(idx) {
		return idx
	}

	switch idx := idx.(type) {
	case *objects.Integer:
		if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
			return objects.NewError(
				node.Token, env.GetFileDescriptorContext(),
				"array index out of bounds: %d",
				idx.Value,
			)
		}

		arr.Elements[idx.Value] = value

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"index operator not supported: %s",
			idx.Type(),
		)
	}

	return value
}

func evalHashAssignmentExpression(
	node *ast.AssignmentExpression,
	hash *objects.Hash,
	index ast.Expression,
	value objects.Object,
	env *objects.Environment,
) objects.Object {
	idx := Eval(index, env)
	if objects.IsError(idx) {
		return idx
	}

	key, ok := idx.(objects.Hashable)
	if !ok {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"invalid type given as hash key: %s",
			idx.Type(),
		)
	}

	hash.Pairs[key.HashKey()] = objects.HashPair{Key: idx, Value: value}

	return value
}

func evalNumberInfixExpression(
	node *ast.InfixExpression,
	left, right objects.Object,
	env *objects.Environment,
) objects.Object {
	leftVal := objects.UnwrapNumberValue(left)
	rightVal := objects.UnwrapNumberValue(right)

	switch node.Operator {
	case "+":
		return objects.WrapNumberValue(leftVal+rightVal, left, right)
	case "-":
		return objects.WrapNumberValue(leftVal-rightVal, left, right)
	case "*":
		return objects.WrapNumberValue(leftVal*rightVal, left, right)
	case "/":
		return objects.WrapNumberValue(leftVal/rightVal, left, right)
	case "^":
		return objects.WrapNumberValue(math.Pow(leftVal, rightVal), left, right)
	case "%":
		return objects.WrapNumberValue(math.Mod(leftVal, rightVal), left, right)
	case "<":
		return objects.NativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return objects.NativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return objects.NativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return objects.NativeBoolToBooleanObject(leftVal != rightVal)
	case "<=":
		return objects.NativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return objects.NativeBoolToBooleanObject(leftVal >= rightVal)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s %s %s",
			left.Type(), node.Operator, right.Type(),
		)
	}
}

func evalStringInfixExpression(
	node *ast.InfixExpression,
	left, right objects.Object,
	env *objects.Environment,
) objects.Object {
	leftVal := left.(*objects.String).Value
	rightVal := right.(*objects.String).Value

	switch node.Operator {
	case "+":
		return &objects.String{Value: leftVal + rightVal}
	case "==":
		return objects.NativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return objects.NativeBoolToBooleanObject(leftVal != rightVal)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s %s %s",
			left.Type(), node.Operator, right.Type(),
		)
	}
}

func evalStringableInfixExpression(
	node *ast.InfixExpression,
	left, right objects.Object,
	env *objects.Environment,
) objects.Object {
	leftVal := objects.StringifyObject(left)
	rightVal := objects.StringifyObject(right)

	switch node.Operator {
	case "+":
		return &objects.String{Value: leftVal + rightVal}

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"unknown operator: %s %s %s",
			left.Type(), node.Operator, right.Type(),
		)
	}
}

func evalExpressions(exps []ast.Expression, env *objects.Environment) []objects.Object {
	var result []objects.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if objects.IsError(evaluated) {
			return []objects.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func evalCallExpression(node *ast.CallExpression, function objects.Object, env *objects.Environment) objects.Object {
	args := evalExpressions(node.Arguments, env)

	if len(args) == 1 && objects.IsError(args[0]) {
		return objects.NewEmptyErrorWithParent(
			args[0].(*objects.Error),
			node.GetToken(),
			env.GetFileDescriptorContext(),
		)
	}

	functionObj, ok := function.(*objects.Function)
	if ok && len(args) < len(functionObj.Parameters) {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"wrong number of arguments. got %d, want %d",
			len(args), len(functionObj.Parameters),
		)
	}

	result := applyFunction(node, function, args, env)
	if objects.IsError(result) {
		return objects.NewEmptyErrorWithParent(
			result.(*objects.Error),
			node.GetToken(),
			env.GetFileDescriptorContext(),
		)
	}

	return result
}

func evalImportStatement(node *ast.ImportStatement, env *objects.Environment) objects.Object {
	if env.GetFileDescriptorContext() == nil {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"import statements can only be used within a file",
		)
	}

	filename := node.Path
	if !strings.HasSuffix(filename, ".zen") {
		filename += ".zen"
	}

	relativePath := filepath.Join(env.GetFileDescriptorContext().Path, filename)
	path, ok := filepath.Abs(relativePath)
	if ok != nil {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"invalid import path: %q",
			path,
		)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"failed to read imported file: %q",
			path,
		)
	}

	lexer := lexer.New(string(content))
	parser := parser.New(lexer, path)

	program := parser.ParseProgram()
	if len(parser.Errors()) > 0 {
		errors := []string{}
		for _, err := range parser.Errors() {
			errors = append(errors, err.String())
		}

		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"failed to parse imported file: %q\n%s",
			path, strings.Join(errors, "\n"),
		)
	}

	newEnv := objects.NewEnvironment(path)
	evaluated := Eval(program, newEnv)
	if evaluated == nil {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"failed to evaluate imported file: %q",
			path,
		)
	}

	if objects.IsError(evaluated) {
		return objects.NewEmptyErrorWithParent(
			evaluated.(*objects.Error),
			node.GetToken(),
			env.GetFileDescriptorContext(),
		)
	}

	hash := objects.CreateImmutableHashFromEnvExports(newEnv)

	if node.Aliased != nil {
		env.SetImmutableForcefully(node.Aliased.Value, hash)
	} else {
		cleanFilename := strings.TrimSuffix(path, ".zen")
		cleanFilename = filepath.Base(cleanFilename)

		env.SetImmutableForcefully(cleanFilename, hash)
	}

	return objects.NULL
}

func evalExportStatement(node *ast.ExportStatement, env *objects.Environment) objects.Object {
	exportedValue := Eval(node.Value, env)
	if objects.IsError(exportedValue) {
		return exportedValue
	}

	err := env.Export(exportedValue)
	if err != nil {
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"failed to export value: %q",
			err,
		)
	}

	return objects.NULL
}

func applyFunction(
	node *ast.CallExpression,
	fn objects.Object,
	args []objects.Object,
	env *objects.Environment,
) objects.Object {
	switch fn := fn.(type) {
	case *objects.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return objects.UnwrapReturnValue(evaluated)
	case *objects.ASTAwareBuiltin:
		for i, arg := range args {
			args[i] = WrapFunctionIfNeeded(arg)
		}

		return captureStdoutForBuiltin(node, fn, args, env)

	default:
		return objects.NewError(
			node.Token, env.GetFileDescriptorContext(),
			"not a function: %s",
			fn.Type(),
		)
	}
}

func extendFunctionEnv(fn *objects.Function, args []objects.Object) *objects.Environment {
	env := objects.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.SetImmutableForcefully(param.Value, args[paramIdx])
	}

	return env
}
