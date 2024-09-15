package evaluator

import (
	"charm/lexer"
	"charm/object"
	"charm/parser"
	"testing"
    "math"
)

func TestEvalIntegerExpression(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"5;", 5},
        {"10;", 10},

        {"-5;", -5},
        {"-10;", -10},

        {"5 + 5 + 5 + 5 - 10;", 10},
        {"2 * 2 * 2 * 2 * 2;", 32},
        {"-50 + 100 + -50;", 0},
        {"5 * 2 + 10;", 20},
        {"5 + 2 * 10;", 25},
        {"20 + 2 * -10;", 0},
        {"50 / 2 * 2 + 10;", 60},
        {"2 * (5 + 10);", 30},
        {"3 * 3 * 3 + 10;", 37},
        {"3 * (3 * 3) + 10;", 37},
        {"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testIntegerObject(t, evaluated, test.expected)
    }
}

func TestEvalFloatExpression(t *testing.T) {
    tests := []struct {
        input string
        expected float64
    } {
        {"5.2;", 5.2},
        {"10.5431;", 10.5431},

        {"-5.2;", -5.2},
        {"-10.5431;", -10.5431},

        {"5.123 + 5.456 + 5.789 + 5.321 - 10.987;", 10.702},
        {"2.345 * 2.678 * 2.123 * 2.987 * 2.111;", 84.067256},
        {"-50.8765 + 100.4321 + -50.1234;", -0.5678},
        {"5.2345 * 2.1234 + 10.5678;", 21.6827},
        {"5.8765 + 2.3456 * 10.4321;", 30.34603376},
        {"20.6789 + 2.3456 * -10.9876;", -5.09361456},
        {"50.9876 / 2.2345 * 2.6789 + 10.4321;", 71.56017413},
        {"2.3456 * (5.7891 + 10.4321);", 38.04844672},
        {"3.1234 * 3.4321 * 3.9876 + 10.5678;", 53.31415878},
        {"3.9876 * (3.8765 * 3.2345) + 10.1234;", 60.12207911},
        {"(5.4321 + 10.9876 * 2.1234 + 15.8765 / 3.1234) * 2.9876 + -10.4321;", 90.68696361},

    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testFloatObject(t, evaluated, test.expected)
    }
}

func TestEvalBooleanExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"true;", true},
        {"false;", false},

        {"!true;", false},
        {"!false;", true},
        {"!5;", false},
        {"!!true;", true},
        {"!!false;", false},
        {"!!5;", true},

        {"true == true;", true},
        {"false == false;", true},
        {"true == false;", false},
        {"true != false;", true},
        {"false != true;", true},
        {"(1 < 2) == true;", true},
        {"(1 < 2) == false;", false},
        {"(1 > 2) == true;", false},
        {"(1 > 2) == false;", true},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testBooleanObject(t, evaluated, test.expected)
    }
}

func TestStringLiteral(t *testing.T) {
    tests := []struct {
        input string
        expected string
    } {
        {`"Hello World!";`, "Hello World!"},
        {`"Hello" + " Concat";`, "Hello Concat"},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)

        str, ok := evaluated.(*object.String)
        if !ok {
            t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
        }

        if str.Value != test.expected {
            t.Errorf("String has wrong value. Expected=%q got=%q", test.expected, str.Value)
        }
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
            "-true;",
            "unknown operator: -BOOLEAN",
        },
        {
            "true + false;",
            "unknown operator: BOOLEAN + BOOLEAN",
        },
        {
            "5; true + false; 5;",
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
        {"if (true) { 10; }", 10},
        {"if (false) { 10; }", nil},
        {"if (1) { 10; }", 10},
        {"if (1 < 2) { 10; }", 10},
        {"if (1 > 2) { 10; }", nil},
        {"if (1 > 2) { 10; } else { 20; }", 20},
        {"if (1 < 2) { 10; } else { 20; }", 10},
        {"if (1 < 2) { 10; 20; 30; } else { 20; }", 30},
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

func TestWhileStatement(t *testing.T) {
    input := `
    x = 10;
    while (x > 0) { x = x - 1; }
    `

    evaluated := evalTest(input)

    int, ok := evaluated.(*object.Integer)
    if !ok {
        t.Fatalf("Expected=%s. Got=%s", object.INTEGER_OBJ, evaluated.Type())
    }

    testIntegerObject(t, int, 0)
}

func TestReturnInWhileLoop(t *testing.T) {
    input := `
    x = 10;
    while (x > 0) { 
        x = x - 1;
        return x;
    }
    `
    evaluated := evalTest(input)

    int, ok := evaluated.(*object.Integer)
    if !ok {
        t.Fatalf("Expected=%s. Got=%s", object.INTEGER_OBJ, evaluated.Type())
    }

    testIntegerObject(t, int, 9)
}

func TestWhileNoEval(t *testing.T) {
    input := `
    x = 10;
    while (x > 120) {
        x = x - 1;
        return x;
    }
    `
    evaluated := evalTest(input)

    testNullObject(t, evaluated)
}

func TestAssignmentStatement(t *testing.T) {
    tests := []struct {
        input string
        expected int64
    } {
        {"a = 5; a;", 5},
        {"a = 5 * 5; a;", 25},
        {"a = 5; b = a; b;", 5},
        {"a = 5; b = a; c = a + b + 5; c;", 15},
    }

    for _, test := range tests {
        testIntegerObject(t, evalTest(test.input), test.expected)
    }
}

func TestFunctionObject(t *testing.T) {
    test := "func(x) { x + 2; };"
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

func TestFunctionStatement(t *testing.T) {
    test := "func z(x) { x + 2; }"
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
        {"identity = func(x) { x; }; identity(5);", 5},
        {"identity = func(x) { return x; }; identity(5);", 5},
        {"double = func(x) { x * 2; }; double(5);", 10},
        {"add = func(x, y) { x + y; }; add(5, 5);", 10},
        {"add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
        {"func(x) { x; }(5);", 5},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        testIntegerObject(t, evaluated, test.expected)
    }
}

func TestClosures(t *testing.T) {
    input := `
    newAdder = func(x) {
        func(y) { x + y; };
    };
    addTwo = newAdder(2);
    addTwo(2);`

    testIntegerObject(t, evalTest(input), 4)
}

func TestArrayLiterals(t *testing.T) {
    input := "[1, 2 * 2, 3 + 3];"
    evaluated := evalTest(input)
    result, ok := evaluated.(*object.Array)
    if !ok {
        t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
    }
    if len(result.Elements) != 3 {
        t.Fatalf("array has wrong num of elements. got=%d",
            len(result.Elements))
    }
    testIntegerObject(t, result.Elements[0], 1)
    testIntegerObject(t, result.Elements[1], 4)
    testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected interface{}
    }{
        {
            "[1, 2, 3][0];", 1,
        },
        {
            "[1, 2, 3][1];", 2,
        },
        {
            "[1, 2, 3][2];", 3,
        },
        {
            "i = 0; [1][i];", 1,
        },
        {
            "[1, 2, 3][1 + 1];", 3,
        },
        {
            "myArray = [1, 2, 3]; myArray[2];", 3,
        },
        {
            "myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6,
        },
        {
            "myArray = [1, 2, 3]; i = myArray[0]; myArray[i];", 2,
        },
        {
            "[1, 2, 3][3];", nil,
        },
        {
            "[1, 2, 3][-1];", nil,
        },
    }
    for _, tt := range tests {
        evaluated := evalTest(tt.input)
        integer, ok := tt.expected.(int)
        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        } else {
            testNullObject(t, evaluated)
        }
    }
}

func TestHashableObjects(t *testing.T) {
    tests := []struct {
        left object.Hashable
        right object.Hashable
        shouldMatch bool
    } {
        {&object.String{Value: "Hello"}, &object.String{Value: "Hello"}, true},
        {&object.String{Value: "Hello"}, &object.String{Value: "World"}, false},
        {&object.Integer{Value: 1234}, &object.Integer{Value: 1234}, true},
        {&object.Integer{Value: 1234}, &object.Integer{Value: 43321}, false},
        {&object.Boolean{Value: true}, &object.Boolean{Value: false}, false},
        {&object.Boolean{Value: true}, &object.Boolean{Value: true}, true},
    }

    for _, test := range tests {
        leftHashCode := test.left.HashCode()
        rightHashCode := test.right.HashCode()

        if test.shouldMatch {
            if leftHashCode != rightHashCode {
                t.Errorf("Hashcodes do not match. left=%s, right=%s", test.left, test.right)
            }
        } else {
            if leftHashCode == rightHashCode {
                t.Errorf("Hashcodes match. left=%s, right=%s", test.left, test.right)
            }
        }
    }
}

func TestHashMapLiteral(t *testing.T) {
    input := `two = "two";
    {
    "one": 10 - 9,
    two: 1 + 1,
    "thr" + "ee": 6 / 2,
    4: 4,
    true: 5,
    false: 6
    };`

    evaluated := evalTest(input)
    hashMapObj, ok := evaluated.(*object.HashMap)
    if !ok {
        t.Fatalf("object not hashMap. Got=%s", evaluated.Type())
    }

    one:= (&object.String{Value: "one"})
    two:= (&object.String{Value: "two"})
    three:= (&object.String{Value: "three"})
    four:= (&object.Integer{Value: 4})
    five:= (&object.Boolean{Value: true})
    six:= (&object.Boolean{Value: false})

    type expectedPair struct {
        expectedObject object.Hashable
        expectedValue int64
    }

    expected := map[uint64]expectedPair {
        one.HashCode() : {one, 1},
        two.HashCode() : {two, 2},
        three.HashCode() : {three, 3},
        four.HashCode() : {four, 4},
        five.HashCode() : {five, 5},
        six.HashCode() : {six, 6},
    }

    if len(expected) != len(hashMapObj.Map) {
        t.Fatalf("length does not match. Expected=%d, Got=%d", len(expected), len(hashMapObj.Map))
    }

    for _, expectedPair := range expected {
        actualPair := hashMapObj.Map[expectedPair.expectedObject.HashCode()]

        // testing keys
        switch expectedObj := expectedPair.expectedObject.(type) {
        case *object.String:
            if !testStringObject(t, actualPair.Key, expectedObj.Value) {
                continue
            }
        case *object.Boolean:
            if !testBooleanObject(t, actualPair.Key, expectedObj.Value) {
                continue
            }
        case *object.Integer:
            if !testIntegerObject(t, actualPair.Key, expectedObj.Value) {
                continue
            }
        default:
            t.Errorf("unexpected val for key: map[%s]: %s. Expected: map[%s]: %d",
                actualPair.Key.Inspect(), actualPair.Value.Type(), 
                expectedPair.expectedObject.Inspect(), expectedPair.expectedValue)
            continue
        }

        // testing values
        if !testIntegerObject(t, actualPair.Value, expectedPair.expectedValue) {
            continue
        }
    }
}

func TestHashIndexExpressions(t *testing.T) {
    tests := []struct {
        input string
        expected any
    } {
        {
            `{"foo": 5}["foo"];`,
            5,
        },
        {
            `{"foo": 5}["bar"];`,
            nil,
        },
        {
            `key = "foo"; {"foo": 5}[key];`,
            5,
        },
        {
            `{}["foo"];`,
            nil,
        },
        {
            `{5: 5}[5];`,
            5,
        },
        {
            `{true: 5}[true];`,
            5,
        },
        {
            `{false: 5}[false];`,
            5,
        },
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)
        integer, ok := test.expected.(int)
        if ok {
            testIntegerObject(t, evaluated, int64(integer))
        } else {
            testNullObject(t, evaluated)
        }
    }
}

// helpers
func evalTest(input string) object.Object {
    lexer := lexer.New(input)
    parser := parser.New(lexer)
    program := parser.ParseProgram()
    environment := object.NewEnvironment()

    return Eval(program, environment)
}

func testIntegerObject(t *testing.T, evaluated object.Object, expected int64) bool{
    intObj, ok := evaluated.(*object.Integer)
    if !ok {
        t.Errorf("object is not an integer. Got=%T (%+v)", evaluated, evaluated)
        return false
    }

    if intObj.Value != expected {
        t.Errorf("object has wrong value. Expected %d, got %d\n", expected, intObj.Value)
        return false
    }

    return true
}

func testFloatObject(t *testing.T, evaluated object.Object, expected float64) bool{
    epsilon := 0.0001

    floatObj, ok := evaluated.(*object.Float)
    if !ok {
        t.Errorf("object is not an float. Got=%T (%+v)", evaluated, evaluated)
        return false
    }

    if math.Abs(floatObj.Value - expected) >= epsilon {
        t.Errorf("object has wrong value. Expected %f, got %f\n", expected, floatObj.Value)
        return false
    }

    return true
}

func testStringObject(t *testing.T, evaluated object.Object, expected string) bool {
    strObj, ok := evaluated.(*object.String)
    if !ok {
        t.Errorf("object is not String. got=%T (%+v)", evaluated, evaluated)
        return false
    }

    if strObj.Value != expected {
        t.Errorf("String has wrong value. Expected=%q got=%q", expected, strObj.Value)
        return false
    }

    return true
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
    boolObj, ok := evaluated.(*object.Boolean)
    if !ok {
        t.Errorf("object is not an boolean. Got=%T (%+v)", evaluated, evaluated)
        return false
    }

    if boolObj.Value != expected {
        t.Errorf("object has wrong value. Expected %t, got %t\n", expected, boolObj.Value)
        return false
    }

    return true
}

func testNullObject(t *testing.T, evaluated object.Object) bool {
    if evaluated != NULL {
        t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
        return false
    }

    return true
}
