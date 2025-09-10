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

func (vm *VM) StackTop() objects.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
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

		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*objects.Integer).Value
			rightValue := right.(*objects.Integer).Value
			result := leftValue + rightValue

			err := vm.push(&objects.Integer{Value: result})
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
