package compiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

const (
	SERIES_HEADER  = "ZENB"
	SERIES_VERSION = uint8(1)

	NULL_CONST = uint8(1)

	INTEGER_CONST = uint8(10)
	FLOAT_CONST   = uint8(11)
	BOOLEAN_CONST = uint8(12)
	STRING_CONST  = uint8(13)

	COMPILED_FUNCTION_CONST = uint8(20)
	COMPILED_IMPORT_CONST   = uint8(21)
)

type Bytecode struct {
	Instructions code.Instructions
	Constants    []objects.Object
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (b *Bytecode) String() string {
	return b.stringFromDepth(0)
}

func (b *Bytecode) stringFromDepth(depth int) string {
	var out bytes.Buffer

	closureDef, _ := code.Lookup(code.OpClosure)
	importDef, _ := code.Lookup(code.OpImport)

	i := 0
	for i < len(b.Instructions) {
		def, err := code.Lookup(code.Opcode(b.Instructions[i]))
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		switch def {
		case closureDef:
			constIndex := binary.BigEndian.Uint16(b.Instructions[i+1 : i+3])
			constant := b.Constants[constIndex]

			if fn, ok := constant.(*objects.CompiledFunction); ok {
				b.printCompiledFunction(&out, def, fn, depth+1)
			}
		case importDef:
			constIndex := binary.BigEndian.Uint16(b.Instructions[i+1 : i+3])
			constant := b.Constants[constIndex]

			if imp, ok := constant.(*objects.CompiledFileImport); ok {
				nestedBytecode := &Bytecode{
					Instructions: imp.OpcodeInstructions,
					Constants:    imp.Constants,
				}

				out.WriteString(nestedBytecode.stringFromDepth(depth + 1))
			}
		}

		i += writeInstructionsToBuffer(&out, i, depth, def, b.Instructions)
	}

	return out.String()
}

func (b *Bytecode) printCompiledFunction(out *bytes.Buffer, closureDef *code.Definition, fn *objects.CompiledFunction, depth int) {
	ins := fn.Instructions()
	if len(ins) == 0 {
		fmt.Fprintf(out, "ERROR: compiled function has no instructions\n")
		return
	}

	x := 0
	for x < len(ins) {
		fnDef, err := code.Lookup(code.Opcode(ins[x]))
		if err != nil {
			fmt.Fprintf(out, "ERROR: %s\n", err)
			break
		}

		if fnDef == closureDef {
			constIndex := binary.BigEndian.Uint16(ins[x+1 : x+3])
			constant := b.Constants[constIndex]

			if nestedFn, ok := constant.(*objects.CompiledFunction); ok {
				b.printCompiledFunction(out, fnDef, nestedFn, depth+1)
			}
		}

		x += writeInstructionsToBuffer(out, x, depth, fnDef, ins)
	}
}

func writeInstructionsToBuffer(out *bytes.Buffer, index, scope int, def *code.Definition, ins code.Instructions) int {
	operands, read := code.ReadOperands(def, ins[index+1:])
	fmt.Fprintf(out, "%04dx%08d %s\n", scope, index, ins.FormatInstruction(def, operands))

	return read + 1
}

func (b *Bytecode) Serialize() []byte {
	buf := &bytes.Buffer{}
	write := func(data any) { binary.Write(buf, binary.BigEndian, data) }

	buf.Write([]byte(SERIES_HEADER))
	buf.WriteByte(SERIES_VERSION)

	write(uint32(len(b.Instructions)))
	buf.Write(b.Instructions)

	b.writeSerializedConstants(buf, write, b.Constants)

	return buf.Bytes()
}

func (b *Bytecode) writeSerializedConstants(buf *bytes.Buffer, write func(data any), constants []objects.Object) {
	write(uint32(len(constants)))
	for _, c := range constants {
		switch v := c.(type) {
		case *objects.Null:
			buf.WriteByte(NULL_CONST)
		case *objects.Integer:
			buf.WriteByte(INTEGER_CONST)
			write(v.Value)
		case *objects.Float:
			buf.WriteByte(FLOAT_CONST)
			write(math.Float64bits(v.Value))
		case *objects.Boolean:
			buf.WriteByte(BOOLEAN_CONST)
			if v.Value {
				buf.WriteByte(1)
			} else {
				buf.WriteByte(0)
			}
		case *objects.String:
			buf.WriteByte(STRING_CONST)
			write(uint32(len(v.Value)))
			buf.WriteString(v.Value)
		case *objects.CompiledFunction:
			buf.WriteByte(COMPILED_FUNCTION_CONST)
			write(uint32(len(v.Name)))
			buf.WriteString(v.Name)
			write(uint32(v.NumLocals))
			write(uint32(v.NumParameters))
			write(uint32(len(v.Instructions())))
			write(v.Instructions())
		case *objects.CompiledFileImport:
			buf.WriteByte(COMPILED_IMPORT_CONST)
			write(uint32(len(v.Name)))
			buf.WriteString(v.Name)
			write(uint32(len(v.Instructions())))
			write(v.Instructions())
			b.writeSerializedConstants(buf, write, v.Constants)

		default:
			panic(fmt.Sprintf("unsupported constant type: %T", v))
		}
	}
}

func Deserialize(data []byte) (*Bytecode, error) {
	r := bytes.NewReader(data)
	if err := verifyBytecodeHeaders(r); err != nil {
		return nil, err
	}

	read := func(data any) error { return binary.Read(r, binary.BigEndian, data) }

	// Reads the instructions set by first getting the length, then reading
	// that many bytes into a byte slice to form the Instructions field
	var instructionLength uint32
	if err := read(&instructionLength); err != nil {
		return nil, err
	}

	ins := make([]byte, instructionLength)
	if _, err := io.ReadFull(r, ins); err != nil {
		return nil, err
	}

	consts, err := deserializeConstants(r, read)
	if err != nil {
		return nil, err
	}

	return &Bytecode{
		Instructions: code.Instructions(ins),
		Constants:    consts,
	}, nil
}

func deserializeConstants(r *bytes.Reader, read func(data any) error) ([]objects.Object, error) {
	var constCount uint32
	if err := read(&constCount); err != nil {
		return nil, err
	}

	consts := make([]objects.Object, 0, constCount)
	for i := uint32(0); i < constCount; i++ {
		tag, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		switch tag {
		case NULL_CONST:
			consts = append(consts, &objects.Null{})
		case INTEGER_CONST:
			var v int64
			if err := read(&v); err != nil {
				return nil, err
			}
			consts = append(consts, &objects.Integer{Value: v})
		case FLOAT_CONST:
			var bits uint64
			if err := read(&bits); err != nil {
				return nil, err
			}
			consts = append(consts, &objects.Float{Value: math.Float64frombits(bits)})
		case BOOLEAN_CONST:
			val, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			consts = append(consts, &objects.Boolean{Value: val == 1})
		case STRING_CONST:
			var strLen uint32
			if err := read(&strLen); err != nil {
				return nil, err
			}

			str := make([]byte, strLen)
			if _, err := io.ReadFull(r, str); err != nil {
				return nil, err
			}

			consts = append(consts, &objects.String{Value: string(str)})
		case COMPILED_FUNCTION_CONST:
			var nameLen uint32
			if err := read(&nameLen); err != nil {
				return nil, err
			}

			nameBytes := make([]byte, nameLen)
			if _, err := io.ReadFull(r, nameBytes); err != nil {
				return nil, err
			}

			var numLocals uint32
			if err := read(&numLocals); err != nil {
				return nil, err
			}

			var numParameters uint32
			if err := read(&numParameters); err != nil {
				return nil, err
			}

			var insLen uint32
			if err := read(&insLen); err != nil {
				return nil, err
			}

			instructions := make([]byte, insLen)
			if _, err := io.ReadFull(r, instructions); err != nil {
				return nil, err
			}

			consts = append(consts, &objects.CompiledFunction{
				Name:               string(nameBytes),
				NumLocals:          int(numLocals),
				NumParameters:      int(numParameters),
				OpcodeInstructions: instructions,
			})
		case COMPILED_IMPORT_CONST:
			var nameLen uint32
			if err := read(&nameLen); err != nil {
				return nil, err
			}

			nameBytes := make([]byte, nameLen)
			if _, err := io.ReadFull(r, nameBytes); err != nil {
				return nil, err
			}

			var insLen uint32
			if err := read(&insLen); err != nil {
				return nil, err
			}

			instructions := make([]byte, insLen)
			if _, err := io.ReadFull(r, instructions); err != nil {
				return nil, err
			}

			nestedConst, err := deserializeConstants(r, read)
			if err != nil {
				return nil, err
			}

			consts = append(consts, &objects.CompiledFileImport{
				Name:               string(nameBytes),
				OpcodeInstructions: instructions,
				Constants:          nestedConst,
			})

		default:
			return nil, fmt.Errorf("unknown constant tag: %d", tag)
		}
	}

	return consts, nil
}

func verifyBytecodeHeaders(r *bytes.Reader) error {
	header := make([]byte, 4)
	if _, err := io.ReadFull(r, header); err != nil {
		return fmt.Errorf("unrecognized bytecode header: %w", err)
	}

	if string(header) != SERIES_HEADER {
		return fmt.Errorf("invalid bytecode header: %q", string(header))
	}

	version, err := r.ReadByte()
	if err != nil {
		return err
	}

	if version != SERIES_VERSION {
		return fmt.Errorf("unsupported bytecode version: %d", version)
	}

	return nil
}
