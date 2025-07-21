package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/seblkma/go-himeji/compiler"
	"github.com/seblkma/go-himeji/object"
	"github.com/seblkma/go-himeji/vm"
	// naming conflicts with go/token
)

func init() {
	// Register the concrete types that will be stored by gob.
	//gob.Register(object.Integer{})
	// See
	// https://stackoverflow.com/questions/54766528/gob-decode-cannot-decode-interface-after-register-type
	gob.Register(&object.Integer{})
}

const HIMEJI_CODES_BIN = "../himeji/codes.bin"

func main() {
	serializedData, err := os.ReadFile(HIMEJI_CODES_BIN)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes from %s\n", len(serializedData), HIMEJI_CODES_BIN)

	// Check raw values
	//v := reflect.ValueOf(serializedData)
	//fmt.Printf("serialized bytecode: %+v\n", v)

	// Deserialize
	var bc *compiler.ByteCode
	decoder := gob.NewDecoder(bytes.NewReader(serializedData))
	if err := decoder.Decode(&bc); err != nil {
		fmt.Println("Error decoding:", err)
		return
	}
	fmt.Printf("Deserialized bytecode: %+v\n", bc)

	machine := vm.New(bc)
	err = machine.Run()
	if err != nil {
		fmt.Printf("Woops! Executing bytecode failed:\n %s\n", err)
		os.Exit(1)
	}

	stackTop := machine.StackTop()
	fmt.Printf("Result: %s\n", stackTop.Inspect())

}
