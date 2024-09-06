package parser

import (
    "charm/ast"
    "charm/lexer"
    "testing"
    "fmt"
)

func TestLetStatements(t *testing.T) {
    input := `
    let x = 5;
    let y = 10;
    let foobar = 838383;
    `

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if program == nil {
        t.Fatal("ParseProgram returned nil\n")
    }

    if len(program.Statements) != 3 {
        t.Fatalf("Expected 3 statements. Got %d\n", len(program.Statements))
    }

    tests := []struct {
        expectedIdentifier string
    } {
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i, testCase := range tests {
        stmt := program.Statements[i]

        if !testLetStatement(t, stmt, testCase.expectedIdentifier) {
            return 
        }
    }
}

func TestReturnStatement(t *testing.T) {
    input := `
    return 5;
    return 10;
    return 923123;
    `

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

    checkParserErrors(t, parser)

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 statements. got=%d\n", len(program.Statements))
    }

    for _, stmt := range program.Statements {
        if stmt.TokenLiteral() != "return" {
            t.Errorf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
        }

        // returnStmt, ok := stmt.(*ast.ReturnStatement)
        _, ok := stmt.(*ast.ReturnStatement)

        if !ok {
            t.Errorf("stmt not *ast.ReturnStatement. got=%T\n", stmt)
            continue
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

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

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

    ident, ok := expStmt.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp is not *ast.Identifier. got=%T", expStmt)
    }

    if ident.Value != "foobar" {
        t.Fatalf("identifier value is not foobar. got=%s", ident.Value)
    }

    if ident.TokenLiteral() != "foobar" {
        t.Fatalf("identifier.TokenLiteral() is not foobar. got=%s", ident.TokenLiteral())
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
        integerValue int64
    } {
        {"!5", "!", 5},
        {"-15", "-", 15},
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

        if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
            return
        }
    }
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

func TestParsingInfixExpressions(t *testing.T) {
    infixTests := []struct {
        input string
        leftValue int64
        operator string
        rightValue int64
    }{
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
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
        if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
            return
        }
        if exp.Operator != tt.operator {
            t.Fatalf("exp.Operator is not '%s'. got=%s",
                tt.operator, exp.Operator)
        }
        if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
            return
        }
    }
}

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
