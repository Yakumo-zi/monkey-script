package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"
	BANG   = "BANG"
	MINUS  = "MINUS"
	STAR   = "STAR"
	SLASH  = "SLASH"

	LT = "LT"
	GT = "GT"
	EQ = "EQ"
	NE = "NE"
	LE = "LE"
	GE = "GE"

	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"
	LPAREN    = "LPAREN"
	RPAREN    = "RPAREN"
	LBRACE    = "LBRACE"
	RBRACE    = "RBRACE"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

func NewToken(ty TokenType, val string) Token {
	return Token{
		Type:    ty,
		Literal: val,
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
