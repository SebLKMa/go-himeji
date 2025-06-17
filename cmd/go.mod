module github.com/seblkma/go-himeji/cmd

replace (
	github.com/seblkma/go-himeji/ast => ../ast
	github.com/seblkma/go-himeji/evaluator => ../evaluator
	github.com/seblkma/go-himeji/lexer => ../lexer
	github.com/seblkma/go-himeji/object => ../object
	github.com/seblkma/go-himeji/parser => ../parser
	github.com/seblkma/go-himeji/repl => ../repl
	github.com/seblkma/go-himeji/token => ../token
)

go 1.22.5

require github.com/seblkma/go-himeji/repl v0.0.0-00010101000000-000000000000

require (
	github.com/seblkma/go-himeji/ast v0.0.0-00010101000000-000000000000 // indirect
	github.com/seblkma/go-himeji/evaluator v0.0.0-00010101000000-000000000000 // indirect
	github.com/seblkma/go-himeji/lexer v0.0.0-00010101000000-000000000000 // indirect
	github.com/seblkma/go-himeji/object v0.0.0-00010101000000-000000000000 // indirect
	github.com/seblkma/go-himeji/parser v0.0.0-00010101000000-000000000000 // indirect
	github.com/seblkma/go-himeji/token v0.0.0-00010101000000-000000000000 // indirect
)
