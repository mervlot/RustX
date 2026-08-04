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
		return NULL

	// Expressions
	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.ArrayLiteral:
		elems := evalExpressions(n.Elements, env)
		if len(elems) == 1 && isError(elems[0]) {
			return elems[0]
		}
		return &object.Array{Elements: elems}

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n, n.Operator, right)

	case *ast.InfixExpression:
		// Handle logical short-circuiting first
		if n.Operator == "&&" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if !isTruthy(left) {
				return left
			}
			return Eval(n.Right, env)
		}
		if n.Operator == "||" {
			left := Eval(n.Left, env)
			if isError(left) {
				return left
			}
			if isTruthy(left) {
				return left
			}
			return Eval(n.Right, env)
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
		return evalAssignExpression(n, env)

	case *ast.ColonAssignExpression:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		ident, ok := n.Left.(*ast.Identifier)
		if !ok {
			return newError(n.GetToken().Line, n.GetToken().Col, "cannot assign to non-identifier")
		}
		env.Set(ident.Value, val)
		return val

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		funcObj := &object.Function{Name: n.Name, Parameters: params, Body: body, Env: env}
		if n.Name != "" {
			env.Set(n.Name, funcObj)
		}
		return funcObj

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
	}

	return NULL
}

func newError(line int, col int, format string, args ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, args...),
		Line:    line,
		Col:     col,
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
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

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL // Initialize to NULL for empty files to prevent pointer issues.

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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError(node.GetToken().Line, node.GetToken().Col, "identifier not found: %s", node.Value)
}

func evalPrefixExpression(node ast.Node, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(node, right)
	default:
		return newError(node.GetToken().Line, node.GetToken().Col, "unknown operator: %s%s", operator, right.Type())
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
		return newError(node.GetToken().Line, node.GetToken().Col, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(node, operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(node, operator, left, right)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Upcast left to float
		lVal := float64(left.(*object.Integer).Value)
		rVal := right.(*object.Float).Value
		return evalCoercedFloatInfixExpression(node, operator, lVal, rVal)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Upcast right to float
		lVal := left.(*object.Float).Value
		rVal := float64(right.(*object.Integer).Value)
		return evalCoercedFloatInfixExpression(node, operator, lVal, rVal)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(node, operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(operator, left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(operator, left, right))
	default:
		return newError(node.GetToken().Line, node.GetToken().Col, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalObjectEquality(operator string, left, right object.Object) bool {
	if left.Type() != right.Type() {
		// Compare float and int values if possible
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
			return newError(node.GetToken().Line, node.GetToken().Col, "division by zero")
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
		return newError(node.GetToken().Line, node.GetToken().Col, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(node ast.Node, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value
	return evalCoercedFloatInfixExpression(node, operator, leftVal, rightVal)
}

func evalCoercedFloatInfixExpression(node ast.Node, operator string, leftVal, rightVal float64) object.Object {
	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0.0 {
			return newError(node.GetToken().Line, node.GetToken().Col, "division by zero")
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
		return newError(node.GetToken().Line, node.GetToken().Col, "unknown operator for float/numeric coercion")
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
		return newError(node.GetToken().Line, node.GetToken().Col, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return newError(node.GetToken().Line, node.GetToken().Col, "post-increment/decrement can only be applied to identifiers")
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newError(node.GetToken().Line, node.GetToken().Col, "identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		oldVal := v.Value
		var newVal int64
		if node.Operator == "++" {
			newVal = oldVal + 1
		} else {
			newVal = oldVal - 1
		}
		env.Update(ident.Value, &object.Integer{Value: newVal})
		return &object.Integer{Value: oldVal}

	case *object.Float:
		oldVal := v.Value
		var newVal float64
		if node.Operator == "++" {
			newVal = oldVal + 1.0
		} else {
			newVal = oldVal - 1.0
		}
		env.Update(ident.Value, &object.Float{Value: newVal})
		return &object.Float{Value: oldVal}

	default:
		return newError(node.GetToken().Line, node.GetToken().Col, "post-increment/decrement can only be applied to numeric types, got %s", val.Type())
	}
}

func evalAssignExpression(node *ast.AssignExpression, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	switch target := node.Left.(type) {
	case *ast.Identifier:
		ok := env.Update(target.Value, val)
		if !ok {
			return newError(node.GetToken().Line, node.GetToken().Col, "identifier not found: %s", target.Value)
		}
		return val

	case *ast.IndexExpression:
		leftObj := Eval(target.Left, env)
		if isError(leftObj) {
			return leftObj
		}
		idxObj := Eval(target.Index, env)
		if isError(idxObj) {
			return idxObj
		}

		arr, ok := leftObj.(*object.Array)
		if !ok {
			return newError(node.GetToken().Line, node.GetToken().Col, "index assignment only supported on arrays, got %s", leftObj.Type())
		}

		idx, ok := idxObj.(*object.Integer)
		if !ok {
			return newError(node.GetToken().Line, node.GetToken().Col, "array index must be integer, got %s", idxObj.Type())
		}

		if idx.Value < 0 || idx.Value >= int64(len(arr.Elements)) {
			return newError(node.GetToken().Line, node.GetToken().Col, "index out of bounds: %d", idx.Value)
		}

		arr.Elements[idx.Value] = val
		return val

	default:
		return newError(node.GetToken().Line, node.GetToken().Col, "cannot assign to expression %s", node.Left.String())
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
	// Create an enclosed scope for the for loop init and body variables
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
		extendedEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(args...)

	default:
		return newError(node.GetToken().Line, node.GetToken().Col, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(f *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(f.Env)

	// Register named function inside its own environment to support recursion!
	if f.Name != "" {
		env.Set(f.Name, f)
	}

	for paramIdx, param := range f.Parameters {
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

func evalIndexExpression(node ast.Node, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(node, left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(node, left, index)
	default:
		return newError(node.GetToken().Line, node.GetToken().Col, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(node ast.Node, array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newError(node.GetToken().Line, node.GetToken().Col, "index out of bounds: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(node ast.Node, str, index object.Object) object.Object {
	strVal := str.(*object.String).Value
	idx := index.(*object.Integer).Value
	max := int64(len(strVal) - 1)

	if idx < 0 || idx > max {
		return newError(node.GetToken().Line, node.GetToken().Col, "index out of bounds: %d", idx)
	}

	return &object.String{Value: string(strVal[idx])}
}

// Builtins map placeholder, will be populated in builtins.go
var builtins = map[string]*object.Builtin{}
