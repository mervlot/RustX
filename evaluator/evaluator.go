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
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)

	case *ast.InfixExpression:
		if n.Operator == "=" || n.Operator == ":=" {
			return evalAssignment(n.Left, n.Operator, n.Right, env)
		}

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
		return evalInfixExpression(n.Operator, left, right)

	case *ast.PostfixExpression:
		return evalPostfixExpression(n.Left, n.Operator, env)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fun := &object.Function{Parameters: params, Env: env, Body: body}
		if n.Name != "" {
			env.Set(n.Name, fun)
		}
		return fun

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

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
		return evalIndexExpression(left, index)
	}

	return NULL
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Value)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError("unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Promote integer to float
		lVal := float64(left.(*object.Integer).Value)
		rVal := right.(*object.Float).Value
		return evalFloatValuesInfixExpression(operator, lVal, rVal)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Promote integer to float
		lVal := left.(*object.Float).Value
		rVal := float64(right.(*object.Integer).Value)
		return evalFloatValuesInfixExpression(operator, lVal, rVal)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(left, right object.Object) bool {
	if left.Type() != right.Type() {
		// Try numeric coercion
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
		return left == right
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
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
			return newError("division by zero")
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	return evalFloatValuesInfixExpression(operator, leftVal, rightVal)
}

func evalFloatValuesInfixExpression(operator string, leftVal, rightVal float64) object.Object {
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError("division by zero")
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
		return newError("unknown operator: FLOAT %s FLOAT", operator)
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
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
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalAssignment(left ast.Expression, operator string, right ast.Expression, env *object.Environment) object.Object {
	val := Eval(right, env)
	if isError(val) {
		return val
	}

	// Handle standard assignments like x = val, x := val, and index assignments array[index] = val
	switch node := left.(type) {
	case *ast.Identifier:
		if operator == ":=" {
			// Shorthand assignment always declares/shadows in current environment scope
			env.Set(node.Value, val)
			return val
		}

		// Standard assign `=` updates the closest defined variable scope
		ok := env.Update(node.Value, val)
		if !ok {
			// If not declared, default to creating/defining in current scope
			env.Set(node.Value, val)
		}
		return val

	case *ast.IndexExpression:
		// Array index assignment
		target := Eval(node.Left, env)
		if isError(target) {
			return target
		}
		idxObj := Eval(node.Index, env)
		if isError(idxObj) {
			return idxObj
		}

		arr, ok := target.(*object.Array)
		if !ok {
			return newError("cannot assign to index of non-array object: %s", target.Type())
		}

		idx, ok := idxObj.(*object.Integer)
		if !ok {
			return newError("array index must be integer, got %s", idxObj.Type())
		}

		i := idx.Value
		if i < 0 || i >= int64(len(arr.Elements)) {
			return newError("array index out of bounds: %d", i)
		}

		arr.Elements[i] = val
		return val
	}

	return newError("invalid assignment target: %T", left)
}

func evalPostfixExpression(left ast.Expression, operator string, env *object.Environment) object.Object {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		return newError("postfix operators %s can only be applied to identifiers", operator)
	}

	val, found := env.Get(ident.Value)
	if !found {
		return newError("identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		var newVal int64
		if operator == "++" {
			newVal = v.Value + 1
		} else {
			newVal = v.Value - 1
		}
		newObj := &object.Integer{Value: newVal}
		env.Update(ident.Value, newObj)
		return newObj
	case *object.Float:
		var newVal float64
		if operator == "++" {
			newVal = v.Value + 1.0
		} else {
			newVal = v.Value - 1.0
		}
		newObj := &object.Float{Value: newVal}
		env.Update(ident.Value, newObj)
		return newObj
	default:
		return newError("cannot increment/decrement non-numeric type: %s", val.Type())
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
	var lastVal object.Object = NULL

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		// Block has its own enclosed scope so variable updates persist but declarations stay scoped
		blockEnv := object.NewEnclosedEnvironment(env)
		val := Eval(we.Body, blockEnv)
		if isError(val) {
			return val
		}

		if val != nil {
			if val.Type() == object.RETURN_VALUE_OBJ {
				return val
			}
		}
		lastVal = val
	}

	return lastVal
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	// For loop header has its own scope so variable initializations don't leak outside
	forEnv := object.NewEnclosedEnvironment(env)

	if fe.Initializer != nil {
		initVal := Eval(fe.Initializer, forEnv)
		if isError(initVal) {
			return initVal
		}
	}

	var lastVal object.Object = NULL

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

		// Body block has its own scope nested under the loop header scope
		bodyEnv := object.NewEnclosedEnvironment(forEnv)
		val := Eval(fe.Body, bodyEnv)
		if isError(val) {
			return val
		}

		if val != nil {
			if val.Type() == object.RETURN_VALUE_OBJ {
				return val
			}
		}
		lastVal = val

		if fe.Post != nil {
			postVal := Eval(fe.Post, forEnv)
			if isError(postVal) {
				return postVal
			}
		}
	}

	return lastVal
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

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
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

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError("array index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newError("string index out of bounds: %d", idx)
	}

	return &object.String{Value: string(strObject.Value[idx])}
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

// newError uses a constant format string to satisfy Go static analysis check
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
