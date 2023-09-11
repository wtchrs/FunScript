package evaluator

import "funscript/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1, got=%d", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}

			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}

			case *object.Hash:
				return &object.Integer{Value: int64(len(arg.Pairs))}

			default:
				return newError("argument to len not supported, got %s", args[0].Type())
			}
		},
	},

	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1, got=%d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to first must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1, got=%d", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to last must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},

	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1, got=%d", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to tail must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElem := make([]object.Object, length-1)
				copy(newElem, arr.Elements[1:])
				return &object.Array{Elements: newElem}
			}

			return NULL
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. expected=2, got=%d", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument0 to push must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElem := make([]object.Object, length+1)
			copy(newElem, arr.Elements)
			newElem[length] = args[1]
			return &object.Array{Elements: newElem}
		},
	},
}
