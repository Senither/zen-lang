package compiler

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/lexer"
	"github.com/senither/zen-lang/objects"
	"github.com/senither/zen-lang/parser"
)

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type CompilationLoop struct {
	startJumpIdx   int
	breakPositions []int
}

type Compiler struct {
	constants []objects.Object

	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int

	loops     []CompilationLoop
	loopIndex int

	file *objects.FileDescriptorContext
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New(path interface{}) *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	WriteBuiltinSymbols(symbolTable)

	var file *objects.FileDescriptorContext
	if pathStr, ok := path.(string); ok {
		file = objects.NewFileDescriptorContext(pathStr)
	}

	return &Compiler{
		constants:   []objects.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
		file:        file,
	}
}

func NewWithState(path interface{}, s *SymbolTable, constants []objects.Object) *Compiler {
	compiler := New(path)

	compiler.symbolTable = s
	compiler.constants = constants

	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	program, ok := node.(*ast.Program)
	if !ok {
		return fmt.Errorf("can only compile program nodes, got %T", node)
	}

	var errors []*objects.Error
	for _, statement := range program.Statements {
		if err := c.compileInstruction(statement); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) == 0 {
		return nil
	}

	var combinedErr bytes.Buffer
	for _, err := range errors {
		combinedErr.WriteString(err.Inspect() + "\n")
	}

	return fmt.Errorf("%s", combinedErr.String())
}

func (c *Compiler) compileInstruction(node ast.Node) *objects.Error {
	switch n := node.(type) {
	case *ast.BlockStatement:
		for _, statement := range n.Statements {
			if err := c.compileInstruction(statement); err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.compileInstruction(n.Expression)
		if err != nil {
			return err
		}

		if c.shouldPopExpression(n.Expression) {
			c.emit(code.OpPop)
		}
	case *ast.VariableStatement:
		symbol := c.symbolTable.Define(n.Name.Value, n.Mutable)

		if fn, ok := n.Value.(*ast.FunctionLiteral); ok {
			if fn.Name != nil {
				return objects.NewError(
					n.Token, c.file,
					"cannot use named function literal in variable statement",
				)
			}

			funcLit := &ast.FunctionLiteral{
				Parameters: fn.Parameters,
				Body:       fn.Body,
				Name:       &ast.Identifier{Value: n.Name.Value},
			}

			err := c.compileFunctionLiteral(funcLit, false)
			if err != nil {
				return err
			}
		} else {
			err := c.compileInstruction(n.Value)
			if err != nil {
				return err
			}
		}

		c.setSymbol(symbol)

	// Expression operators
	case *ast.PrefixExpression:
		err := c.compilePrefixExpression(n)
		if err != nil {
			return err
		}
	case *ast.InfixExpression:
		err := c.compileInfixExpression(n)
		if err != nil {
			return err
		}
	case *ast.SuffixExpression:
		ident, ok := n.Left.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				n.Token, c.file,
				"unsupported suffix expression left side: %T",
				n.Left,
			)
		}

		symbol, ok := c.symbolTable.Resolve(ident.Value)
		if !ok {
			return objects.NewError(
				n.Token, c.file,
				"undefined variable %s",
				ident.Value,
			)
		}

		c.loadSymbol(symbol)
		c.emit(code.OpConstant, c.addConstant(&objects.Integer{Value: 1}))

		switch n.Operator {
		case "++":
			c.emit(code.OpAdd)
		case "--":
			c.emit(code.OpSub)

		default:
			return objects.NewError(
				n.Token, c.file,
				"unknown operator %s",
				n.Operator,
			)
		}

		c.setSymbol(symbol)
	case *ast.IndexExpression:
		err := c.compileInstruction(n.Left)
		if err != nil {
			return err
		}

		err = c.compileInstruction(n.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)
	case *ast.ChainExpression:
		err := c.compileChainExpression(n, false)
		if err != nil {
			return err
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(n.Value)
		if !ok {
			return objects.NewError(
				n.Token, c.file,
				"undefined variable %s",
				n.Value,
			)
		}

		c.loadSymbol(symbol)
	case *ast.AssignmentExpression:
		err := c.compileAssignmentExpression(n)
		if err != nil {
			return err
		}

	// Expression types
	case *ast.NullLiteral:
		c.emit(code.OpNull)
	case *ast.StringLiteral:
		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: n.Value}))
	case *ast.IntegerLiteral:
		c.emit(code.OpConstant, c.addConstant(&objects.Integer{Value: n.Value}))
	case *ast.FloatLiteral:
		c.emit(code.OpConstant, c.addConstant(&objects.Float{Value: n.Value}))
	case *ast.BooleanLiteral:
		if n.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.ArrayLiteral:
		for _, element := range n.Elements {
			err := c.compileInstruction(element)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpArray, len(n.Elements))
	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for key := range n.Pairs {
			keys = append(keys, key)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, key := range keys {
			err := c.compileInstruction(key)
			if err != nil {
				return err
			}

			err = c.compileInstruction(n.Pairs[key])
			if err != nil {
				return err
			}
		}

		c.emit(code.OpHash, len(n.Pairs)*2)
	case *ast.IfExpression:
		err := c.compileConditionalIfExpression(n)
		if err != nil {
			return err
		}
	case *ast.FunctionLiteral:
		err := c.compileFunctionLiteral(n, true)
		if err != nil {
			return err
		}
	case *ast.ReturnStatement:
		err := c.compileInstruction(n.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := c.compileInstruction(n.Function)
		if err != nil {
			return err
		}

		err = c.compileFunctionArguments(n)
		if err != nil {
			return err
		}

	case *ast.WhileExpression:
		err := c.compileWhileExpression(n)
		if err != nil {
			return err
		}
	case *ast.BreakStatement:
		if c.loopIndex == 0 {
			return objects.NewError(
				n.Token, c.file,
				"break statement not within a loop",
			)
		}

		loop := &c.loops[c.loopIndex-1]
		pos := c.emit(code.OpJump, 9999)
		loop.breakPositions = append(loop.breakPositions, pos)
	case *ast.ContinueStatement:
		if c.loopIndex == 0 {
			return objects.NewError(
				n.Token, c.file,
				"continue statement not within a loop",
			)
		}

		loop := c.loops[c.loopIndex-1]
		c.emit(code.OpJump, loop.startJumpIdx)

	// Import & Export statements
	case *ast.ImportStatement:
		err := c.compileImportStatement(n)
		if err != nil {
			return err
		}
	case *ast.ExportStatement:
		err := c.compileExportStatement(n)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++

	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) enterLoop() int {
	loop := CompilationLoop{
		startJumpIdx:   len(c.currentInstructions()),
		breakPositions: []int{},
	}

	c.loops = append(c.loops, loop)
	c.loopIndex++

	return loop.startJumpIdx
}

func (c *Compiler) leaveLoop(endIdx int) CompilationLoop {
	loop := c.loops[c.loopIndex-1]

	for _, breakPos := range loop.breakPositions {
		c.changeInstructionOperandAt(breakPos, endIdx)
	}

	c.loops = c.loops[:len(c.loops)-1]
	c.loopIndex--

	return loop
}

func (c *Compiler) addConstant(obj objects.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) addInstruction(ins []byte) int {
	newInstructionPos := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return newInstructionPos
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) replaceInstruction(pos int, newInstructions []byte) {
	ins := c.currentInstructions()

	for i := 0; i < len(newInstructions); i++ {
		ins[pos+i] = newInstructions[i]
	}
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) changeInstructionOperandAt(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	return len(c.currentInstructions()) > 0 &&
		c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].previousInstruction = previous
}

func (c *Compiler) shouldPopExpression(expr ast.Expression) bool {
	switch expr := expr.(type) {
	case *ast.AssignmentExpression, *ast.SuffixExpression, *ast.WhileExpression:
		return false
	case *ast.ChainExpression:
		if _, ok := expr.Right.(*ast.AssignmentExpression); ok {
			return false
		}
		return true
	case *ast.FunctionLiteral:
		return expr.Name == nil

	default:
		return true
	}
}

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) *objects.Error {
	err := c.compileInstruction(node.Right)
	if err != nil {
		return err
	}

	switch node.Operator {
	case "!":
		c.emit(code.OpBang)
	case "-":
		c.emit(code.OpMinus)

	default:
		return objects.NewError(
			node.Token, c.file,
			"unknown operator %s",
			node.Operator,
		)
	}

	return nil
}

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) *objects.Error {
	err := c.compileInfixExpressionOperands(node)
	if err != nil {
		return err
	}

	switch node.Operator {
	case "+":
		c.emit(code.OpAdd)
	case "-":
		c.emit(code.OpSub)
	case "*":
		c.emit(code.OpMul)
	case "/":
		c.emit(code.OpDiv)
	case "^":
		c.emit(code.OpPow)
	case "%":
		c.emit(code.OpMod)
	case "==":
		c.emit(code.OpEqual)
	case "!=":
		c.emit(code.OpNotEqual)
	case ">", "<":
		c.emit(code.OpGreaterThan)
	case ">=", "<=":
		c.emit(code.OpGreaterThanOrEqual)
	case "&&":
		c.emit(code.OpAnd)
	case "||":
		c.emit(code.OpOr)

	default:
		return objects.NewError(
			node.Token, c.file,
			"unknown operator %s",
			node.Operator,
		)
	}

	return nil
}

func (c *Compiler) compileInfixExpressionOperands(node *ast.InfixExpression) *objects.Error {
	// Reverse order for < and <= to so we're able to
	// treat them as > and >= during evaluation
	if node.Operator == "<" || node.Operator == "<=" {
		err := c.compileInstruction(node.Right)
		if err != nil {
			return err
		}

		err = c.compileInstruction(node.Left)
		if err != nil {
			return err
		}

		return nil
	}

	err := c.compileInstruction(node.Left)
	if err != nil {
		return err
	}

	err = c.compileInstruction(node.Right)
	if err != nil {
		return err
	}

	return nil
}

func (c *Compiler) compileConditionalIfExpression(node *ast.IfExpression) *objects.Error {
	err := c.compileInstruction(node.Condition)
	if err != nil {
		return err
	}

	// Emits a jump instruction with an invalid operand that we'll change
	// later on when we know where in the stack to jump to.
	jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

	err = c.compileInstruction(node.Consequence)
	if err != nil {
		return err
	}

	if c.lastInstructionIs(code.OpPop) {
		c.removeLastPop()
	}

	var jumpPos int = -1
	if !c.lastInstructionIs(code.OpJump) {
		jumpPos = c.emit(code.OpJump, 9999)
	}

	afterConsequencePos := len(c.currentInstructions())
	c.changeInstructionOperandAt(jumpNotTruthyPos, afterConsequencePos)

	if node.Intermediary != nil {
		err := c.compileConditionalIfExpression(node.Intermediary)
		if err != nil {
			return err
		}
	} else if node.Alternative == nil {
		c.emit(code.OpNull)
	} else {
		err := c.compileInstruction(node.Alternative)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}
	}

	if jumpPos >= 0 {
		afterAlternativePos := len(c.currentInstructions())
		c.changeInstructionOperandAt(jumpPos, afterAlternativePos)
	}

	return nil
}

func (c *Compiler) compileFunctionLiteral(node *ast.FunctionLiteral, constructNamed bool) *objects.Error {
	var symbol *Symbol
	if constructNamed && node.Name != nil {
		sym := c.symbolTable.Define(node.Name.Value, false)
		symbol = &sym
	}

	c.enterScope()

	if node.Name != nil {
		c.symbolTable.DefineFunctionName(node.Name.Value, false)
	}

	for _, param := range node.Parameters {
		c.symbolTable.Define(param.Value, false)
	}

	err := c.compileInstruction(node.Body)
	if err != nil {
		return err
	}

	if c.lastInstructionIs(code.OpPop) {
		c.replaceLastPopWithReturn()
	}

	if !c.lastInstructionIs(code.OpReturnValue) {
		c.emit(code.OpReturn)
	}

	freeSymbols := c.symbolTable.FreeSymbols
	numLocals := c.symbolTable.numDefinitions
	instructions := c.leaveScope()

	for _, sym := range freeSymbols {
		c.loadSymbol(sym)
	}

	cfName := ""
	if node.Name != nil {
		cfName = node.Name.Value
	}

	compiledFn := &objects.CompiledFunction{
		Name:               cfName,
		OpcodeInstructions: instructions,
		NumLocals:          numLocals,
		NumParameters:      len(node.Parameters),
	}

	c.emit(code.OpClosure, c.addConstant(compiledFn), len(freeSymbols))

	if symbol != nil {
		c.setSymbol(*symbol)
	}

	return nil
}

func (c *Compiler) compileChainExpression(node *ast.ChainExpression, inner bool) *objects.Error {
	leftIdent, ok := node.Left.(*ast.Identifier)
	if !ok {
		return objects.NewError(
			node.Token, c.file,
			"unsupported chain expression left side: %T",
			node.Left,
		)
	}

	if inner {
		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: leftIdent.Value}))
		c.emit(code.OpIndex)
	} else {
		c.compileInstruction(node.Left)
	}

	switch right := node.Right.(type) {
	case *ast.Identifier:
		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: right.Value}))
	case *ast.ChainExpression:
		return c.compileChainExpression(right, true)
	case *ast.IndexExpression:
		ident, ok := right.Left.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				node.Token, c.file,
				"unsupported call expression function in chain: %T",
				right.Left,
			)
		}

		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: ident.Value}))
		c.emit(code.OpIndex)

		index, ok := right.Index.(*ast.IntegerLiteral)
		if !ok {
			return objects.NewError(
				node.Token, c.file,
				"unsupported index expression index in chain: %T",
				right.Index,
			)
		}

		c.emit(code.OpConstant, c.addConstant(&objects.Integer{Value: index.Value}))
	case *ast.CallExpression:
		ident, ok := right.Function.(*ast.Identifier)
		if !ok {
			return objects.NewError(
				node.Token, c.file,
				"unsupported call expression function in chain: %T",
				right.Function,
			)
		}

		if !inner {
			symbol, ok := c.symbolTable.Resolve(fmt.Sprintf("%s.%s", node.Left, ident.Value))
			if ok {
				c.loadSymbol(symbol)

				return c.compileFunctionArguments(right)
			}
		}

		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: ident.Value}))
		c.emit(code.OpIndex)

		return c.compileFunctionArguments(right)
	case *ast.AssignmentExpression:
		return c.compileChainAssignment(right)

	default:
		return objects.NewError(
			node.Token, c.file,
			"unsupported chain expression right side: %T",
			right,
		)
	}

	c.emit(code.OpIndex)
	return nil
}

func (c *Compiler) compileChainAssignment(assign *ast.AssignmentExpression) *objects.Error {
	innerAssign, ok := assign.Right.(*ast.AssignmentExpression)
	if !ok {
		if index, ok := assign.Left.(*ast.IndexExpression); ok {
			return c.compileChainIndexAssignment(assign, index)
		}

		return objects.NewError(
			assign.Token, c.file,
			"unsupported chain assignment structure: %T",
			assign.Right,
		)
	}

	assignmentPath, err := c.retrieveAssignmentKeys(innerAssign.Left, innerAssign)
	if err != nil {
		return err
	}

	if len(assignmentPath) == 0 {
		return objects.NewError(
			innerAssign.Token, c.file,
			"empty assignment key path in chain",
		)
	}

	for i := 0; i < len(assignmentPath)-1; i++ {
		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: assignmentPath[i]}))
		c.emit(code.OpIndex)
	}

	finalKey := assignmentPath[len(assignmentPath)-1]
	c.emit(code.OpConstant, c.addConstant(&objects.String{Value: finalKey}))

	err = c.compileInstruction(innerAssign.Right)
	if err != nil {
		return err
	}

	c.emit(code.OpIndexAssign)
	return nil
}

func (c *Compiler) compileChainIndexAssignment(assign *ast.AssignmentExpression, index *ast.IndexExpression) *objects.Error {
	propIdent, ok := index.Left.(*ast.Identifier)
	if !ok {
		return objects.NewError(
			index.Token, c.file,
			"unsupported index expression left in chain assignment: %T",
			index.Left,
		)
	}

	c.emit(code.OpConstant, c.addConstant(&objects.String{Value: propIdent.Value}))
	c.emit(code.OpIndex)

	err := c.compileInstruction(index.Index)
	if err != nil {
		return err
	}

	err = c.compileInstruction(assign.Right)
	if err != nil {
		return err
	}

	c.emit(code.OpIndexAssign)
	return nil
}

func (c *Compiler) retrieveAssignmentKeys(exp ast.Expression, innerAssign *ast.AssignmentExpression) ([]string, *objects.Error) {
	switch v := exp.(type) {
	case *ast.Identifier:
		return []string{v.Value}, nil
	case *ast.ChainExpression:
		leftKeys, err := c.retrieveAssignmentKeys(v.Left, innerAssign)
		if err != nil {
			return nil, err
		}

		rightKeys, err := c.retrieveAssignmentKeys(v.Right, innerAssign)
		if err != nil {
			return nil, err
		}

		return append(leftKeys, rightKeys...), nil
	case *ast.AssignmentExpression:
		return c.retrieveAssignmentKeys(v.Left, innerAssign)

	default:
		return nil, objects.NewError(
			innerAssign.Token, c.file,
			"invalid assignment key path element in chain: %T",
			v,
		)
	}
}

func (c *Compiler) compileAssignmentExpression(node *ast.AssignmentExpression) *objects.Error {
	switch left := node.Left.(type) {
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(left.Value)
		if !ok {
			return objects.NewError(
				node.Token, c.file,
				"assignment to undeclared variable: %s",
				left.Value,
			)
		}

		if !symbol.Mutable {
			return objects.NewError(
				node.Token, c.file,
				"cannot modify immutable variable: %s",
				left.Value,
			)
		}

		err := c.compileInstruction(node.Right)
		if err != nil {
			return err
		}

		c.setSymbol(symbol)

	case *ast.IndexExpression:
		err := c.compileInstruction(left.Left)
		if err != nil {
			return err
		}

		err = c.compileInstruction(left.Index)
		if err != nil {
			return err
		}

		err = c.compileInstruction(node.Right)
		if err != nil {
			return err
		}

		c.emit(code.OpIndexAssign)

	default:
		return objects.NewError(
			node.Token, c.file,
			"left hand side of assignment is not a valid expression: %T",
			left,
		)
	}

	return nil
}

func (c *Compiler) compileFunctionArguments(node *ast.CallExpression) *objects.Error {
	for _, arg := range node.Arguments {
		if err := c.compileInstruction(arg); err != nil {
			return objects.NewEmptyErrorWithParent(err, node.Token, c.file)
		}
	}

	c.emit(code.OpCall, len(node.Arguments))

	return nil
}

func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) *objects.Error {
	startJumpIdx := c.enterLoop()

	err := c.compileInstruction(node.Condition)
	if err != nil {
		return objects.NewEmptyErrorWithParent(err, node.Token, c.file)
	}

	jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

	err = c.compileInstruction(node.Body)
	if err != nil {
		return err
	}

	c.emit(code.OpJump, startJumpIdx)

	endJumpIdx := len(c.currentInstructions())

	c.emit(code.OpLoopEnd)
	c.leaveLoop(endJumpIdx)

	c.changeInstructionOperandAt(jumpNotTruthyPos, endJumpIdx)

	return nil
}

func (c *Compiler) compileImportStatement(node *ast.ImportStatement) *objects.Error {
	if c.file == nil {
		return objects.NewError(
			node.Token, c.file,
			"cannot use import statement without a file context",
		)
	}

	filename := node.Path
	if !strings.HasSuffix(filename, ".zen") {
		filename += ".zen"
	}

	relativePath := filepath.Join(c.file.Path, filename)
	path, ok := filepath.Abs(relativePath)
	if ok != nil {
		return objects.NewError(
			node.Token, c.file,
			"invalid import path: %q",
			path,
		)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return objects.NewError(
			node.Token, c.file,
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
			node.Token, c.file,
			"failed to parse imported file: %q\n%s",
			path, strings.Join(errors, "\n"),
		)
	}

	var name string
	if node.Aliased != nil {
		name = node.Aliased.Value
	} else {
		cleanFilename := strings.TrimSuffix(path, ".zen")
		name = filepath.Base(cleanFilename)
	}

	symbol := c.symbolTable.Define(name, false)

	importCompiler := New(path)
	err = importCompiler.Compile(program)
	if err != nil {
		return objects.NativeErrorToErrorObject(err)
	}

	c.emit(code.OpImport, c.addConstant(&objects.CompiledFileImport{
		Name:               name,
		Constants:          importCompiler.constants,
		OpcodeInstructions: importCompiler.currentInstructions(),
	}))

	c.setSymbol(symbol)

	return nil
}

func (c *Compiler) compileExportStatement(node *ast.ExportStatement) *objects.Error {
	switch v := node.Value.(type) {
	case *ast.Identifier:
		err := c.compileInstruction(v)
		if err != nil {
			return err
		}

		c.emit(code.OpExport)
	case *ast.FunctionLiteral:
		if v.Name == nil {
			return objects.NewError(
				node.Token, c.file,
				"cannot use unnamed function literal in export statement",
			)
		}

		symbol := c.symbolTable.Define(v.Name.Value, false)

		err := c.compileFunctionLiteral(v, false)
		if err != nil {
			return err
		}

		c.setSymbol(symbol)
		c.loadSymbol(symbol)

		c.emit(code.OpExport)

	default:
		return objects.NewError(
			node.Token, c.file,
			"cannot export expression of type %T",
			node.Value,
		)
	}

	return nil
}

func (c *Compiler) loadSymbol(symbol Symbol) {
	switch symbol.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, symbol.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, symbol.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, symbol.Index)
	case GlobalBuiltinScope:
		c.emit(code.OpGetGlobalBuiltin, symbol.Index)
	case FreeScope:
		c.emit(code.OpGetFree, symbol.Index)
	case FunctionScope:
		c.emit(code.OpCurrentClosure)
	}
}

func (c *Compiler) setSymbol(symbol Symbol) {
	if symbol.Scope == GlobalScope {
		c.emit(code.OpSetGlobal, symbol.Index)
	} else {
		c.emit(code.OpSetLocal, symbol.Index)
	}
}
