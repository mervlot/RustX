package evaluator

import (
	"fmt"
	"math"
	"rustx/ast"
	"rustx/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func newPosError(node ast.Node, format string, a ...interface{}) *object.Error {
	tok := node.GetToken()
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

// Eval evaluates an AST node in the given environment.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n.Statements, env)

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
		var val object.Object = NULL
		if n.Value != nil {
			val = Eval(n.Value, env)
			if isError(val) {
				return val
			}
		}
		env.Set(n.Name.Value, val)
		return NULL

	case *ast.AssignStatement:
		return evalAssignStatement(n, env)

	case *ast.WhileStatement:
		return evalWhileStatement(n, env)

	case *ast.ForStatement:
		return evalForStatement(n, env)

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
		return evalIndexExpression(left, index, n)

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right, n)

	case *ast.PostfixExpression:
		return evalPostfixExpression(n, env)

	case *ast.InfixExpression:
		// Logical && and || support short-circuit evaluation
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
		return evalInfixExpression(n.Operator, left, right, n)

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.Identifier:
		return evalIdentifier(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fn := &object.Function{Name: n.Name, Parameters: params, Body: body, Env: env}
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
		return applyFunction(function, args, n)
	}

	return NULL
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range stmts {
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

func evalAssignStatement(node *ast.AssignStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	if node.IsDefine {
		// x := value
		if node.Name == nil {
			return newPosError(node, "Invalid left hand side in definition")
		}
		env.Set(node.Name.Value, val)
		return NULL
	}

	// Regular assignment: left can be identifier or index expression
	if node.Name != nil {
		updated := env.Update(node.Name.Value, val)
		if !updated {
			// If not defined, set in current environment
			env.Set(node.Name.Value, val)
		}
		return NULL
	}

	// Index assignment: arr[idx] = val
	if idxExpr, ok := node.Left.(*ast.IndexExpression); ok {
		left := Eval(idxExpr.Left, env)
		if isError(left) {
			return left
		}
		arr, ok := left.(*object.Array)
		if !ok {
			return newPosError(idxExpr, "Index assignment only supported for Array, got %s", left.Type())
		}
		idxObj := Eval(idxExpr.Index, env)
		if isError(idxObj) {
			return idxObj
		}
		intIdx, ok := idxObj.(*object.Integer)
		if !ok {
			return newPosError(idxExpr, "Array index must be an integer, got %s", idxObj.Type())
		}
		idx := intIdx.Value
		if idx < 0 || idx >= int64(len(arr.Elements)) {
			return newPosError(idxExpr, "Array index out of bounds: index %d with length %d", idx, len(arr.Elements))
		}
		arr.Elements[idx] = val
		return NULL
	}

	return newPosError(node, "Invalid assignment target")
}

func evalWhileStatement(ws *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(ws.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		// Enclosed environment for while loop body
		loopEnv := object.NewEnclosedEnvironment(env)
		result = evalBlockStatement(ws.Body, loopEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(fs *ast.ForStatement, env *object.Environment) object.Object {
	// For loop creates its own scope for init variables
	forEnv := object.NewEnclosedEnvironment(env)

	if fs.Init != nil {
		initRes := Eval(fs.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var result object.Object = NULL

	for {
		if fs.Condition != nil {
			condRes := Eval(fs.Condition, forEnv)
			if isError(condRes) {
				return condRes
			}
			if !isTruthy(condRes) {
				break
			}
		}

		// Body executes within forEnv
		result = evalBlockStatement(fs.Body, forEnv)
		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}

		if fs.Post != nil {
			postRes := Eval(fs.Post, forEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return result
}

func evalPostfixExpression(pe *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := pe.Left.(*ast.Identifier)
	if !ok {
		return newPosError(pe, "operator %s only supported on variables", pe.Operator)
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		return newPosError(pe, "identifier not found: %s", ident.Value)
	}

	switch v := val.(type) {
	case *object.Integer:
		oldVal := v.Value
		var newVal int64
		if pe.Operator == "++" {
			newVal = oldVal + 1
		} else if pe.Operator == "--" {
			newVal = oldVal - 1
		} else {
			return newPosError(pe, "unknown postfix operator: %s", pe.Operator)
		}
		env.Update(ident.Value, &object.Integer{Value: newVal})
		return &object.Integer{Value: oldVal}
	case *object.Float:
		oldVal := v.Value
		var newVal float64
		if pe.Operator == "++" {
			newVal = oldVal + 1.0
		} else if pe.Operator == "--" {
			newVal = oldVal - 1.0
		} else {
			return newPosError(pe, "unknown postfix operator: %s", pe.Operator)
		}
		env.Update(ident.Value, &object.Float{Value: newVal})
		return &object.Float{Value: oldVal}
	default:
		return newPosError(pe, "operator %s not supported for type %s", pe.Operator, val.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object, node ast.Node) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, node)
	default:
		return newPosError(node, "unknown operator: %s%s", operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, node ast.Node) object.Object {
	switch v := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -v.Value}
	case *object.Float:
		return &object.Float{Value: -v.Value}
	default:
		return newPosError(node, "unknown operator: -%s", right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right, node)
	}
	if (left.Type() == object.INTEGER_OBJ || left.Type() == object.FLOAT_OBJ) &&
		(right.Type() == object.INTEGER_OBJ || right.Type() == object.FLOAT_OBJ) {
		return evalFloatInfixExpression(operator, left, right, node)
	}
	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left, right, node)
	}
	if operator == "==" {
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	}
	if operator == "!=" {
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	}
	return newPosError(node, "type mismatch: %s %s %s", left.Type(), operator, right.Type())
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
	case *object.Null:
		return true
	default:
		return left == right
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
			return newPosError(node, "division by zero")
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return newPosError(node, "modulo by zero")
		}
		return &object.Integer{Value: leftVal % rightVal}
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
		return newPosError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	var leftVal, rightVal float64
	if l, ok := left.(*object.Integer); ok {
		leftVal = float64(l.Value)
	} else {
		leftVal = left.(*object.Float).Value
	}

	if r, ok := right.(*object.Integer); ok {
		rightVal = float64(r.Value)
	} else {
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
		if rightVal == 0 {
			return newPosError(node, "division by zero")
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
		return nativeBoolToBooleanObject(math.Abs(leftVal-rightVal) < 1e-9)
	case "!=":
		return nativeBoolToBooleanObject(math.Abs(leftVal-rightVal) >= 1e-9)
	default:
		return newPosError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
		return newPosError(node, "unknown operator: %s %s %s", left.Type(), operator, right.Type())
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

func evalIndexExpression(left, index object.Object, node ast.Node) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index, node)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringIndexExpression(left, index, node)
	default:
		return newPosError(node, "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object, node ast.Node) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newPosError(node, "Array index out of bounds: index %d with length %d", idx, len(arrayObject.Elements))
	}

	return arrayObject.Elements[idx]
}

func evalStringIndexExpression(str, index object.Object, node ast.Node) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return newPosError(node, "String index out of bounds: index %d with length %d", idx, len(strObject.Value))
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newPosError(node, "identifier not found: %s", node.Value)
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
	switch f := fn.(type) {
	case *object.Function:
		if len(args) != len(f.Parameters) {
			return newPosError(node, "wrong number of arguments: expected %d, got %d", len(f.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(f, args)
		evaluated := Eval(f.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(args...)

	default:
		return newPosError(node, "not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
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
		if intObj, ok := obj.(*object.Integer); ok {
			return intObj.Value != 0
		}
		if floatObj, ok := obj.(*object.Float); ok {
			return floatObj.Value != 0.0
		}
		if strObj, ok := obj.(*object.String); ok {
			return len(strObj.Value) > 0
		}
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
