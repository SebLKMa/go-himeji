package opcodes

import (
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2}},
}

func Lookup(op byte) (*Definition, error) {
	def, found := definitions[Opcode(op)]
	if !found {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int) []byte {
	// Get the opcode definition
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	// Get the instruction length from opcode definition
	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// Make the instruction
	// First byte of insruction is the opcode
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	fmt.Printf("instruction: %v\n", instruction)
	// Make the rest of the instruction
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		fmt.Printf("i:%d, width:%d, offset:%d\n", i, width, offset)
		switch width {
		case 2: // operands starts at 2
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += width
	}
	fmt.Printf("instruction: %v\n", instruction)
	return instruction
}
