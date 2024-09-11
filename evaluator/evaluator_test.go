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
        {"false == false", true},
        {"true == false", false},
        {"true != false", true},
        {"false != true", true},
        {"(1 < 2) == true", true},
        {"(1 < 2) == false", false},
        {"(1 > 2) == true", false},
        {"(1 > 2) == false", true},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testBooleanObject(t, evaluated, test.expected)
    }
}

func TestErrorHandling(t *testing.T) {
    tests := []struct {
        input string
        expectedMessage string
    } {
        {
            "5 + true;",
            "type mismatch: INTEGER + BOOLEAN",
        },
        {
            "5 + true; 5;",
            "type mismatch: INTEGER + BOOLEAN",
        },
        {
            "-true",
            "unknown operator: -BOOLEAN",
        },
        {
            "true + false;",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            "5; true + false; 5",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            "if (10 > 1) { true + false; }",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            `if (10 > 1) {
            if (10 > 1) {
            return true + false;
            }
            return 1;
            }`, "unknown operator: BOOLEAN + BOOLEAN",
        },
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)

        error, ok := evaluated.(*object.Error)
        if !ok {
            t.Errorf("No error returned. Expected message %s. Got %T(%+v)", test.expectedMessage, evaluated, evaluated)
            continue
        }

        if error.Message != test.expectedMessage {
            t.Errorf("wrong message received. Expected=%s. Got=%s", test.expectedMessage, error.Message)
            continue
        }
    }
}

// TODO: refactor tests to run every sub-test in its own goroutine
func TestIfElseExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected any
    } {
        {"if (true) { 10 }", 10},
        {"if (false) { 10 }", nil},
        {"if (1) { 10 }", 10},
        {"if (1 < 2) { 10 }", 10},
        {"if (1 > 2) { 10 }", nil},
        {"if (1 > 2) { 10 } else { 20 }", 20},
        {"if (1 < 2) { 10 } else { 20 }", 10},
        {"if (1 < 2) { 10; 20; 30; } else { 20 }", 30},
    }

    for _, test := range tests {
        evalulated := evalTest(test.input)
        integer, ok := test.expected.(int)
        if ok {
            testIntegerObject(t, evalulated, int64(integer))
        } else {
            testNullObject(t, evalulated)
        }
    }
}

func TestReturnStatements(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"return 10;", 10},
        {"return 10; 11;", 10},
        {"return 2 * 5; 9;", 10},
        {"9; return 2 * 5; 9;", 10},
        {`if (10 > 1) {
            if (10 > 1) {
                return 10;
            }
            return 1;
            }`, 10,
        },
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testIntegerObject(t, evaluated, test.expected)
    }
}

func TestLetStatement(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"let a = 5; a;", 5},
        {"let a = 5 * 5; a;", 25},
        {"let a = 5; let b = a; b;", 5},
        {"let a = 5; let b = a; let c = a + b + 5; c;", 15},
    }

    for _, test := range tests {
        testIntegerObject(t, evalTest(test.input), test.expected)
    }
}

func TestFunctionObject(t *testing.T) {
    test := "fn(x) {x + 2; };"
    evaluated := evalTest(test)

    funcObj, ok := evaluated.(*object.Function)
    if !ok {
        t.Errorf("object not Function. Got type=%T(%+v)", evaluated, evaluated)
        return
    }

    params := funcObj.Parameters
    if len(params) != 1 {
        t.Fatalf("Function has more than 1 parameter. Got=%d)", len(params))
    }
    if params[0].String() != "x" {
        t.Fatalf("Function paramter not X. Got=%s)", params[0].String())
    }

    expectedBody := "(x + 2)"
    if funcObj.Body.String() != expectedBody {
        t.Fatalf("body is %q. got %q", expectedBody, funcObj.Body.String())
    }
}

func TestFunctionCall(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"let identity = fn(x) { x; }; identity(5);", 5},
        {"let identity = fn(x) { return x; }; identity(5);", 5},
        {"let double = fn(x) { x * 2; }; double(5);", 10},
        {"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
        {"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
        {"fn(x) { x; }(5)", 5},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testIntegerObject(t, evaluated, test.expected)
    }
}

func TestClosures(t *testing.T) {
    input := `
    let newAdder = fn(x) {
        fn(y) { x + y };
    };
    let addTwo = newAdder(2);
    addTwo(2);`

    testIntegerObject(t, evalTest(input), 4)
}

// helpers
func evalTest(input string) object.Object {
    lexer := lexer.New(input)
    parser := parser.New(lexer)
    program := parser.ParseProgram()
    environment := object.NewEnvironment()

    return Eval(program, environment)
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
        t.Errorf("object is not an boolean. Got=%T (%+v)", evaluated, evaluated)
    }

    if boolObj.Value != expected {
        t.Errorf("object has wrong value. Expected %t, got %t\n", expected, boolObj.Value)
    }
}

func testNullObject(t *testing.T, evaluated object.Object) {
    if evaluated != NULL {
        t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
    }
}
