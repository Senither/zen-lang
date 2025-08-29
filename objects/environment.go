package objects

import "fmt"

type Environment struct {
	store map[string]EnvironmentStateItem
	outer *Environment
}

type EnvironmentStateItem struct {
	value   Object
	mutable bool
}

func NewEnvironment() *Environment {
	return &Environment{
		store: make(map[string]EnvironmentStateItem),
		outer: nil,
	}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
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

func (e *Environment) Set(name string, val Object, mutable bool) Object {
	item, ok := e.GetStateItem(name)
	if ok {
		if !item.mutable {
			return &Error{Message: fmt.Sprintf("Cannot modify immutable variable '%s'", name)}
		}

		mutable = item.mutable
	}

	e.store[name] = EnvironmentStateItem{
		value:   val,
		mutable: mutable,
	}

	return val
}

func (e *Environment) SetImmutableForcefully(name string, val Object) Object {
	e.store[name] = EnvironmentStateItem{
		value:   val,
		mutable: false,
	}

	return val
}
