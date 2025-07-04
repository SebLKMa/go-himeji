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
func Eval(n ast.Node, env *object.Environment) object.Object {
	switch node := n.(type) {
	case *ast.Program:
		//return evalStatements(node.Statements)
		return evalProgram(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		// Allocates new Integer values
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		// No need to allocate new Boolean objects since they represent the same TRUE or FALSE values
		return toBooleanObjectInstance(node.Value)
	case *ast.PrefixExpression:
		rhs := Eval(node.Right, env)
		if isError(rhs) {
			return rhs
		}
		return evalPrefixExpression(node.Operator, rhs)
	case *ast.InfixExpression:
		lhs := Eval(node.Left, env)
		if isError(lhs) {
			return lhs
		}
		rhs := Eval(node.Right, env)
		if isError(rhs) {
			return rhs
		}
		return evalInfixExpression(node.Operator, lhs, rhs)
	case *ast.BlockStatement:
		//return evalStatements(node.Statements)
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		valExpr := Eval(node.Value, env)
		if isError(valExpr) {
			return valExpr
		}
		return &object.ReturnValue{Value: valExpr}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.LetStatement:
		valExpr := Eval(node.Value, env)
		if isError(valExpr) {
			return valExpr
		}
		env.Set(node.Name.Value, valExpr)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		// Function arguments have to be evaluated before passing them as args
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return executeFunction(fn, args)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

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
	case lhs.Type() == object.STRING_OBJ && rhs.Type() == object.STRING_OBJ:
		// Both operands are String objects
		return evalInfixStringExpression(op, lhs, rhs)
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

func evalInfixStringExpression(op string, lhs, rhs object.Object) object.Object {
	leftValue := lhs.(*object.String).Value
	rightValue := rhs.(*object.String).Value

	switch op {
	case "+":
		return &object.String{Value: leftValue + rightValue}
	default:
		return newError("unknown operator: %s %s %s", lhs.Type(), op, rhs.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if obj, found := env.Get(node.Value); found {
		return obj
	}

	if builtin, found := builtins[node.Value]; found {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
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

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.TrueBlock, env)
	} else if ie.FalseBlock != nil {
		return Eval(ie.FalseBlock, env)
	} else {
		return NULL
	}
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object
	for _, e := range exprs {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}
	return results
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrObj.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrObj.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObj := hash.(*object.Hashes)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for nodeKey, nodeValue := range node.Pairs {
		key := Eval(nodeKey, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}

		value := Eval(nodeValue, env)
		if isError(value) {
			return value
		}

		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}

	return &object.Hashes{Pairs: pairs}
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

func executeFunction(fn object.Object, args []object.Object) object.Object {
	switch function := fn.(type) {
	case *object.Function:
		scopedEnv := scopeFunctionEnv(function, args)
		// Recursively Eval until the last function body
		// Unbox it so that evalBlockStatement won’t stop evaluating statements in “outer” functions
		executed := Eval(function.Body, scopedEnv)
		return unboxReturnValue(executed)
	case *object.Builtin:
		return function.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func scopeFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	newScope := object.NewInnerEnvironment(fn.Env)

	for i, p := range fn.Parameters {
		newScope.Set(p.Value, args[i])
	}

	return newScope
}

func unboxReturnValue(obj object.Object) object.Object {
	if unboxed, ok := obj.(*object.ReturnValue); ok {
		return unboxed.Value
	}
	return obj
}
