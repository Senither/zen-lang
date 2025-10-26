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

type Compiler struct {
	constants []objects.Object

	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
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

		c.emit(code.OpPop)
	case *ast.VariableStatement:
		err := c.Compile(n.Value)
		if err != nil {
			return err
		}

		symbol := c.symbolTable.Define(n.Name.Value, n.Mutable)
		if symbol.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, symbol.Index)
		} else {
			c.emit(code.OpSetLocal, symbol.Index)
		}

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

		if symbol.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, symbol.Index)
		} else {
			c.emit(code.OpGetLocal, symbol.Index)
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
		c.enterScope()

		for _, param := range n.Parameters {
			c.symbolTable.Define(param.Value, false)
		}

		err := c.Compile(n.Body)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		compiledFn := &objects.CompiledFunction{
			OpcodeInstructions: instructions,
			NumLocals:          numLocals,
			NumParameters:      len(n.Parameters),
		}

		c.emit(code.OpConstant, c.addConstant(compiledFn))
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

	jumpPos := c.emit(code.OpJump, 9999)

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

	afterAlternativePos := len(c.currentInstructions())
	c.changeInstructionOperandAt(jumpPos, afterAlternativePos)

	return nil
}

func (c *Compiler) compileChainExpression(node *ast.ChainExpression, inner bool) error {
	switch left := node.Left.(type) {
	case *ast.Identifier:
		if inner {
			c.emit(code.OpConstant, c.addConstant(&objects.String{Value: left.Value}))
			c.emit(code.OpIndex)
		} else {
			c.Compile(node.Left)
		}

	default:
		return fmt.Errorf("unsupported chain expression left side: %T", left)
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

		c.emit(code.OpConstant, c.addConstant(&objects.String{Value: ident.Value}))
		c.emit(code.OpIndex)

		return c.compileFunctionArguments(right)

	default:
		return fmt.Errorf("unsupported chain expression right side: %T", right)
	}

	c.emit(code.OpIndex)
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
