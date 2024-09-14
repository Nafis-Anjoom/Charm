package evaluator

import (
	"charm/lexer"
	"charm/object"
	"charm/parser"
	"testing"
)

func TestBuiltinLenFunction(t *testing.T) {
    tests := []struct {
        input string
        expected any
    } {
        {`len("");`, 0},
        {`len("four");`, 4},
        {`len("hello world");`, 11},
        {`len(1);`, "argument to `len` not supported, got INTEGER"},
        {`len("one", "two");`, "wrong number of arguments. got=2, want=1"},
        {`len([1, "string", true, 4]);`, 4},
        {`len([]);`, 0},
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
        {`x = [1, 2, 3]; push(x, 4);`, []int{1, 2, 3, 4}},
        {`x = []; push(x, 4);`, []int{4}},
        {`x = "string"; push(x, 4);`, "argument to `push` must be ARRAY, got STRING"},
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
        {`x = [1, 2, 3]; pop(x);`, []int64{1, 2}, 3},
        {`x = [1]; pop(x);`, []int64{}, 1},
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

func TestBuiltinKeysFunction(t *testing.T) {
    input := `two = "two";
    keys({
    "one": 10 - 9,
    two: 1 + 1,
    "thr" + "ee": 6 / 2,
    4: 4,
    true: 5,
    false: 6
    });`

    expectedKeys := map[any]int {
        "one": 0,
        "two": 0,
        "three": 0,
        4: 1,
        true: 0,
        false: 0,
    }

    evaluated := evalTest(input)

    keysArrayobj, ok := evaluated.(*object.Array)
    if !ok {
        t.Fatalf("object not Array. Got=%s", evaluated.Type())
    }

    keys := keysArrayobj.Elements
    if len(keys) != len(expectedKeys) {
        t.Fatalf("Incorrect length. Expected=%d. Got=%d", len(expectedKeys), len(keys))
    }

    for _, key := range keys {
        switch keyVal := key.(type) {
        case *object.String:
            if _, ok = expectedKeys[keyVal.Value]; !ok {
                t.Fatalf("unexpected key: %s", keyVal.Inspect())
            }
            expectedKeys[keyVal.Value] = 1
        case *object.Boolean:
            if _, ok = expectedKeys[keyVal.Value]; !ok {
                t.Fatalf("unexpected key: %s", keyVal.Inspect())
            }
            expectedKeys[keyVal.Value] = 1
        case *object.Integer:
            if _, ok = expectedKeys[int(keyVal.Value)]; !ok {
                t.Fatalf("unexpected key: %s", keyVal.Inspect())
            }
            expectedKeys[keyVal.Value] = 1
        }
    }

    for key, val := range expectedKeys {
        if val != 1 {
            t.Fatalf("key missing: %s", key)
        }
    }
}

func TestBuiltinDeleteFunction(t *testing.T) {
    input := `
    x = {
    "one": 10 - 9,
    "two": 1 + 1,
    };
    
    delete(x, "two");
    delete(x, "three");
    `

    lexer := lexer.New(input)
    parser := parser.New(lexer)
    program := parser.ParseProgram()

    env := object.NewEnvironment()
    Eval(program, env)

    x, ok := env.Get("x")
    if !ok {
        t.Fatal("x not found")
    }

    hashMap, ok := x.(*object.HashMap)
    if !ok {
        t.Fatalf("x is not an hashMap: %s", x.Type())
    }

    if len(hashMap.Map) != 1 {
        t.Fatalf("Incorrect length. Expected %d. Got %d", 1, len(hashMap.Map))
    }

    _, ok = hashMap.Map[(&object.String{Value: "one"}).HashCode()]
    if !ok {
        t.Fatalf("Missing entry")
    }
}
