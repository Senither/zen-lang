package compiler

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestBytecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input: "true == false",
			expected: []string{
				"0000x00000000 OpTrue",
				"0000x00000001 OpFalse",
				"0000x00000002 OpEqual",
				"0000x00000003 OpPop",
			},
		},
		{
			input: "1 + 2",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpAdd",
				"0000x00000007 OpPop",
			},
		},
		{
			input: "2.5 + 3f",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpAdd",
				"0000x00000007 OpPop",
			},
		},
		{
			input: "1 + 2 * 3 - 4 / 5 % 6",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpConstant 2",
				"0000x00000009 OpMul",
				"0000x00000010 OpAdd",
				"0000x00000011 OpConstant 3",
				"0000x00000014 OpConstant 4",
				"0000x00000017 OpDiv",
				"0000x00000018 OpConstant 5",
				"0000x00000021 OpMod",
				"0000x00000022 OpSub",
				"0000x00000023 OpPop",
			},
		},
		{
			input: "'Hello, World!'",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpPop",
			},
		},
		{
			input: "[1, 2, 3]",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpConstant 2",
				"0000x00000009 OpArray 3",
				"0000x00000012 OpPop",
			},
		},
		{
			input: "{ 'key': 'value' }",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpHash 2",
				"0000x00000009 OpPop",
			},
		},
		{
			input: "var a = 10; var b = 20; a + b",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpSetGlobal 0",
				"0000x00000006 OpConstant 1",
				"0000x00000009 OpSetGlobal 1",
				"0000x00000012 OpGetGlobal 0",
				"0000x00000015 OpGetGlobal 1",
				"0000x00000018 OpAdd",
				"0000x00000019 OpPop",
			},
		},
		{
			input: "var arr = [1, 2, 3]; arr[0] + arr[1]",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpConstant 2",
				"0000x00000009 OpArray 3",
				"0000x00000012 OpSetGlobal 0",
				"0000x00000015 OpGetGlobal 0",
				"0000x00000018 OpConstant 3",
				"0000x00000021 OpIndex",
				"0000x00000022 OpGetGlobal 0",
				"0000x00000025 OpConstant 4",
				"0000x00000028 OpIndex",
				"0000x00000029 OpAdd",
				"0000x00000030 OpPop",
			},
		},
		{
			input: "var obj = { 'key': 'value' }; obj['key']",
			expected: []string{
				"0000x00000000 OpConstant 0",
				"0000x00000003 OpConstant 1",
				"0000x00000006 OpHash 2",
				"0000x00000009 OpSetGlobal 0",
				"0000x00000012 OpGetGlobal 0",
				"0000x00000015 OpConstant 2",
				"0000x00000018 OpIndex",
				"0000x00000019 OpPop",
			},
		},
		{
			input: "if (false) { 10 } else if (true) { 20 } else { 30 }",
			expected: []string{
				"0000x00000000 OpFalse",
				"0000x00000001 OpJumpNotTruthy 10",
				"0000x00000004 OpConstant 0",
				"0000x00000007 OpJump 23",
				"0000x00000010 OpTrue",
				"0000x00000011 OpJumpNotTruthy 20",
				"0000x00000014 OpConstant 1",
				"0000x00000017 OpJump 23",
				"0000x00000020 OpConstant 2",
				"0000x00000023 OpPop",
			},
		},
		{
			input: "func () { return 5 + 10 }",
			expected: []string{
				"0001x00000000 OpConstant 0",
				"0001x00000003 OpConstant 1",
				"0001x00000006 OpAdd",
				"0001x00000007 OpReturnValue",
				"0000x00000000 OpClosure 2 0",
				"0000x00000004 OpPop",
			},
		},
		{
			input: "func () { 5 + 10 }",
			expected: []string{
				"0001x00000000 OpConstant 0",
				"0001x00000003 OpConstant 1",
				"0001x00000006 OpAdd",
				"0001x00000007 OpReturnValue",
				"0000x00000000 OpClosure 2 0",
				"0000x00000004 OpPop",
			},
		},
	}

	for _, test := range tests {
		program := parse(test.input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Compilation failed: %v", err)
		}

		bytecode := strings.TrimRight(compiler.Bytecode().String(), "\n")
		expected := strings.Join(test.expected, "\n")
		if bytecode != expected {
			t.Errorf("Stringified bytecode mismatch.\ngot:\n%s\nwant:\n%s", bytecode, expected)
		}
	}
}

func TestBytecodeSerializeDeserialize(t *testing.T) {
	tests := []string{
		"1 + 2",
		"2.5 + 3f",
		"'Hello, World!'",
		"[1, 2, 3]",
		"{ 'key': 'value' }",
		"1 + 2 * 3 - 4 / 5 % 6",
		"var a = 10; var b = 20; a + b",
		"var mut x = 5; x = x + 10; x",
		"var arr = [1, 2, 3]; arr[0] + arr[1]",
		"var obj = { 'key': 'value' }; obj['key']",
		"if (true) { 10 } else if (false) { 20 } else { 30 }",
		"while (false) { 1 }",
		"func () { return 5 + 10 }",
		"func () { 5 + 10 }",
		"func (a) { a + 10 }(5)",
		"func (a, b) { a + b }(5, 10)",
	}

	for _, input := range tests {
		program := parse(input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Compilation failed: %v", err)
		}

		bytecode := compiler.Bytecode()
		serialized := bytecode.Serialize()
		deserialized, err := Deserialize(serialized)
		if err != nil {
			t.Fatalf("Deserialize failed: %v", err)
		}

		if !bytes.Equal(deserialized.Instructions, bytecode.Instructions) {
			t.Errorf(
				"Instructions mismatch after Deserialize.\ninput: %s\ngot:\n%v\nwant:\n%v",
				input, deserialized.Instructions, bytecode.Instructions,
			)
		}

		if len(deserialized.Constants) != len(bytecode.Constants) {
			t.Fatalf("Constants length mismatch. got %d, want %d", len(deserialized.Constants), len(bytecode.Constants))
		}

		for i, constantObject := range bytecode.Constants {
			deserializedConstant := deserialized.Constants[i]
			if reflect.TypeOf(constantObject) != reflect.TypeOf(deserializedConstant) {
				t.Errorf("Constant %d type mismatch. got %T, want %T", i, deserializedConstant, constantObject)
			}

			switch v := constantObject.(type) {
			case *objects.Null:
				// nothing to check
			case *objects.Integer:
				if v.Value != deserializedConstant.(*objects.Integer).Value {
					t.Errorf(
						"Integer constant %d value mismatch. got %d, want %d",
						i, deserializedConstant.(*objects.Integer).Value, v.Value,
					)
				}
			case *objects.Float:
				if v.Value != deserializedConstant.(*objects.Float).Value {
					t.Errorf(
						"Float constant %d value mismatch. got %f, want %f",
						i, deserializedConstant.(*objects.Float).Value, v.Value,
					)
				}
			case *objects.Boolean:
				if v.Value != deserializedConstant.(*objects.Boolean).Value {
					t.Errorf(
						"Boolean constant %d value mismatch. got %v, want %v",
						i, deserializedConstant.(*objects.Boolean).Value, v.Value,
					)
				}
			case *objects.String:
				if v.Value != deserializedConstant.(*objects.String).Value {
					t.Errorf(
						"String constant %d value mismatch. got %v, want %v",
						i, deserializedConstant.(*objects.String).Value, v.Value,
					)
				}
			case *objects.CompiledFunction:
				if !reflect.DeepEqual(v, deserializedConstant) {
					t.Errorf(
						"CompiledFunction constant %d value mismatch. got %v, want %v",
						i, deserializedConstant, v,
					)
				}

			default:
				t.Errorf("Unsupported constant type %T", v)
			}
		}

	}
}
