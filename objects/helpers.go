package objects

import (
	"fmt"
	"path/filepath"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/tokens"
)

var (
	NULL  = &Null{}
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}

	BREAK    = &Break{}
	CONTINUE = &Continue{}
)

type FileDescriptorContext struct {
	Name string
	Path string
}

func NewFileDescriptorContext(path string) *FileDescriptorContext {
	name := filepath.Base(path)

	return &FileDescriptorContext{
		Name: name,
		Path: path[:len(path)-len(name)-1],
	}
}

func NativeErrorToErrorObject(err error) *Error {
	return &Error{Message: err.Error()}
}

func NewError(token tokens.Token, fileCtx *FileDescriptorContext, format string, a ...any) *Error {
	err := Error{
		Message: fmt.Sprintf(format, a...),
		Line:    token.Line,
		Column:  token.Column,
	}

	if fileCtx != nil {
		err.File = fileCtx.Name
		err.Path = fileCtx.Path
	}

	return &err
}

func NewEmptyErrorWithParent(parent *Error, token tokens.Token, fileCtx *FileDescriptorContext) *Error {
	err := NewError(token, fileCtx, "")
	err.Parent = parent

	return err
}

func IsError(obj Object) bool {
	return obj != nil && obj.Type() == ERROR_OBJ
}

func IsTruthy(obj Object) bool {
	switch obj := obj.(type) {
	case *Boolean:
		return obj.Value
	case *Null:
		return false

	default:
		return true
	}
}

func IsNumber(t ObjectType) bool {
	for _, numberType := range GetNumberTypes() {
		if t == numberType {
			return true
		}
	}

	return false
}

func GetNumberTypes() []ObjectType {
	return []ObjectType{INTEGER_OBJ, FLOAT_OBJ}
}

func WrapNumberValue(value float64, left, right Object) Object {
	if left.Type() == FLOAT_OBJ || right.Type() == FLOAT_OBJ {
		return &Float{Value: value}
	}

	if float64(int64(value)) == value {
		return &Integer{Value: int64(value)}
	}

	return &Float{Value: value}
}

func UnwrapNumberValue(obj Object) float64 {
	switch n := obj.(type) {
	case *Integer:
		return float64(n.Value)
	case *Float:
		return n.Value

	default:
		return 0
	}
}

func IsStringable(obj Object) bool {
	switch obj.(type) {
	case *Integer, *Float, *Boolean:
		return true

	default:
		return false
	}
}

func StringifyObject(obj Object) string {
	switch obj := obj.(type) {
	case *String:
		return obj.Value
	case *Integer:
		return fmt.Sprintf("%d", obj.Value)
	case *Float:
		return fmt.Sprintf("%g", obj.Value)
	case *Boolean:
		if obj.Value {
			return "true"
		}
		return "false"

	default:
		return obj.Inspect()
	}
}

func UnwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func NativeBoolToBooleanObject(input bool) *Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func Equals(left, right Object) *Boolean {
	if left.Type() != right.Type() {
		return FALSE
	}

	switch left := left.(type) {
	case *String:
		return NativeBoolToBooleanObject(left.Value == right.(*String).Value)
	case *Integer:
		return NativeBoolToBooleanObject(left.Value == right.(*Integer).Value)
	case *Float:
		return NativeBoolToBooleanObject(left.Value == right.(*Float).Value)
	case *Array:
		rightArr := right.(*Array)
		if len(left.Elements) != len(rightArr.Elements) {
			return FALSE
		}

		for i, leftElem := range left.Elements {
			rightElem := rightArr.Elements[i]
			if Equals(leftElem, rightElem) != TRUE {
				return FALSE
			}
		}

		return TRUE
	case *Hash:
		rightHash := right.(*Hash)
		if len(left.Pairs) != len(rightHash.Pairs) {
			return FALSE
		}

		for key, leftPair := range left.Pairs {
			rightPair, ok := rightHash.Pairs[key]
			if !ok || Equals(leftPair.Value, rightPair.Value) != TRUE {
				return FALSE
			}
		}

		return TRUE
	case *Null:
		return TRUE

	default:
		return NativeBoolToBooleanObject(left == right)
	}
}

func CreateImmutableHashFromEnvExports(env *Environment) *ImmutableHash {
	hashPairs := []HashPair{}

	for key, val := range env.GetExports() {
		hashPairs = append(hashPairs, HashPair{
			Key:   &String{Value: key},
			Value: val,
		})
	}

	hash := BuildImmutableHash(hashPairs...)

	return hash
}

func BuildImmutableHash(args ...HashPair) *ImmutableHash {
	pairs := make(map[HashKey]HashPair)

	for _, arg := range args {
		hash := arg.Key.(Hashable)
		pairs[hash.HashKey()] = arg
	}

	return &ImmutableHash{Value: Hash{Pairs: pairs}}
}

func WrapBuiltinFunctionInASTAwareMap(name string, fn *Builtin) HashPair {
	return HashPair{
		Key:   &String{Value: name},
		Value: BuiltinToASTAwareBuiltin(fn),
	}
}

func Copy(obj Object) Object {
	switch original := obj.(type) {
	case *String:
		return &String{Value: original.Value}
	case *Integer:
		return &Integer{Value: original.Value}
	case *Float:
		return &Float{Value: original.Value}
	case *Boolean:
		return original
	case *Null:
		return NULL
	case *Array:
		copied := &Array{Elements: make([]Object, len(original.Elements))}

		for i, element := range original.Elements {
			copied.Elements[i] = Copy(element)
		}

		return copied
	case *Hash:
		copied := &Hash{Pairs: make(map[HashKey]HashPair)}

		for key, pair := range original.Pairs {
			copied.Pairs[key] = HashPair{
				Key:   Copy(pair.Key),
				Value: Copy(pair.Value),
			}
		}

		return copied
	case *CompiledFileImport:
		copied := &CompiledFileImport{
			Name:               original.Name,
			OpcodeInstructions: make(code.Instructions, len(original.OpcodeInstructions)),
			Constants:          make([]Object, len(original.Constants)),
		}

		copy(copied.OpcodeInstructions, original.OpcodeInstructions)
		for i, constant := range original.Constants {
			copied.Constants[i] = Copy(constant)
		}

		return copied
	case *CompiledFunction:
		copied := &CompiledFunction{
			Name:               original.Name,
			OpcodeInstructions: make(code.Instructions, len(original.OpcodeInstructions)),
			NumLocals:          original.NumLocals,
			NumParameters:      original.NumParameters,
		}

		copy(copied.OpcodeInstructions, original.OpcodeInstructions)

		return copied

	default:
		panic(fmt.Sprintf("unsupported object type for copy: %s", obj.Type()))
	}
}
