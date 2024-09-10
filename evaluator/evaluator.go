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
        if node.Value {
            return TRUE
        }
        return FALSE
    case *ast.PrefixExpression:
        return evalPrefixExpression(node)
        
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

func evalPrefixExpression(exp *ast.PrefixExpression) object.Object {
    switch exp.Operator {
    case "!":
        right := Eval(exp.Right)
        return evalBangExpression(right)
    default:
        return NULL
    }
}

// Integers are truthy. Nulls are falsy
func evalBangExpression(right object.Object) object.Object {
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
