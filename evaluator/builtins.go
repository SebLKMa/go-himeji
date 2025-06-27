package evaluator

import (
	"github.com/seblkma/go-himeji/object"
)

// A map of built-in functions
var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		// This function returns the first array element, use REPL to test, e.g. first(myArr)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				//arr := args[0].(*object.Array)
				if len(arg.Elements) > 0 {
					return arg.Elements[0]
				}
				return NULL
			default:
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"last": &object.Builtin{
		// This function returns the last array element, use REPL to test, e.g. last(myArr)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				count := len(arg.Elements)
				if count > 0 {
					return arg.Elements[count-1]
				}
				return NULL
			default:
				return newError("argument to `last` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"tail": &object.Builtin{
		// This function returns a new array containing all array elements except the first, use REPL to test, e.g. tail(myArr)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				count := len(arg.Elements)
				if count > 0 {
					newElements := make([]object.Object, count-1, count-1)
					copy(newElements, arg.Elements[1:count]) // all except the first
					return &object.Array{Elements: newElements}
				}
				return NULL
			default:
				return newError("argument to `tail` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
	"push": &object.Builtin{
		// This function appends returns a new array, use REPL to test, e.g. tail(myArr)
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				count := len(arg.Elements)
				if count > 0 {
					newElements := make([]object.Object, count+1, count+1)
					copy(newElements, arg.Elements) // all
					newElements[count] = args[1]    // assign second arg to last array element
					return &object.Array{Elements: newElements}
				}
				return NULL
			default:
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}
		},
	},
}
