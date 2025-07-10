module github.com/seblkma/go-himeji/compiler

go 1.22.5

replace (
	github.com/seblkma/go-himeji/ast => ../ast
	//github.com/seblkma/go-himeji/lexer => ../lexer
	github.com/seblkma/go-himeji/object => ../object
	//github.com/seblkma/go-himeji/parser => ../parser
	//github.com/seblkma/go-himeji/token => ../token
    github.com/seblkma/go-himeji/opcodes => ../opcodes
)