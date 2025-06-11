package parser

import (
	"fmt"
	"strings"
)

var traceLevel int = 0

const traceIndent string = "\t"

func indentLevel() string {
	return strings.Repeat(traceIndent, traceLevel-1)
}

func tracePrint(fs string) {
	fmt.Printf("%s%s\n", indentLevel(), fs)
}

func increIdent() { traceLevel = traceLevel + 1 }
func decreIdent() { traceLevel = traceLevel - 1 }

func trace(msg string) string {
	increIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	tracePrint("END " + msg)
	decreIdent()
}
