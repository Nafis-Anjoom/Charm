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
        {input: "let x = 5;", expectedIdentifier: "x", expectedValue: 5},
        {input: "let y = 10;", expectedIdentifier: "y", expectedValue: 10},
        {input: "let foobar = 838383;", expectedIdentifier: "foobar", expectedValue: 838383},
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

        val := stmt.(*ast.LetStatement).Value
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

func TestError(t *testing.T) {
    tests := []struct {
        inputStatement string
        outputError string
    } {
        {"let x 1234;", "expected next token to be =, got INT instead"},
        {"let = 5;", "expected next token to be IDENT, got = instead"},
        {"let 838383;", "expected next token to be IDENT, got INT instead"},
    }

    for i, test := range tests {
        lexer := lexer.New(test.inputStatement)
        parser := New(lexer)
        parser.ParseProgram()

        errors := parser.GetErrors()

        if len(errors) == 0 {
            t.Fatalf("expected parser to produce error\n")
        } else {
            if errors[0] != test.outputError {
                t.Fatalf(`[%d]:expected message "%s", got "%s"\n`, i, test.outputError, errors[0])
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

func TestParsingPrefixExpressions(t *testing.T) {
    prefixTests := []struct {
        input string
        operator string
        value any
    } {
        {"!5", "!", 5},
        {"-15", "-", 15},

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

        {"true == true", true, "==", true},
        {"true != false", true, "!=", false},
        {"false == false", false, "==", false},
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
            "-a * b",
            "((-a) * b)",
        },
        {
            "!-a",
            "(!(-a))",
        },
        {
            "a + b + c",
            "((a + b) + c)",
        },
        {
            "a + b - c",
            "((a + b) - c)",
        },
        {
            "a * b * c",
            "((a * b) * c)",
        },
        {
            "a * b / c",
            "((a * b) / c)",
        },
        {
            "a + b / c",
            "(a + (b / c))",
        },
        {
            "a + b * c + d / e - f",
            "(((a + (b * c)) + (d / e)) - f)",
        },
        {
            "3 + 4; -5 * 5",
            "(3 + 4)((-5) * 5)",
        },
        {
            "5 > 4 == 3 < 4",
            "((5 > 4) == (3 < 4))",
        },
        {
            "5 < 4 != 3 > 4",
            "((5 < 4) != (3 > 4))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "true",
            "true",
        },
        {
            "false",
            "false",
        },
        {
            "3 > 5 == false",
            "((3 > 5) == false)",
        },
        {
            "3 < 5 == true",
            "((3 < 5) == true)",
        },
        {
            "1 + (2 + 3) + 4",
            "((1 + (2 + 3)) + 4)",
        },
        {
            "(5 + 5) * 2",
            "((5 + 5) * 2)",
        },
        {
            "2 / (5 + 5)",
            "(2 / (5 + 5))",
        },
        {
            "-(5 + 5)",
            "(-(5 + 5))",
        },
        {
            "!(true == true)",
            "(!(true == true))",
        },
        {
            "a + add(b * c) + d",
            "((a + add((b * c))) + d)",
        },
        {
            "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
            "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
        },
        {
            "add(a + b + c * d / f + g)",
            "add((((a + b) + ((c * d) / f)) + g))",
        },
    }

    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParserErrors(t, p)
        actual := program.String()
        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }
}


func TestIfExpression(t *testing.T) {
    input := `if (x < y) { x }`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }
    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
            stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    if exp.Alternative != nil {
        t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
    }
}

func TestIfElseExpression(t *testing.T) {
    input := `if (x < y) { x } else { y }`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }
    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
    }
    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
            stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n", len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Consequence.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
    if exp.Alternative == nil {
        t.Fatalf("Alternative.Statements[0] is not ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
    }

    if !testIdentifier(t, alternative.Expression, "y") {
        return
    }
}

func TestParseFunctionLiteral(t *testing.T) {
    input := `fn(x, y) {x + y;}`

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)
    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
            program.Statements[0])
    }

    function, ok := stmt.Expression.(*ast.FunctionLiteral)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
            stmt.Expression)
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

func TestParseFunctionParameters(t *testing.T) {
    tests := []struct {
        input string
        expected []string
    } {
        {input: "fn() {}", expected: []string{}},
        {input: "fn(x) {}", expected: []string{"x"}},
        {input: "fn(x, y, z) {}", expected: []string{"x", "y", "z"}},
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
    if stmt.TokenLiteral() != "let" {
        t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
        return false
    }

    letStmt, ok := stmt.(*ast.LetStatement)
    if !ok {
        t.Errorf("stmt not of type *ast.LetStatement. got=%q", stmt.TokenLiteral())
        return false
    }

    if letStmt.Identifier.Value != name {
        t.Errorf("letStmt.Identifier.value not %s. got=%s", name, letStmt.Identifier.Value)
        return false
    }

    if letStmt.Identifier.TokenLiteral() != name {
        t.Errorf("stmt.Identifier not %s. got=%s", name, letStmt.Identifier)
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

