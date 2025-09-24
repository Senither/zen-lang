package vm

import (
	"fmt"
	"math"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
)

const STACK_SIZE = 2048
const GLOBALS_SIZE = 65536

var (
	NULL  = &objects.Null{}
	TRUE  = &objects.Boolean{Value: true}
	FALSE = &objects.Boolean{Value: false}
)

type VM struct {
	constants    []objects.Object
	instructions code.Instructions

	stack []objects.Object
	sp    int

	globals []objects.Object
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack: make([]objects.Object, STACK_SIZE),
		sp:    0,

		globals: make([]objects.Object, GLOBALS_SIZE),
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []objects.Object) *VM {
	vm := New(bytecode)

	vm.globals = globals

	return vm
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

		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpPow, code.OpMod:
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

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
			err := vm.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpBang:
			err := vm.executeBangOperator()
			if err != nil {
				return err
			}
		case code.OpMinus:
			err := vm.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(vm.instructions[ip+1:]))
			ip += 2

			condition := vm.pop()
			if !isTruthy(condition) {
				ip = pos - 1
			}

		case code.OpNull:
			err := vm.push(NULL)
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

	switch {
	case isNumber(leftType) && isNumber(rightType):
		return vm.executeBinaryNumberOperation(op, left, right)
	case leftType == objects.STRING_OBJ && rightType == objects.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)

	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
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
	case code.OpPow:
		result = math.Pow(leftValue, rightValue)
	case code.OpMod:
		result = math.Mod(leftValue, rightValue)
	default:
		return fmt.Errorf("unknown number operator: %d", op)
	}

	return vm.push(wrapNumberValue(result, left, right))
}

func (vm *VM) executeBinaryStringOperation(op code.Opcode, left, right objects.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*objects.String).Value
	rightValue := right.(*objects.String).Value

	return vm.push(&objects.String{Value: leftValue + rightValue})
}

func (vm *VM) executeComparison(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	if isNumber(leftType) && isNumber(rightType) {
		return vm.executeComparisonNumberOperation(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, leftType, rightType)
	}
}

func (vm *VM) executeComparisonNumberOperation(op code.Opcode, left, right objects.Object) error {
	leftValue := unwrapNumberValue(left)
	rightValue := unwrapNumberValue(right)

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftValue > rightValue))
	case code.OpGreaterThanOrEqual:
		return vm.push(nativeBoolToBooleanObject(leftValue >= rightValue))
	default:
		return fmt.Errorf("unknown number operator: %d", op)
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case TRUE:
		return vm.push(FALSE)
	case FALSE:
		return vm.push(TRUE)
	case NULL:
		return vm.push(TRUE)

	default:
		return vm.push(FALSE)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if !isNumber(operand.Type()) {
		return fmt.Errorf("unknown operator: -%s", operand.Type())
	}

	return vm.push(wrapNumberValue(-unwrapNumberValue(operand), operand, operand))
}

func nativeBoolToBooleanObject(input bool) *objects.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj objects.Object) bool {
	switch obj := obj.(type) {
	case *objects.Boolean:
		return obj.Value
	case *objects.Null:
		return false

	default:
		return true
	}
}

func isNumber(t objects.ObjectType) bool {
	return t == objects.INTEGER_OBJ || t == objects.FLOAT_OBJ
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
