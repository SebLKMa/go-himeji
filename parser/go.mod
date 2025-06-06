module github.com/seblkma/go-himeji/parser

replace (
	github.com/seblkma/go-himeji/ast => ../ast
	github.com/seblkma/go-himeji/lexer => ../lexer
	github.com/seblkma/go-himeji/token => ../token
)

go 1.22.5

require (
	github.com/seblkma/go-himeji/ast v0.0.0-00010101000000-000000000000
	github.com/seblkma/go-himeji/lexer v0.0.0-00010101000000-000000000000
	github.com/seblkma/go-himeji/token v0.0.0-00010101000000-000000000000
)
