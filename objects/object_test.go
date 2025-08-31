package objects

import (
	"testing"
)

func TestStringHashKey(t *testing.T) {
	hello1 := &String{Value: "Hello, World"}
	hello2 := &String{Value: "Hello, World"}
	greeting1 := &String{Value: "Greetings, Zen"}
	greeting2 := &String{Value: "Greetings, Zen"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if greeting1.HashKey() != greeting2.HashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.HashKey() == greeting1.HashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestIntegerHashKey(t *testing.T) {
	int1 := &Integer{Value: 42}
	int2 := &Integer{Value: 42}
	int3 := &Integer{Value: 100}

	if int1.HashKey() != int2.HashKey() {
		t.Errorf("integers with same value have different hash keys")
	}

	if int1.HashKey() == int3.HashKey() {
		t.Errorf("integers with different values have same hash keys")
	}
}

func TestFloatHashKey(t *testing.T) {
	float1 := &Float{Value: 3.14159}
	float2 := &Float{Value: 3.14159}
	float3 := &Float{Value: 42.2468}

	if float1.HashKey() != float2.HashKey() {
		t.Errorf("floats with same value have different hash keys")
	}

	if float1.HashKey() == float3.HashKey() {
		t.Errorf("floats with different values have same hash keys")
	}
}

func TestBooleanHashKey(t *testing.T) {
	bool1 := &Boolean{Value: true}
	bool2 := &Boolean{Value: true}
	bool3 := &Boolean{Value: false}

	if bool1.HashKey() != bool2.HashKey() {
		t.Errorf("booleans with same value have different hash keys")
	}

	if bool1.HashKey() == bool3.HashKey() {
		t.Errorf("booleans with different values have same hash keys")
	}
}
