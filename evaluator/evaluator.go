package evaluator

import (
	"fmt"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/errors"
	"github.com/username/rustx/object"
	"github.com/username/rustx/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func init() {
	registerBuiltins()
}

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

	case *ast.ShortAssignStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		env.Set(n.Name.Value, val)

	case *ast.AssignmentStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}
		// Handle direct update in the environment chain
		if ident, ok := n.Name.(*ast.Identifier); ok {
			_, ok := env.Update(ident.Value, val)
			if !ok {
				return newErrorWithToken("identifier %s is not defined", ident.GetToken(), ident.Value)
			}
		} else if indexExp, ok := n.Name.(*ast.IndexExpression); ok {
			// Handle array index assignments, e.g. arr[0] = 5
			left := Eval(indexExp.Left, env)
			if isError(left) {
				return left
			}
			idx := Eval(indexExp.Index, env)
			if isError(idx) {
				return idx
			}
			arr, okArr := left.(*object.Array)
			idxObj, okIdx := idx.(*object.Integer)
			if !okArr {
				return newErrorWithToken("left-hand side of indexing assignment must be Array, got %s", indexExp.Left.GetToken(), left.Type())
			}
			if !okIdx {
				return newErrorWithToken("index must be Integer, got %s", indexExp.Index.GetToken(), idx.Type())
			}
			if idxObj.Value < 0 || idxObj.Value >= int64(len(arr.Elements)) {
				return newErrorWithToken("index out of bounds: %d", indexExp.Index.GetToken(), idxObj.Value)
			}
			arr.Elements[idxObj.Value] = val
		} else {
			return newErrorWithToken("invalid left-hand side of assignment", n.Name.GetToken())
		}

	case *ast.FunctionStatement:
		// Named function declaration
		fn := Eval(n.Value, env)
		if isError(fn) {
			return fn
		}
		env.Set(n.Name.Value, fn)

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

	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right, n.GetToken())

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
		return evalInfixExpression(n.Operator, left, right, n.GetToken())

	case *ast.PostfixExpression:
		// Restricted only to identifiers in current scope/environment
		ident, ok := n.Left.(*ast.Identifier)
		if !ok {
			return newErrorWithToken("postfix %s can only be applied to identifiers", n.GetToken(), n.Operator)
		}
		valObj, ok := env.Get(ident.Value)
		if !ok {
			return newErrorWithToken("identifier %s is not defined", ident.GetToken(), ident.Value)
		}
		switch valObj.Type() {
		case object.INTEGER_OBJ:
			curr := valObj.(*object.Integer).Value
			var next int64
			if n.Operator == "++" {
				next = curr + 1
			} else {
				next = curr - 1
			}
			env.Update(ident.Value, &object.Integer{Value: next})
			return &object.Integer{Value: curr} // post-increment returns the pre-incremented value
		case object.FLOAT_OBJ:
			curr := valObj.(*object.Float).Value
			var next float64
			if n.Operator == "++" {
				next = curr + 1.0
			} else {
				next = curr - 1.0
			}
			env.Update(ident.Value, &object.Float{Value: next})
			return &object.Float{Value: curr}
		default:
			return newErrorWithToken("postfix %s cannot be applied to type %s", n.GetToken(), n.Operator, valObj.Type())
		}

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
		return evalIndexExpression(left, index, n.GetToken())

	case *ast.IfExpression:
		return evalIfExpression(n, env)

	case *ast.WhileExpression:
		return evalWhileExpression(n, env)

	case *ast.ForExpression:
		return evalForExpression(n, env)

	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		return &object.Function{Parameters: params, Body: body, Env: env}

	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args, n.GetToken())
	}

	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL // prevent Go nil issues in empty files
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

func evalPrefixExpression(operator string, right object.Object, tok token.Token) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right, tok)
	default:
		return newErrorWithToken("unknown operator: %s%s", tok, operator, right.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object, tok token.Token) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newErrorWithToken("unknown operator: -%s", tok, right.Type())
	}
}

func evalInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right, tok)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right, tok)
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		// Coerce left to Float
		coercedLeft := &object.Float{Value: float64(left.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, coercedLeft, right, tok)
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		// Coerce right to Float
		coercedRight := &object.Float{Value: float64(right.(*object.Integer).Value)}
		return evalFloatInfixExpression(operator, left, coercedRight, tok)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right, tok)
	case operator == "==":
		return nativeBoolToBooleanObject(evalObjectEquality(left, right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!evalObjectEquality(left, right))
	default:
		return newErrorWithToken("type mismatch: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
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
			return newErrorWithToken("division by zero", tok)
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
		return newErrorWithToken("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalFloatInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
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
			return newErrorWithToken("division by zero", tok)
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
		return newErrorWithToken("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(operator string, left, right object.Object, tok token.Token) object.Object {
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
		return newErrorWithToken("unknown operator: %s %s %s", tok, left.Type(), operator, right.Type())
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

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newErrorWithToken("identifier not found: %s", node.GetToken(), node.Value)
}

func evalIndexExpression(left, index object.Object, tok token.Token) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayNavIndexExpression(left, index, tok)
	case left.Type() == object.STRING_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalStringNavIndexExpression(left, index, tok)
	default:
		return newErrorWithToken("index operator not supported: %s[%s]", tok, left.Type(), index.Type())
	}
}

func evalArrayNavIndexExpression(array, index object.Object, tok token.Token) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if idx < 0 || idx > max {
		return newErrorWithToken("index out of bounds: %d", tok, idx)
	}

	return arrayObject.Elements[idx]
}

func evalStringNavIndexExpression(str, index object.Object, tok token.Token) object.Object {
	stringObject := str.(*object.String)
	idx := index.(*object.Integer).Value
	max := int64(len(stringObject.Value) - 1)

	if idx < 0 || idx > max {
		return newErrorWithToken("index out of bounds: %d", tok, idx)
	}

	return &object.String{Value: string(stringObject.Value[idx])}
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
	var lastResult object.Object = NULL

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		lastResult = Eval(we.Body, env)
		if lastResult != nil {
			if lastResult.Type() == object.RETURN_VALUE_OBJ || lastResult.Type() == object.ERROR_OBJ {
				return lastResult
			}
		}
	}

	return lastResult
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	// Create an enclosed scope so that loop variables do not pollute outer scope
	forEnv := object.NewEnclosedEnvironment(env)

	if fe.Init != nil {
		initRes := Eval(fe.Init, forEnv)
		if isError(initRes) {
			return initRes
		}
	}

	var lastResult object.Object = NULL

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

		lastResult = Eval(fe.Body, forEnv)
		if lastResult != nil {
			if lastResult.Type() == object.RETURN_VALUE_OBJ || lastResult.Type() == object.ERROR_OBJ {
				return lastResult
			}
		}

		if fe.Post != nil {
			postRes := Eval(fe.Post, forEnv)
			if isError(postRes) {
				return postRes
			}
		}
	}

	return lastResult
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

func applyFunction(fn object.Object, args []object.Object, tok token.Token) object.Object {
	switch f := fn.(type) {
	case *object.Function:
		extendEnv := extendFunctionEnv(f, args)
		// Register name inside the local environment to support recursive calls/self-reference
		// Fl.Name is stored in the function litteral but we can retrieve it or use closure.
		// Since we want named recursion, if we have standard named functions we've registered them in outer scope.
		evaluated := Eval(f.Body, extendEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return f.Fn(args...)
	default:
		return newErrorWithToken("not a function: %s", tok, fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
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
		// String and Numbers are truthy unless they have explicit empty/zero values (though traditionally truthy)
		// Let's keep a standard behavior: non-null, non-false is true.
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	// Must use a format that conforms to Go static analysis.
	// Always use dynamic formatting internally or safe errors
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func newErrorWithToken(format string, tok token.Token, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	err := &errors.RuntimeError{
		Message: msg,
		Line:    tok.Line,
		Column:  tok.Column,
	}
	return &object.Error{Message: err.Error()}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
