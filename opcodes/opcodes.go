package opcodes

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2}},
	OpAdd:      {Name: "OpAdd", OperandWidths: []int{}}, // empty slice, no operand
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

func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}
