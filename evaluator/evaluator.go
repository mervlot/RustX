package evaluator

import (
	"fmt"
	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
	"math"
	"os"
	"strings"
	"time"
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

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node, right)

	case *ast.InfixExpression:
		if node.Operator == "&&" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return left
			}
			return Eval(node.Right, env)
		}
		if node.Operator == "||" {
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return left
			}
			return Eval(node.Right, env)
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node, left, right)

	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		return evalPostfixExpression(node, left, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
		if node.Name != "" {
			env.Set(node.Name, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(node, function, args)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.AssignExpression:
		return evalAssignExpression(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(node, left, index)

	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

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
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch node.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right)
	default:
		return newErrorAt(node.Token, "unknown operator: %s%s", node.Operator, right.Type())
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
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newErrorAt(node.Token, "unknown operator: -%s", right.Type())
}

func evalInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(node, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(node, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		rightFloat := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, left, rightFloat)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		leftFloat := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, leftFloat, right)
	case node.Operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case node.Operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, left, right)
	case left.Type() != right.Type():
		return newErrorAt(node.Token, "type mismatch: %s %s %s", left.Type(), node.Operator, right.Type())
	default:
		return newErrorAt(node.Token, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalIntegerInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch node.Operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newErrorAt(node.Token, "division by zero")
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
	case "&&":
		return nativeBoolToBooleanObject(leftVal != 0 && rightVal != 0)
	case "||":
		return nativeBoolToBooleanObject(leftVal != 0 || rightVal != 0)
	default:
		return newErrorAt(node.Token, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalBooleanInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch node.Operator {
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case "&&":
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case "||":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	default:
		return newErrorAt(node.Token, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalFloatInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch node.Operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
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
		return newErrorAt(node.Token, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalStringInfixExpression(node *ast.InfixExpression, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch node.Operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newErrorAt(node.Token, "unknown operator: %s %s %s", left.Type(), node.Operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, left object.Object, env *object.Environment) object.Object {
	if node.Operator == "++" {
		if left.Type() == object.INTEGER_OBJ {
			val := left.(*object.Integer).Value
			newObj := &object.Integer{Value: val + 1}
			if ident, ok := node.Left.(*ast.Identifier); ok {
				env.Update(ident.Value, newObj)
			}
			return left
		}
	}
	return newErrorAt(node.Token, "unknown operator: %s%s", left.Type(), node.Operator)
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newErrorAt(node.Token, "identifier not found: %s", node.Value)
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

func applyFunction(node *ast.CallExpression, fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(node, args...)
	default:
		return newErrorAt(node.Token, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
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
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	// For loops create a new scope
	childEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		Eval(fe.Init, childEnv)
	}

	var result object.Object = NULL

	for {
		if fe.Condition != nil {
			condition := Eval(fe.Condition, childEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = Eval(fe.Body, childEnv)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

		if fe.Post != nil {
			Eval(fe.Post, childEnv)
		}
	}

	return result
}

func evalAssignExpression(ae *ast.AssignExpression, env *object.Environment) object.Object {
	val := Eval(ae.Value, env)
	if isError(val) {
		return val
	}

	if ident, ok := ae.Left.(*ast.Identifier); ok {
		if ae.Token.Literal == ":=" {
			env.Set(ident.Value, val)
		} else {
			if !env.Update(ident.Value, val) {
				return newErrorAt(ident.Token, "identifier not found: %s", ident.Value)
			}
		}
		return val
	}

	if indexExp, ok := ae.Left.(*ast.IndexExpression); ok {
		left := Eval(indexExp.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(indexExp.Index, env)
		if isError(index) {
			return index
		}

		if array, ok := left.(*object.Array); ok {
			idxObj, ok := index.(*object.Integer)
			if !ok {
				return newErrorAt(indexExp.Index.GetToken(), "array index must be an integer")
			}
			idx := idxObj.Value
			if idx < 0 || idx >= int64(len(array.Elements)) {
				return newErrorAt(indexExp.Index.GetToken(), "index out of bounds")
			}
			array.Elements[idx] = val
			return val
		}
		return newErrorAt(indexExp.Left.GetToken(), "target is not an array")
	}

	return newErrorAt(ae.Token, "invalid assignment target")
}

func evalIndexExpression(node *ast.IndexExpression, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newErrorAt(node.Token, "index operator not supported: %s", left.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string(stringObject.Value[idx])}
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
		if obj.Type() == object.INTEGER_OBJ {
			return obj.(*object.Integer).Value != 0
		}
		return true
	}
}

func newErrorAt(tok token.Token, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		Line:    tok.Line,
		Col:     tok.Col,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

type BuiltinFunction func(node *ast.CallExpression, args ...object.Object) object.Object

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newErrorAt(node.Token, "argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"println": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			for i, arg := range args {
				if i > 0 {
					fmt.Print(" ")
				}
				fmt.Print(arg.Inspect())
			}
			fmt.Println()
			return NULL
		},
	},
	"print": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Print(arg.Inspect())
			}
			return NULL
		},
	},
	"push": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newErrorAt(node.Token, "argument to `push` must be ARRAY, got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			newElements := append(arr.Elements, args[1])
			return &object.Array{Elements: newElements}
		},
	},
	// strings
	"split": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorAt(node.Token, "arguments to `split` must be STRING")
			}
			parts := strings.Split(s.Value, sep.Value)
			objs := make([]object.Object, len(parts))
			for i, p := range parts {
				objs[i] = &object.String{Value: p}
			}
			return &object.Array{Elements: objs}
		},
	},
	"join": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=2", len(args))
			}
			arr, ok1 := args[0].(*object.Array)
			sep, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorAt(node.Token, "arguments to `join` must be (ARRAY, STRING)")
			}
			parts := make([]string, len(arr.Elements))
			for i, e := range arr.Elements {
				parts[i] = e.Inspect()
			}
			return &object.String{Value: strings.Join(parts, sep.Value)}
		},
	},
	"to_upper": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newErrorAt(node.Token, "argument to `to_upper` must be STRING")
			}
			return &object.String{Value: strings.ToUpper(s.Value)}
		},
	},
	"to_lower": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			s, ok := args[0].(*object.String)
			if !ok {
				return newErrorAt(node.Token, "argument to `to_lower` must be STRING")
			}
			return &object.String{Value: strings.ToLower(s.Value)}
		},
	},
	"contains": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=2", len(args))
			}
			s, ok1 := args[0].(*object.String)
			sub, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorAt(node.Token, "arguments to `contains` must be STRING")
			}
			return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
		},
	},
	// math
	"abs": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
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
			return newErrorAt(node.Token, "argument to `abs` must be INTEGER or FLOAT")
		},
	},
	"pow": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=2", len(args))
			}
			f1, ok1 := args[0].(*object.Float)
			f2, ok2 := args[1].(*object.Float)
			if !ok1 || !ok2 {
				return newErrorAt(node.Token, "arguments to `pow` must be FLOAT")
			}
			return &object.Float{Value: math.Pow(f1.Value, f2.Value)}
		},
	},
	"sin": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			f, ok := args[0].(*object.Float)
			if !ok {
				return newErrorAt(node.Token, "argument to `sin` must be FLOAT")
			}
			return &object.Float{Value: math.Sin(f.Value)}
		},
	},
	"cos": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			f, ok := args[0].(*object.Float)
			if !ok {
				return newErrorAt(node.Token, "argument to `cos` must be FLOAT")
			}
			return &object.Float{Value: math.Cos(f.Value)}
		},
	},
	"tan": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			f, ok := args[0].(*object.Float)
			if !ok {
				return newErrorAt(node.Token, "argument to `tan` must be FLOAT")
			}
			return &object.Float{Value: math.Tan(f.Value)}
		},
	},
	// time
	"now": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			return &object.Integer{Value: time.Now().Unix()}
		},
	},
	"sleep": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			ms, ok := args[0].(*object.Integer)
			if !ok {
				return newErrorAt(node.Token, "argument to `sleep` must be INTEGER")
			}
			time.Sleep(time.Duration(ms.Value) * time.Millisecond)
			return NULL
		},
	},
	"read_file": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=1", len(args))
			}
			filename, ok := args[0].(*object.String)
			if !ok {
				return newErrorAt(node.Token, "argument to `read_file` must be STRING")
			}
			content, err := os.ReadFile(filename.Value)
			if err != nil {
				return newErrorAt(node.Token, "could not read file %s: %s", filename.Value, err)
			}
			return &object.String{Value: string(content)}
		},
	},
	"write_file": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorAt(node.Token, "wrong number of arguments. got=%d, want=2", len(args))
			}
			filename, ok1 := args[0].(*object.String)
			content, ok2 := args[1].(*object.String)
			if !ok1 || !ok2 {
				return newErrorAt(node.Token, "arguments to `write_file` must be STRING")
			}
			err := os.WriteFile(filename.Value, []byte(content.Value), 0644)
			if err != nil {
				return newErrorAt(node.Token, "could not write to file %s: %s", filename.Value, err)
			}
			return NULL
		},
	},
	"exit": {
		Fn: func(node *ast.CallExpression, args ...object.Object) object.Object {
			code := 0
			if len(args) > 0 {
				if c, ok := args[0].(*object.Integer); ok {
					code = int(c.Value)
				}
			}
			os.Exit(code)
			return NULL
		},
	},
}
