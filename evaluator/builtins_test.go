package evaluator

import (
	"charm/object"
	"charm/lexer"
	"charm/parser"
	"testing"
)

func TestBuiltinLenFunction(t *testing.T) {
    tests := []struct {
        input string
        expected any
    } {
        {`len("")`, 0},
        {`len("four")`, 4},
        {`len("hello world")`, 11},
        {`len(1)`, "argument to `len` not supported, got INTEGER"},
        {`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
        {`len([1, "string", true, 4])`, 4},
        {`len([])`, 0},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)

        switch expected := test.expected.(type) {
        case int:
            testIntegerObject(t, evaluated, int64(expected))
        case string:
            errObj, ok := evaluated.(*object.Error)
            if !ok {
                t.Errorf("object is not Error. Got=%T (%+v)", evaluated, evaluated)
                continue
            }

            if errObj.Message != expected {
                t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
            }
        }
    }
}

func TestBuiltinPushFunction(t *testing.T) {
    tests := []struct {
        input string
        expected any
    } {
        {`let x = [1, 2, 3]; push(x, 4);`, []int{1, 2, 3, 4}},
        {`let x = []; push(x, 4);`, []int{4}},
        {`let x = "string"; push(x, 4);`, "argument to `push` must be ARRAY, got STRING"},
    }

    for _, test := range tests {
        evaluated := evalTest(test.input)

        switch expected := test.expected.(type) {
        case []int:
            arrayObj, ok := evaluated.(*object.Array)
            if !ok {
                t.Errorf("object not Array. Got=%T (%+v)", evaluated, evaluated)
                continue
            }
            if len(expected) != len(arrayObj.Elements) {
                t.Errorf("array length incorrect. Expected=%d. Got=%d", len(expected), len(arrayObj.Elements))
            }

            for i := range len(expected) {
                if !testIntegerObject(t, arrayObj.Elements[i], int64(expected[i])) {
                    continue
                }
            }
        case string:
            errObj, ok := evaluated.(*object.Error)
            if !ok {
                t.Errorf("object is not Error. Got=%T (%+v)", evaluated, evaluated)
                continue
            }

            if errObj.Message != expected {
                t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
            }
        }
    }
}

func TestBuiltinPopFunction(t *testing.T) {
    tests := []struct {
        input string
        expectedArray []int64
        expectedRetValue int64
    } {
        {`let x = [1, 2, 3]; pop(x);`, []int64{1, 2}, 3},
        {`let x = [1]; pop(x);`, []int64{}, 1},
    }

    for _, test := range tests {
        lexer := lexer.New(test.input)
        parser := parser.New(lexer)
        program := parser.ParseProgram()
        environment := object.NewEnvironment()

        evaluated := Eval(program, environment)

        intObj, ok := evaluated.(*object.Integer)
        if !ok {
            t.Errorf("object not Integer. Got=%T (%+v)", evaluated, evaluated)
            continue
        }

        if !testIntegerObject(t, intObj, test.expectedRetValue) {
            t.Errorf("incorrect return value. Expected=%d. Got=%d", test.expectedRetValue, intObj.Value)
            continue
        }

        poppedArray, ok := environment.Get("x")
        if !ok {
            t.Errorf("variable does not exist: %s", "x")
            continue
        }

        arrayObj, ok := poppedArray.(*object.Array)
        if !ok {
            t.Errorf("object not Array. Got=%T (%+v)", evaluated, evaluated)
            continue
        }

        if len(test.expectedArray) != len(arrayObj.Elements) {
            t.Errorf("array length incorrect. Expected=%d. Got=%d", len(test.expectedArray), len(arrayObj.Elements))
            continue
        }

        for i := range len(test.expectedArray) {
            if !testIntegerObject(t, arrayObj.Elements[i], test.expectedArray[i]) {
                continue
            }
        }
    }
}
