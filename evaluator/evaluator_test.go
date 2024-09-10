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

        {"-5", -5},
        {"-10", -10},

        {"5 + 5 + 5 + 5 - 10", 10},
        {"2 * 2 * 2 * 2 * 2", 32},
        {"-50 + 100 + -50", 0},
        {"5 * 2 + 10", 20},
        {"5 + 2 * 10", 25},
        {"20 + 2 * -10", 0},
        {"50 / 2 * 2 + 10", 60},
        {"2 * (5 + 10)", 30},
        {"3 * 3 * 3 + 10", 37},
        {"3 * (3 * 3) + 10", 37},
        {"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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

        {"!true", false},
        {"!false", true},
        {"!5", false},
        {"!!true", true},
        {"!!false", false},
        {"!!5", true},

        {"true == true", true},
        // {"false == false", true},
        // {"true == false", false},
        // {"true != false", true},
        // {"false != true", true},
        // {"(1 < 2) == true", true},
        // {"(1 < 2) == false", false},
        // {"(1 > 2) == true", false},
        // {"(1 > 2) == false", true},
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
