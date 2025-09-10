package compiler

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type Compiler struct {
	instructions code.Instructions
	constants    []objects.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []objects.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch n := node.(type) {
	case *ast.Program:
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
	case *ast.InfixExpression:
		err := c.Compile(n.Left)
		if err != nil {
			return err
		}

		err = c.Compile(n.Right)
		if err != nil {
			return err
		}

		switch n.Operator {
		case "+":
			c.emit(code.OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", n.Operator)
		}

	case *ast.IntegerLiteral:
		c.emit(code.OpConstant, c.addConstant(&objects.Integer{Value: n.Value}))
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

	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)

	return posNewInstruction
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []objects.Object
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
