package lexer

import (
	"charm/token"
)

// position points to the character in the input that corresponds to ch rune
// readPostion refers to next character to read
type Lexer struct {
    input []rune
    position int
    readPosition int
    ch rune
}

func New(input string) *Lexer {
    lexer := &Lexer{input: []rune(input)}
    lexer.readChar()

    return lexer
}

func (lexer *Lexer) readChar() {
    if lexer.readPosition >= len(lexer.input) {
        lexer.ch = 0
    } else {
        lexer.ch = lexer.input[lexer.readPosition]
    }

    lexer.position = lexer.readPosition
    lexer.readPosition += 1
}

func (lexer *Lexer) peekChar() rune {
    if lexer.readPosition >= len(lexer.input) {
        return 0
    } else {
        return lexer.input[lexer.readPosition]
    }
}

func (lexer *Lexer) NextToken() token.Token {
    var tok token.Token

    lexer.skipWhitespace()

    switch lexer.ch {
        case '=':
            if lexer.peekChar() == '=' {
                lexer.readChar()
                tok.Type = token.EQ
                tok.Literal = "=="
            } else {
                tok = newToken(token.ASSIGN, lexer.ch)
            }
        case ';':
            tok = newToken(token.SEMICOLON, lexer.ch)
        case '(':
            tok = newToken(token.LPAREN, lexer.ch)
        case ')':
            tok = newToken(token.RPAREN, lexer.ch)
        case '{':
            tok = newToken(token.LBRACE, lexer.ch)
        case '}':
            tok = newToken(token.RBRACE, lexer.ch)
        case '[':
            tok = newToken(token.LBRACKET, lexer.ch)
        case ']':
            tok = newToken(token.RBRACKET, lexer.ch)
        case ',':
            tok = newToken(token.COMMA, lexer.ch)
        case '+':
            tok = newToken(token.PLUS, lexer.ch)
        case '-':
            tok = newToken(token.MINUS, lexer.ch)
        case '/':
            tok = newToken(token.SLASH, lexer.ch)
        case '*':
            tok = newToken(token.ASTERISK, lexer.ch)
        case '!':
            if lexer.peekChar() == '=' {
                lexer.readChar()
                tok = token.Token{Type: token.NOT_EQ, Literal:  "!="}
            } else {
                tok = newToken(token.BANG, lexer.ch)
            }
        case '<':
            if lexer.peekChar() == '=' {
                lexer.readChar()
                tok = token.Token{Type: token.LT_EQ, Literal: "<="}
            } else {
                tok = newToken(token.LT, lexer.ch)
            }
        case '>':
            if lexer.peekChar() == '=' {
                lexer.readChar()
                tok = token.Token{Type: token.GT_EQ, Literal: ">="}
            } else {
                tok = newToken(token.GT, lexer.ch)
            }
        case '"':
            stringLiteral := lexer.readString()
            tok = token.Token{Type: token.STRING, Literal: stringLiteral }
            lexer.readChar()
            return tok
        case 0:
            tok.Literal = ""
            tok.Type = token.EOF
        default:
            if isLetter(lexer.ch) {
                tok.Literal = lexer.readIdentifier()
                tok.Type = token.LookupIdentifier(tok.Literal)
                return tok
            } else if isDigit(lexer.ch) {
                tok.Literal = lexer.readNumber()
                tok.Type = token.INT
                return tok
            }else {
                tok = newToken(token.ILLEGAL, lexer.ch)
            }
    }

    lexer.readChar()
    return tok
}

func (lexer *Lexer) skipWhitespace() {
    for lexer.ch == ' ' || lexer.ch == '\t' || lexer.ch == '\n' || lexer.ch == '\r' {
        lexer.readChar()
    }
}

func isLetter(ch rune) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
    return '0' <= ch && ch <= '9'
}

// does it support alphanumeric identifier?
func (lexer *Lexer) readIdentifier() string {
    position := lexer.position

    for isLetter(lexer.ch) {
        lexer.readChar()
    }

    return string(lexer.input[position: lexer.position])
}

// Challenge: add support for escape characters such as \t \n
func (lexer *Lexer) readString() string {
    lexer.readChar()
    startPosition := lexer.position

    for  lexer.ch != 0 && lexer.ch != '"' {
        lexer.readChar()
    }
    return string(lexer.input[startPosition: lexer.position])
}

// TODO: support floats
func (lexer *Lexer) readNumber() string {
    position := lexer.position

    for isDigit(lexer.ch) {
        lexer.readChar()
    }

    return string(lexer.input[position: lexer.position])
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
    return token.Token{Type: tokenType, Literal: string(ch)}
}
