package ast

import (
	"bytes"
	"charm/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
    String() string
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

func (program *Program) String() string {
    var out bytes.Buffer

    for _, s := range program.Statements {
        out.WriteString(s.String())
    }

    return out.String()
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
func (id *Identifier) String() string {
    return id.Value
}

type LetStatement struct {
	Token      token.Token
	Identifier *Identifier
	Value      Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) String() string {
    var out bytes.Buffer

    out.WriteString(ls.TokenLiteral() + " ")
    out.WriteString(ls.Identifier.String())
    out.WriteString(" = ")

    if ls.Value != nil {
        out.WriteString(ls.Value.String())
    }

    out.WriteString(";")
    return out.String()
}


type ReturnStatement struct {
	Token      token.Token
	ReturnValue      Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
    var out bytes.Buffer

    out.WriteString(rs.TokenLiteral() + " ")
    if rs.ReturnValue != nil {
        out.WriteString(rs.ReturnValue.String())
    }
    out.WriteString(";")
    return out.String()
}

// Expression statement are special statements. It accomodates expressions written as statements.
// For example: "3 + x;" is a valid statement.
type ExpressionStatement struct {
    Token token.Token
    Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string {
    if es.Expression != nil {
        return es.Expression.String()
    }
    return ""
}

type IntegerLiteral struct {
    Token token.Token
    Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
    return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
    return il.Token.Literal
}

type PrefixExpression struct {
    Token token.Token
    Operator string
    Right Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
    return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(pe.Operator)
    out.WriteString(pe.Right.String())
    out.WriteString(")")

    return out.String()
}

type InfixExpression struct {
    Token token.Token
    Left Expression
    Operator string
    Right Expression
}
func (oe *InfixExpression) expressionNode() {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(oe.Left.String())
    out.WriteString(" " + oe.Operator + " ")
    out.WriteString(oe.Right.String())
    out.WriteString(")")

    return out.String()
}

type BooleanLiteral struct {
    Token token.Token
    Value bool
}
func (b *BooleanLiteral) expressionNode() {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) String() string { return b.Token.Literal }

type StringLiteral struct {
    Token token.Token
    Value string
}
func (s *StringLiteral) expressionNode() {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string { return s.Token.Literal }

type BlockStatement struct {
    Token token.Token
    Statements []Statement
}
func (bs *BlockStatement) statementNode() {}
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}
func (bs *BlockStatement) String() string {
    var out bytes.Buffer
    for _, s := range bs.Statements {
        out.WriteString(s.String())
    }
    return out.String()
}

type IfExpression struct {
    token.Token
    Condition Expression
    Consequence *BlockStatement
    Alternative *BlockStatement
}
func (ie *IfExpression) expressionNode() {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
    var out bytes.Buffer
    out.WriteString("if")
    out.WriteString(ie.Condition.String())
    out.WriteString(" ")
    out.WriteString(ie.Consequence.String())

    if ie.Alternative != nil {
        out.WriteString("else ")
    out.WriteString(ie.Alternative.String())
    }
    return out.String()
}

type FunctionLiteral struct {
    Token token.Token
    Parameters []*Identifier
    Body *BlockStatement
}
func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
    var out bytes.Buffer

    params := []string {}
    for _, p := range fl.Parameters {
        params = append(params, p.String())
    }

    out.WriteString(fl.Token.Literal)
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") ")
    out.WriteString(fl.Body.String())

    return out.String()
}

type CallExpression struct {
    Token token.Token
    FunctionLiteral Expression
    Arguments []Expression
}
func (ce *CallExpression) expressionNode() {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
    var out bytes.Buffer

    args := []string{}
    for _, a := range ce.Arguments {
        args = append(args, a.String())
    }

    out.WriteString(ce.FunctionLiteral.String())
    out.WriteString("(")
    out.WriteString(strings.Join(args, ", "))
    out.WriteString(")")

    return out.String()
}

type ArrayLiteral struct {
    Token token.Token
    Elements []Expression
}

func (a *ArrayLiteral) expressionNode() {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
    var out bytes.Buffer

    elements := []string{}
    for _, exp := range a.Elements {
        elements = append(elements, exp.String())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")
    
    return out.String()
}

type IndexExpression struct {
    Token token.Token
    Left Expression
    Index Expression
}

func (i *IndexExpression) expressionNode() {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
    var out bytes.Buffer

    out.WriteString("(")
    out.WriteString(i.Left.String())
    out.WriteString("[")
    out.WriteString(i.Index.String())
    out.WriteString("])")

    return out.String()
}

type HashMapLiteral struct {
    Token token.Token
    Map map[Expression]Expression
}
func (hm *HashMapLiteral) expressionNode() {}
func (hm *HashMapLiteral) TokenLiteral() string { return hm.Token.Literal }
func (hm *HashMapLiteral) String() string {
    var out bytes.Buffer

    pairs := []string{}

    for key, value := range hm.Map {
        pairs = append(pairs, key.String() + ":" + value.String())
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}
