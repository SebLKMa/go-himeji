package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/seblkma/go-himeji/lexer"
	tk "github.com/seblkma/go-himeji/token" // naming conflicts with go/token
)

const PROMPT = ">>"

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

		for tok := l.NextToken(); tok.Type != tk.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
