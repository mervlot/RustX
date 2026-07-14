package evaluator

import (
	"os"
	"strings"
	"time"

	"github.com/username/rustx/object"
)

func init() {
	// Add more builtins to the global map in evaluator.go
	// Since they are in the same package, we can just add to it if it were exported,
	// but it's not. I'll merge this into evaluator.go or make it a separate registration function.
}

// For now, I will add more builtins directly in evaluator.go to keep it simple.
// But I'll provide the registration here if I refactor.

func RegisterBuiltins(env *object.Environment) {
    // This could be used for stdlib-like registration
}

// Let's add some more to the builtins map in evaluator.go via a separate file if we want,
// but in Go, maps in the same package can be updated in init() if they are package-level.

func extendBuiltins() {
	builtins["sleep"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.INTEGER_OBJ {
				return newError("argument to `sleep` must be INTEGER, got %s", args[0].Type())
			}
			ms := args[0].(*object.Integer).Value
			time.Sleep(time.Duration(ms) * time.Millisecond)
			return NULL
		},
	}

	builtins["now"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().UnixMilli()}
		},
	}

	builtins["read_file"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `read_file` must be STRING, got %s", args[0].Type())
			}
			path := args[0].(*object.String).Value
			content, err := os.ReadFile(path)
			if err != nil {
				return newError("could not read file %s: %s", path, err)
			}
			return &object.String{Value: string(content)}
		},
	}

	builtins["write_file"] = &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return newError("arguments to `write_file` must be STRING, got %s and %s", args[0].Type(), args[1].Type())
			}
			path := args[0].(*object.String).Value
			content := args[1].(*object.String).Value
			err := os.WriteFile(path, []byte(content), 0644)
			if err != nil {
				return newError("could not write file %s: %s", path, err)
			}
			return NULL
		},
	}

    // String utilities
    builtins["to_upper"] = &object.Builtin{
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
                return newError("to_upper expects 1 string argument")
            }
            return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
        },
    }

    builtins["to_lower"] = &object.Builtin{
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
                return newError("to_lower expects 1 string argument")
            }
            return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
        },
    }

    builtins["split"] = &object.Builtin{
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 2 || args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
                return newError("split expects 2 string arguments")
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
    }
}

func init() {
    extendBuiltins()
}
