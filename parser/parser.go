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
    INDEX // arr[index]
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
    token.LPAREN: CALL,
    token.LBRACKET: INDEX,
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
    parser.registerPrefix(token.STRING, parser.parseStringLiteral)
    parser.registerPrefix(token.IF, parser.parseIfExpression)
    parser.registerPrefix(token.FUNCTION, parser.parseFunctionLiteral)
    parser.registerPrefix(token.LBRACKET, parser.parseArrayLiteral)
    parser.registerPrefix(token.LBRACE, parser.parseHashMapLiteral)
    parser.registerPrefix(token.LPAREN, parser.parseGroupedExpression)

    parser.infixParseFns = make(map[token.TokenType]infixParseFn)
    parser.registerInfix(token.PLUS, parser.parseInfixExpression)
    parser.registerInfix(token.MINUS, parser.parseInfixExpression)
    parser.registerInfix(token.SLASH, parser.parseInfixExpression)
    parser.registerInfix(token.ASTERISK, parser.parseInfixExpression)
    parser.registerInfix(token.EQ, parser.parseInfixExpression)
    parser.registerInfix(token.NOT_EQ, parser.parseInfixExpression)
    parser.registerInfix(token.LT, parser.parseInfixExpression)
    parser.registerInfix(token.GT, parser.parseInfixExpression)
    parser.registerInfix(token.LPAREN, parser.parseCallExpression)
    parser.registerInfix(token.LBRACKET, parser.parseIndexExpression)

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

    parser.nextToken()
    stmt.Value = parser.parseExpression(LOWEST)

    if parser.peekToken.Type == token.SEMICOLON {
        parser.nextToken()
    }

    return stmt
}

func (parser *Parser) parseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token: parser.currToken}

    parser.nextToken()

    stmt.ReturnValue = parser.parseExpression(LOWEST)

    if parser.peekToken.Type == token.SEMICOLON {
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

func (parser *Parser) parseIfExpression() ast.Expression {
    exp := &ast.IfExpression{Token: parser.currToken}

    if !parser.expectPeek(token.LPAREN) {
        fmt.Println("left paren not found for if statement")
        return nil
    }

    parser.nextToken()
    exp.Condition = parser.parseExpression(LOWEST)

    if !parser.expectPeek(token.RPAREN) {
        return nil
    }

    if !parser.expectPeek(token.LBRACE) {
        return nil
    }

    exp.Consequence = parser.parseBlockStatement()

    if parser.peekToken.Type == token.ELSE {
        parser.nextToken()

        if !parser.expectPeek(token.LBRACE) {
            return nil
        }

        exp.Alternative = parser.parseBlockStatement()
    }
    return exp
}

func (parser *Parser) parseBlockStatement() *ast.BlockStatement {
    block := &ast.BlockStatement{Token: parser.currToken}
    block.Statements = []ast.Statement{}

    parser.nextToken()

    for parser.currToken.Type != token.RBRACE && parser.currToken.Type != token.EOF {
        stmt := parser.parseStatement()

        if stmt != nil {
            block.Statements = append(block.Statements, stmt)
        }
        parser.nextToken()
    }
    return block
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

func (parser *Parser) parseStringLiteral() ast.Expression {
    return &ast.StringLiteral{Token: parser.currToken, Value: parser.currToken.Literal}
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

func (parser *Parser) parseGroupedExpression() ast.Expression {
    parser.nextToken()

    exp := parser.parseExpression(LOWEST)

    if !parser.expectPeek(token.RPAREN) {
        return nil
    }

    return exp
}

func (parser *Parser) parseFunctionLiteral() ast.Expression {
    function := &ast.FunctionLiteral{Token: parser.currToken}

    if !parser.expectPeek(token.LPAREN) {
        return nil
    }

    function.Parameters = parser.parseFunctionParameters()

    if !parser.expectPeek(token.LBRACE) {
        return nil
    }

    function.Body = parser.parseBlockStatement()

    return function
}

func (parser *Parser) parseFunctionParameters() []*ast.Identifier {
    parameters := []*ast.Identifier{}

    if parser.peekToken.Type == token.RPAREN {
        parser.nextToken()
        return parameters
    }

    if !parser.expectPeek(token.IDENT) {
        return nil
    }

    identifier := &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
    parameters = append(parameters, identifier)

    for parser.peekToken.Type == token.COMMA {
        parser.nextToken()
        parser.nextToken()
        identifier := &ast.Identifier{Token: parser.currToken, Value: parser.currToken.Literal}
        parameters = append(parameters, identifier)
    }

    if !parser.expectPeek(token.RPAREN) {
        return nil
    }

    return parameters
}

func (parser *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    callExp := &ast.CallExpression{Token: parser.currToken, FunctionLiteral: function}
    callExp.Arguments = parser.parseCallArguments()

    return callExp
}

func (parser *Parser) parseArrayLiteral() ast.Expression {
    elements := []ast.Expression{}

    if parser.peekToken.Type == token.RBRACKET {
        parser.nextToken()
         return &ast.ArrayLiteral{Elements: elements}
    }

    parser.nextToken()
    var expr ast.Expression
    expr = parser.parseExpression(LOWEST)
    elements = append(elements, expr)

    for parser.peekToken.Type == token.COMMA {
        parser.nextToken()
        parser.nextToken()
        expr = parser.parseExpression(LOWEST)
        elements = append(elements, expr)
    }

    if !parser.expectPeek(token.RBRACKET) {
        return nil
    }

    return &ast.ArrayLiteral{Elements: elements}
}

func (parser *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if parser.peekToken.Type == token.RPAREN {
        parser.nextToken()
        return args
    }

    parser.nextToken()
    args = append(args, parser.parseExpression(LOWEST))

    for parser.peekToken.Type == token.COMMA {
        parser.nextToken()
        parser.nextToken()
        args = append(args, parser.parseExpression(LOWEST))
    }

    // TODO: return error instead of nil, then handle error in caller
    if !parser.expectPeek(token.RPAREN) {
        return nil
    }
    return args
}

func (parser *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
    indexExpr := &ast.IndexExpression{
        Token: parser.currToken,
        Left: left,
    }

    parser.nextToken()
    index := parser.parseExpression(LOWEST)

    if !parser.expectPeek(token.RBRACKET) {
        return nil
    }

    indexExpr.Index = index
    return indexExpr
}

// TODO: Better error reporting for incorrect grammar of hashmaps
func (parser *Parser) parseHashMapLiteral() ast.Expression {
    hashMap := &ast.HashMapLiteral{
        Token: parser.currToken,
        Map: make(map[ast.Expression]ast.Expression),
    }

    for parser.peekToken.Type != token.RBRACE {
        parser.nextToken()
        keyExpr := parser.parseExpression(LOWEST)

        if !parser.expectPeek(token.COLON) {
            return nil
        }

        parser.nextToken()
        valueExpr := parser.parseExpression(LOWEST)
        if valueExpr == nil {
            return nil
        }

        hashMap.Map[keyExpr] = valueExpr
        
        if parser.peekToken.Type != token.RBRACE && !parser.expectPeek(token.COMMA) {
            return nil
        }
    }

    if !parser.expectPeek(token.RBRACE) {
        return nil
    }

    return hashMap
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
