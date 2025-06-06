module github.com/seblkma/go-himeji/repl

replace (
	github.com/seblkma/go-himeji/lexer => ../lexer
	github.com/seblkma/go-himeji/token => ../token
)

go 1.22.5

require github.com/seblkma/go-himeji/lexer v0.0.0-00010101000000-000000000000

require github.com/seblkma/go-himeji/token v0.0.0-00010101000000-000000000000 // indirect
