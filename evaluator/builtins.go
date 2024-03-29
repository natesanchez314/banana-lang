package evaluator

import "banana/object"

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("Wrong number of args, got=%d, expected=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Val: int64(len(arg.Val))}
			default:
				return newError("Arg to `len` not supported, got %s", arg.Type())
			}
		},
	},
}