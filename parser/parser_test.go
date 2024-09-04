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
        let foobar = 838383; `

    lexer := lexer.New(input)
    parser := New(lexer)
    program := parser.ParseProgram()

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

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
    if stmt.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
        return false
    }
    t.Log("token")

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
