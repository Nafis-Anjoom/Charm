package evaluator

import (
	"charm/ast"
	"charm/object"
)

var (
    TRUE = &object.Boolean{Value: true}
    FALSE = &object.Boolean{Value: false}
    NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
    case *ast.Program:
        return evalStatements(node.Statements)
    case *ast.ExpressionStatement:
        return Eval(node.Expression)
    case *ast.IntegerLiteral:
        return &object.Integer{Value: node.Value}
    case *ast.BooleanLiteral:
        return nativeBooltoBoolObject(node.Value)
    case *ast.PrefixExpression:
        return evalPrefixExpression(node)
    case *ast.InfixExpression:
        return evalInfixExpression(node)
    }
    return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
    var result object.Object

    for _, stmt := range stmts {
        result = Eval(stmt)
    }

    return result
}

// DNOTE: return NULL when encountering unknown prefix operator 
func evalPrefixExpression(exp *ast.PrefixExpression) object.Object {
    switch exp.Operator {
    case "!":
        right := Eval(exp.Right)
        return evalBangPrefixExpression(right)
    case "-":
        right := Eval(exp.Right)
        return evalMinusPrefixExpression(right)
    default:
        return NULL
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
        return NULL
    }
    value := right.(*object.Integer).Value
    return &object.Integer{Value: -value}
}

func evalInfixExpression(exp *ast.InfixExpression) object.Object {
    right := Eval(exp.Right)
    left := Eval(exp.Left)

    switch {
    case right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ:
        rightValue := right.(*object.Integer).Value
        leftValue := left.(*object.Integer).Value
        return evalIntegerInfixExpression(leftValue, exp.Operator, rightValue)
    case exp.Operator == "==":
        return nativeBooltoBoolObject(left == right)
    case exp.Operator == "!=":
        return nativeBooltoBoolObject(left != right)
    default:
        return NULL
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

func nativeBooltoBoolObject(boolean bool) *object.Boolean {
    if boolean {
        return TRUE
    }
    return FALSE
}
