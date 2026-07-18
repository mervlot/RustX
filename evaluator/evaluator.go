package evaluator

import (
	"fmt"
	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(n, env)

	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return val

	case *ast.AssignStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		_, ok := env.Update(n.Name.Value, val)
		if !ok {
			return newError(n, "identifier not found: %s", n.Name.Value)
		}
		return val

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.ArrayLiteral:
		elements := evalExpressions(n.Elements, env)
		if len(elements) > 0 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(n.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(n, left, index)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n, n.Operator, right)

	case *ast.InfixExpression:
		// Short-circuit logic for && and ||
		if n.Operator == "&&" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return FALSE
			}
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}
		if n.Operator == "||" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return TRUE
			}
			right := Eval(n.Right, env)
			if isError(right) {
				return right
			}
			return nativeBoolToBooleanObject(isTruthy(right))
		}

		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n, n.Operator, left, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(n, n.Left, n.Operator, env)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fn := &object.Function{Parameters: params, Body: body, Env: env}
		if n.Name != "" {
			env.Set(n.Name, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) > 0 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(n, function, args)
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	// Initialize result to NULL to ensure consistent behavior in empty files
	var result object.Object = NULL

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	// Block statements introduce local scope for control flows / functions
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node, "identifier not found: %s", node.Value)
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

func evalIndexExpression(node ast.Node, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(node, left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(node, left, index)
	default:
		return newError(node, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(node ast.Node, array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError(node, "index out of bounds: length %d, got %d", len(arrayObject.Elements), idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(node ast.Node, str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return newError(node, "index out of bounds: length %d, got %d", len(stringObject.Value), idx)
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func evalPrefixExpression(node ast.Node, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right)
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
		if isTruthy(right) {
			return FALSE
		}
		return TRUE
	}
}

func evalMinusPrefixOperatorExpression(node ast.Node, right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError(node, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(node, operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Coerce left to float
		coercedLeft := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, operator, coercedLeft, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Coerce right to float
		coercedRight := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, operator, left, coercedRight)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newError(node, "type mismatch or unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
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

func evalFloatInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
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
		if rightVal == 0.0 {
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

func evalStringInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
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

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		return false
	}

	switch l := left.(type) {
	case *object.Integer:
		return l.Value == right.(*object.Integer).Value
	case *object.Float:
		return l.Value == right.(*object.Float).Value
	case *object.Boolean:
		return l.Value == right.(*object.Boolean).Value
	case *object.String:
		return l.Value == right.(*object.String).Value
	default:
		return left == right
	}
}

func evalPostfixExpression(node ast.Node, left ast.Expression, operator string, env *object.Environment) object.Object {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return newError(node, "postfix operator can only be applied to identifiers")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(node, "identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		var newVal int64
		if operator == "++" {
			newVal = v.Value + 1
		} else if operator == "--" {
			newVal = v.Value - 1
		} else {
			return newError(node, "unknown postfix operator: %s", operator)
		}
		newObj := &object.Integer{Value: newVal}
		env.Update(ident.Value, newObj)
		return newObj

	case *object.Float:
		var newVal float64
		if operator == "++" {
			newVal = v.Value + 1.0
		} else if operator == "--" {
			newVal = v.Value - 1.0
		} else {
			return newError(node, "unknown postfix operator: %s", operator)
		}
		newObj := &object.Float{Value: newVal}
		env.Update(ident.Value, newObj)
		return newObj

	default:
		return newError(node, "postfix operator cannot be applied to type: %s", val.Type())
	}
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
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	var result object.Object = NULL

	// For loop introduces a block environment
	forEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		initRes := Eval(fe.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

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
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
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

func applyFunction(node ast.Node, fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(args...)

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

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		switch o := obj.(type) {
		case *object.Boolean:
			return o.Value
		default:
			return true
		}
	}
}

func newError(node ast.Node, format string, a ...interface{}) *object.Error {
	tok := node.GetToken()
	msg := fmt.Sprintf(format, a...)
	return &object.Error{
		Message: fmt.Sprintf("%s at line %d column %d", msg, tok.Line, tok.Col),
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
