package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/seblkma/go-himeji/cmd/common"
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

func printParserErrors(errors []string) {
	for _, msg := range errors {
		fmt.Printf("\t" + msg + "\n")
	}
}

func outputFile(filePath string, newExtension string) string {
	// Get the directory
	dir := filepath.Dir(filePath)

	// Get the base name without the extension
	baseName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// Construct the new file path
	newFilePath := filepath.Join(dir, baseName+newExtension)

	//fmt.Printf("Original path: %s\n", filePath)
	//fmt.Printf("New path: %s\n", newFilePath)

	//noExtPath := "/path/to/my/document"
	//newExtPath := filepath.Join(filepath.Dir(noExtPath), strings.TrimSuffix(filepath.Base(noExtPath), filepath.Ext(noExtPath))+newExtension)
	//fmt.Printf("Original path (no ext): %s\n", noExtPath)
	//fmt.Printf("New path (no ext): %s\n", newExtPath)
	return newFilePath
}

func main() {
	// Get the first command line arg (zero index)
	inputFile := common.GetCmdArg(1, os.Args)
	if inputFile == "" {
		fmt.Println("Please provide source file. Example:")
		fmt.Printf("%s codes.txt\n", os.Args[0])
		os.Exit(1)
	}

	src, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %s\n%v\n", inputFile, err)
		os.Exit(1)
	}

	outFile := outputFile(inputFile, ".bin")
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

	file, err := os.Create(outFile)
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
	fmt.Printf("%d bytes written to %s\n", n, outFile)

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
