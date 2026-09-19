package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				if arg != nil {
					fmt.Print(arg.Inspect())
				}
			}
			return NULL
		},
	},
	"println": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				if arg != nil {
					fmt.Print(arg.Inspect())
				}
			}
			fmt.Println()
			return NULL
		},
	},
	"len": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for len: got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(tok, "argument to len not supported, got %s", args[0].Type())
			}
		},
	},
	"contains": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments for contains: got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			substr, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to contains must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			return nativeBoolToBooleanObject(strings.Contains(str.Value, substr.Value))
		},
	},
	"split": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments for split: got=%d, want=2", len(args))
			}
			str, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to split must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			parts := strings.Split(str.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments for join: got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to join must be ARRAY and STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for to_upper: got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(tok, "argument to to_upper must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for to_lower: got=%d, want=1", len(args))
			}
			str, ok := args[0].(*object.String)
			if !ok {
				return newError(tok, "argument to to_lower must be STRING, got %s", args[0].Type())
			}
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
	"abs": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for abs: got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value < 0 {
					return &object.Integer{Value: -arg.Value}
				}
				return arg
			case *object.Float:
				return &object.Float{Value: math.Abs(arg.Value)}
			default:
				return newError(tok, "argument to abs must be INTEGER or FLOAT, got %s", args[0].Type())
			}
		},
	},
	"pow": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments for pow: got=%d, want=2", len(args))
			}
			baseVal, isBaseNum := getFloatValue(args[0])
			expVal, isExpNum := getFloatValue(args[1])
			if !isBaseNum || !isExpNum {
				return newError(tok, "arguments to pow must be numeric, got %s and %s", args[0].Type(), args[1].Type())
			}
			res := math.Pow(baseVal, expVal)
			if args[0].Type() == object.INTEGER_OBJ && args[1].Type() == object.INTEGER_OBJ {
				return &object.Integer{Value: int64(res)}
			}
			return &object.Float{Value: res}
		},
	},
	"sin": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for sin: got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(tok, "argument to sin must be numeric, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for cos: got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(tok, "argument to cos must be numeric, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for tan: got=%d, want=1", len(args))
			}
			val, ok := getFloatValue(args[0])
			if !ok {
				return newError(tok, "argument to tan must be numeric, got %s", args[0].Type())
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},
	"now": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError(tok, "wrong number of arguments for now: got=%d, want=0", len(args))
			}
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for sleep: got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(tok, "argument to sleep must be INTEGER (ms), got %s", args[0].Type())
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"read_file": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(tok, "wrong number of arguments for read_file: got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError(tok, "argument to read_file must be STRING path, got %s", args[0].Type())
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return newError(tok, "failed to read file %q: %s", path.Value, err.Error())
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(tok token.Token, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(tok, "wrong number of arguments for write_file: got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(tok, "arguments to write_file must be STRING path and STRING content, got %s and %s", args[0].Type(), args[1].Type())
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(tok, "failed to write file %q: %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}

func getFloatValue(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
