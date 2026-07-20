package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
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
		return &object.Null{}

	case *ast.AssignmentStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		if n.Token.Literal == ":=" {
			env.Set(n.Name.Value, val)
		} else {
			_, updated := env.Update(n.Name.Value, val)
			if !updated {
				return newRuntimeError(n, "identifier not found: %s", n.Name.Value)
			}
		}
		return &object.Null{}

	case *ast.WhileStatement:
		var result object.Object = &object.Null{}
		for {
			cond := Eval(n.Condition, env)
			if isError(cond) {
				return cond
			}
			if !isTruthy(cond) {
				break
			}
			res := Eval(n.Body, env)
			if isError(res) {
				return res
			}
			if res.Type() == object.RETURN_VALUE_OBJ {
				return res
			}
			result = res
		}
		return result

	case *ast.ForStatement:
		loopEnv := object.NewEnclosedEnvironment(env)
		if n.Init != nil {
			initRes := Eval(n.Init, loopEnv)
			if isError(initRes) {
				return initRes
			}
		}

		var result object.Object = &object.Null{}
		for {
			if n.Condition != nil {
				cond := Eval(n.Condition, loopEnv)
				if isError(cond) {
					return cond
				}
				if !isTruthy(cond) {
					break
				}
			}

			res := Eval(n.Body, loopEnv)
			if isError(res) {
				return res
			}
			if res.Type() == object.RETURN_VALUE_OBJ {
				return res
			}
			result = res

			if n.Post != nil {
				postRes := Eval(n.Post, loopEnv)
				if isError(postRes) {
					return postRes
				}
			}
		}
		return result

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(n.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
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

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right, n)

	case *ast.InfixExpression:
		if n.Operator == "&&" || n.Operator == "||" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			return evalLogicalExpression(n, left, env)
		}

		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right, n)

	case *ast.PostfixExpression:
		ident, ok := n.Left.(*ast.Identifier)
		if !ok {
			return newRuntimeError(n, "postfix operators can only be applied to identifiers")
		}
		val, ok := env.Get(ident.Value)
		if !ok {
			return newRuntimeError(n, "identifier not found: %s", ident.Value)
		}

		var newVal object.Object
		switch v := val.(type) {
		case *object.Integer:
			if n.Operator == "++" {
				newVal = &object.Integer{Value: v.Value + 1}
			} else {
				newVal = &object.Integer{Value: v.Value - 1}
			}
		case *object.Float:
			if n.Operator == "++" {
				newVal = &object.Float{Value: v.Value + 1.0}
			} else {
				newVal = &object.Float{Value: v.Value - 1.0}
			}
		default:
			return newRuntimeError(n, "postfix operator %s cannot be applied to type %s", n.Operator, val.Type())
		}

		env.Update(ident.Value, newVal)
		return newVal

	case *ast.IfExpression:
		cond := Eval(n.Condition, env)
		if isError(cond) {
			return cond
		}
		if isTruthy(cond) {
			return Eval(n.Consequence, env)
		} else if n.Alternative != nil {
			return Eval(n.Alternative, env)
		}
		return &object.Null{}

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fn := &object.Function{Parameters: params, Body: body, Env: env}
		if n.Name != nil {
			env.Set(n.Name.Value, fn)
		}
		return fn

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(n, function, args)
	}

	return &object.Null{}
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = &object.Null{}

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

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
	var result object.Object = &object.Null{}

	// Block statements create their own lexical scope, so updates to outer variables affect outer variables,
	// but local variables are scoped inside.
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, stmt := range block.Statements {
		result = Eval(stmt, blockEnv)

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
	return &object.Boolean{Value: input}
}

func isTruthy(obj object.Object) bool {
	switch o := obj.(type) {
	case *object.Boolean:
		return o.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newRuntimeError(node ast.Node, format string, a ...interface{}) *object.Error {
	tok := node.GetToken()
	msg := fmt.Sprintf(format, a...)
	return &object.Error{
		Message: fmt.Sprintf("RuntimeError:\n%s\nat line %d column %d", msg, tok.Line, tok.Col),
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newRuntimeError(node, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(operator string, right object.Object, node ast.Node) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, node)
	default:
		return newRuntimeError(node, "unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	if b, ok := right.(*object.Boolean); ok {
		return &object.Boolean{Value: !b.Value}
	}
	if _, ok := right.(*object.Null); ok {
		return &object.Boolean{Value: true}
	}
	return &object.Boolean{Value: false}
}

func evalMinusPrefixOperatorExpression(right object.Object, node ast.Node) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newRuntimeError(node, "unknown operator: -%s", right.Type())
	}
}

func evalLogicalExpression(node *ast.InfixExpression, left object.Object, env *object.Environment) object.Object {
	leftVal := isTruthy(left)

	if node.Operator == "&&" {
		if !leftVal {
			return &object.Boolean{Value: false}
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return &object.Boolean{Value: isTruthy(right)}
	}

	if node.Operator == "||" {
		if leftVal {
			return &object.Boolean{Value: true}
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return &object.Boolean{Value: isTruthy(right)}
	}

	return newRuntimeError(node, "unknown operator: %s", node.Operator)
}

func evalInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left.(*object.Integer), right.(*object.Integer), node)
	case left.Type() == object.FLOAT_OBJ || right.Type() == object.FLOAT_OBJ:
		if left.Type() == object.INTEGER_OBJ || left.Type() == object.FLOAT_OBJ {
			if right.Type() == object.INTEGER_OBJ || right.Type() == object.FLOAT_OBJ {
				return evalFloatInfixExpression(operator, left, right, node)
			}
		}
		return newRuntimeError(node, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left.(*object.String), right.(*object.String), node)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case left.Type() != right.Type():
		return newRuntimeError(node, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newRuntimeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		// Try coercing float and int if necessary
		if left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ {
			return float64(left.(*object.Integer).Value) == right.(*object.Float).Value
		}
		if left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ {
			return left.(*object.Float).Value == float64(right.(*object.Integer).Value)
		}
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

func evalIntegerInfixExpression(operator string, left, right *object.Integer, node ast.Node) object.Object {
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
			return newRuntimeError(node, "division by zero")
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
		return newRuntimeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	var leftVal, rightVal float64

	if left.Type() == object.INTEGER_OBJ {
		leftVal = float64(left.(*object.Integer).Value)
	} else if left.Type() == object.FLOAT_OBJ {
		leftVal = left.(*object.Float).Value
	}

	if right.Type() == object.INTEGER_OBJ {
		rightVal = float64(right.(*object.Integer).Value)
	} else if right.Type() == object.FLOAT_OBJ {
		rightVal = right.(*object.Float).Value
	}

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newRuntimeError(node, "division by zero")
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
		return newRuntimeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right *object.String, node ast.Node) object.Object {
	leftVal := left.Value
	rightVal := right.Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newRuntimeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
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
		return evalArrayAtIndex(node, left.(*object.Array), index.(*object.Integer))
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringAtIndex(node, left.(*object.String), index.(*object.Integer))
	default:
		return newRuntimeError(node, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayAtIndex(node ast.Node, array *object.Array, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(array.Elements) - 1)
	if idx < 0 || idx > max {
		return newRuntimeError(node, "index out of bounds: array length %d, got index %d", len(array.Elements), idx)
	}
	return array.Elements[idx]
}

func evalStringAtIndex(node ast.Node, str *object.String, index *object.Integer) object.Object {
	idx := index.Value
	max := int64(len(str.Value) - 1)
	if idx < 0 || idx > max {
		return newRuntimeError(node, "index out of bounds: string length %d, got index %d", len(str.Value), idx)
	}
	return &object.String{Value: string(str.Value[idx])}
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
		return newRuntimeError(node, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		if i < len(args) {
			env.Set(param.Name.Value, args[i])
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
