package objects

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/code"
)

type ObjectType string

const (
	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	STRING_OBJ  = "STRING"
	INTEGER_OBJ = "INTEGER"
	FLOAT_OBJ   = "FLOAT"
	BOOLEAN_OBJ = "BOOLEAN"

	ARRAY_OBJ          = "ARRAY"
	HASH_OBJ           = "HASH"
	IMMUTABLE_HASH_OBJ = "IMMUTABLE_HASH"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	BREAK_OBJ    = "BREAK"
	CONTINUE_OBJ = "CONTINUE"

	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"

	COMPILED_FUNCTION_OBJ = "COMPILED_FUNCTION"
	CLOSURE_OBJ           = "CLOSURE"

	IMPORTED_CLOSURE_OBJ     = "IMPORTED_CLOSURE"
	COMPILED_FILE_IMPORT_OBJ = "COMPILED_FILE_IMPORT"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type CompiledInstructionsObject interface {
	Type() ObjectType
	Inspect() string
	Instructions() code.Instructions
}

type Callable interface {
	Call(args ...Object) Object
	ParametersCount() int
	Inspect() string
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf("%v", s.Value) }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey { return HashKey{Type: i.Type(), Value: uint64(i.Value)} }

type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%f", f.Value) }
func (f *Float) HashKey() HashKey {
	h := fnv.New64a()
	fmt.Fprintf(h, "%f", f.Value)

	return HashKey{Type: f.Type(), Value: h.Sum64()}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64 = 0
	if b.Value {
		value = 1
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	keys := make([]HashKey, 0, len(h.Pairs))
	for k := range h.Pairs {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Type == keys[j].Type {
			return keys[i].Value < keys[j].Value
		}
		return keys[i].Type < keys[j].Type
	})

	pairs := []string{}
	for _, k := range keys {
		pair := h.Pairs[k]
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type ImmutableHash struct {
	Value Hash
}

func (h *ImmutableHash) Type() ObjectType { return IMMUTABLE_HASH_OBJ }
func (h *ImmutableHash) Inspect() string  { return h.Value.Inspect() }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
	Path    string
	File    string
	Line    int
	Column  int
	Parent  *Error
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	var out bytes.Buffer

	if e.Parent != nil {
		out.WriteString(e.Parent.Inspect())
	}

	file := "<unknown>"
	if e.Path != "" && e.File != "" {
		file = e.Path + string(os.PathSeparator) + e.File
	}

	if e.Parent == nil {
		out.WriteString(fmt.Sprintf("%s\n    at %s:%d:%d", e.Message, file, e.Line, e.Column))
	} else {
		out.WriteString(fmt.Sprintf("\n    at %s:%d:%d", file, e.Line, e.Column))
	}

	seen := make(map[string]struct{})
	unique := []string{}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if _, exists := seen[line]; !exists {
			seen[line] = struct{}{}
			unique = append(unique, line)
		}
	}

	return strings.Join(unique, "\n")
}

type Break struct{}

func (b *Break) Type() ObjectType { return BREAK_OBJ }
func (b *Break) Inspect() string  { return "break" }

type Continue struct{}

func (c *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (c *Continue) Inspect() string  { return "continue" }

type Function struct {
	Name       *ast.Identifier
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	if f.Name != nil {
		out.WriteString(f.Name.String())
	}

	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuiltinFunction func(args ...Object) (Object, error)
type BuiltinDefinition struct {
	Name    string
	Builtin *Builtin
}

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

type ASTAwareBuiltinFunction func(node *ast.CallExpression, env *Environment, args ...Object) Object
type ASTAwareBuiltin struct {
	Fn ASTAwareBuiltinFunction
}

func (b *ASTAwareBuiltin) Type() ObjectType { return BUILTIN_OBJ }
func (b *ASTAwareBuiltin) Inspect() string  { return "builtin function" }

type CompiledFunction struct {
	Name               string
	OpcodeInstructions code.Instructions
	NumLocals          int
	NumParameters      int
}

func (cf *CompiledFunction) Type() ObjectType { return COMPILED_FUNCTION_OBJ }
func (cf *CompiledFunction) Inspect() string {
	if len(cf.Name) > 0 {
		return fmt.Sprintf("CompiledFunction[%s|%p]", cf.Name, cf)
	} else {
		return fmt.Sprintf("CompiledFunction[%p]", cf)
	}
}
func (cf *CompiledFunction) Instructions() code.Instructions { return cf.OpcodeInstructions }

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() ObjectType                { return CLOSURE_OBJ }
func (c *Closure) Inspect() string                 { return fmt.Sprintf("Closure[%p]", c) }
func (c *Closure) Instructions() code.Instructions { return c.Fn.OpcodeInstructions }

type ImportedClosure struct {
	Closure   *Closure
	Constants []Object
	Globals   []Object
}

func (ic *ImportedClosure) Type() ObjectType                { return IMPORTED_CLOSURE_OBJ }
func (ic *ImportedClosure) Inspect() string                 { return fmt.Sprintf("ImportedClosure[%p]", ic) }
func (ic *ImportedClosure) Instructions() code.Instructions { return ic.Closure.Fn.OpcodeInstructions }

type CompiledFileImport struct {
	Name               string
	OpcodeInstructions code.Instructions
	Constants          []Object
}

func (cfi *CompiledFileImport) Type() ObjectType { return COMPILED_FILE_IMPORT_OBJ }
func (cfi *CompiledFileImport) Inspect() string {
	return fmt.Sprintf("CompiledFileImport[%s|%p]", cfi.Name, cfi)
}
func (cfi *CompiledFileImport) Instructions() code.Instructions { return cfi.OpcodeInstructions }
