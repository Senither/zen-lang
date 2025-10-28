package vm

import (
	"fmt"
	"math"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/compiler"
	"github.com/senither/zen-lang/objects"
)

const (
	MAX_FRAMES   = 1024
	STACK_SIZE   = 2048
	GLOBALS_SIZE = 65536
)

type VM struct {
	constants []objects.Object

	stack []objects.Object
	// The stack pointer, this always points to the next value.
	// Top of stack is stack[sp-1]
	sp int

	globals []objects.Object

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &objects.CompiledFunction{OpcodeInstructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MAX_FRAMES)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,

		stack: make([]objects.Object, STACK_SIZE),
		sp:    0,

		globals: make([]objects.Object, GLOBALS_SIZE),

		frames:      frames,
		framesIndex: 1,
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
	var (
		ip  int
		ins code.Instructions
		op  code.Opcode
	)

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = code.Opcode(ins[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

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
			err := vm.push(objects.TRUE)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(objects.FALSE)
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
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !objects.IsTruthy(condition) {
				vm.currentFrame().ip = pos - 1
			}

		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArray(vm.sp-numElements, vm.sp)
			vm.sp -= numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHash:
			numElements := int(code.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2

			hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}

			vm.sp -= numElements

			err = vm.push(hash)
			if err != nil {
				return err
			}
		case code.OpNull:
			err := vm.push(objects.NULL)
			if err != nil {
				return err
			}

		case code.OpCall:
			numArgs := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			err := vm.executeCall(int(numArgs))
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			returnValue := vm.pop()

			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(objects.NULL)
			if err != nil {
				return err
			}

		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(ins[ip+1:])
			vm.currentFrame().ip += 1

			definition := objects.Builtins[builtinIndex]

			err := vm.push(definition.Builtin)
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

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := left.Type()

	switch {
	case objects.IsNumber(leftType) && objects.IsNumber(rightType):
		return vm.executeBinaryNumberOperation(op, left, right)
	case leftType == objects.STRING_OBJ && rightType == objects.STRING_OBJ:
		return vm.executeBinaryStringOperation(op, left, right)

	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
	}
}

func (vm *VM) executeBinaryNumberOperation(op code.Opcode, left, right objects.Object) error {
	leftValue := objects.UnwrapNumberValue(left)
	rightValue := objects.UnwrapNumberValue(right)

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

	return vm.push(objects.WrapNumberValue(result, left, right))
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

	if objects.IsNumber(leftType) && objects.IsNumber(rightType) {
		return vm.executeComparisonNumberOperation(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(objects.NativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(objects.NativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, leftType, rightType)
	}
}

func (vm *VM) executeComparisonNumberOperation(op code.Opcode, left, right objects.Object) error {
	leftValue := objects.UnwrapNumberValue(left)
	rightValue := objects.UnwrapNumberValue(right)

	switch op {
	case code.OpEqual:
		return vm.push(objects.NativeBoolToBooleanObject(leftValue == rightValue))
	case code.OpNotEqual:
		return vm.push(objects.NativeBoolToBooleanObject(leftValue != rightValue))
	case code.OpGreaterThan:
		return vm.push(objects.NativeBoolToBooleanObject(leftValue > rightValue))
	case code.OpGreaterThanOrEqual:
		return vm.push(objects.NativeBoolToBooleanObject(leftValue >= rightValue))
	default:
		return fmt.Errorf("unknown number operator: %d", op)
	}
}

func (vm *VM) executeBangOperator() error {
	operand := vm.pop()

	switch operand {
	case objects.TRUE:
		return vm.push(objects.FALSE)
	case objects.FALSE:
		return vm.push(objects.TRUE)
	case objects.NULL:
		return vm.push(objects.TRUE)

	default:
		return vm.push(objects.FALSE)
	}
}

func (vm *VM) executeMinusOperator() error {
	operand := vm.pop()

	if !objects.IsNumber(operand.Type()) {
		return fmt.Errorf("unknown operator: -%s", operand.Type())
	}

	return vm.push(objects.WrapNumberValue(-objects.UnwrapNumberValue(operand), operand, operand))
}

func (vm *VM) executeIndexExpression(left, index objects.Object) error {
	switch {
	case left.Type() == objects.ARRAY_OBJ && index.Type() == objects.INTEGER_OBJ:
		return vm.executeArrayIndex(left, index)
	case left.Type() == objects.HASH_OBJ:
		return vm.executeHashIndex(left, index)

	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (vm *VM) executeArrayIndex(array, index objects.Object) error {
	arrayObj := array.(*objects.Array)
	i := index.(*objects.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if i < 0 || i > max {
		return vm.push(objects.NULL)
	}

	return vm.push(arrayObj.Elements[i])
}

func (vm *VM) executeHashIndex(hash, index objects.Object) error {
	hashObj := hash.(*objects.Hash)

	key, ok := index.(objects.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(objects.NULL)
	}

	return vm.push(pair.Value)
}

func (vm *VM) buildArray(startIndex, endIndex int) objects.Object {
	elements := make([]objects.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &objects.Array{Elements: elements}
}

func (vm *VM) buildHash(startIndex, endIndex int) (objects.Object, error) {
	pairs := make(map[objects.HashKey]objects.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]

		pair := objects.HashPair{Key: key, Value: value}

		hashable, ok := key.(objects.Hashable)
		if !ok {
			return nil, fmt.Errorf("key is not hashable: %T", key)
		}

		pairs[hashable.HashKey()] = pair
	}

	return &objects.Hash{Pairs: pairs}, nil
}

func (vm *VM) executeCall(numArgs int) error {
	callee := vm.stack[vm.sp-1-numArgs]

	switch callee := callee.(type) {
	case *objects.CompiledFunction:
		return vm.callFunction(callee, numArgs)
	case *objects.Builtin:
		return vm.callBuiltin(callee, numArgs)

	default:
		return fmt.Errorf("calling non-function and non-builtin")
	}
}

func (vm *VM) callFunction(fn *objects.CompiledFunction, numArgs int) error {
	if numArgs != fn.NumParameters {
		return fmt.Errorf("wrong number of arguments: got %d, want %d", numArgs, fn.NumParameters)
	}

	frame := NewFrame(fn, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *objects.Builtin, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]

	result, err := builtin.Fn(args...)
	vm.sp = vm.sp - numArgs - 1

	if err != nil {
		result = objects.NativeErrorToErrorObject(err)
	}

	return vm.push(result)
}
