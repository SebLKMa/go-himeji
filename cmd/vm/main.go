package main

import (
	"fmt"
	"os"
	"os/user"

	repl "github.com/seblkma/go-himeji/replcompiler"
)

const PROGLANG = "Himeji"

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Guten Tag %s, welcome to the %s programming language!\n", user.Name, PROGLANG)
	repl.Start(os.Stdin, os.Stdout)
}
