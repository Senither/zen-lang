package objects

import (
	"testing"

	"github.com/senither/zen-lang/ast"
	"github.com/senither/zen-lang/tokens"
)

func TestNewEnvironment(t *testing.T) {
	env := NewEnvironment("path/to/file.zen")

	if env == nil {
		t.Fatal("Expected NewEnvironment to return a non-nil Environment")
	}

	if env.file == nil {
		t.Fatal("Expected Environment to have a non-nil FileDescriptorContext")
	}

	if env.file.Path != "path/to" {
		t.Fatalf("Expected FullPath to be 'path/to', got '%s'", env.file.Path)
	}

	if env.file.Name != "file.zen" {
		t.Fatalf("Expected Name to be 'file.zen', got '%s'", env.file.Name)
	}

	if len(env.store) != 0 {
		t.Fatalf("Expected store to be empty, got length %d", len(env.store))
	}

	if len(env.exports) != 0 {
		t.Fatalf("Expected exports to be empty, got length %d", len(env.exports))
	}
}

func TestNewEnclosedEnvironment(t *testing.T) {
	outerEnv := NewEnvironment("outer/file.zen")
	enclosedEnv := NewEnclosedEnvironment(outerEnv)

	if enclosedEnv.outer != outerEnv {
		t.Fatal("Expected enclosed Environment to reference the correct outer Environment")
	}
}

func TestHas(t *testing.T) {
	env := NewEnvironment(nil)
	env.store["var1"] = EnvironmentStateItem{value: nil, mutable: true}

	if !env.Has("var1") {
		t.Fatal("Expected Has to return true for existing variable 'var1'")
	}

	if env.Has("var2") {
		t.Fatal("Expected Has to return false for non-existing variable 'var2'")
	}
}

func TestGet(t *testing.T) {
	env := NewEnvironment(nil)
	expectedValue := &Integer{Value: 42}
	env.store["var1"] = EnvironmentStateItem{value: expectedValue, mutable: true}

	value, ok := env.Get("var1")
	if !ok {
		t.Fatal("Expected Get to return true for existing variable 'var1'")
	}

	if value != expectedValue {
		t.Fatalf("Expected Get to return the correct value for 'var1', got %v", value)
	}

	_, ok = env.Get("var2")
	if ok {
		t.Fatal("Expected Get to return false for non-existing variable 'var2'")
	}
}

func TestGetStateItem(t *testing.T) {
	env := NewEnvironment(nil)
	expectedValue := &Integer{Value: 100}
	env.store["var1"] = EnvironmentStateItem{value: expectedValue, mutable: false}

	item, ok := env.GetStateItem("var1")
	if !ok {
		t.Fatal("Expected GetStateItem to return true for existing variable 'var1'")
	}

	if item.value != expectedValue {
		t.Fatalf("Expected GetStateItem to return the correct value for 'var1', got %v", item.value)
	}

	_, ok = env.GetStateItem("var2")
	if ok {
		t.Fatal("Expected GetStateItem to return false for non-existing variable 'var2'")
	}
}

func TestSet(t *testing.T) {
	env := NewEnvironment(nil)
	val := &Integer{Value: 10}

	result := env.Set(nil, "var1", val, true)
	if result != val {
		t.Fatal("Expected Set to return true when setting a new variable 'var1'")
	}

	retrievedVal, ok := env.GetStateItem("var1")
	if !ok || retrievedVal.value != val {
		t.Fatal("Expected Get to return the correct value for 'var1' after Set")
	}

	if !retrievedVal.mutable {
		t.Fatal("Expected 'var1' to be mutable after Set")
	}

	result = env.Set(nil, "var2", val, false)
	if result != val {
		t.Fatal("Expected Set to return true when setting a new variable 'var1'")
	}

	retrievedVal, ok = env.GetStateItem("var2")
	if !ok || retrievedVal.value != val {
		t.Fatal("Expected Get to return the correct value for 'var2' after Set")
	}

	if retrievedVal.mutable {
		t.Fatal("Expected 'var2' to be immutable after Set")
	}
}

func TestAssignWithMutable(t *testing.T) {
	env := NewEnvironment(nil)
	val := &Integer{Value: 20}
	newVal := &Integer{Value: 30}
	env.Set(nil, "var1", val, true)

	result := env.Assign(nil, "var1", newVal)
	if result != newVal {
		t.Fatal("Expected Assign to return the new value when assigning to mutable variable 'var1'")
	}

	retrievedVal, ok := env.GetStateItem("var1")
	if !ok || retrievedVal.value != newVal {
		t.Fatal("Expected Get to return the updated value for 'var1' after Assign")
	}
}

func TestAssignWithImmutable(t *testing.T) {
	env := NewEnvironment(nil)
	val := &Integer{Value: 20}
	newVal := &Integer{Value: 30}
	env.Set(nil, "var1", val, false)

	result := env.Assign(&ast.AssignmentExpression{
		Token: tokens.Token{Line: 1, Column: 1, Literal: "="},
	}, "var1", newVal)

	errorObj, ok := result.(*Error)
	if !ok {
		t.Fatal("Expected Assign to return an Error when assigning to immutable variable 'var1'")
	}

	expectedMessage := "cannot modify immutable variable: var1"
	if errorObj.Message != expectedMessage {
		t.Fatalf("Expected error message '%s', got '%s'", expectedMessage, errorObj.Message)
	}
}

func TestSetImmutableForcefully(t *testing.T) {
	env := NewEnvironment(nil)
	val := &Integer{Value: 50}
	env.Set(nil, "var1", val, false)

	retrievedVal, ok := env.GetStateItem("var1")
	if !ok || retrievedVal.value != val {
		t.Fatal("Expected Get to return the correct value for 'var1' after Set")
	}

	if retrievedVal.mutable {
		t.Fatal("Expected 'var1' to be immutable after Set")
	}

	newVal := &Integer{Value: 60}
	result := env.SetImmutableForcefully("var1", newVal)

	if result != newVal {
		t.Fatal("Expected SetImmutableForcefully to return the new value when updating 'var1'")
	}

	retrievedVal, ok = env.GetStateItem("var1")
	if !ok || retrievedVal.value != newVal {
		t.Fatal("Expected Get to return the updated value for 'var1' after SetImmutableForcefully")
	}
}

func TestExports(t *testing.T) {
	env := NewEnvironment(nil)

	err := env.Export(&Function{Name: &ast.Identifier{Value: "myFunction"}})
	if err != nil {
		t.Fatalf("Expected Export to succeed, got error: %v", err)
	}

	exports := env.GetExports()
	if len(exports) != 1 {
		t.Fatalf("Expected 1 export, got %d", len(exports))
	}

	if _, ok := exports["myFunction"]; !ok {
		t.Fatal("Expected 'myFunction' to be in exports")
	}

	err = env.Export(&Function{})
	if err == nil {
		t.Fatal("Expected Export to fail for unnamed function")
	}

	if err.Error() != "cannot export unnamed function" {
		t.Fatalf("Unexpected error message for unnamed function: %s", err.Error())
	}

	err = env.Export(&Integer{Value: 10})
	if err == nil {
		t.Fatal("Expected Export to fail for non-function object")
	}

	if err.Error() != "cannot export object of type INTEGER" {
		t.Fatalf("Unexpected error message for non-function object: %s", err.Error())
	}
}

func TestGetFileDescriptorContext(t *testing.T) {
	env := NewEnvironment("dir/sample.zen")
	fileDesc := env.GetFileDescriptorContext()

	if fileDesc.Path != "dir" {
		t.Fatalf("Expected Path to be 'dir', got '%s'", fileDesc.Path)
	}

	if fileDesc.Name != "sample.zen" {
		t.Fatalf("Expected Name to be 'sample.zen', got '%s'", fileDesc.Name)
	}

	enclosedEnv := NewEnclosedEnvironment(env)
	enclosedFileDesc := enclosedEnv.GetFileDescriptorContext()

	if enclosedFileDesc.Path != "dir" {
		t.Fatalf("Expected Path to be 'dir' in enclosed environment, got '%s'", enclosedFileDesc.Path)
	}

	if enclosedFileDesc.Name != "sample.zen" {
		t.Fatalf("Expected Name to be 'sample.zen' in enclosed environment, got '%s'", enclosedFileDesc.Name)
	}
}
