package evaluator

import (
	"fmt"
	"math"
	"os"
	"rustx/object"
	"strings"
	"time"
)

var builtins = map[string]*object.Builtin{
	// Printing
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},

	// Common
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments to push. got=%d, want=2", len(args))
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("first argument to push must be ARRAY, got %s", args[0].Type())
			}
			arr.Elements = append(arr.Elements, args[1])
			return arr
		},
	},

	// Strings standard library
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("str_contains requires 2 arguments (s, substr)")
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_contains must be strings")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	},
	"str_split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("str_split requires 2 arguments (s, sep)")
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("arguments to str_split must be strings")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"str_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("str_upper requires 1 argument (s)")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_upper must be string")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("str_lower requires 1 argument (s)")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_lower must be string")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"str_trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("str_trim requires 1 argument (s)")
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError("argument to str_trim must be string")
			}
			return &object.String{Value: strings.TrimSpace(s.Value)}
		},
	},
	"str_replace": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 3 {
				return newError("str_replace requires 3 arguments (s, old, new)")
			}
			s, ok1 := args[0].(*object.String)
			oldStr, ok2 := args[1].(*object.String)
			newStr, ok3 := args[2].(*object.String)
			if !ok1 || !ok2 || !ok3 {
				return newError("arguments to str_replace must be strings")
			}
			return &object.String{Value: strings.ReplaceAll(s.Value, oldStr.Value, newStr.Value)}
		},
	},

	// Math standard library
	"math_abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_abs requires 1 argument")
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
				return newError("math_abs argument must be number")
			}
		},
	},
	"math_sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_sqrt requires 1 argument")
			}
			var val float64
			switch arg := args[0].(type) {
			case *object.Integer:
				val = float64(arg.Value)
			case *object.Float:
				val = arg.Value
			default:
				return newError("math_sqrt argument must be number")
			}
			if val < 0 {
				return newError("cannot take square root of negative number")
			}
			return &object.Float{Value: math.Sqrt(val)}
		},
	},
	"math_pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_pow requires 2 arguments (base, exp)")
			}
			var b, e float64
			if bi, ok := args[0].(*object.Integer); ok {
				b = float64(bi.Value)
			} else if bf, ok := args[0].(*object.Float); ok {
				b = bf.Value
			} else {
				return newError("math_pow base must be number")
			}

			if ei, ok := args[1].(*object.Integer); ok {
				e = float64(ei.Value)
			} else if ef, ok := args[1].(*object.Float); ok {
				e = ef.Value
			} else {
				return newError("math_pow exponent must be number")
			}

			res := math.Pow(b, e)
			// Return integer if both were integer and res is integer
			if _, ok1 := args[0].(*object.Integer); ok1 {
				if _, ok2 := args[1].(*object.Integer); ok2 {
					if res == math.Trunc(res) && !math.IsInf(res, 0) {
						return &object.Integer{Value: int64(res)}
					}
				}
			}
			return &object.Float{Value: res}
		},
	},
	"math_min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_min requires 2 arguments")
			}
			if i1, ok1 := args[0].(*object.Integer); ok1 {
				if i2, ok2 := args[1].(*object.Integer); ok2 {
					if i1.Value < i2.Value {
						return i1
					}
					return i2
				}
			}
			f1 := toFloat(args[0])
			f2 := toFloat(args[1])
			return &object.Float{Value: math.Min(f1, f2)}
		},
	},
	"math_max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("math_max requires 2 arguments")
			}
			if i1, ok1 := args[0].(*object.Integer); ok1 {
				if i2, ok2 := args[1].(*object.Integer); ok2 {
					if i1.Value > i2.Value {
						return i1
					}
					return i2
				}
			}
			f1 := toFloat(args[0])
			f2 := toFloat(args[1])
			return &object.Float{Value: math.Max(f1, f2)}
		},
	},
	"math_floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_floor requires 1 argument")
			}
			f := toFloat(args[0])
			return &object.Integer{Value: int64(math.Floor(f))}
		},
	},
	"math_ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("math_ceil requires 1 argument")
			}
			f := toFloat(args[0])
			return &object.Integer{Value: int64(math.Ceil(f))}
		},
	},

	// Time standard library
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"time_millis": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano() / int64(time.Millisecond)}
		},
	},
	"time_sleep": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("time_sleep requires 1 argument (milliseconds)")
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError("time_sleep argument must be integer")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("read_file requires 1 argument (filepath)")
			}
			filepath, ok := args[0].(*object.String)
			if !ok {
				return newError("read_file filepath must be a string")
			}
			data, err := os.ReadFile(filepath.Value)
			if err != nil {
				return newError("read_file error: %s", err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("write_file requires 2 arguments (filepath, content)")
			}
			filepath, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError("write_file arguments must be strings")
			}
			err := os.WriteFile(filepath.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError("write_file error: %s", err.Error())
			}
			return TRUE
		},
	},
	"file_exists": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("file_exists requires 1 argument (filepath)")
			}
			filepath, ok := args[0].(*object.String)
			if !ok {
				return newError("file_exists argument must be string")
			}
			_, err := os.Stat(filepath.Value)
			return nativeBoolToBooleanObject(!os.IsNotExist(err))
		},
	},
}

func toFloat(obj object.Object) float64 {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value)
	case *object.Float:
		return o.Value
	default:
		return 0
	}
}
