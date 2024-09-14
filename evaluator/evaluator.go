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
    case *ast.StringLiteral:
        return &object.String{Value: node.Value}
    case *ast.BooleanLiteral:
        return nativeBooltoBoolObject(node.Value)
    case *ast.PrefixExpression:
        return evalPrefixExpression(node, env)
    case *ast.InfixExpression:
        return evalInfixExpression(node, env)
    case *ast.IfStatement:
        return evalIfStatement(node, env)
    case *ast.WhileStatement:
        return evalWhileStatemnt(node, env)
    case *ast.BlockStatement:
        return evalStatements(node.Statements, env)
    case *ast.ReturnStatement:
        returnValue := Eval(node.ReturnValue, env)
        if isError(returnValue) {
            return returnValue
        }
        return &object.ReturnValue{Value: returnValue}
    case *ast.AssignmentStatement:
        value := Eval(node.Value, env)
        if isError(value) {
            return value
        }
        env.Set(node.Identifier.Value, value)
        return value
    case *ast.FunctionStatement:
        return evalFunctionStatement(node, env)
    case *ast.Identifier:
        return evalIdentifier(node, env)
    case *ast.FunctionLiteral:
        return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
    case *ast.CallExpression:
        return evalCallExpression(node, env)
    case *ast.ArrayLiteral:
        return evalArrayLiteral(node, env)
    case *ast.IndexExpression:
        return evalIndexExpression(node, env)
    case *ast.HashMapLiteral:
        return evalHashMapLiteral(node, env)
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
            fmt.Println(result.Message)
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
    case right.Type() == object.STRING_OBJ && left.Type() == object.STRING_OBJ && exp.Operator == "+":
        rightValue := right.(*object.String).Value
        leftValue := left.(*object.String).Value
        return &object.String{Value: leftValue + rightValue}
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

func evalIfStatement(node *ast.IfStatement, env *object.Environment) object.Object {
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

func evalWhileStatemnt(node *ast.WhileStatement, env *object.Environment) object.Object {
    condition := Eval(node.Condition, env)
    
    if isError(condition) {
        return condition
    }

    var evaluated object.Object = NULL
    for isTruthy(condition) {
        evaluated = Eval(node.Body, env)

        if evaluated.Type() == object.RETURN_VALUE_OBJ {
            return evaluated
        }

        condition = Eval(node.Condition, env)
    }

    return evaluated
}

func evalFunctionStatement(stmt *ast.FunctionStatement, env *object.Environment) object.Object {
    function := &object.Function{
        Parameters: stmt.FunctionLiteral.Parameters,
        Body: stmt.FunctionLiteral.Body,
        Env: env,
    }

    env.Set(stmt.Identifier.Value, function)

    return function
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
    arguments := evalExpressions(node.Arguments, env)
    if len(arguments) == 1 && arguments[0].Type() == object.ERROR_OBJ {
        return arguments[0]
    }

    obj := Eval(node.FunctionLiteral, env)

    switch functionObj := obj.(type) {
    case *object.Function:
        if len(arguments) > len(functionObj.Parameters) {
            return newError("too many arguments")
        }
        if len(arguments) < len(functionObj.Parameters) {
            return newError("not enough arguments")
        }

        enclosedEnv := object.NewEnclosedEnvironment(functionObj.Env)

        for i, param := range functionObj.Parameters {
            enclosedEnv.Set(param.Value, arguments[i])
        }

        evaluated := Eval(functionObj.Body, enclosedEnv)
        return unwrapReturnValue(evaluated)
    case *object.Builtin:
        return functionObj.Fn(arguments...)
    default:
        return newError("not a function: %s", obj.Type())
    }
}

func evalIdentifier(Identifier *ast.Identifier, env *object.Environment) object.Object {
    if val, ok := env.Get(Identifier.Value); ok {
        return val
    }

    if builtin, ok := builtins[Identifier.Value]; ok {
        return builtin
    }

    return newError("identifier not found: %s", Identifier.Value)
}

func evalIndexExpression(indexExpr *ast.IndexExpression, env *object.Environment) object.Object {
    obj := Eval(indexExpr.Left, env)
    // what happens if this is NULL?
    indexObj := Eval(indexExpr.Index, env)

    if isError(indexObj) || indexObj == NULL {
        return newError("error evaluating index: %s", indexObj.Inspect())
    }

    switch obj := obj.(type) {
    case *object.Array:
        return evalArrayIndexExpression(obj, indexObj, env)
    case *object.HashMap:
        return evalHashMapIndexExpression(obj, indexObj, env)
    default:
        return newError("index operator not supported: %s", obj.Type())
    }
}

func evalArrayIndexExpression(arrayObj *object.Array, indexObj object.Object, env *object.Environment) object.Object {
    index, ok := indexObj.(*object.Integer)
    if !ok {
        return newError("not an integer: %T", index)
    }

    if int(index.Value) >= len(arrayObj.Elements) || index.Value < 0 {
        return NULL
    }

    return arrayObj.Elements[index.Value]
}

func evalHashMapIndexExpression(hashMapObj *object.HashMap, indexObj object.Object, env *object.Environment) object.Object {
    HashObj, ok := indexObj.(object.Hashable)
    if !ok {
        return newError("unusable as haskey: %s", indexObj.Type())
    }

    hashCode := HashObj.HashCode()

    pair, ok := hashMapObj.Map[hashCode]
    if !ok {
        return NULL
    }
    return pair.Value
}

func evalHashMapLiteral(hashMap *ast.HashMapLiteral, env *object.Environment) object.Object {
    hashMapObj := &object.HashMap{
        Map : make(map[uint64]object.Pair),
    }

    for key, val := range hashMap.Map {
        keyObj := Eval(key, env)
        if isError(keyObj) {
            return keyObj
        }
        
        hashableKey, ok := keyObj.(object.Hashable)
        if !ok {
            return newError("Object not hashable: %s", keyObj.Type())
        }

        valObj := Eval(val, env)
        if isError(valObj) {
            return valObj
        }

        hashMapObj.Map[hashableKey.HashCode()] = object.Pair{Key: keyObj, Value: valObj}
    }

    return hashMapObj
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
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}

func evalExpressions(args []ast.Expression, env *object.Environment) []object.Object {
    evaledArgs := []object.Object{}
    for _, arg := range args {
        evaluated := Eval(arg, env)
        if isError(evaluated) {
            return []object.Object{ evaluated }
        }
        evaledArgs = append(evaledArgs, evaluated)
    }
    return evaledArgs
}

func evalArrayLiteral(node *ast.ArrayLiteral, env *object.Environment) object.Object {
    elements := []object.Object{}
    var evaluated object.Object

    for _, exp := range node.Elements {
        evaluated = Eval(exp, env)
        if evaluated != nil && evaluated.Type() == object.ERROR_OBJ {
            return evaluated
        }

        elements = append(elements, evaluated)
    }

    return &object.Array{Elements: elements}
}

func unwrapReturnValue(obj object.Object) object.Object {
    if returnValue, ok := obj.(*object.ReturnValue); ok {
        return returnValue.Value
    }
    return obj
}
