package parser

import (
	"fmt"
	"strings"
)

var wantTrace bool = false // decide if you want to see parsing trace
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
	if !wantTrace {
		return ""
	}
	increIdent()
	tracePrint("BEGIN " + msg)
	return msg
}

func untrace(msg string) {
	if !wantTrace {
		return
	}
	tracePrint("END " + msg)
	decreIdent()
}
