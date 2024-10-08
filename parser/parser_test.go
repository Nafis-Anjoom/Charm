package parser

import (
	"charm/ast"
	"charm/lexer"
	"fmt"
	"testing"
)

func TestLetStatements(t *testing.T) {
    tests := []struct {
        input string
        expectedIdentifier string
        expectedValue any
    } {
        {input: "x = 5;", expectedIdentifier: "x", expectedValue: 5},
        {input: "y = 10;", expectedIdentifier: "y", expectedValue: 10},
        {input: "foobar = 838383;", expectedIdentifier: "foobar", expectedValue: 838383},
    }

    for _, test := range tests {

        lexer := lexer.New(test.input)
        parser := New(lexer)
        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if program == nil {
            t.Fatal("ParseProgram returned nil\n")
        }

        if len(program.Statements) != 1 {
            t.Fatalf("Expected 1 statements. Got %d\n", len(program.Statements))
        }

        stmt := program.Statements[0]
        if !testLetStatement(t, stmt, test.expectedIdentifier) {
            return 
        }

        val := stmt.(*ast.AssignmentStatement).Value
        if !testLiteralExpression(t, val, test.expectedValue) {
            return
        }
    }
}

func TestReturnStatement(t *testing.T) {
    tests := []struct {
        input string
        expectedReturnValue any
    } {
        
        {input: "return 5;", expectedReturnValue: "5"},
        {input: "return foobar;", expectedReturnValue: "foobar"},
        {input: "return 10 + 10;", expectedReturnValue: "(10 + 10)"},
    }

    for _, test := range tests {
        lexer := lexer.New(test.input)
        parser := New(lexer)
        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain 1 statements. got=%d\n", len(program.Statements))
        }

        stmt := program.Statements[0]
        if stmt.TokenLiteral() != "return" {
            t.Errorf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
        }

        returnStmt, ok := stmt.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("stmt not *ast.ReturnStatement. got=%T\n", stmt)
            continue
        }

        if returnStmt.ReturnValue.String() != test.expectedReturnValue {
            t.Fatalf("Expected %s. Got=%s\n", test.expectedReturnValue, returnStmt.ReturnValue.String())
        }
    }
}

func TestWhileStatement(t *testing.T) {
    input := `while (x < 10) { x; }`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Statements does not contain 1 statements. got=%d\n", len(program.Statements))
    }

    stmt := program.Statements[0]
    if stmt.TokenLiteral() != "while" {
        t.Fatalf("stmt.TokenLiteral not 'while'. got=%q", stmt.TokenLiteral())
    }

    whileStmt, ok := stmt.(*ast.WhileStatement)
    if !ok {
        t.Fatalf("stmt not *ast.WhileStatement. got=%T\n", stmt)
    }

    if !testInfixExpression(t, whileStmt.Condition, "x", "<", 10) {
        t.FailNow()
    }

    if len(whileStmt.Body.Statements) != 1 {
        t.Fatalf("Body length wrong. Expected=1. Got=%d", len(whileStmt.Body.Statements))
    }

    exprStmt, ok := whileStmt.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("body is not an expression statement. Got=%T", whileStmt.Body.Statements[0])
    }

    if exprStmt.Expression.String() != "x" {
        t.Fatalf(`body does not match. Expected=x. Got=%s`, exprStmt.Expression.String())
    }
}

func TestError(t *testing.T) {
    tests := []struct {
        inputStatement string
        outputError string
    } {
        {"x 1234;", "expected next token to be ;, got INT instead"},
        {"x = 1234 5;", "expected next token to be ;, got INT instead"},
        {" = 5;", "unexpected token: '='"},
    }

    for i, test := range tests {
        lexer := lexer.New(test.inputStatement)
        parser := New(lexer)
        program := parser.ParseProgram()

        errors := parser.GetErrors()

        if len(errors) == 0 {
            fmt.Println(program.Statements[1].String())
            t.Errorf("[%d]:expected parser to produce error\n", i)
            continue
        } else {
            if errors[0] != test.outputError {
                t.Errorf(`[%d]:expected message "%s", got "%s"\n`, i, test.outputError, errors[0])
                continue
            }
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"
    expected := "foobar"

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }

    expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    if !testIdentifier(t, expStmt.Expression, expected) {
        return
    }
}

func TestBooleanLiteralExpression(t *testing.T) {
    tests := []struct {
        input string
        expected bool
    } {
        {"true;", true},
        {"false;", false},
    }

    for _, test := range tests {
        lexer := lexer.New(test.input)
        parser := New(lexer)
        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
        }

        expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
        }

        if !testBooleanLiteral(t, expStmt.Expression, test.expected) {
            return 
        }
    }
}


func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }

    expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    intLiteral, ok := expStmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("exp is not *ast.IntegerLiteral. got=%T", expStmt)
    }

    if intLiteral.Value != 5 {
        t.Fatalf("intLiteral.value is not 5. got=%d", intLiteral.Value)
    }

    if intLiteral.TokenLiteral() != "5" {
        t.Fatalf("intLiteral.TokenLiteral() is not 5. got=%s", intLiteral.TokenLiteral())
    }
}

func TestFloatLiteralExpression(t *testing.T) {
    input := "5.4321;"

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
    }

    expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }

    floatLiteral, ok := expStmt.Expression.(*ast.FloatLiteral)
    if !ok {
        t.Fatalf("exp is not *ast.FloatLiteral. got=%T", expStmt)
    }

    if floatLiteral.Value != float64(5.4321) {
        t.Fatalf("floatLiteral.value is not 5.4321. got=%f", floatLiteral.Value)
    }

    if floatLiteral.TokenLiteral() != "5.4321" {
        t.Fatalf("intLiteral.TokenLiteral() is not 5.4321. got=%s", floatLiteral.TokenLiteral())
    }
}

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        input string
        operator string
        value any
    } {
        {"!5;", "!", 5},
        {"-15;", "-", 15},

        {"!true;", "!", true},
        {"!false;", "!", false},
    }

    for _, tt := range prefixTests {
        lexer := lexer.New(tt.input)
        parser := New(lexer)
        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                1, len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        exp, ok := stmt.Expression.(*ast.PrefixExpression)
        if !ok {
            t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
        }

        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s'. got=%s",
                tt.operator, exp.Operator)
        }

        if !testLiteralExpression(t, exp.Right, tt.value) {
            return
        }
    }
}

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        leftValue any
        operator string
        rightValue any
    }{
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},

        {"true == true;", true, "==", true},
        {"true != false;", true, "!=", false},
        {"false == false;", false, "==", false},
    }
    for _, tt := range infixTests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()

        checkParserErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
                1, len(program.Statements))
        }
        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }
        exp, ok := stmt.Expression.(*ast.InfixExpression)
        if !ok {
            t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
        }
        if !testLiteralExpression(t, exp.Left, tt.leftValue) {
            return
        }
        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s'. got=%s",
                tt.operator, exp.Operator)
        }
        if !testLiteralExpression(t, exp.Right, tt.rightValue) {
            return
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input string
        expected string
    }{
        {
            "-a * b;",
            "((-a) * b)",
        },
        {
            "!-a;",
            "(!(-a))",
        },
        {
            "a + b + c;",
            "((a + b) + c)",
        },
        {
            "a + b - c;",
            "((a + b) - c)",
        },
        {
            "a * b * c;",
            "((a * b) * c)",
        },
        {
            "a * b / c;",
            "((a * b) / c)",
        },
        {
            "a + b / c;",
            "(a + (b / c))",
        },
        {
            "a + b * c + d / e - f;",
            "(((a + (b * c)) + (d / e)) - f)",
        },
        {
            "3 + 4; -5 * 5;",
            "(3 + 4)((-5) * 5)",
        },
        {
            "5 > 4 == 3 < 4;",
            "((5 > 4) == (3 < 4))",
        },
        {
            "5 < 4 != 3 > 4;",
            "((5 < 4) != (3 > 4))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5;",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5;",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "true;",
            "true",
        },
        {
            "false;",
            "false",
        },
        {
            "3 > 5 == false;",
            "((3 > 5) == false)",
        },
        {
            "3 < 5 == true;",
            "((3 < 5) == true)",
        },
        {
            "1 + (2 + 3) + 4;",
            "((1 + (2 + 3)) + 4)",
        },
        {
            "(5 + 5) * 2;",
            "((5 + 5) * 2)",
        },
        {
            "2 / (5 + 5);",
            "(2 / (5 + 5))",
        },
        {
            "-(5 + 5);",
            "(-(5 + 5))",
        },
        {
            "!(true == true);",
            "(!(true == true))",
        },
        {
            "a + add(b * c) + d;",
            "((a + add((b * c))) + d)",
        },
        {
            "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8));",
            "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
        },
        {
            "add(a + b + c * d / f + g);",
            "add((((a + b) + ((c * d) / f)) + g))",
        },
        {
            "a * [1, 2, 3, 4][b * c] * d;",
            "((a * ([1, 2, 3, 4][(b * c)])) * d)",
        },
        {
            "add(a * b[2], b[1], 2 * [1, 2][1]);",
            "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
        },
    }

    for i, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()

        checkParserErrors(t, p)

        actual := program.String()
        if actual != tt.expected {
            t.Errorf("[%d]: expected=%q, got=%q", i, tt.expected, actual)
        }
    }
}

func TestParsingIndexExpressions(t *testing.T) {
    input := "myArray[1 + 1];"
    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()

    checkParserErrors(t, p)

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    indexExp, ok := stmt.Expression.(*ast.IndexExpression)
    if !ok {
        t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
    }
    if !testIdentifier(t, indexExp.Left, "myArray") {
        return
    }
    if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
        return
    }
}

func TestIfStatement(t *testing.T) {
    input := `if (x < y) { x; }`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }
    ifStmt, ok := program.Statements[0].(*ast.IfStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
    }

    if !testInfixExpression(t, ifStmt.Condition, "x", "<", "y") {
        return
    }

    if len(ifStmt.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(ifStmt.Consequence.Statements))
    }

    consequence, ok := ifStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", ifStmt.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    if ifStmt.Alternative != nil {
        t.Errorf("ifStmt.Alternative.Statements was not nil. got=%+v", ifStmt.Alternative)
    }
}

func TestIfElseExpression(t *testing.T) {
    input := `if (x < y) { x; } else { y; }`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }
    ifStmt, ok := program.Statements[0].(*ast.IfStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not *ast.IfStatement. got=%T", program.Statements[0])
    }

    if !testInfixExpression(t, ifStmt.Condition, "x", "<", "y") {
        return
    }

    if len(ifStmt.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(ifStmt.Consequence.Statements))
    }

    consequence, ok := ifStmt.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", ifStmt.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    alternative, ok := ifStmt.Alternative.Statements[0].(*ast.ExpressionStatement)
    if ifStmt.Alternative == nil {
        t.Fatalf("Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", ifStmt.Alternative.Statements[0])
    }

    if !testIdentifier(t, alternative.Expression, "y") {
        return
    }
}

func TestParseFunctionLiteral(t *testing.T) {
    input := `z = func(x, y) {x + y;};`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.AssignmentStatement. got=%T",
            program.Statements[0])
    }

    function, ok := stmt.Value.(*ast.FunctionLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
            stmt.Value)
    }

    if len(function.Parameters) != 2 {
        t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
            len(function.Parameters))
    }
    testLiteralExpression(t, function.Parameters[0], "x")
    testLiteralExpression(t, function.Parameters[1], "y")
    if len(function.Body.Statements) != 1 {
        t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
            len(function.Body.Statements))
    }
    bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
            function.Body.Statements[0])
    }
    testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestParseFunctionStatement(t *testing.T) {
    input := `func z(x, y) {x + y;}`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.FunctionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
            program.Statements[0])
    }

    function := stmt.FunctionLiteral

    if len(function.Parameters) != 2 {
        t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
            len(function.Parameters))
    }
    testLiteralExpression(t, function.Parameters[0], "x")
    testLiteralExpression(t, function.Parameters[1], "y")
    if len(function.Body.Statements) != 1 {
        t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
            len(function.Body.Statements))
    }
    bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
            function.Body.Statements[0])
    }
    testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestParseFunctionParameters(t *testing.T) {
    tests := []struct {
        input string
        expected []string
    } {
        {input: "func() {};", expected: []string{}},
        {input: "func(x) {};", expected: []string{"x"}},
        {input: "func(x, y, z) {};", expected: []string{"x", "y", "z"}},
    }

    for _, test := range tests {
        lexer := lexer.New(test.input)
        parser := New(lexer)
        program := parser.ParseProgram()

        checkParserErrors(t, parser)

        expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("statement not expression statement\n")
        }

        fnLit, ok := expStmt.Expression.(*ast.FunctionLiteral)
        if !ok {
            t.Fatalf("Statement.Expression is not of type functionLiteral. It is of type=%T\n", expStmt.Expression)
        }

        if len(fnLit.Parameters) != len(test.expected) {
            t.Fatalf("Expected %d parameters. Got=%d\n", len(test.expected), len(fnLit.Parameters))
        }

        for i, ident := range fnLit.Parameters {
            if !testIdentifier(t, ident, test.expected[i]) {
                return
            }
        }
    }
}

func TestParseCallExpression(t *testing.T) {
    input := "add(1, 2 * 3, 4 + 5);"

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 1 {
        t.Fatalf("length of statements not 1. Got %d\n", len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.statements[0] is not of type *ast.ExpressionStatement. Got=%T\n",
            program.Statements[0])
    }

    callExp, ok := stmt.Expression.(*ast.CallExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not of type *ast.CallExpression. Got=%T\n",
            stmt.Expression)
    }

    if !testIdentifier(t, callExp.FunctionLiteral, "add") {
        return
    }

    if len(callExp.Arguments) != 3 {
        t.Fatalf("wrong length of arguments. got=%d", len(callExp.Arguments))
    }

    // input := "add(1, 2 * 3, 4 + 5);"
    testIntegerLiteral(t, callExp.Arguments[0], 1)
    testInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
    testInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
    input := `"hello world";`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    stmt := program.Statements[0].(*ast.ExpressionStatement)
    literal, ok := stmt.Expression.(*ast.StringLiteral)
    if !ok {
        t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
    }
    if literal.Value != "hello world" {
        t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
    }
}

func TestArrayLiteral(t *testing.T) {
    input := "[1, 2 * 2, 3 + 3];"

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    array, ok := stmt.Expression.(*ast.ArrayLiteral)

    if !ok {
        t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
    }
    if len(array.Elements) != 3 {
        t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
    }
    testIntegerLiteral(t, array.Elements[0], 1)
    testInfixExpression(t, array.Elements[1], 2, "*", 2)
    testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingHashMapLiteralsStringKeys(t *testing.T) {
    input := `{"one": 1, "two": 2, "three": 3};`

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()

    checkParserErrors(t, p)

    stmt := program.Statements[0].(*ast.ExpressionStatement)
    hash, ok := stmt.Expression.(*ast.HashMapLiteral)
    if !ok {
        t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
    }
    if len(hash.Map) != 3 {
        t.Errorf("hash.Map has wrong length. got=%d", len(hash.Map))
    }

    expected := map[string]int64{
        "one": 1,
        "two": 2,
        "three": 3,
    }
    for key, value := range hash.Map {
        literal, ok := key.(*ast.StringLiteral)
        if !ok {
            t.Errorf("key is not ast.StringLiteral. got=%T", key)
        }
        expectedValue := expected[literal.String()]
        testIntegerLiteral(t, value, expectedValue)
    }
}

func TestParsingHashMapLiteralsStringKeysInvalid(t *testing.T) {
    input := `{"one": 1, "two": 2, "three": }`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    if len(parser.errors) == 0 {
        t.Fatalf("Expected Error. Got=%T", program.Statements[0])
    }
}

// helper functions 
func checkParserErrors(t *testing.T, p *Parser) {
    errors := p.GetErrors()
    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))
    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integ, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
        return false
    }
    if integ.Value != value {
        t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
        return false
    }
    if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
        t.Errorf("integ.TokenLiteral not %d. got=%s", value,
            integ.TokenLiteral())
        return false
    }
    return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
    ident, ok := exp.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp is not *ast.Identifier. got=%T", exp)
        return false
    }

    if ident.Value != value {
        t.Fatalf("identifier value is not %s. got=%s", value, ident.Value)
        return false
    }

    if ident.TokenLiteral() != value {
        t.Fatalf("identifier.TokenLiteral() is not %s. got=%s", value, ident.TokenLiteral())
        return false
    }

    return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
    boolLit, ok := exp.(*ast.BooleanLiteral)
    if !ok {
        t.Fatalf("exp is not *ast.BooleanLiteral. got=%T", exp)
        return false
    }

    if boolLit.Value != value {
        t.Fatalf("boolean value is not %t. got=%t", value, boolLit.Value)
        return false
    }

    return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
    switch v := expected.(type) {
    case int:
        return testIntegerLiteral(t, exp, int64(v))
    case int64:
        return testIntegerLiteral(t, exp,v)
    case string:
        return testIdentifier(t, exp, v)
    case bool:
        return testBooleanLiteral(t, exp, v)
}
    t.Errorf("type of exp not handled. Got=%T", exp)
    return false
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
    assignStmt, ok := stmt.(*ast.AssignmentStatement)
    if !ok {
        t.Errorf("stmt not of type *ast.AssignmentStatement. got=%q", stmt.TokenLiteral())
        return false
    }

    if assignStmt.Identifier.Value != name {
        t.Errorf("assignStmt.Identifier.value not %s. got=%s", name, assignStmt.Identifier.Value)
        return false
    }

    if assignStmt.Identifier.TokenLiteral() != name {
        t.Errorf("stmt.Identifier not %s. got=%s", name, assignStmt.Identifier)
        return false
    }

    return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
    opExp, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
        return false
    }

    if !testLiteralExpression(t, opExp.Left, left) {
        return false
    }

    if opExp.Operator != operator {
        t.Errorf("exp.Operator is not %s. got=%q", operator, opExp.Operator)
        return false
    }

    if !testLiteralExpression(t, opExp.Right, right) {
        return false
    }
    return true
}

func stringToExpression(input string) ast.Expression {
    lexer := lexer.New(input)
    parser := New(lexer)
    exp := parser.parseExpression(LOWEST)

    return exp
}

