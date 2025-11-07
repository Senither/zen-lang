package vm

import (
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type Frame struct {
	closure     *objects.Closure
	ip          int
	basePointer int
}

func NewFrame(obj *objects.Closure, basePointer int) *Frame {
	return &Frame{
		closure:     obj,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.closure.Instructions()
}
