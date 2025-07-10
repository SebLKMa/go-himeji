package opcodes

import "testing"

// GOFLAGS="-count=1" go test -run TestMake
func TestMake(t *testing.T) {
	testInputs := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}}, // fmt.Printf("Byte as hex-> %x,%x\n", byte(255), byte(254)) -> ff,fe
	}

	for _, ti := range testInputs {
		instruction := Make(ti.op, ti.operands...)

		if len(instruction) != len(ti.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d", len(ti.expected), len(instruction))
		}

		for i, b := range ti.expected {
			if instruction[i] != ti.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d", i, b, instruction[i])
			}
		}
	}

}
