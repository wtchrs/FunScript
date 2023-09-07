package evaluator

import "funscript/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. expected=1, got=%d", len(args))
			}

			str, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}

			return &object.Integer{Value: int64(len(str.Value))}
		},
	},
}
