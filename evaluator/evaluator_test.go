package evaluator

import (
	"charm/lexer"
	"charm/object"
	"charm/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"5", 5},
        {"10", 10},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testIntegerObject(t, evaluated, test.expected)
    }
}

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"true;", true},
        {"false;", false},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testBooleanObject(t, evaluated, test.expected)
    }
}

func TestBangOperator(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testBooleanObject(t, evaluated, test.expected)
    }
}

// helpers
func evalTest(input string) object.Object {
    lexer := lexer.New(input)
    parser := parser.New(lexer)
    program := parser.ParseProgram()

    return Eval(program)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) {
    intObj, ok := evaluated.(*object.Integer)
    if !ok {
        t.Fatalf("object is not an integer. Got=%T (%+v)", evaluated, evaluated)
    }

    if intObj.Value != expected {
        t.Fatalf("object has wrong value. Expected %d, got %d\n", expected, intObj.Value)
    }
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) {
    boolObj, ok := evaluated.(*object.Boolean)
    if !ok {
        t.Fatalf("object is not an boolean. Got=%T (%+v)", evaluated, evaluated)
    }

    if boolObj.Value != expected {
        t.Fatalf("object has wrong value. Expected %t, got %t\n", expected, boolObj.Value)
    }
}
