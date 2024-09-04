package parser

import (
	"charm/ast"
	"charm/lexer"
	"testing"
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
    
    parserErrors := parser.GetErrors()
    if len(parserErrors) != 0 {
        t.Errorf("parser has %d errors\n", len(parserErrors))
        for _, msg := range parserErrors {
            t.Errorf("parser error: %q", msg)
        }
        t.FailNow()
    }

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
        t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
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
