package compiler

import (
	"fmt"
	"sort"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
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
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()
	WriteBuiltinSymbols(symbolTable)

	return &Compiler{
		constants:   []objects.Object{},
		symbolTable: symbolTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

func NewWithState(s *SymbolTable, constants []objects.Object) *Compiler {
	compiler := New()

	compiler.symbolTable = s
	compiler.constants = constants

	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
		for _, statement := range n.Statements {
			if err := c.Compile(statement); err != nil {
				return err
			}
		}
	case *ast.BlockStatement:
		for _, statement := range n.Statements {
			if err := c.Compile(statement); err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(n.Expression)
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
				return fmt.Errorf("cannot use named function literal in variable statement")
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
			err := c.Compile(n.Value)
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
			return fmt.Errorf("unsupported suffix expression left side: %T", n.Left)
		}

		symbol, ok := c.symbolTable.Resolve(ident.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", ident.Value)
		}

		c.loadSymbol(symbol)
		c.emit(code.OpConstant, c.addConstant(&objects.Integer{Value: 1}))

		switch n.Operator {
		case "++":
			c.emit(code.OpAdd)
		case "--":
			c.emit(code.OpSub)

		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}

		c.setSymbol(symbol)
	case *ast.IndexExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}

		err = c.Compile(n.Index)
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
			return fmt.Errorf("undefined variable %s", n.Value)
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
			err := c.Compile(element)
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
			err := c.Compile(key)
			if err != nil {
				return err
			}

			err = c.Compile(n.Pairs[key])
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
		err := c.Compile(n.ReturnValue)
		if err != nil {
			return err
		}

		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := c.Compile(n.Function)
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
			return fmt.Errorf("break statement not within a loop")
		}

		loop := &c.loops[c.loopIndex-1]
		pos := c.emit(code.OpJump, 9999)
		loop.breakPositions = append(loop.breakPositions, pos)
	case *ast.ContinueStatement:
		if c.loopIndex == 0 {
			return fmt.Errorf("continue statement not within a loop")
		}

		loop := c.loops[c.loopIndex-1]
		c.emit(code.OpJump, loop.startJumpIdx)
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
	case *ast.FunctionLiteral:
		return expr.Name == nil

	default:
		return true
	}
}

func (c *Compiler) compilePrefixExpression(node *ast.PrefixExpression) error {
	err := c.Compile(node.Right)
	if err != nil {
		return err
	}

	switch node.Operator {
	case "!":
		c.emit(code.OpBang)
	case "-":
		c.emit(code.OpMinus)
	default:
		return fmt.Errorf("unknown operator %s", node.Operator)
	}

	return nil
}

func (c *Compiler) compileInfixExpression(node *ast.InfixExpression) error {
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
	default:
		return fmt.Errorf("unknown operator %s", node.Operator)
	}

	return nil
}

func (c *Compiler) compileInfixExpressionOperands(node *ast.InfixExpression) error {
	// Reverse order for < and <= to so we're able to
	// treat them as > and >= during evaluation
	if node.Operator == "<" || node.Operator == "<=" {
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		err = c.Compile(node.Left)
		if err != nil {
			return err
		}

		return nil
	}

	err := c.Compile(node.Left)
	if err != nil {
		return err
	}

	err = c.Compile(node.Right)
	if err != nil {
		return err
	}

	return nil
}

func (c *Compiler) compileConditionalIfExpression(node *ast.IfExpression) error {
	err := c.Compile(node.Condition)
	if err != nil {
		return err
	}

	// Emits a jump instruction with an invalid operand that we'll change
	// later on when we know where in the stack to jump to.
	jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

	err = c.Compile(node.Consequence)
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
		err := c.Compile(node.Alternative)
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

func (c *Compiler) compileFunctionLiteral(node *ast.FunctionLiteral, constructNamed bool) error {
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

	err := c.Compile(node.Body)
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

	compiledFn := &objects.CompiledFunction{
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

func (c *Compiler) compileChainExpression(node *ast.ChainExpression, inner bool) error {
	leftIdent, ok := node.Left.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("unsupported chain expression left side: %T", node.Left)
	}

	if inner {
		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: leftIdent.Value}))
		c.emit(code.OpIndex)
	} else {
		c.Compile(node.Left)
	}

	switch right := node.Right.(type) {
	case *ast.Identifier:
		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: right.Value}))
	case *ast.ChainExpression:
		return c.compileChainExpression(right, true)
	case *ast.IndexExpression:
		ident, ok := right.Left.(*ast.Identifier)
		if !ok {
			return fmt.Errorf("unsupported call expression function in chain: %T", right.Left)
		}

		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: ident.Value}))
		c.emit(code.OpIndex)

		index, ok := right.Index.(*ast.IntegerLiteral)
		if !ok {
			return fmt.Errorf("unsupported index expression index in chain: %T", right.Index)
		}

		c.emit(code.OpConstant, c.addConstant(&objects.Integer{Value: index.Value}))
	case *ast.CallExpression:
		ident, ok := right.Function.(*ast.Identifier)
		if !ok {
			return fmt.Errorf("unsupported call expression function in chain: %T", right.Function)
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

	default:
		return fmt.Errorf("unsupported chain expression right side: %T", right)
	}

	c.emit(code.OpIndex)
	return nil
}

func (c *Compiler) compileAssignmentExpression(node *ast.AssignmentExpression) error {
	switch left := node.Left.(type) {
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(left.Value)
		if !ok {
			return fmt.Errorf("assignment to undeclared variable: %s", left.Value)
		}

		err := c.Compile(node.Right)
		if err != nil {
			return err
		}

		c.setSymbol(symbol)

	case *ast.IndexExpression:
		err := c.Compile(left.Left)
		if err != nil {
			return err
		}

		err = c.Compile(left.Index)
		if err != nil {
			return err
		}

		err = c.Compile(node.Right)
		if err != nil {
			return err
		}

		c.emit(code.OpIndexAssign)

	default:
		return fmt.Errorf("left hand side of assignment is not a valid expression: %T", left)
	}

	return nil
}

func (c *Compiler) compileFunctionArguments(node *ast.CallExpression) error {
	for _, arg := range node.Arguments {
		if err := c.Compile(arg); err != nil {
			return err
		}
	}

	c.emit(code.OpCall, len(node.Arguments))

	return nil
}

func (c *Compiler) compileWhileExpression(node *ast.WhileExpression) error {
	startJumpIdx := c.enterLoop()

	err := c.Compile(node.Condition)
	if err != nil {
		return err
	}

	jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

	err = c.Compile(node.Body)
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
