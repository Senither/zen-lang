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

type VMSettings struct {
	CaptureStdout bool
}

type VM struct {
	constants []objects.Object

	stack []objects.Object
	// The stack pointer, this always points to the next value.
	// Top of stack is stack[sp-1]
	sp int

	globals []objects.Object

	frames      []*Frame
	framesIndex int

	settings VMSettings
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &objects.CompiledFunction{OpcodeInstructions: bytecode.Instructions}
	mainClosure := &objects.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MAX_FRAMES)
	frames[0] = mainFrame

	settings := VMSettings{
		CaptureStdout: false,
	}

	return &VM{
		constants: bytecode.Constants,

		stack: make([]objects.Object, STACK_SIZE),
		sp:    0,

		globals: make([]objects.Object, GLOBALS_SIZE),

		frames:      frames,
		framesIndex: 1,

		settings: settings,
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, globals []objects.Object) *VM {
	vm := New(bytecode)

	vm.globals = globals

	return vm
}

func (vm *VM) EnableStdoutCapture() {
	vm.settings.CaptureStdout = true
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

		err := vm.executeInstructions(op, ins, ip)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vm *VM) executeInstructions(op code.Opcode, ins code.Instructions, ip int) error {
	switch op {
	case code.OpConstant:
		constIndex := code.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		return vm.push(vm.constants[constIndex])
	case code.OpClosure:
		constIndex := code.ReadUint16(ins[ip+1:])
		numFree := code.ReadUint8(ins[ip+3:])
		vm.currentFrame().ip += 3

		err := vm.pushClosure(int(constIndex), int(numFree))
		if err != nil {
			return err
		}

	case code.OpPop:
		vm.pop()

	case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv, code.OpPow, code.OpMod:
		return vm.executeBinaryOperation(op)

	case code.OpTrue:
		return vm.push(objects.TRUE)
	case code.OpFalse:
		return vm.push(objects.FALSE)

	case code.OpEqual, code.OpNotEqual, code.OpGreaterThan, code.OpGreaterThanOrEqual:
		return vm.executeComparison(op)

	case code.OpBang:
		return vm.executeBangOperator()
	case code.OpMinus:
		return vm.executeMinusOperator()

	case code.OpSetGlobal:
		globalIndex := code.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		vm.globals[globalIndex] = vm.pop()
	case code.OpGetGlobal:
		globalIndex := code.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		return vm.push(vm.globals[globalIndex])
	case code.OpSetLocal:
		localIndex := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		frame := vm.currentFrame()

		vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
	case code.OpGetLocal:
		localIndex := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		frame := vm.currentFrame()

		return vm.push(vm.stack[frame.basePointer+int(localIndex)])
	case code.OpGetFree:
		freeIndex := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		currentClosure := vm.currentFrame().closure

		err := vm.push(currentClosure.Free[freeIndex])
		if err != nil {
			return err
		}
	case code.OpCurrentClosure:
		currentClosure := vm.currentFrame().closure
		err := vm.push(currentClosure)
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
	case code.OpLoopEnd:
		// Nothing needs to happen here, this is simply a marker for
		// the end of loops that Jump operands are able to point
		// to, so we don't pop the result off the stack.

	case code.OpIndex:
		index := vm.pop()
		left := vm.pop()

		return vm.executeIndexExpression(left, index)
	case code.OpIndexAssign:
		value := vm.pop()
		index := vm.pop()
		left := vm.pop()

		return vm.executeIndexAssignment(left, index, value)
	case code.OpArray:
		numElements := int(code.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip += 2

		array := vm.buildArray(vm.sp-numElements, vm.sp)
		vm.sp -= numElements

		return vm.push(array)
	case code.OpHash:
		numElements := int(code.ReadUint16(ins[ip+1:]))
		vm.currentFrame().ip += 2

		hash, err := vm.buildHash(vm.sp-numElements, vm.sp)
		if err != nil {
			return err
		}

		vm.sp -= numElements

		return vm.push(hash)

	case code.OpNull:
		return vm.push(objects.NULL)

	case code.OpCall:
		numArgs := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		return vm.executeCall(int(numArgs))
	case code.OpReturnValue:
		returnValue := vm.pop()

		frame := vm.popFrame()
		vm.sp = frame.basePointer - 1

		return vm.push(returnValue)
	case code.OpReturn:
		frame := vm.popFrame()
		vm.sp = frame.basePointer - 1

		return vm.push(objects.NULL)
	case code.OpGetBuiltin:
		builtinIndex := code.ReadUint8(ins[ip+1:])
		vm.currentFrame().ip += 1

		definition := objects.Builtins[builtinIndex]

		return vm.push(definition.Builtin)
	case code.OpGetGlobalBuiltin:
		builtinIndex := code.ReadUint16(ins[ip+1:])
		vm.currentFrame().ip += 2

		scopeIdx := uint8(builtinIndex >> 8)
		builtIdx := uint8(builtinIndex & 0xFF)

		definition := objects.Globals[scopeIdx].Builtins[builtIdx]

		return vm.push(definition.Builtin)

	default:
		return fmt.Errorf("unsupported opcode in compiled function: %d", op)
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

func (vm *VM) pushClosure(constIndex int, numFree int) error {
	constant := vm.constants[constIndex]
	fn, ok := constant.(*objects.CompiledFunction)
	if !ok {
		return fmt.Errorf("not a function: %T", constant)
	}

	free := make([]objects.Object, numFree)
	for i := 0; i < numFree; i++ {
		free[i] = vm.stack[vm.sp-numFree+i]
	}

	vm.sp = vm.sp - numFree

	closure := &objects.Closure{Fn: fn, Free: free}

	return vm.push(closure)
}

func (vm *VM) pop() objects.Object {
	if vm.sp == 0 {
		return nil
	}

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

func (vm *VM) executeIndexAssignment(left, index, value objects.Object) error {
	switch obj := left.(type) {
	case *objects.Array:
		idx, ok := index.(*objects.Integer)
		if !ok {
			return fmt.Errorf("index operator not supported: %T", index)
		}

		if idx.Value < 0 || idx.Value >= int64(len(obj.Elements)) {
			return fmt.Errorf("array index out of bounds: %d", idx.Value)
		}

		obj.Elements[idx.Value] = value
	case *objects.Hash:
		key, ok := index.(objects.Hashable)
		if !ok {
			return fmt.Errorf("unusable as hash key: %T", index)
		}

		obj.Pairs[key.HashKey()] = objects.HashPair{Key: index, Value: value}

	default:
		return fmt.Errorf("index assignment not supported: %T", left)
	}

	return nil
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
	case *objects.Closure:
		return vm.callClosure(callee, numArgs)
	case *objects.Builtin:
		return vm.callBuiltin(callee, numArgs)

	default:
		return fmt.Errorf("calling non-function and non-builtin")
	}
}

func (vm *VM) callClosure(cl *objects.Closure, numArgs int) error {
	if numArgs != cl.Fn.NumParameters {
		return fmt.Errorf("wrong number of arguments. got %d, want %d", numArgs, cl.Fn.NumParameters)
	}

	frame := NewFrame(cl, vm.sp-numArgs)
	vm.pushFrame(frame)

	vm.sp = frame.basePointer + cl.Fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *objects.Builtin, numArgs int) error {
	args := vm.stack[vm.sp-numArgs : vm.sp]
	vm.sp = vm.sp - numArgs - 1

	for i, arg := range args {
		args[i] = WrapClosuresIfNeeded(vm, arg)
	}

	if vm.settings.CaptureStdout {
		return vm.push(captureStdoutForBuiltin(builtin, args))
	}

	result, err := builtin.Fn(args...)
	if err != nil {
		result = objects.NativeErrorToErrorObject(err)
	}

	return vm.push(result)
}
