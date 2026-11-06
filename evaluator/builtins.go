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
	"print": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Print(strings.Join(parts, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			var parts []string
			for _, arg := range args {
				parts = append(parts, arg.Inspect())
			}
			fmt.Println(strings.Join(parts, " "))
			return NULL
		},
	},
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for len. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError(0, 0, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for push. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError(0, 0, "first argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			arr.Elements = append(arr.Elements, args[1])
			return arr
		},
	},

	// Strings
	"str_contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for str_contains. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to str_contains must be STRING")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	},
	"str_split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for str_split. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to str_split must be STRING")
			}
			parts := strings.Split(s.Value, sep.Value)
			elements := make([]object.Object, len(parts))
			for i, p := range parts {
				elements[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: elements}
		},
	},
	"str_join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for str_join. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to str_join must be (ARRAY, STRING)")
			}
			parts := make([]string, len(arr.Elements))
			for i, el := range arr.Elements {
				parts[i] = el.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"str_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for str_upper. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to str_upper must be STRING")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"str_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for str_lower. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to str_lower must be STRING")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"str_trim": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for str_trim. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to str_trim must be STRING")
			}
			return &object.String{Value: strings.TrimSpace(s.Value)}
		},
	},

	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for abs. got=%d, want=1", len(args))
			}
			switch val := args[0].(type) {
			case *object.Integer:
				if val.Value < 0 {
					return &object.Integer{Value: -val.Value}
				}
				return val
			case *object.Float:
				return &object.Float{Value: math.Abs(val.Value)}
			default:
				return newError(0, 0, "argument to abs must be numeric")
			}
		},
	},
	"sqrt": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for sqrt. got=%d, want=1", len(args))
			}
			switch val := args[0].(type) {
			case *object.Integer:
				return &object.Float{Value: math.Sqrt(float64(val.Value))}
			case *object.Float:
				return &object.Float{Value: math.Sqrt(val.Value)}
			default:
				return newError(0, 0, "argument to sqrt must be numeric")
			}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for pow. got=%d, want=2", len(args))
			}
			var base, exp float64
			switch b := args[0].(type) {
			case *object.Integer:
				base = float64(b.Value)
			case *object.Float:
				base = b.Value
			default:
				return newError(0, 0, "base to pow must be numeric")
			}
			switch e := args[1].(type) {
			case *object.Integer:
				exp = float64(e.Value)
			case *object.Float:
				exp = e.Value
			default:
				return newError(0, 0, "exponent to pow must be numeric")
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"min": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for min. got=%d, want=2", len(args))
			}
			v1, ok1 := args[0].(*object.Integer)
			v2, ok2 := args[1].(*object.Integer)
			if ok1 && ok2 {
				if v1.Value < v2.Value {
					return v1
				}
				return v2
			}
			f1 := toFloat(args[0])
			f2 := toFloat(args[1])
			if !math.IsNaN(f1) && !math.IsNaN(f2) {
				return &object.Float{Value: math.Min(f1, f2)}
			}
			return newError(0, 0, "arguments to min must be numeric")
		},
	},
	"max": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for max. got=%d, want=2", len(args))
			}
			v1, ok1 := args[0].(*object.Integer)
			v2, ok2 := args[1].(*object.Integer)
			if ok1 && ok2 {
				if v1.Value > v2.Value {
					return v1
				}
				return v2
			}
			f1 := toFloat(args[0])
			f2 := toFloat(args[1])
			if !math.IsNaN(f1) && !math.IsNaN(f2) {
				return &object.Float{Value: math.Max(f1, f2)}
			}
			return newError(0, 0, "arguments to max must be numeric")
		},
	},
	"round": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for round. got=%d, want=1", len(args))
			}
			f := toFloat(args[0])
			if math.IsNaN(f) {
				return newError(0, 0, "argument to round must be numeric")
			}
			return &object.Integer{Value: int64(math.Round(f))}
		},
	},
	"floor": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for floor. got=%d, want=1", len(args))
			}
			f := toFloat(args[0])
			if math.IsNaN(f) {
				return newError(0, 0, "argument to floor must be numeric")
			}
			return &object.Integer{Value: int64(math.Floor(f))}
		},
	},
	"ceil": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for ceil. got=%d, want=1", len(args))
			}
			f := toFloat(args[0])
			if math.IsNaN(f) {
				return newError(0, 0, "argument to ceil must be numeric")
			}
			return &object.Integer{Value: int64(math.Ceil(f))}
		},
	},

	// Time
	"time_now": {
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixNano()}
		},
	},
	"sleep_ms": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for sleep_ms. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newError(0, 0, "argument to sleep_ms must be integer")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},

	// Filesystem basics
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(0, 0, "wrong number of arguments for read_file. got=%d, want=1", len(args))
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return newError(0, 0, "argument to read_file must be STRING")
			}
			data, err := os.ReadFile(path.Value)
			if err != nil {
				return newError(0, 0, "failed to read file '%s': %s", path.Value, err.Error())
			}
			return &object.String{Value: string(data)}
		},
	},
	"write_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError(0, 0, "wrong number of arguments for write_file. got=%d, want=2", len(args))
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newError(0, 0, "arguments to write_file must be (STRING, STRING)")
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return newError(0, 0, "failed to write file '%s': %s", path.Value, err.Error())
			}
			return TRUE
		},
	},
}

func toFloat(obj object.Object) float64 {
	switch v := obj.(type) {
	case *object.Integer:
		return float64(v.Value)
	case *object.Float:
		return v.Value
	default:
		return math.NaN()
	}
}
