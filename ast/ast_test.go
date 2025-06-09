package ast

import (
	"testing"

	tk "github.com/seblkma/go-himeji/token"
)

// GOFLAGS="-count=1" go test -run TestProgramString
func TestProgramString(t *testing.T) {
	testStr := "let myVar = anotherVar;"

	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: tk.Token{Type: tk.LET, Literal: "let"},
				Name: &Identifier{
					Token: tk.Token{Type: tk.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: tk.Token{Type: tk.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != testStr {
		t.Errorf("\nexpected program string: %s\nbut got:%s\n", testStr, program.String())
	}

}
