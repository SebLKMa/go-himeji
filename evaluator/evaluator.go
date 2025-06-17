package evaluator

import (
	"github.com/seblkma/go-himeji/ast"
	"github.com/seblkma/go-himeji/object"
	//hparser "github.com/seblkma/go-himeji/parser"
)

// Singletons to be referenced by all Boolean objects
var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval evaluates an AST node to our value Object representation
func Eval(n ast.Node) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		// Allocates new Integer values
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		// No need to allocate new Boolean objects since they represent the same TRUE or FALSE values
		return toBooleanObjectInstance(node.Value)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func toBooleanObjectInstance(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}
