package evaluator

import (
	"fmt"
	"github.com/username/rustx/object"
	"os"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
				fmt.Print(" ")
			}
			fmt.Println()
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `to_upper` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `read_file` must be STRING, got %s", args[0].Type())
			}
			content, err := os.ReadFile(args[0].(*object.String).Value)
			if err != nil {
				return newError("could not read file: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
}
