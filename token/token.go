package token

const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	IDENT     = "IDENT"
	INT       = "INT"
	STRING    = "STRING"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	BANG      = "!"
	LT        = "<"
	LTE       = "<="
	GT        = ">"
	GTE       = ">="
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"
	EQ        = "=="
	NOT_EQ    = "!="
	FUNCTION  = "FUNCTION"
	RETURN    = "RETURN"
	LET       = "LET"
	IF        = "IF"
	ELSE      = "ELSE"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
)

type TokenType string

// TODO: Add positional information to Token

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"return": RETURN,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
