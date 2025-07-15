package vm

import (
	"fmt"

	"github.com/seblkma/go-himeji/compiler"
	"github.com/seblkma/go-himeji/object"
	"github.com/seblkma/go-himeji/opcodes"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions opcodes.Instructions

	stack    []object.Object
	stackptr int // Always point to the next free slot. Top of the stack is stack[sp-1]
	// Incremented and decremented as the stack grows or shrinks.
}

func (vm *VM) push(o object.Object) error {
	if vm.stackptr >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.stackptr] = o
	vm.stackptr++
	return nil
}

func New(bytecode *compiler.ByteCode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack:    make([]object.Object, StackSize),
		stackptr: 0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.stackptr == 0 {
		return nil
	}
	return vm.stack[vm.stackptr-1]
}

func (vm *VM) Run() error {
	// Fetch instructions
	for insptr := 0; insptr < len(vm.instructions); insptr++ {
		// Decoded each instruction as opcode. Not using opcodes.Lookup for performance reasons.
		op := opcodes.Opcode(vm.instructions[insptr])

		switch op {
		case opcodes.OpConstant:
			// Decode the operands, the byte right after the opcode at insptr+1.
			// Not using opcodes.ReadOperands for performance reasons.
			constIndex := opcodes.ReadUint16(vm.instructions[insptr+1:])
			insptr += 2 // increment the correct size - the no. of bytes read to decode operands
			// next iteration the loops starts at opcode

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
