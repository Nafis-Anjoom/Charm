package lexer

import (
	"charm/token"
	"testing"
)

func TestNextToken(t *testing.T) {
    input := `five = 5;
    ten = 10;
    add = fn(x, y) {
    x + y;
    };
    result = add(five, ten);
    !-/*5;
    5 < 10 > 5;
    if (5 < 10) {
    return true;
    } else {
    return false;
    }
    10 == 10;
    10 != 9;

    1 <= 10;
    10 >= 0;
    "hello"
    "hello, world!@#"
    gg[0]
    {"foo": "bar"}
    while (true)
    `   
    tests := []struct {
        expectedType token.TokenType
        expectedLiteral string
    }{
        {token.IDENT, "five"},
        {token.ASSIGN, "="},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
        {token.IDENT, "ten"},
        {token.ASSIGN, "="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
        {token.IDENT, "add"},
        {token.ASSIGN, "="},
        {token.FUNCTION, "fn"},
        {token.LPAREN, "("},
        {token.IDENT, "x"},
        {token.COMMA, ","},
        {token.IDENT, "y"},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        {token.IDENT, "x"},
        {token.PLUS, "+"},
        {token.IDENT, "y"},
        {token.SEMICOLON, ";"},
        {token.RBRACE, "}"},
        {token.SEMICOLON, ";"},
        {token.IDENT, "result"},
        {token.ASSIGN, "="},
        {token.IDENT, "add"},
        {token.LPAREN, "("},
        {token.IDENT, "five"},
        {token.COMMA, ","},
        {token.IDENT, "ten"},
        {token.RPAREN, ")"},
        {token.SEMICOLON, ";"},
        {token.BANG, "!"},
        {token.MINUS, "-"},
        {token.SLASH, "/"},
        {token.ASTERISK, "*"},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
        {token.INT, "5"},
        {token.LT, "<"},
        {token.INT, "10"},
        {token.GT, ">"},
        {token.INT, "5"},
        {token.SEMICOLON, ";"},
        {token.IF, "if"},
        {token.LPAREN, "("},
        {token.INT, "5"},
        {token.LT, "<"},
        {token.INT, "10"},
        {token.RPAREN, ")"},
        {token.LBRACE, "{"},
        {token.RETURN, "return"},
        {token.TRUE, "true"},
        {token.SEMICOLON, ";"},
        {token.RBRACE, "}"},
        {token.ELSE, "else"},
        {token.LBRACE, "{"},
        {token.RETURN, "return"},
        {token.FALSE, "false"},
        {token.SEMICOLON, ";"},
        {token.RBRACE, "}"},
        {token.INT, "10"},
        {token.EQ, "=="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
        {token.INT, "10"},
        {token.NOT_EQ, "!="},
        {token.INT, "9"},
        {token.SEMICOLON, ";"},
        {token.INT, "1"},
        {token.LT_EQ, "<="},
        {token.INT, "10"},
        {token.SEMICOLON, ";"},
        {token.INT, "10"},
        {token.GT_EQ, ">="},
        {token.INT, "0"},
        {token.SEMICOLON, ";"},
        {token.STRING, "hello"},
        {token.STRING, "hello, world!@#"},
        {token.IDENT, "gg"},
        {token.LBRACKET, "["},
        {token.INT, "0"},
        {token.RBRACKET, "]"},
        {token.LBRACE, "{"},
        {token.STRING, "foo"},
        {token.COLON, ":"},
        {token.STRING, "bar"},
        {token.RBRACE, "}"},
        {token.WHILE, "while"},
        {token.LPAREN, "("},
        {token.TRUE, "true"},
        {token.RPAREN, ")"},
        {token.EOF, ""},
    }

    lexer := New(input)

    for i, testCase := range tests {
        token := lexer.NextToken()

        if token.Type != testCase.expectedType {
            t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, testCase.expectedType, token.Type)
        }

        if token.Literal != testCase.expectedLiteral {
            t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, testCase.expectedLiteral, token.Literal)
        }

    }
}
