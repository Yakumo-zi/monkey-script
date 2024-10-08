package lexer

import (
	"interpreter/token"
)

type Lexer struct {
	input        string
	position     int // current position
	readPosition int // next position
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
func (l *Lexer) readIdentifer() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}
func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	start := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[start:l.position]
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()
	switch l.ch {
	case '+':
		tok = token.NewToken(token.PLUS, string(l.ch))
	case '-':
		tok = token.NewToken(token.MINUS, string(l.ch))
	case '*':
		tok = token.NewToken(token.STAR, string(l.ch))
	case '/':
		tok = token.NewToken(token.SLASH, string(l.ch))
	case '(':
		tok = token.NewToken(token.LPAREN, string(l.ch))
	case ')':
		tok = token.NewToken(token.RPAREN, string(l.ch))
	case '{':
		tok = token.NewToken(token.LBRACE, string(l.ch))
	case '}':
		tok = token.NewToken(token.RBRACE, string(l.ch))
	case ';':
		tok = token.NewToken(token.SEMICOLON, string(l.ch))
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewToken(token.EQ, "==")
		} else {

			tok = token.NewToken(token.ASSIGN, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewToken(token.LE, "<=")
		} else {

			tok = token.NewToken(token.LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewToken(token.GE, ">=")
		} else {

			tok = token.NewToken(token.GT, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.NewToken(token.NE, "!=")
		} else {

			tok = token.NewToken(token.BANG, string(l.ch))
		}
	case ',':
		tok = token.NewToken(token.COMMA, string(l.ch))
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = token.NewToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = token.NewToken(token.RBRACKET, string(l.ch))
	case ':':
		tok = token.NewToken(token.COLON, string(l.ch))
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifer()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.NewToken(token.ILLEGAL, string(l.ch))
		}

	}
	l.readChar()
	return tok
}
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
