package evaluator

import (
	"fmt"
	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
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

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

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
		return evalPrefixExpression(node.Operator, right, node)

	case *ast.InfixExpression:
		if node.Operator == "&&" {
			return evalAndExpression(node.Left, node.Right, env)
		}
		if node.Operator == "||" {
			return evalOrExpression(node.Left, node.Right, env)
		}

		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right, env, node)

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

		return applyFunction(function, args, node)

	case *ast.WhileExpression:
		return evalWhileExpression(node, env)

	case *ast.ForExpression:
		return evalForExpression(node, env)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

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
		return evalIndexExpression(left, index, node)
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

func evalPrefixExpression(operator string, right object.Object, node ast.Node) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, node)
	default:
		return newError(node, "unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, node ast.Node) object.Object {
	if right.Type() == object.INTEGER_OBJ {
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
	if right.Type() == object.FLOAT_OBJ {
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	}
	return newError(node, "unknown operator: -%s", right.Type())
}

func evalInfixExpression(operator string, left, right object.Object, env *object.Environment, node ast.Node) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, node)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, node)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		rightFloat := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, left, rightFloat, node)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		leftFloat := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, leftFloat, right, node)
	case operator == "==":
		return evalEqualityExpression(left, right)
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalEqualityExpression(left, right).Value)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, node)
	case operator == "=" || operator == ":=":
		return evalAssignmentExpression(operator, left, right, env, node)
	case left.Type() != right.Type():
		return newError(node, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalEqualityExpression(left, right object.Object) *object.Boolean {
	if left.Type() != right.Type() {
		return FALSE
	}

	switch left.Type() {
	case object.INTEGER_OBJ:
		return nativeBoolToBooleanObject(left.(*object.Integer).Value == right.(*object.Integer).Value)
	case object.FLOAT_OBJ:
		return nativeBoolToBooleanObject(left.(*object.Float).Value == right.(*object.Float).Value)
	case object.STRING_OBJ:
		return nativeBoolToBooleanObject(left.(*object.String).Value == right.(*object.String).Value)
	case object.BOOLEAN_OBJ:
		return nativeBoolToBooleanObject(left.(*object.Boolean).Value == right.(*object.Boolean).Value)
	case object.NULL_OBJ:
		return TRUE
	default:
		return nativeBoolToBooleanObject(left == right)
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError(node, "division by zero")
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
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError(node, "division by zero")
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
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalAndExpression(left, right ast.Expression, env *object.Environment) object.Object {
	leftRes := Eval(left, env)
	if isError(leftRes) {
		return leftRes
	}
	if !isTruthy(leftRes) {
		return FALSE
	}
	rightRes := Eval(right, env)
	if isError(rightRes) {
		return rightRes
	}
	return nativeBoolToBooleanObject(isTruthy(rightRes))
}

func evalOrExpression(left, right ast.Expression, env *object.Environment) object.Object {
	leftRes := Eval(left, env)
	if isError(leftRes) {
		return leftRes
	}
	if isTruthy(leftRes) {
		return TRUE
	}
	rightRes := Eval(right, env)
	if isError(rightRes) {
		return rightRes
	}
	return nativeBoolToBooleanObject(isTruthy(rightRes))
}

func evalAssignmentExpression(operator string, left, right object.Object, env *object.Environment, node ast.Node) object.Object {
	infix, ok := node.(*ast.InfixExpression)
	if !ok {
		return newError(node, "invalid assignment")
	}

	ident, ok := infix.Left.(*ast.Identifier)
	if !ok {
		return newError(node, "left hand side of assignment must be an identifier")
	}

	if operator == "=" {
		_, updated := env.Update(ident.Value, right)
		if !updated {
			return newError(node, "identifier %s not defined", ident.Value)
		}
	} else {
		env.Set(ident.Value, right)
	}

	return right
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

	return newError(node, "identifier not found: %s", node.Value)
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
		return true
	}
}

func newError(node ast.Node, format string, a ...interface{}) *object.Error {
	tok := node.GetToken()
	msg := fmt.Sprintf(format, a...)
	return &object.Error{Message: fmt.Sprintf("%s at line %d column %d", msg, tok.Line, tok.Column)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
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

func applyFunction(fn object.Object, args []object.Object, node ast.Node) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError(node, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx < len(args) {
			env.Set(param.Value, args[paramIdx])
		} else {
			env.Set(param.Value, NULL)
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

		if isTruthy(condition) {
			result = Eval(we.Body, env)
			if result != nil {
				rt := result.Type()
				if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
					return result
				}
			}
		} else {
			break
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
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

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node, "postfix operator must be used on identifier")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(node, "identifier not found: %s", ident.Value)
	}

	if val.Type() != object.INTEGER_OBJ {
		return newError(node, "postfix operator can only be used on integers")
	}

	currVal := val.(*object.Integer).Value
	var newVal int64
	if node.Operator == "++" {
		newVal = currVal + 1
	} else {
		newVal = currVal - 1
	}

	newObj := &object.Integer{Value: newVal}
	env.Update(ident.Value, newObj)

	return newObj
}

func evalIndexExpression(left, index object.Object, node ast.Node) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, node)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, node)
	default:
		return newError(node, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, node ast.Node) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, node ast.Node) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

// Builtins
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
	// Strings
	"to_upper": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "to_upper requires 1 string argument"}
			}
			return &object.String{Value: strings.ToUpper(args[0].(*object.String).Value)}
		},
	},
	"to_lower": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "to_lower requires 1 string argument"}
			}
			return &object.String{Value: strings.ToLower(args[0].(*object.String).Value)}
		},
	},
	"contains": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "contains requires 2 string arguments"}
			}
			return nativeBoolToBooleanObject(strings.Contains(args[0].(*object.String).Value, args[1].(*object.String).Value))
		},
	},
	"split": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "split requires 2 string arguments"}
			}
			res := strings.Split(args[0].(*object.String).Value, args[1].(*object.String).Value)
			elements := []object.Object{}
			for _, s := range res {
				elements = append(elements, &object.String{Value: s})
			}
			return &object.Array{Elements: elements}
		},
	},
	"join": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 || args[0].Type() != object.ARRAY_OBJ || args[1].Type() != object.STRING_OBJ {
				return &object.Error{Message: "join requires (array, string) arguments"}
			}
			arr := args[0].(*object.Array)
			sep := args[1].(*object.String).Value
			elements := []string{}
			for _, e := range arr.Elements {
				elements = append(elements, e.Inspect())
			}
			return &object.String{Value: strings.Join(elements, sep)}
		},
	},
	// Math
	"abs": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Message: "abs requires 1 argument"}
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
			return &object.Error{Message: "abs requires numeric argument"}
		},
	},
	"pow": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return &object.Error{Message: "pow requires 2 arguments"}
			}
			v1 := 0.0
			v2 := 0.0
			if args[0].Type() == object.INTEGER_OBJ {
				v1 = float64(args[0].(*object.Integer).Value)
			} else if args[0].Type() == object.FLOAT_OBJ {
				v1 = args[0].(*object.Float).Value
			}
			if args[1].Type() == object.INTEGER_OBJ {
				v2 = float64(args[1].(*object.Integer).Value)
			} else if args[1].Type() == object.FLOAT_OBJ {
				v2 = args[1].(*object.Float).Value
			}
			return &object.Float{Value: math.Pow(v1, v2)}
		},
	},
	"sin": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "sin requires 1 argument"}}
			val := 0.0
			if args[0].Type() == object.INTEGER_OBJ { val = float64(args[0].(*object.Integer).Value) } else if args[0].Type() == object.FLOAT_OBJ { val = args[0].(*object.Float).Value }
			return &object.Float{Value: math.Sin(val)}
		},
	},
	"cos": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "cos requires 1 argument"}}
			val := 0.0
			if args[0].Type() == object.INTEGER_OBJ { val = float64(args[0].(*object.Integer).Value) } else if args[0].Type() == object.FLOAT_OBJ { val = args[0].(*object.Float).Value }
			return &object.Float{Value: math.Cos(val)}
		},
	},
	"tan": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 { return &object.Error{Message: "tan requires 1 argument"}}
			val := 0.0
			if args[0].Type() == object.INTEGER_OBJ { val = float64(args[0].(*object.Integer).Value) } else if args[0].Type() == object.FLOAT_OBJ { val = args[0].(*object.Float).Value }
			return &object.Float{Value: math.Tan(val)}
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
				return &object.Error{Message: "sleep requires 1 integer argument (ms)"}
			}
			time.Sleep(time.Duration(args[0].(*object.Integer).Value) * time.Millisecond)
			return NULL
		},
	},
	// FS
	"read_file": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 || args[0].Type() != object.STRING_OBJ {
				return &object.Error{Message: "read_file requires 1 string argument"}
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
				return &object.Error{Message: "write_file requires (string, string) arguments"}
			}
			err := os.WriteFile(args[0].(*object.String).Value, []byte(args[1].(*object.String).Value), 0644)
			if err != nil {
				return &object.Error{Message: err.Error()}
			}
			return TRUE
		},
	},
}
