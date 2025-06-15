package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/seblkma/go-himeji/lexer"
	"github.com/seblkma/go-himeji/parser"
	// naming conflicts with go/token
)

const PROMPT = ">>"

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")

		//for tok := l.NextToken(); tok.Type != tk.EOF; tok = l.NextToken() {
		//	fmt.Printf("%+v\n", tok)
		//}
	}
}
