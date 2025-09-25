package compiler

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/senither/zen-lang/objects"
)

func TestBytecodeString(t *testing.T) {
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
		"func add(a, b) { return a + b }; add(2, 3)",
	}

	for _, input := range tests {
		program := parse(input)

		compiler := New()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Compilation failed: %v", err)
		}

		bytecode := compiler.Bytecode()

		if bytecode.String() != bytecode.Instructions.String() {
			t.Errorf("expected String() to return instructions string, got %q", bytecode.String())
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
		"func add(a, b) { return a + b }; add(2, 3)",
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

			default:
				t.Errorf("Unsupported constant type %T", v)
			}
		}

	}
}
