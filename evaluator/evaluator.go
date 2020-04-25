package evaluator

import (
	"fmt"
	"hummus-lang/ast"
	"hummus-lang/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

var predefs = map[string]*object.Predef{
	"print": &object.Predef{
		Function: func(args ...object.Object) object.Object {
			fmt.Print(args[0].Printable())
			return Null
		},
	},
	"printLine": &object.Predef{
		Function: func(args ...object.Object) object.Object {
			fmt.Println(args[0].Printable())
			return Null
		},
	},
	"len": &object.Predef{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("len: expected exactly 1 argument. given %d", len(args))
			}
			switch x := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(x.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(x.Elements))}
			default:
				return newError("len: can only take length of strings and arrays")
			}
		},
	},
	"head": &object.Predef{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("head: expected exactly 1 argument. given %d", len(args))
			}
			switch x := args[0].(type) {
			case *object.Array:
				if len(x.Elements) == 0 {
					return newError("head: can not take head of empty array")
				} else {
					return x.Elements[0]
				}
			case *object.String:
				if len(x.Value) == 0 {
					return newError("head: can not take head of empty string")
				} else {
					return &object.String{Value: fmt.Sprintf("%c", x.Value[0])}
				}
			default:
				return newError("head: can not take head of `%s`", args[0].Type())
			}
		},
	},
	"tail": &object.Predef{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("head: expected exactly 1 argument. given %d", len(args))
			}
			switch x := args[0].(type) {
			case *object.Array:
				if len(x.Elements) == 0 {
					return newError("tail: can not take tail of empty array")
				} else {
					tailElements := x.Elements[1:]
					return &object.Array{Elements: tailElements}
				}
			case *object.String:
				if len(x.Value) == 0 {
					return newError("tail: can not take tail of empty string")
				} else {
					return &object.String{Value: x.Value[1:]}
				}
			default:
				return newError("head: can not take head of `%s`", args[0].Type())
			}
		},
	},
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ErrorObj
	} else {
		return false
	}
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

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
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueObj || rt == object.ErrorObj {
				return result
			}
		}

	}

	return result
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return newError("unknown operator: unary - not defined for `%s`", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
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
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: binary %s not defined for `%s` + `%s`", operator, left.Type(), right.Type())
	}
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftString := left.(*object.String).Value
	rightString := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftString + rightString}
	default:
		return newError("unknown operator: binary %s not defined for `%s` and `%s`", operator, left.Type(), right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s not defined for type `%s`", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.StringObj && right.Type() == object.StringObj:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: can not %s `%s` and `%s`", operator, left.Type(), right.Type())
	default:
		return newError("unknown operator: binary %s not defined for `%s` and `%s`", operator, left.Type(), right.Type())
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
		return Null
	}
}

func evalIndexExpression(ie *ast.IndexExpression, env *object.Environment) object.Object {
	leftEval := Eval(ie.Left, env)
	if isError(leftEval) {
		return leftEval
	}

	rightEval := Eval(ie.Right, env)
	if isError(rightEval) {
		return rightEval
	}

	if leftEval.Type() == object.ArrayObj && rightEval.Type() == object.IntegerObj {
		array := leftEval.(*object.Array).Elements
		idx := int(rightEval.(*object.Integer).Value)

		if idx < 0 || idx > len(array)-1 {
			//throw new ArrayIndexOutOfBoundsException()
			return newError("index expression: index out of array bounds")
		}

		return array[idx]
	} else if leftEval.Type() == object.StringObj && rightEval.Type() == object.IntegerObj {
		str := leftEval.(*object.String).Value
		idx := int(rightEval.(*object.Integer).Value)

		if idx < 0 || idx > len(str)-1 {
			return newError("index expression: index out of string bounds")
		}

		return &object.String{Value: fmt.Sprintf("%c", str[idx])}
	}

	return newError("index expression: can not take index of type `%s` with `%s`", leftEval.Type(), rightEval.Type())
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}

	val, ok = predefs[node.Value]
	if ok {
		return val
	}

	return newError("unknown reference on line %d: %s", node.Token.Line, node.Value)

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
	switch fn.(type) {
	case *object.Function:
		actual := fn.(*object.Function)
		if len(actual.Parameters) != len(args) {
			return newError("incorrect number of arguments: need %d, got %d", len(actual.Parameters), len(args))
		}
		extendedEnv := extendFunctionEnv(actual, args)
		evaluated := Eval(actual.Body, extendedEnv)
		return unwrapReturnValue(evaluated)

	case *object.Predef:
		actual := fn.(*object.Predef)
		return actual.Function(args...)

	default:
		return newError("applyFunction: unknown function; got %s", fn.Type())
	}

}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	} else {
		return obj
	}
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}

		return &object.Array{Elements: elements}

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	default:
		fmt.Printf("Eval: unknown expression type encountered. Expression type: %T\n", node)
	}

	return nil
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	} else {
		return False
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}
