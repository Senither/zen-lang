package vm

import (
	"fmt"
	"math"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
)

const STACK_SIZE = 2048

var (
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
)

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

		case code.OpTrue:
			err := vm.push(TRUE)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(FALSE)
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

	if isNumber(leftType) && isNumber(rightType) {
		return vm.executeBinaryNumberOperation(op, left, right)
	}

	return nil
}

func (vm *VM) executeBinaryNumberOperation(op code.Opcode, left, right objects.Object) error {
	leftValue := unwrapNumberValue(left)
	rightValue := unwrapNumberValue(right)

	var result float64

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
		result = math.Mod(leftValue, rightValue)
	default:
		return fmt.Errorf("unknown number operator: %d", op)
	}

	return vm.push(wrapNumberValue(result, left, right))
}

func isNumber(obj objects.ObjectType) bool {
	switch obj {
	case objects.INTEGER_OBJ, objects.FLOAT_OBJ:
		return true
	default:
		return false
	}
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
