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

	case *ast.BlockStatement:
		return evalBlockStatement(n, env)

	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		if val == nil {
			return &object.ReturnValue{Value: NULL}
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return NULL

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.Boolean:
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
		return evalPrefixExpression(n, n.Operator, right)

	case *ast.InfixExpression:
		// Short-circuiting for logical AND / OR
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
		return evalPostfixExpression(n, env)

	case *ast.AssignExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}

		if n.Operator == ":=" {
			// Shorthand assign always sets in the CURRENT environment
			env.Set(n.Name.Value, val)
			return val
		} else {
			// Regular assign (=) updates recursively
			ok := env.Update(n.Name.Value, val)
			if !ok {
				return nodeError(n, "identifier '%s' is not defined", n.Name.Value)
			}
			return val
		}

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
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
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(n, function, args)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
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

func evalIndexExpression(node ast.Node, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		array := left.(*object.Array)
		idx := index.(*object.Integer).Value
		max := int64(len(array.Elements) - 1)
		if idx < 0 || idx > max {
			return nodeError(node, "array index out of bounds: %d", idx)
		}
		return array.Elements[idx]
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		str := left.(*object.String)
		idx := index.(*object.Integer).Value
		max := int64(len(str.Value) - 1)
		if idx < 0 || idx > max {
			return nodeError(node, "string index out of bounds: %d", idx)
		}
		return &object.String{Value: string(str.Value[idx])}
	default:
		return nodeError(node, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return nodeError(node, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(node ast.Node, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right)
	default:
		return nodeError(node, "unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(node ast.Node, right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return nodeError(node, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(node, operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Cast int to float
		lFloat := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, operator, lFloat, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Cast int to float
		rFloat := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(node, operator, left, rFloat)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return nodeError(node, "type mismatch or unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
			return nodeError(node, "division by zero")
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
		return nodeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
			return nodeError(node, "division by zero")
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
		return nodeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return nodeError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
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
	case *object.Null:
		return true
	default:
		return false
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return nodeError(node, "postfix operators can only be applied to identifiers")
	}
	val, ok := env.Get(ident.Value)
	if !ok {
		return nodeError(node, "identifier not found: %s", ident.Value)
	}
	switch v := val.(type) {
	case *object.Integer:
		var newVal *object.Integer
		if node.Operator == "++" {
			newVal = &object.Integer{Value: v.Value + 1}
		} else {
			newVal = &object.Integer{Value: v.Value - 1}
		}
		env.Update(ident.Value, newVal)
		return newVal
	case *object.Float:
		var newVal *object.Float
		if node.Operator == "++" {
			newVal = &object.Float{Value: v.Value + 1.0}
		} else {
			newVal = &object.Float{Value: v.Value - 1.0}
		}
		env.Update(ident.Value, newVal)
		return newVal
	default:
		return nodeError(node, "postfix operators can only be applied to integer or float, got %s", v.Type())
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
	// Create a sub-environment for the entire loop so variables defined in init are scoped correctly
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

func applyFunction(node ast.Node, fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		if len(args) != len(f.Parameters) {
			return nodeError(node, "wrong number of arguments. got=%d, want=%d", len(args), len(f.Parameters))
		}
		extendedEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(args...)

	default:
		return nodeError(node, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Name.Value, args[paramIdx])
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
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func nodeError(node ast.Node, format string, a ...interface{}) *object.Error {
	tok := node.GetToken()
	msg := fmt.Sprintf(format, a...)
	formatted := fmt.Sprintf("RuntimeError: %s at line %d column %d", msg, tok.Line, tok.Column)
	return &object.Error{Message: formatted}
}
