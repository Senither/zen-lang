package vm

import (
	"github.com/senither/zen-lang/code"
	"github.com/senither/zen-lang/objects"
)

type Frame struct {
	obj objects.CompiledInstructionsObject
	ip  int
}

func NewFrame(obj objects.CompiledInstructionsObject) *Frame {
	return &Frame{obj: obj, ip: -1}
}

func (f *Frame) Instructions() code.Instructions {
	return f.obj.Instructions()
}
