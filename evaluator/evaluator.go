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
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.VarStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		if val == nil {
			val = NULL
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.WhileStatement:
		return evalWhileStatement(node, env)

	case *ast.ForStatement:
		return evalForStatement(node, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(node.Value)

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

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right, node)

	case *ast.PostfixExpression:
		return evalPostfixExpression(node, env)

	case *ast.InfixExpression:
		if node.Operator == ":=" {
			ident, ok := node.Left.(*ast.Identifier)
			if !ok {
				tok := node.GetToken()
				return &object.Error{Message: "left side of := must be an identifier", Line: tok.Line, Column: tok.Column}
			}
			val := Eval(node.Right, env)
			if isError(val) {
				return val
			}
			env.Set(ident.Value, val)
			return val
		}

		if node.Operator == "=" {
			val := Eval(node.Right, env)
			if isError(val) {
				return val
			}

			if ident, ok := node.Left.(*ast.Identifier); ok {
				if _, updated := env.Update(ident.Value, val); !updated {
					tok := node.GetToken()
					return &object.Error{
						Message: fmt.Sprintf("identifier not found for assignment: %s", ident.Value),
						Line:    tok.Line,
						Column:  tok.Column,
					}
				}
				return val
			}

			if indexExp, ok := node.Left.(*ast.IndexExpression); ok {
				left := Eval(indexExp.Left, env)
				if isError(left) {
					return left
				}
				idx := Eval(indexExp.Index, env)
				if isError(idx) {
					return idx
				}

				arr, ok := left.(*object.Array)
				if !ok {
					tok := indexExp.GetToken()
					return &object.Error{Message: fmt.Sprintf("index assignment not supported on %s", left.Type()), Line: tok.Line, Column: tok.Column}
				}

				iObj, ok := idx.(*object.Integer)
				if !ok {
					tok := indexExp.GetToken()
					return &object.Error{Message: fmt.Sprintf("array index must be INTEGER, got %s", idx.Type()), Line: tok.Line, Column: tok.Column}
				}

				i := iObj.Value
				if i < 0 || i >= int64(len(arr.Elements)) {
					tok := indexExp.GetToken()
					return &object.Error{Message: fmt.Sprintf("index out of bounds: %d", i), Line: tok.Line, Column: tok.Column}
				}

				arr.Elements[i] = val
				return val
			}

			tok := node.GetToken()
			return &object.Error{Message: "invalid left-hand side in assignment", Line: tok.Line, Column: tok.Column}
		}

		// Short-circuiting logical operations
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
		return evalInfixExpression(node.Operator, left, right, node)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		fn := &object.Function{
			Name:       node.Name,
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
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
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	// Automatic entry point call if zero-arg main() function is defined
	if mainObj, ok := env.Get("main"); ok {
		if fn, isFn := mainObj.(*object.Function); isFn && len(fn.Parameters) == 0 {
			mainRes := applyFunction(program, fn, []object.Object{})
			if isError(mainRes) {
				return mainRes
			}
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	// Block statements get an enclosed environment for block scoping
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

func evalWhileStatement(node *ast.WhileStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for {
		condition := Eval(node.Condition, env)
		if isError(condition) {
			return condition
		}

		if !isTruthy(condition) {
			break
		}

		result = evalBlockStatement(node.Body, env)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalForStatement(node *ast.ForStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	forEnv := object.NewEnclosedEnvironment(env)

	if node.Init != nil {
		initRes := Eval(node.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	for {
		if node.Condition != nil {
			condition := Eval(node.Condition, forEnv)
			if isError(condition) {
				return condition
			}
			if !isTruthy(condition) {
				break
			}
		}

		result = evalBlockStatement(node.Body, forEnv)
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}

		if node.Post != nil {
			postRes := Eval(node.Post, forEnv)
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

func evalPrefixExpression(operator string, right object.Object, node ast.Node) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, node)
	default:
		tok := node.GetToken()
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s%s", operator, right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
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
	tok := node.GetToken()
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	case *object.Float:
		return &object.Float{Value: -right.Value}
	default:
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: -%s", right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
	}
}

func evalPostfixExpression(node *ast.PostfixExpression, env *object.Environment) object.Object {
	ident, ok := node.Left.(*ast.Identifier)
	if !ok {
		tok := node.GetToken()
		return &object.Error{Message: "postfix operator requires an identifier", Line: tok.Line, Column: tok.Column}
	}

	val, ok := env.Get(ident.Value)
	if !ok {
		tok := node.GetToken()
		return &object.Error{Message: fmt.Sprintf("identifier not found: %s", ident.Value), Line: tok.Line, Column: tok.Column}
	}

	switch v := val.(type) {
	case *object.Integer:
		var newVal int64
		if node.Operator == "++" {
			newVal = v.Value + 1
		} else {
			newVal = v.Value - 1
		}
		newObj := &object.Integer{Value: newVal}
		env.Update(ident.Value, newObj)
		return v
	case *object.Float:
		var newVal float64
		if node.Operator == "++" {
			newVal = v.Value + 1.0
		} else {
			newVal = v.Value - 1.0
		}
		newObj := &object.Float{Value: newVal}
		env.Update(ident.Value, newObj)
		return v
	default:
		tok := node.GetToken()
		return &object.Error{Message: fmt.Sprintf("postfix operator %s not supported for %s", node.Operator, v.Type()), Line: tok.Line, Column: tok.Column}
	}
}

func evalInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	tok := node.GetToken()

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, node)
	case (left.Type() == object.FLOAT_OBJ || left.Type() == object.INTEGER_OBJ) &&
		(right.Type() == object.FLOAT_OBJ || right.Type() == object.INTEGER_OBJ):
		return evalFloatInfixExpression(operator, left, right, node)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, node)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	case left.Type() != right.Type():
		return &object.Error{
			Message: fmt.Sprintf("type mismatch: %s %s %s", left.Type(), operator, right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
	default:
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	tok := node.GetToken()

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &object.Error{Message: "division by zero", Line: tok.Line, Column: tok.Column}
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		if rightVal == 0 {
			return &object.Error{Message: "modulo by zero", Line: tok.Line, Column: tok.Column}
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
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	var leftVal, rightVal float64

	if l, ok := left.(*object.Float); ok {
		leftVal = l.Value
	} else {
		leftVal = float64(left.(*object.Integer).Value)
	}

	if r, ok := right.(*object.Float); ok {
		rightVal = r.Value
	} else {
		rightVal = float64(right.(*object.Integer).Value)
	}

	tok := node.GetToken()

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return &object.Error{Message: "division by zero", Line: tok.Line, Column: tok.Column}
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
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, node ast.Node) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	tok := node.GetToken()

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return &object.Error{
			Message: fmt.Sprintf("unknown operator: %s %s %s", left.Type(), operator, right.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
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
	case *object.Null:
		return true
	default:
		return left == right
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	tok := node.GetToken()
	return &object.Error{
		Message: fmt.Sprintf("identifier not found: %s", node.Value),
		Line:    tok.Line,
		Column:  tok.Column,
	}
}

func evalIndexExpression(left, index object.Object, node ast.Node) object.Object {
	tok := node.GetToken()
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arrayObject := left.(*object.Array)
		idx := index.(*object.Integer).Value
		max := int64(len(arrayObject.Elements) - 1)
		if idx < 0 || idx > max {
			return &object.Error{Message: fmt.Sprintf("index out of bounds: %d", idx), Line: tok.Line, Column: tok.Column}
		}
		return arrayObject.Elements[idx]
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		strObject := left.(*object.String)
		idx := index.(*object.Integer).Value
		max := int64(len(strObject.Value) - 1)
		if idx < 0 || idx > max {
			return &object.Error{Message: fmt.Sprintf("index out of bounds: %d", idx), Line: tok.Line, Column: tok.Column}
		}
		return &object.String{Value: string(strObject.Value[idx])}
	default:
		return &object.Error{
			Message: fmt.Sprintf("index operator not supported: %s[%s]", left.Type(), index.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
	}
}

func applyFunction(node ast.Node, fn object.Object, args []object.Object) object.Object {
	tok := node.GetToken()
	switch function := fn.(type) {
	case *object.Function:
		if len(args) != len(function.Parameters) {
			return &object.Error{
				Message: fmt.Sprintf("wrong argument count: expected %d, got %d", len(function.Parameters), len(args)),
				Line:    tok.Line,
				Column:  tok.Column,
			}
		}
		extendedEnv := extendFunctionEnv(function, args)
		evaluated := Eval(function.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Builtin:
		res := function.Fn(args...)
		if err, ok := res.(*object.Error); ok {
			err.Line = tok.Line
			err.Column = tok.Column
		}
		return res

	default:
		return &object.Error{
			Message: fmt.Sprintf("not a function: %s", fn.Type()),
			Line:    tok.Line,
			Column:  tok.Column,
		}
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
		if integer, ok := obj.(*object.Integer); ok {
			return integer.Value != 0
		}
		if float, ok := obj.(*object.Float); ok {
			return float.Value != 0.0
		}
		if str, ok := obj.(*object.String); ok {
			return len(str.Value) > 0
		}
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
