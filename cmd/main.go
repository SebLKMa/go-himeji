package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/seblkma/go-himeji/repl"
)

const PROGLANG = "himeji"

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Guten Tag %s, welcome to the %s programming language!\n", user.Name, PROGLANG)
	repl.Start(os.Stdin, os.Stdout)
}
