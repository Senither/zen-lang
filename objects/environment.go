package objects

import (
	"fmt"

	"github.com/senither/zen-lang/ast"
)

type Environment struct {
	store   map[string]EnvironmentStateItem
	exports map[string]Object
	outer   *Environment
	file    *FileDescriptorContext
}

type EnvironmentStateItem struct {
	value   Object
	mutable bool
}

func NewEnvironment(fullFilePath interface{}) *Environment {
	env := &Environment{
		store:   make(map[string]EnvironmentStateItem),
		exports: make(map[string]Object),
		outer:   nil,
	}

	if _, ok := fullFilePath.(string); !ok {
		return env
	}

	fullPath := fullFilePath.(string)
	env.file = NewFileDescriptorContext(fullPath)

	return env
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment(nil)
	env.outer = outer

	return env
}

func (e *Environment) IsRoot() bool {
	return e.outer == nil
}

func (e *Environment) Has(name string) bool {
	_, ok := e.GetStateItem(name)
	return ok
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.GetStateItem(name)
	return val.value, ok
}

func (e *Environment) GetStateItem(name string) (EnvironmentStateItem, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.GetStateItem(name)
	}

	return val, ok
}

func (e *Environment) Set(node ast.Node, name string, val Object, mutable bool) Object {
	item, ok := e.GetStateItem(name)
	if ok {
		if !item.mutable {
			return NewError(
				node.GetToken(),
				e.GetFileDescriptorContext(),
				"cannot modify immutable variable: %s",
				name,
			)
		}

		mutable = item.mutable
	}

	e.store[name] = EnvironmentStateItem{
		value:   val,
		mutable: mutable,
	}

	return val
}

func (e *Environment) Assign(node ast.Node, name string, new Object) Object {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Assign(node, name, new)
	}

	if !val.mutable {
		return NewError(
			node.GetToken(),
			e.GetFileDescriptorContext(),
			"cannot modify immutable variable: %s",
			name,
		)
	}

	val.value = new
	e.store[name] = val

	return val.value
}

func (e *Environment) SetImmutableForcefully(name string, val Object) Object {
	e.store[name] = EnvironmentStateItem{
		value:   val,
		mutable: false,
	}

	return val
}

func (e *Environment) Export(val Object) error {
	if e.outer != nil {
		return e.outer.Export(val)
	}

	switch val := val.(type) {
	case *Function:
		if val.Name == nil {
			return fmt.Errorf("cannot export unnamed function")
		}

		e.exports[val.Name.Value] = val

	default:
		return fmt.Errorf("cannot export object of type %s", val.Type())
	}

	return nil
}

func (e *Environment) GetExports() map[string]Object {
	return e.exports
}

func (e *Environment) GetFileDescriptorContext() *FileDescriptorContext {
	if e.file == nil && e.outer != nil {
		return e.outer.GetFileDescriptorContext()
	}

	return e.file
}
