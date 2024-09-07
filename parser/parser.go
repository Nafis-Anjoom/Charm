package parser

import (
	"charm/ast"
	"charm/lexer"
	"charm/token"
	"fmt"
	"strconv"
)

type (
    prefixParseFn func() ast.Expression
    infixParseFn func(ast.Expression) ast.Expression
)

const (
    _ int = iota
    LOWEST
    EQUALS // ==
    LESSGREATER // > or <
    SUM // + or -
    PRODUCT // * or /
    PREFIX // -X or !X
    CALL // myFunction(X)
)

var precedences = map[token.TokenType]int {
    token.EQ: EQUALS,
    token.NOT_EQ: EQUALS,
    token.LT: LESSGREATER,
    token.GT: LESSGREATER,
    token.PLUS: SUM,
    token.MINUS: SUM,
    token.SLASH: PRODUCT,
    token.ASTERISK: PRODUCT,
}

type Parser struct {
    lexer *lexer.Lexer

    currToken token.Token
    peekToken token.Token

    errors []string

    prefixParseFns map[token.TokenType]prefixParseFn
    infixParseFns map[token.TokenType]infixParseFn
}

func New(lexer *lexer.Lexer) *Parser {
    parser := &Parser{
        lexer: lexer,
        errors: []string{},
    }

    parser.nextToken()
    parser.nextToken()

    parser.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    parser.registerPrefix(token.IDENT, parser.parseIdentifier)
    parser.registerPrefix(token.INT, parser.parseIntegerLiteral)
    parser.registerPrefix(token.BANG, parser.parsePrefixExpression)
    parser.registerPrefix(token.MINUS, parser.parsePrefixExpression)
    parser.registerPrefix(token.TRUE, parser.parseBooleanLiteral)
    parser.registerPrefix(token.FALSE, parser.parseBooleanLiteral)

    parser.infixParseFns = make(map[token.TokenType]infixParseFn)
    parser.registerInfix(token.PLUS, parser.parseInfixExpression)
    parser.registerInfix(token.MINUS, parser.parseInfixExpression)
    parser.registerInfix(token.SLASH, parser.parseInfixExpression)
    parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
    parser.registerInfix(token.EQ, parser.parseInfixExpression)
    parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
    parser.registerInfix(token.LT, parser.parseInfixExpression)
    parser.registerInfix(token.GT, parser.parseInfixExpression)

    return parser
}

func (parser *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    parser.prefixParseFns[tokenType] = fn
}


func (parser *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    parser.infixParseFns[tokenType] = fn
}

func (parser *Parser) Err(expectedType token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead",
        expectedType, parser.peekToken.Type)
    parser.errors = append(parser.errors, msg)
}

func (parser *Parser) GetErrors() []string {
    return parser.errors
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
        case token.RETURN:
            return parser.parseReturnStatement()
        default:
            return parser.parseExpressionStatement()
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

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: parser.currToken}

    parser.nextToken()

    // skipping expression for now
    for parser.currToken.Type != token.SEMICOLON {
        parser.nextToken()
    }

    return stmt
}

func (parser *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token: parser.currToken}

    stmt.Expression = parser.parseExpression(LOWEST)

    if parser.peekToken.Type == token.SEMICOLON {
        parser.nextToken()
    }

    return stmt
}

func (parser *Parser) parseExpression(precedence int) ast.Expression {
    prefix := parser.prefixParseFns[parser.currToken.Type]

    if prefix == nil {
        parser.noPrefixFnError(parser.currToken.Type)
        return nil
    }

    leftExp := prefix()

    for parser.peekToken.Type != token.SEMICOLON && precedence < parser.peekPrecedence() {
        infix := parser.infixParseFns[parser.peekToken.Type]

        if infix == nil {
            return leftExp
        }

        parser.nextToken()

        leftExp = infix(leftExp)
    }

    return leftExp
}

func (parser *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
}

func (parser *Parser) expectPeek(expectedType token.TokenType) bool {
    if parser.peekToken.Type != expectedType {
        parser.Err(expectedType)
        return false
    }
    parser.nextToken()
    return true
}

func (parser *Parser) parseIntegerLiteral() ast.Expression {
    intLit := &ast.IntegerLiteral{Token: parser.currToken}

    value, err := strconv.ParseInt(parser.currToken.Literal, 0, 64)
    if err != nil {
        // abstract away in a function to append to errors
        msg := fmt.Sprintf("could not parse %q as integer", parser.currToken.Literal)
        parser.errors = append(parser.errors, msg)

        return nil
    }

    intLit.Value = value
    return intLit
}

func (parser *Parser) parseBooleanLiteral() ast.Expression {
    boolLit := &ast.BooleanLiteral{Token: parser.currToken}

    if parser.currToken.Type == token.TRUE {
        boolLit.Value = true
    } else {
        boolLit.Value = false
    }

    return boolLit
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token: parser.currToken,
        Operator: parser.currToken.Literal,
    }

    parser.nextToken()
    expression.Right = parser.parseExpression(PREFIX)

    return expression
}

func (parser *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression {
        Token: parser.currToken,
        Operator: parser.currToken.Literal,
        Left: left,
    }

    precedence := parser.currPrecedence()
    parser.nextToken()
    expression.Right = parser.parseExpression(precedence)

    return expression
}

func (parser *Parser) noPrefixFnError(token token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function for %s found", token)
    parser.errors = append(parser.errors, msg)
}

func (parser *Parser) peekPrecedence() int {
    if p, ok := precedences[parser.peekToken.Type]; ok {
        return p
    }

    return LOWEST
}


func (parser *Parser) currPrecedence() int {
    if p, ok := precedences[parser.currToken.Type]; ok {
        return p
    }

    return LOWEST
}
