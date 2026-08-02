package evaluator

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ShorthandAssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		tok := node.GetToken()
		updated := env.Update(node.Name.Value, val)
		if !updated {
			return newError("identifier %s is not defined", tok, node.Name.Value)
		}
		return NULL

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(node.GetToken(), left, index)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.GetToken(), node.Operator, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		// Handle short-circuiting for && and ||
		if node.Operator == "&&" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}
		if node.Operator == "||" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.GetToken(), node.Operator, left, right)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fun := &object.Function{Parameters: params, Env: env, Body: body}
		if node.Name != nil {
			env.Set(node.Name.Value, fun)
		}
		return fun

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(node.GetToken(), function, args)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL // Prevents Go nil pointer issues

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.GetToken(), node.Value)
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(tok token.Token, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(tok, right)
	default:
		return newError("unknown operator: %s%s", tok, operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		if isTruthy(right) {
			return FALSE
		}
		return TRUE
	}
}

func evalMinusPrefixOperatorExpression(tok token.Token, right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	default:
		return newError("unknown operator: -%s", tok, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError("postfix operator can only be applied to identifiers", node.GetToken())
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError("identifier not found: %s", node.GetToken(), ident.Value)
	}

	var newVal object.Object
	switch val := val.(type) {
	case *object.Integer:
		if node.Operator == "++" {
			newVal = &object.Integer{Value: val.Value + 1}
		} else {
			newVal = &object.Integer{Value: val.Value - 1}
		}
	case *object.Float:
		if node.Operator == "++" {
			newVal = &object.Float{Value: val.Value + 1.0}
		} else {
			newVal = &object.Float{Value: val.Value - 1.0}
		}
	default:
		return newError("postfix operator can only be applied to numbers, got %s", node.GetToken(), val.Type())
	}

	env.Update(ident.Value, newVal)
	return val // Return value before modification (post-increment/decrement)
}

func evalInfixExpression(tok token.Token, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(tok, operator, left.(*object.Integer), right.(*object.Integer))
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(tok, operator, left.(*object.Float), right.(*object.Float))
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Cast integer to float
		fl := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(tok, operator, fl, right.(*object.Float))
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Cast integer to float
		fr := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(tok, operator, left.(*object.Float), fr)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(tok, operator, left.(*object.String), right.(*object.String))
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(tok, operator, left.(*object.Boolean), right.(*object.Boolean))
	default:
		return newError("type mismatch: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(tok token.Token, operator string, left, right *object.Integer) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero", tok)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(tok token.Token, operator string, left, right *object.Float) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError("division by zero", tok)
		}
		return &object.Float{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(tok token.Token, operator string, left, right *object.String) object.Object {
	switch operator {
	case "+":
		return &object.String{Value: left.Value + right.Value}
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return newError("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(tok token.Token, operator string, left, right *object.Boolean) object.Object {
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left.Value == right.Value)
	case "!=":
		return nativeBoolToBooleanObject(left.Value != right.Value)
	default:
		return newError("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalIndexExpression(tok token.Token, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(tok, left.(*object.Array), index.(*object.Integer))
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(tok, left.(*object.String), index.(*object.Integer))
	default:
		return newError("index operator not supported: %s", tok, left.Type())
	}
}

func evalArrayIndexExpression(tok token.Token, array *object.Array, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(array.Elements) - 1)

	if idx < 0 || idx > max {
		return newError("index out of bounds: array length is %d, got index %d", tok, len(array.Elements), idx)
	}

	return array.Elements[idx]
}

func evalStringIndexExpression(tok token.Token, str *object.String, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(str.Value) - 1)

	if idx < 0 || idx > max {
		return newError("index out of bounds: string length is %d, got index %d", tok, len(str.Value), idx)
	}

	return &object.String{Value: string(str.Value[idx])}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = Eval(we.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	// Create an enclosed environment for variables defined in the initializer
	forEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		initRes := Eval(fe.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = NULL

	for {
		if fe.Condition != nil {
			condition := Eval(fe.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fe.Body, forEnv)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

		if fe.Post != nil {
			postRes := Eval(fe.Post, forEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return result
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func applyFunction(tok token.Token, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return fn.Fn(tok.Line, tok.Column, args...)

	default:
		return newError("not a function: %s", tok, fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Name.Value, args[paramIdx])
		} else {
			env.Set(param.Name.Value, NULL)
		}
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		// Any non-null/non-false is truthy in RustX Starter
		return true
	}
}

// To prevent Go static analysis warnings about format strings:
// We use a helper function that takes line/col/args and produces a constant-looking format message.
func newError(format string, tok token.Token, args ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, args...)
	return &object.Error{Message: msg, Line: tok.Line, Column: tok.Column}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

// Global Built-in Functions
var builtins = map[string]*object.Builtin{
	"print": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Print(strings.Join(out, " "))
			return NULL
		},
	},
	"println": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			var out []string
			for _, arg := range args {
				out = append(out, arg.Inspect())
			}
			fmt.Println(strings.Join(out, " "))
			return NULL
		},
	},

	// Math Builtins
	"abs": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("abs() expected 1 argument, got %d", len(args)), Line: line, Column: col}
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
				return &object.Error{Message: "abs() argument must be an integer or float", Line: line, Column: col}
			}
		},
	},
	"pow": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("pow() expected 2 arguments, got %d", len(args)), Line: line, Column: col}
			}
			var base, exp float64
			switch arg := args[0].(type) {
			case *object.Integer:
				base = float64(arg.Value)
			case *object.Float:
				base = arg.Value
			default:
				return &object.Error{Message: "pow() first argument must be a number", Line: line, Column: col}
			}
			switch arg := args[1].(type) {
			case *object.Integer:
				exp = float64(arg.Value)
			case *object.Float:
				exp = arg.Value
			default:
				return &object.Error{Message: "pow() second argument must be a number", Line: line, Column: col}
			}
			return &object.Float{Value: math.Pow(base, exp)}
		},
	},
	"sin": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("sin() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			val, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "sin() argument must be a number", Line: line, Column: col}
			}
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("cos() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			val, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "cos() argument must be a number", Line: line, Column: col}
			}
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("tan() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			val, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "tan() argument must be a number", Line: line, Column: col}
			}
			return &object.Float{Value: math.Tan(val)}
		},
	},

	// String Builtins
	"len": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("len() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return &object.Error{Message: "argument to len() not supported", Line: line, Column: col}
			}
		},
	},
	"contains": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("contains() expected 2 arguments, got %d", len(args)), Line: line, Column: col}
			}
			s1, ok1 := args[0].(*object.String)
			s2, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "contains() arguments must be strings", Line: line, Column: col}
			}
			return nativeBoolToBooleanObject(strings.Contains(s1.Value, s2.Value))
		},
	},
	"split": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("split() expected 2 arguments, got %d", len(args)), Line: line, Column: col}
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "split() arguments must be strings", Line: line, Column: col}
			}
			parts := strings.Split(s.Value, sep.Value)
			var elements []object.Object
			for _, p := range parts {
				elements = append(elements, &object.String{Value: p})
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("join() expected 2 arguments, got %d", len(args)), Line: line, Column: col}
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "join() arguments must be (Array, String)", Line: line, Column: col}
			}
			var strParts []string
			for _, el := range arr.Elements {
				strParts = append(strParts, el.Inspect())
			}
			return &object.String{Value: strings.Join(strParts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("to_upper() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "to_upper() argument must be a string", Line: line, Column: col}
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("to_lower() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "to_lower() argument must be a string", Line: line, Column: col}
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},

	// Time Builtins
	"sleep": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("sleep() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			sec, ok := getFloatVal(args[0])
			if !ok {
				return &object.Error{Message: "sleep() argument must be a number of seconds", Line: line, Column: col}
			}
			time.Sleep(time.Duration(sec * float64(time.Second)))
			return NULL
		},
	},
	"now": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 0 {
				return &object.Error{Message: fmt.Sprintf("now() expected 0 arguments, got %d", len(args)), Line: line, Column: col}
			}
			return &object.Integer{Value: time.Now().Unix()}
		},
	},

	// Filesystem Builtins
	"read_file": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: fmt.Sprintf("read_file() expected 1 argument, got %d", len(args)), Line: line, Column: col}
			}
			path, ok := args[0].(*object.String)
			if !ok {
				return &object.Error{Message: "read_file() argument must be a string path", Line: line, Column: col}
			}
			content, err := os.ReadFile(path.Value)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("read_file() error: %s", err.Error()), Line: line, Column: col}
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(line, col int, args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: fmt.Sprintf("write_file() expected 2 arguments, got %d", len(args)), Line: line, Column: col}
			}
			path, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return &object.Error{Message: "write_file() arguments must be (String, String)", Line: line, Column: col}
			}
			err := os.WriteFile(path.Value, []byte(content.Value), 0644)
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("write_file() error: %s", err.Error()), Line: line, Column: col}
			}
			return NULL
		},
	},
}

func getFloatVal(obj object.Object) (float64, bool) {
	switch o := obj.(type) {
	case *object.Integer:
		return float64(o.Value), true
	case *object.Float:
		return o.Value, true
	default:
		return 0, false
	}
}
