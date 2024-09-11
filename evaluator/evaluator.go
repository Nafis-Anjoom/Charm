package evaluator

import (
	"charm/ast"
	"charm/object"
    "fmt"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalProgram(node, env)
    case *ast.ExpressionStatement:
        return Eval(node.Expression, env)
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
    case *ast.BooleanLiteral:
        return nativeBooltoBoolObject(node.Value)
    case *ast.PrefixExpression:
        return evalPrefixExpression(node, env)
    case *ast.InfixExpression:
        return evalInfixExpression(node, env)
    case *ast.IfExpression:
        return evalIfExpression(node, env)
    case *ast.BlockStatement:
        return evalStatements(node.Statements, env)
    case *ast.ReturnStatement:
        returnValue := Eval(node.ReturnValue, env)
        if isError(returnValue) {
            return returnValue
        }
        return &object.ReturnValue{Value: returnValue}
    case *ast.LetStatement:
        value := Eval(node.Value, env)
        if isError(value) {
            return value
        }
        env.Set(node.Identifier.Value, value)
        return value
    case *ast.Identifier:
        val, found := env.Get(node.Value)
        if !found {
            return newError("identifier not found: %s", node.Value)
        }
        return val
    case *ast.FunctionLiteral:
        return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
    }
    return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
    var result object.Object

    for _, stmt := range program.Statements {
        result = Eval(stmt, env)
        switch result := result.(type) {
        case *object.ReturnValue:
            return result.Value
        case *object.Error:
            return result
        }
    }
    return result
}


func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
    var result object.Object

    for _, stmt := range stmts {
        result = Eval(stmt, env)
        if result != nil {
            if  result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
                return result
            }
        }
    }

    return result
}

// DNOTE: return NULL when encountering unknown prefix operator 
func evalPrefixExpression(exp *ast.PrefixExpression, env *object.Environment) object.Object {
    right := Eval(exp.Right, env)
    if isError(right) {
        return right
    }

    switch exp.Operator {
    case "!":
        return evalBangPrefixExpression(right)
    case "-":
        return evalMinusPrefixExpression(right)
    default:
        return newError("unknown operator: %s%s",exp.Operator, right.Type())
    }
}

// DNOTE: Integers are truthy. Nulls are falsy
func evalBangPrefixExpression(right object.Object) object.Object {
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

// DNOTE: if right expression not integers, then return NULL
func evalMinusPrefixExpression(right object.Object) object.Object {
    if right.Type() != object.INTEGER_OBJ {
        return newError("unknown operator: -%s", right.Type())
    }
    value := right.(*object.Integer).Value
    return &object.Integer{Value: -value}
}

func evalInfixExpression(exp *ast.InfixExpression, env *object.Environment) object.Object {
    right := Eval(exp.Right, env)
    if isError(right) {
        return right
    }
    left := Eval(exp.Left, env)
    if isError(left) {
        return left
    }

    switch {
    case right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ:
        rightValue := right.(*object.Integer).Value
        leftValue := left.(*object.Integer).Value
        return evalIntegerInfixExpression(leftValue, exp.Operator, rightValue)
    case exp.Operator == "==":
        return nativeBooltoBoolObject(left == right)
    case exp.Operator == "!=":
        return nativeBooltoBoolObject(left != right)
    case left.Type() != right.Type():
        return newError("type mismatch: %s %s %s", left.Type(), exp.Operator, right.Type())
    default:
        return newError("unknown operator: %s %s %s", left.Type(), exp.Operator, right.Type())
    }
}

func evalIntegerInfixExpression(left int64, operator string, right int64) object.Object {
    switch operator {
    case "+":
        return &object.Integer{Value: left + right}
    case "-":
        return &object.Integer{Value: left - right}
    case "*":
        return &object.Integer{Value: left * right}
    case "/":
        return &object.Integer{Value: left / right}
    case "<":
        return nativeBooltoBoolObject(left < right)
    case ">":
        return nativeBooltoBoolObject(left > right)
    case "==":
        return nativeBooltoBoolObject(left == right)
    case "!=":
        return nativeBooltoBoolObject(left != right)
    default:
        return NULL
    }
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(node.Condition, env)
    if isError(condition) {
        return condition
    }

    if isTruthy(condition) {
        return Eval(node.Consequence, env)
    } else if node.Alternative != nil {
        return Eval(node.Alternative, env)
    } else {
        return NULL
    }
}

func nativeBooltoBoolObject(boolean bool) *object.Boolean {
    if boolean {
        return TRUE
    }
    return FALSE
}

func isTruthy(obj object.Object) bool {
    switch obj := obj.(type) {
    case *object.Boolean:
        return obj.Value
    case *object.Null:
        return false
    default:
        return true
    }
}

func newError(format string, arguments ...any) *object.Error {
    return &object.Error{Message: fmt.Sprintf(format, arguments...)}
}

func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.RETURN_VALUE_OBJ
    }
    return false
}
