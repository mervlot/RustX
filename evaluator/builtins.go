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
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
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
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `to_lower` must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newError("arguments to `contains` must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			return nativeBoolToBooleanObject(strings.Contains(args[0].(*object.String).Value, args[1].(*object.String).Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newError("arguments to `split` must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			res := strings.Split(args[0].(*object.String).Value, args[1].(*object.String).Value)
			elements := make([]object.Object, len(res))
			for i, s := range res {
				elements[i] = &object.String{Value: s}
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ || args[1].Type() != object.STRING_OBJ {
				return newError("arguments to `join` must be ARRAY and STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			array := args[0].(*object.Array)
			sep := args[1].(*object.String).Value
			elements := make([]string, len(array.Elements))
			for i, e := range array.Elements {
				elements[i] = e.Inspect()
			}
			return &object.String{Value: strings.Join(elements, sep)}
		},
	},
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() == object.INTEGER_OBJ {
				v := args[0].(*object.Integer).Value
				if v < 0 {
					return &object.Integer{Value: -v}
				}
				return args[0]
			}
			if args[0].Type() == object.FLOAT_OBJ {
				return &object.Float{Value: math.Abs(args[0].(*object.Float).Value)}
			}
			return newError("argument to `abs` must be INTEGER or FLOAT, got %s", args[0].Type())
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			v1 := 0.0
			v2 := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				v1 = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				v1 = args[0].(*object.Float).Value
			} else {
				return newError("arguments to `pow` must be number, got %s", args[0].Type())
			}
			if args[1].Type() == object.INTEGER_OBJ {
				v2 = float64(args[1].(*object.Integer).Value)
			} else if args[1].Type() == object.FLOAT_OBJ {
				v2 = args[1].(*object.Float).Value
			} else {
				return newError("arguments to `pow` must be number, got %s", args[1].Type())
			}
			return &object.Float{Value: math.Pow(v1, v2)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			v := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				v = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				v = args[0].(*object.Float).Value
			} else {
				return newError("argument to `sin` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(v)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			v := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				v = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				v = args[0].(*object.Float).Value
			} else {
				return newError("argument to `cos` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(v)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			v := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				v = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				v = args[0].(*object.Float).Value
			} else {
				return newError("argument to `tan` must be number, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(v)}
		},
	},
	"now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `sleep` must be INTEGER (ms), got %s", args[0].Type())
			}
			ms := args[0].(*object.Integer).Value
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return NULL
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
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newError("arguments to `write_file` must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			err := os.WriteFile(args[0].(*object.String).Value, []byte(args[1].(*object.String).Value), 0644)
			if err != nil {
				return newError("could not write file: %s", err.Error())
			}
			return NULL
		},
	},
}
