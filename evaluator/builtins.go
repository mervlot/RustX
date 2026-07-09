package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newBuiltinError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			for i, arg := range args {
				fmt.Print(arg.Inspect())
				if i < len(args)-1 {
					fmt.Print(" ")
				}
			}
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for i, arg := range args {
				fmt.Print(arg.Inspect())
				if i < len(args)-1 {
					fmt.Print(" ")
				}
			}
			fmt.Println()
			return NULL
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newBuiltinError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			arr := args[0].(*object.Array)
			newElements := append(arr.Elements, args[1])
			return &object.Array{Elements: newElements}
		},
	},
	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newBuiltinError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() == object.INTEGER_OBJ {
				val := args[0].(*object.Integer).Value
				if val < 0 {
					return &object.Integer{Value: -val}
				}
				return args[0]
			}
			if args[0].Type() == object.FLOAT_OBJ {
				val := args[0].(*object.Float).Value
				return &object.Float{Value: math.Abs(val)}
			}
			return newBuiltinError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newBuiltinError("wrong number of arguments. got=%d, want=2", len(args))
			}
			x, ok1 := getFloat(args[0])
			y, ok2 := getFloat(args[1])
			if !ok1 || !ok2 {
				return newBuiltinError("arguments to `pow` must be numeric")
			}
			return &object.Float{Value: math.Pow(x, y)}
		},
	},
	// Strings
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return newBuiltinError("argument to `to_upper` must be a single STRING")
			}
			return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return newBuiltinError("argument to `to_lower` must be a single STRING")
			}
			return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newBuiltinError("arguments to `split` must be (STRING, STRING)")
			}
			s := args[0].(*object.String).Value
			sep := args[1].(*object.String).Value
			parts := strings.Split(s, sep)
			objs := make([]object.Object, len(parts))
			for i, p := range parts {
				objs[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: objs}
		},
	},
	// Time
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.INTEGER_OBJ {
				return newBuiltinError("argument to `sleep` must be INTEGER (ms)")
			}
			ms := args[0].(*object.Integer).Value
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return NULL
		},
	},
	// Filesystem
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return newBuiltinError("argument to `read_file` must be STRING (path)")
			}
			path := args[0].(*object.String).Value
			content, err := os.ReadFile(path)
			if err != nil {
				return newBuiltinError("failed to read file: %s", err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newBuiltinError("arguments to `write_file` must be (STRING path, STRING content)")
			}
			path := args[0].(*object.String).Value
			content := args[1].(*object.String).Value
			err := os.WriteFile(path, []byte(content), 0644)
			if err != nil {
				return newBuiltinError("failed to write file: %s", err.Error())
			}
			return TRUE
		},
	},
}

func newBuiltinError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func getFloat(obj object.Object) (float64, bool) {
	if obj.Type() == object.INTEGER_OBJ {
		return float64(obj.(*object.Integer).Value), true
	}
	if obj.Type() == object.FLOAT_OBJ {
		return obj.(*object.Float).Value, true
	}
	return 0, false
}
