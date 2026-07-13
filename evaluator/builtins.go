package evaluator

import (
	"fmt"
	"github.com/username/rustx/object"
	"math"
	"os"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return &object.Error{Message: fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type())}
			}
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
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return &object.Error{Message: fmt.Sprintf("argument to `push` must be ARRAY, got %s", args[0].Type())}
			}

			arr := args[0].(*object.Array)
			newElements := make([]object.Object, len(arr.Elements)+1)
			copy(newElements, arr.Elements)
			newElements[len(arr.Elements)] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(args))}
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
			return &object.Error{Message: fmt.Sprintf("argument to `abs` must be number, got %s", args[0].Type())}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(args))}
			}
			var base, exp float64
			if args[0].Type() == object.INTEGER_OBJ {
				base = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				base = args[0].(*object.Float).Value
			} else {
				return &object.Error{Message: "first argument to `pow` must be number"}
			}

			if args[1].Type() == object.INTEGER_OBJ {
				exp = float64(args[1].(*object.Integer).Value)
			} else if args[1].Type() == object.FLOAT_OBJ {
				exp = args[1].(*object.Float).Value
			} else {
				return &object.Error{Message: "second argument to `pow` must be number"}
			}

			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			val := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				val = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				val = args[0].(*object.Float).Value
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			val := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				val = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				val = args[0].(*object.Float).Value
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			val := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				val = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				val = args[0].(*object.Float).Value
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	// Strings
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "arguments to `contains` must be strings"}
			}
			return &object.Boolean{Value: strings.Contains(args[0].(*object.String).Value, args[1].(*object.String).Value)}
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "arguments to `split` must be strings"}
			}
			parts := strings.Split(args[0].(*object.String).Value, args[1].(*object.String).Value)
			objs := make([]object.Object, len(parts))
			for i, p := range parts {
				objs[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: objs}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "wrong number of arguments"}
			}
			if args[0].Type() != object.ARRAY_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "arguments to `join` must be ARRAY and STRING"}
			}
			arr := args[0].(*object.Array)
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, args[1].(*object.String).Value)}
		},
	},
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "invalid arguments"}
			}
			return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "invalid arguments"}
			}
			return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
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
				return &object.Error{Message: "invalid arguments"}
			}
			time.Sleep(time.Duration(args[0].(*object.Integer).Value) * time.Millisecond)
			return NULL
		},
	},
	// Filesystem
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "invalid arguments"}
			}
			content, err := os.ReadFile(args[0].(*object.String).Value)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "invalid arguments"}
			}
			err := os.WriteFile(args[0].(*object.String).Value, []byte(args[1].(*object.String).Value), 0644)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			return TRUE
		},
	},
}
