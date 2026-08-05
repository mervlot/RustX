package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

func init() {
	builtins = map[string]*object.Builtin{
		"print": {
			Fn: func(args ...object.Object) object.Object {
				var parts []string
				for _, arg := range args {
					if s, ok := arg.(*object.String); ok {
						parts = append(parts, s.Value)
					} else {
						parts = append(parts, arg.Inspect())
					}
				}
				fmt.Print(strings.Join(parts, " "))
				return NULL
			},
		},
		"println": {
			Fn: func(args ...object.Object) object.Object {
				var parts []string
				for _, arg := range args {
					if s, ok := arg.(*object.String); ok {
						parts = append(parts, s.Value)
					} else {
						parts = append(parts, arg.Inspect())
					}
				}
				fmt.Println(strings.Join(parts, " "))
				return NULL
			},
		},
		"len": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("len: expected 1 argument, got %d", len(args))}
				}
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return &object.Error{Message: fmt.Sprintf("len: argument not supported, got %s", args[0].Type())}
				}
			},
		},
		"contains": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Error{Message: fmt.Sprintf("contains: expected 2 arguments, got %d", len(args))}
				}
				switch container := args[0].(type) {
				case *object.String:
					substr, ok := args[1].(*object.String)
					if !ok {
						return &object.Error{Message: fmt.Sprintf("contains: second argument must be string when first argument is string, got %s", args[1].Type())}
					}
					return nativeBoolToBooleanObject(strings.Contains(container.Value, substr.Value))
				case *object.Array:
					for _, elem := range container.Elements {
						if evalObjectEquality("==", elem, args[1]) {
							return TRUE
						}
					}
					return FALSE
				default:
					return &object.Error{Message: fmt.Sprintf("contains: first argument must be string or array, got %s", args[0].Type())}
				}
			},
		},
		"split": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Error{Message: fmt.Sprintf("split: expected 2 arguments, got %d", len(args))}
				}
				str, ok1 := args[0].(*object.String)
				sep, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Error{Message: fmt.Sprintf("split: both arguments must be strings, got %s and %s", args[0].Type(), args[1].Type())}
				}
				parts := strings.Split(str.Value, sep.Value)
				var elems []object.Object
				for _, p := range parts {
					elems = append(elems, &object.String{Value: p})
				}
				return &object.Array{Elements: elems}
			},
		},
		"join": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Error{Message: fmt.Sprintf("join: expected 2 arguments, got %d", len(args))}
				}
				arr, ok1 := args[0].(*object.Array)
				sep, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Error{Message: fmt.Sprintf("join: expected (array, string), got (%s, %s)", args[0].Type(), args[1].Type())}
				}
				var parts []string
				for _, elem := range arr.Elements {
					if s, ok := elem.(*object.String); ok {
						parts = append(parts, s.Value)
					} else {
						parts = append(parts, elem.Inspect())
					}
				}
				return &object.String{Value: strings.Join(parts, sep.Value)}
			},
		},
		"to_upper": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("to_upper: expected 1 argument, got %d", len(args))}
				}
				str, ok := args[0].(*object.String)
				if !ok {
					return &object.Error{Message: fmt.Sprintf("to_upper: argument must be string, got %s", args[0].Type())}
				}
				return &object.String{Value: strings.ToUpper(str.Value)}
			},
		},
		"to_lower": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("to_lower: expected 1 argument, got %d", len(args))}
				}
				str, ok := args[0].(*object.String)
				if !ok {
					return &object.Error{Message: fmt.Sprintf("to_lower: argument must be string, got %s", args[0].Type())}
				}
				return &object.String{Value: strings.ToLower(str.Value)}
			},
		},
		"abs": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("abs: expected 1 argument, got %d", len(args))}
				}
				switch val := args[0].(type) {
				case *object.Integer:
					v := val.Value
					if v < 0 {
						v = -v
					}
					return &object.Integer{Value: v}
				case *object.Float:
					return &object.Float{Value: math.Abs(val.Value)}
				default:
					return &object.Error{Message: fmt.Sprintf("abs: argument must be numeric, got %s", args[0].Type())}
				}
			},
		},
		"pow": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Error{Message: fmt.Sprintf("pow: expected 2 arguments, got %d", len(args))}
				}
				var base, exponent float64
				switch b := args[0].(type) {
				case *object.Integer:
					base = float64(b.Value)
				case *object.Float:
					base = b.Value
				default:
					return &object.Error{Message: fmt.Sprintf("pow: base must be numeric, got %s", args[0].Type())}
				}
				switch e := args[1].(type) {
				case *object.Integer:
					exponent = float64(e.Value)
				case *object.Float:
					exponent = e.Value
				default:
					return &object.Error{Message: fmt.Sprintf("pow: exponent must be numeric, got %s", args[1].Type())}
				}
				return &object.Float{Value: math.Pow(base, exponent)}
			},
		},
		"sin": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("sin: expected 1 argument, got %d", len(args))}
				}
				var val float64
				switch v := args[0].(type) {
				case *object.Integer:
					val = float64(v.Value)
				case *object.Float:
					val = v.Value
				default:
					return &object.Error{Message: fmt.Sprintf("sin: argument must be numeric, got %s", args[0].Type())}
				}
				return &object.Float{Value: math.Sin(val)}
			},
		},
		"cos": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("cos: expected 1 argument, got %d", len(args))}
				}
				var val float64
				switch v := args[0].(type) {
				case *object.Integer:
					val = float64(v.Value)
				case *object.Float:
					val = v.Value
				default:
					return &object.Error{Message: fmt.Sprintf("cos: argument must be numeric, got %s", args[0].Type())}
				}
				return &object.Float{Value: math.Cos(val)}
			},
		},
		"tan": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("tan: expected 1 argument, got %d", len(args))}
				}
				var val float64
				switch v := args[0].(type) {
				case *object.Integer:
					val = float64(v.Value)
				case *object.Float:
					val = v.Value
				default:
					return &object.Error{Message: fmt.Sprintf("tan: argument must be numeric, got %s", args[0].Type())}
				}
				return &object.Float{Value: math.Tan(val)}
			},
		},
		"sleep": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("sleep: expected 1 argument, got %d", len(args))}
				}
				ms, ok := args[0].(*object.Integer)
				if !ok {
					return &object.Error{Message: fmt.Sprintf("sleep: argument must be integer milliseconds, got %s", args[0].Type())}
				}
				time.Sleep(time.Duration(ms.Value) * time.Millisecond)
				return NULL
			},
		},
		"now": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 0 {
					return &object.Error{Message: fmt.Sprintf("now: expected 0 arguments, got %d", len(args))}
				}
				return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
			},
		},
		"read_file": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Error{Message: fmt.Sprintf("read_file: expected 1 argument, got %d", len(args))}
				}
				path, ok := args[0].(*object.String)
				if !ok {
					return &object.Error{Message: fmt.Sprintf("read_file: argument must be string filepath, got %s", args[0].Type())}
				}
				data, err := os.ReadFile(path.Value)
				if err != nil {
					return &object.Error{Message: fmt.Sprintf("read_file error: %s", err.Error())}
				}
				return &object.String{Value: string(data)}
			},
		},
		"write_file": {
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Error{Message: fmt.Sprintf("write_file: expected 2 arguments, got %d", len(args))}
				}
				path, ok1 := args[0].(*object.String)
				content, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Error{Message: fmt.Sprintf("write_file: expected (string, string), got (%s, %s)", args[0].Type(), args[1].Type())}
				}
				err := os.WriteFile(path.Value, []byte(content.Value), 0644)
				if err != nil {
					return &object.Error{Message: fmt.Sprintf("write_file error: %s", err.Error())}
				}
				return TRUE
			},
		},
	}
}
