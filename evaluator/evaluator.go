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
	case *ast.PrefixExpression:
		rhs := Eval(node.Right)
		return evalPrefixExpression(node.Operator, rhs)
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

func evalBangOperatorExpression(rhs object.Object) object.Object {
	switch rhs {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(rhs object.Object) object.Object {
	if rhs.Type() != object.INTEGER_OBJ {
		return NULL
	}

	// Returns the minus value
	value := rhs.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalPrefixExpression(op string, rhs object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(rhs)
	case "-":
		return evalMinusPrefixOperatorExpression(rhs)
	default:
		return NULL
	}
}
