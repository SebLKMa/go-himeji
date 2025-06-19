package evaluator

import (
	"fmt"

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
		//return evalStatements(node.Statements)
		return evalProgram(node)
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
		if isError(rhs) {
			return rhs
		}
		return evalPrefixExpression(node.Operator, rhs)
	case *ast.InfixExpression:
		lhs := Eval(node.Left)
		if isError(lhs) {
			return lhs
		}
		rhs := Eval(node.Right)
		if isError(rhs) {
			return rhs
		}
		return evalInfixExpression(node.Operator, lhs, rhs)
	case *ast.BlockStatement:
		//return evalStatements(node.Statements)
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		valExpr := Eval(node.Value)
		if isError(valExpr) {
			return valExpr
		}
		return &object.ReturnValue{Value: valExpr}

	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		//if returnValue, ok := result.(*object.ReturnValue); ok {
		//	return returnValue.Value
		//}
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
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
		return newError("unknown operator: -%s", rhs.Type())
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
		return newError("unknown operator: %s%s", op, rhs.Type())
	}
}

func evalInfixExpression(op string, lhs, rhs object.Object) object.Object {
	switch {
	case lhs.Type() == object.INTEGER_OBJ && rhs.Type() == object.INTEGER_OBJ:
		// Both operands are Integer objects
		return evalInfixIntegerExpression(op, lhs, rhs)
	case op == "==":
		// Reaching here means lhs and rhs are pointers to the Boolean singleton instance(s)
		return toBooleanObjectInstance(lhs == rhs)
	case op == "!=":
		// Reaching here means lhs and rhs are pointers to the Boolean singleton instance(s)
		return toBooleanObjectInstance(lhs != rhs)
	case lhs.Type() != rhs.Type():
		return newError("type mismatch: %s %s %s", lhs.Type(), op, rhs.Type())
	default:
		return newError("unknown operator: %s %s %s", lhs.Type(), op, rhs.Type())
	}
}

func evalInfixIntegerExpression(op string, lhs, rhs object.Object) object.Object {
	leftValue := lhs.(*object.Integer).Value
	rightValue := rhs.(*object.Integer).Value

	switch op {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return toBooleanObjectInstance(leftValue < rightValue)
	case ">":
		return toBooleanObjectInstance(leftValue > rightValue)
	case "==":
		return toBooleanObjectInstance(leftValue == rightValue)
	case "!=":
		return toBooleanObjectInstance(leftValue != rightValue)
	case ">=":
		return toBooleanObjectInstance(leftValue >= rightValue)
	default:
		return newError("unknown operator: %s %s %s", lhs.Type(), op, rhs.Type())
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true // handles if (n)
	}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.TrueBlock)
	} else if ie.FalseBlock != nil {
		return Eval(ie.FalseBlock)
	} else {
		return NULL
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
