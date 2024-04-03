package compiler

import (
	"banana/ast"
	"banana/code"
	"banana/object"
)

type Compiler struct {
	instructions code.Instructions
	constants []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants: []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.Instructions,
		Constants: c.constants,
	}
}

type ByteCode struct {
	Instructions code.Instructions
	Constants []object.Object
}