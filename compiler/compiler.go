package compiler

import (
	"fmt"
	"sort"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type Compiler struct {
	instructions code.Instructions
	constants    []objects.Object

	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction

	symbolTable *SymbolTable
}

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []objects.Object{},

		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},

		symbolTable: NewSymbolTable(),
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
		c.emit(code.OpSetGlobal, symbol.Index)

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

		c.emit(code.OpGetGlobal, symbol.Index)

	// Expression types
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
	}

	return nil
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

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)

	return posNewInstruction
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	c.previousInstruction = c.lastInstruction
	c.lastInstruction = EmittedInstruction{Opcode: op, Position: pos}
}

func (c *Compiler) replaceInstruction(pos int, newInstructions []byte) {
	for i := range newInstructions {
		c.instructions[pos+i] = newInstructions[i]
	}
}

func (c *Compiler) changeInstructionOperandAt(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	return c.lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
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

	afterConsequencePos := len(c.instructions)
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

	afterAlternativePos := len(c.instructions)
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
	default:
		return fmt.Errorf("unsupported chain expression right side: %T", right)
	}

	c.emit(code.OpIndex)

	return nil
}
