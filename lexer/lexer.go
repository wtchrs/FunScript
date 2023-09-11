package lexer

import (
	"errors"
	"funscript/token"
	"strings"
)

var escapeMap = map[string]byte{
	"\\n":  '\n',
	"\\r":  '\r',
	"\\t":  '\t',
	"\\v":  '\v',
	"\\\"": '"',
	"\\'":  '\'',
	"\\\\": '\\',
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			literal := l.readTwoCharToken()
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = token.New(token.ASSIGN, l.ch)
		}
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case '-':
		tok = token.New(token.MINUS, l.ch)
	case '*':
		tok = token.New(token.ASTERISK, l.ch)
	case '/':
		tok = token.New(token.SLASH, l.ch)
	case '!':
		if l.peekChar() == '=' {
			literal := l.readTwoCharToken()
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = token.New(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			literal := l.readTwoCharToken()
			tok = token.Token{Type: token.LTE, Literal: literal}
		} else {
			tok = token.New(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			literal := l.readTwoCharToken()
			tok = token.Token{Type: token.GTE, Literal: literal}
		} else {
			tok = token.New(token.GT, l.ch)
		}
	case '{':
		tok = token.New(token.LBRACE, l.ch)
	case '}':
		tok = token.New(token.RBRACE, l.ch)
	case '[':
		tok = token.New(token.LBRACKET, l.ch)
	case ']':
		tok = token.New(token.RBRACKET, l.ch)
	case ':':
		tok = token.New(token.COLON, l.ch)
	case '"':
		tok = l.readStringToken()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
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

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readStringToken() token.Token {
	position := l.position + 1
	escape := false
	for {
		l.readChar()
		if !escape && l.ch == '"' {
			break
		}
		if l.ch == 0 || l.ch == '\n' || l.ch == '\r' {
			return token.Token{Type: token.ILLEGAL, Literal: l.input[position-1 : l.position+1]}
		}
		escape = !escape && l.ch == '\\'
	}
	s, err := unescapeString(l.input[position:l.position])
	if err != nil {
		return token.Token{Type: token.ILLEGAL, Literal: l.input[position-1 : l.position+1]}
	}
	return token.Token{Type: token.STRING, Literal: s}
}

func (l *Lexer) readTwoCharToken() string {
	ch := l.ch
	l.readChar()
	return string(ch) + string(l.ch)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func unescapeString(input string) (string, error) {
	var out strings.Builder
	for i := 0; i < len(input); i++ {
		if input[i] == '\\' {
			if i+1 < len(input) {
				if replacement, ok := escapeMap[input[i:i+2]]; ok {
					out.WriteByte(replacement)
					i++
				} else {
					return out.String(), errors.New("Unknown escape sequence: " + string(replacement))
				}
			} else {
				return out.String(), errors.New("Incomplete escape sequence: " + input[i:])
			}
		} else {
			out.WriteByte(input[i])
		}
	}
	return out.String(), nil
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
