package ast

import "charm/token"

type Node interface {
    TokenLiteral() string
}

type Statement interface {
    Node
    statementNode()
}

type Expression interface {
    Node
    expressionNode()
}

type Program struct {
    Statements []Statement
}

func (program *Program) TokenLiteral() string {
    if len(program.Statements) > 0 {
        return program.Statements[0].TokenLiteral()
    } else {
        return ""
    }
}

// implements Expression node because it returns value
type Identifier struct {
    Token token.Token
    Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
    return id.Token.Literal
}

type LetStatement struct {
    Token token.Token
    Identifier *Identifier
    Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
    return ls.Token.Literal
}
