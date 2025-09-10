package vm

import (
	"fmt"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
)

const STACK_SIZE = 2048

type VM struct {
	constants    []objects.Object
	instructions code.Instructions

	stack []objects.Object
	sp    int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]objects.Object, STACK_SIZE),
		sp:    0,
	}
}

func (vm *VM) LastPoppedStackElem() objects.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}

		case code.OpPop:
			vm.pop()

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpMod:
			err := vm.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (vm *VM) push(obj objects.Object) error {
	if vm.sp >= STACK_SIZE {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

func (vm *VM) pop() objects.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if leftType == objects.INTEGER_OBJ && rightType == objects.INTEGER_OBJ {
		return vm.executeBinaryIntegerOperation(op, left, right)
	}

	return nil
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right objects.Object) error {
	leftValue := left.(*objects.Integer).Value
	rightValue := right.(*objects.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	case code.OpMod:
		result = leftValue % rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return vm.push(&objects.Integer{Value: result})
}
