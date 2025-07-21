package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/seblkma/go-himeji/compiler"
	"github.com/seblkma/go-himeji/lexer"
	"github.com/seblkma/go-himeji/object"
	"github.com/seblkma/go-himeji/parser"
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

func printParserErrors(errors []string) {
	for _, msg := range errors {
		fmt.Printf("\t" + msg + "\n")
	}
}

func getArg(index int, userArgs []string) string {

	// Access arguments excluding the program name
	if len(os.Args[1:]) < 1 {
		fmt.Println("Missing command line args")
		return ""
	}

	//userArgs := os.Args[1:]
	//fmt.Println("User-provided arguments:", userArgs)

	// Access individual arguments by index
	//if len(os.Args) > 1 {
	//	fmt.Println("First user argument:", os.Args[1])
	//}

	// Iterate through user-provided arguments
	fmt.Println("Iterating through user arguments:")
	for i, arg := range userArgs {
		//fmt.Printf("Argument %d: %s\n", i+1, arg)
		if i == index {
			return arg
		}
	}
	return ""
}

func main() {
	// Get the first command line arg (zero index)
	src_file := getArg(1, os.Args)
	if src_file == "" {
		fmt.Println("Please profile source file. Example:")
		fmt.Printf("%s ../himeji/codes.txt\n", os.Args[0])
		os.Exit(1)
	}

	src, err := os.ReadFile(src_file)
	if err != nil {
		fmt.Printf("Error reading file: %s\n%v\n", src_file, err)
		os.Exit(1)
	}

	line := string(src)
	fmt.Printf("Source codes:\n%s\n", line)

	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(p.Errors())
		os.Exit(1)
	}

	comp := compiler.New()
	err = comp.Compile(program)
	if err != nil {
		fmt.Printf("Woops! Compilation failed:\n %s\n", err)
		os.Exit(1)
	}

	file, err := os.Create(HIMEJI_CODES_BIN)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	bc := comp.ByteCode()
	// Serialize
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(bc); err != nil {
		fmt.Println("Error encoding:", err)
		return
	}
	serializedData := buffer.Bytes()

	n, err := file.Write(serializedData)
	if err != nil {
		fmt.Println("file write failed:", err)
	}
	fmt.Printf("%d bytes wriiten.\n", n)

	// Deserialize
	/*
		var bc *compiler.ByteCode
		decoder := gob.NewDecoder(bytes.NewReader(serializedData))
		if err := decoder.Decode(&bc); err != nil {
			fmt.Println("Error decoding:", err)
			return
		}
		fmt.Printf("Deserialized bytecode: %+v\n", bc)
	*/
}
