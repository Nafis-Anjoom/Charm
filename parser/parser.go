package parser

import (
	"charm/ast"
	"charm/lexer"
	"charm/token"
)

type Parser struct {
    lexer *lexer.Lexer

    currToken token.Token
    peekToken token.Token
}

func New(lexer *lexer.Lexer) *Parser {
    parser := &Parser{lexer: lexer}

    parser.nextToken()
    parser.nextToken()

    return parser
}

func (parser *Parser) nextToken() {
    parser.currToken = parser.peekToken
    parser.peekToken = parser.lexer.NextToken()
}

func (parser *Parser) ParseProgram() *ast.Program {
    program := ast.Program{}
    program.Statements = []ast.Statement{}

    for parser.currToken.Type != token.EOF {
        stmt := parser.parseStatement()

        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }

        parser.nextToken()
    }
    return &program
}


func (parser *Parser) parseStatement() ast.Statement {
    switch parser.currToken.Type {
        case token.LET:
            return parser.parseLetStatement()
        default:
            return nil
    }
}

func (parser *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: parser.currToken}

    if !parser.expectPeek(token.IDENT) {
        return nil
    }

    stmt.Identifier = &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal } 

    if !parser.expectPeek(token.ASSIGN) {
        return nil
    }

    // skipping expression for now
    for parser.currToken.Type != token.SEMICOLON {
        parser.nextToken()
    }

    return stmt
}

func (parser *Parser) expectPeek(expectedType token.TokenType) bool {
    if parser.peekToken.Type != expectedType {
        return false
    }
    parser.nextToken()
    return true
}

