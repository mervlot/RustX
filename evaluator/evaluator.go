package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type Evaluator struct {
	filename string
}

func New(filename string) *Evaluator {
	return &Evaluator{filename: filename}
}

func (ev *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {

	// Statements
	case *ast.Program:
		return ev.evalProgram(n, env)

	case *ast.BlockStatement:
		return ev.evalBlockStatement(n, env)

	case *ast.ExpressionStatement:
		return ev.Eval(n.Expression, env)

	case *ast.ReturnStatement:
		val := ev.Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.VarStatement:
		val := ev.Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)
		return val

	case *ast.AssignStatement:
		// Standard assignment
		val := ev.Eval(n.Value, env)
		if isError(val) {
			return val
		}

		switch target := n.Name.(type) {
		case *ast.Identifier:
			ok := env.Update(target.Value, val)
			if !ok {
				// Fallback to definition if not yet declared (though usually error)
				env.Set(target.Value, val)
			}
			return val

		case *ast.IndexExpression:
			// Left hand side is array
			arrObj := ev.Eval(target.Left, env)
			if isError(arrObj) {
				return arrObj
			}

			idxObj := ev.Eval(target.Index, env)
			if isError(idxObj) {
				return idxObj
			}

			array, ok1 := arrObj.(*object.Array)
			index, ok2 := idxObj.(*object.Integer)

			if !ok1 {
				return ev.newError(target.Left.GetToken(), "LHS of assignment must be an ARRAY for index assignment, got %s", arrObj.Type())
			}
			if !ok2 {
				return ev.newError(target.Index.GetToken(), "Index must be an INTEGER, got %s", idxObj.Type())
			}

			idx := index.Value
			if idx < 0 || idx >= int64(len(array.Elements)) {
				return ev.newError(target.Index.GetToken(), "Index out of bounds: %d", idx)
			}

			array.Elements[idx] = val
			return val

		default:
			return ev.newError(n.GetToken(), "Invalid assignment target")
		}

	// Expressions
	case *ast.Identifier:
		return ev.evalIdentifier(n, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(n.Value)

	case *ast.StringLiteral:
		return &object.String{Value: n.Value}

	case *ast.ArrayLiteral:
		elements := ev.evalExpressions(n.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		left := ev.Eval(n.Left, env)
		if isError(left) {
			return left
		}
		index := ev.Eval(n.Index, env)
		if isError(index) {
			return index
		}
		return ev.evalIndexExpression(n, left, index)

	case *ast.PrefixExpression:
		right := ev.Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return ev.evalPrefixExpression(n, n.Operator, right)

	case *ast.InfixExpression:
		left := ev.Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := ev.Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return ev.evalInfixExpression(n, n.Operator, left, right)

	case *ast.PostfixExpression:
		return ev.evalPostfixExpression(n, env)

	case *ast.IfExpression:
		return ev.evalIfExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		fn := &object.Function{Parameters: params, Env: env, Body: body}
		if n.Name != nil {
			env.Set(n.Name.Value, fn)
		}
		return fn

	case *ast.CallExpression:
		function := ev.Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := ev.evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return ev.applyFunction(n, function, args)

	case *ast.WhileExpression:
		return ev.evalWhileExpression(n, env)

	case *ast.ForExpression:
		return ev.evalForExpression(n, env)
	}

	return NULL
}

func (ev *Evaluator) evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL // Safe default result

	for _, statement := range program.Statements {
		result = ev.Eval(statement, env)

		switch r := result.(type) {
		case *object.ReturnValue:
			return r.Value
		case *object.Error:
			return r
		}
	}

	return result
}

func (ev *Evaluator) evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	// Block statements introduce local scope
	blockEnv := object.NewEnclosedEnvironment(env)

	for _, statement := range block.Statements {
		result = ev.Eval(statement, blockEnv)

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

func (ev *Evaluator) evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return ev.newError(node.GetToken(), "Identifier not found: %s", node.Value)
}

func (ev *Evaluator) evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := ev.Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func (ev *Evaluator) evalIndexExpression(node *ast.IndexExpression, left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return ev.evalArrayIndexExpression(node, left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return ev.evalStringIndexExpression(node, left, index)
	default:
		return ev.newError(node.GetToken(), "index operator not supported: %s[%s]", left.Type(), index.Type())
	}
}

func (ev *Evaluator) evalArrayIndexExpression(node *ast.IndexExpression, array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return ev.newError(node.Index.GetToken(), "index out of range: %d", idx)
	}

	return arrayObject.Elements[idx]
}

func (ev *Evaluator) evalStringIndexExpression(node *ast.IndexExpression, str, index object.Object) object.Object {
	strObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(strObject.Value) - 1)

	if idx < 0 || idx > max {
		return ev.newError(node.Index.GetToken(), "index out of range: %d", idx)
	}

	return &object.String{Value: string(strObject.Value[idx])}
}

func (ev *Evaluator) evalPrefixExpression(node *ast.PrefixExpression, operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return ev.evalBangOperatorExpression(right)
	case "-":
		return ev.evalMinusPrefixOperatorExpression(node, right)
	default:
		return ev.newError(node.GetToken(), "unknown operator: %s%s", operator, right.Type())
	}
}

func (ev *Evaluator) evalBangOperatorExpression(right object.Object) object.Object {
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

func (ev *Evaluator) evalMinusPrefixOperatorExpression(node *ast.PrefixExpression, right object.Object) object.Object {
	switch right.Type() {
	case object.INTEGER_OBJ:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	case object.FLOAT_OBJ:
		value := right.(*object.Float).Value
		return &object.Float{Value: -value}
	default:
		return ev.newError(node.GetToken(), "unknown operator: -%s", right.Type())
	}
}

func (ev *Evaluator) evalInfixExpression(node *ast.InfixExpression, operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return ev.evalIntegerInfixExpression(node, operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return ev.evalFloatInfixExpression(node, operator, left, right)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Cast Int to Float
		rCast := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return ev.evalFloatInfixExpression(node, operator, left, rCast)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Cast Int to Float
		lCast := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return ev.evalFloatInfixExpression(node, operator, lCast, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return ev.evalStringInfixExpression(node, operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return ev.evalBooleanInfixExpression(node, operator, left, right)
	default:
		if operator == "==" {
			return nativeBoolToBooleanObject(evalObjectEquality(left, right))
		}
		if operator == "!=" {
			return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
		}
		return ev.newError(node.GetToken(), "type mismatch: %s %s %s", left.Type(), operator, right.Type())
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

func (ev *Evaluator) evalIntegerInfixExpression(node *ast.InfixExpression, operator string, left, right object.Object) object.Object {
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
			return ev.newError(node.Right.GetToken(), "division by zero")
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
		return ev.newError(node.GetToken(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (ev *Evaluator) evalFloatInfixExpression(node *ast.InfixExpression, operator string, left, right object.Object) object.Object {
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
			return ev.newError(node.Right.GetToken(), "division by zero")
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
		return ev.newError(node.GetToken(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (ev *Evaluator) evalStringInfixExpression(node *ast.InfixExpression, operator string, left, right object.Object) object.Object {
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
		return ev.newError(node.GetToken(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (ev *Evaluator) evalBooleanInfixExpression(node *ast.InfixExpression, operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case "||":
		return nativeBoolToBooleanObject(leftVal || rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return ev.newError(node.GetToken(), "unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (ev *Evaluator) evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		return ev.newError(node.Left.GetToken(), "post increment/decrement operator must be applied to an identifier")
	}

	val, found := env.Get(ident.Value)
	if !found {
		return ev.newError(ident.GetToken(), "Identifier not found: %s", ident.Value)
	}

	switch node.Operator {
	case "++":
		switch v := val.(type) {
		case *object.Integer:
			newVal := &object.Integer{Value: v.Value + 1}
			env.Update(ident.Value, newVal)
			return newVal
		case *object.Float:
			newVal := &object.Float{Value: v.Value + 1.0}
			env.Update(ident.Value, newVal)
			return newVal
		default:
			return ev.newError(node.GetToken(), "post increment can only be applied to INTEGER or FLOAT, got %s", val.Type())
		}
	case "--":
		switch v := val.(type) {
		case *object.Integer:
			newVal := &object.Integer{Value: v.Value - 1}
			env.Update(ident.Value, newVal)
			return newVal
		case *object.Float:
			newVal := &object.Float{Value: v.Value - 1.0}
			env.Update(ident.Value, newVal)
			return newVal
		default:
			return ev.newError(node.GetToken(), "post decrement can only be applied to INTEGER or FLOAT, got %s", val.Type())
		}
	default:
		return ev.newError(node.GetToken(), "unknown postfix operator: %s", node.Operator)
	}
}

func (ev *Evaluator) evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := ev.Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return ev.Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return ev.Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func (ev *Evaluator) evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	var lastResult object.Object = NULL

	for {
		condition := ev.Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		lastResult = ev.Eval(we.Body, env)
		if lastResult != nil {
			rt := lastResult.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return lastResult
			}
		}
	}

	return lastResult
}

func (ev *Evaluator) evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	// For-expression executes in its own scope enclosing outer scope
	forEnv := object.NewEnclosedEnvironment(env)

	// Execute Initializer if present
	if fe.Init != nil {
		initRes := ev.Eval(fe.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var lastResult object.Object = NULL

	for {
		// Evaluate condition if present
		if fe.Condition != nil {
			condition := ev.Eval(fe.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		// Execute Body
		lastResult = ev.Eval(fe.Body, forEnv)
		if lastResult != nil {
			rt := lastResult.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return lastResult
			}
		}

		// Execute Post Statement if present
		if fe.Post != nil {
			postRes := ev.Eval(fe.Post, forEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return lastResult
}

func (ev *Evaluator) applyFunction(node *ast.CallExpression, fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		extendEnv := ev.extendFunctionEnv(f, args)
		evaluated := ev.Eval(f.Body, extendEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		return f.Fn(args...)

	default:
		return ev.newError(node.GetToken(), "not a function: %s", fn.Type())
	}
}

func (ev *Evaluator) extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
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
		return true
	}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func (ev *Evaluator) newError(tok token.Token, format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
		File:    ev.filename,
		Line:    tok.Line,
		Column:  tok.Column,
	}
}
