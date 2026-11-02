package lexer

import (
	"bytes"
	"rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

func (l *Lexer) Input() string {
	return l.input
}

func (l *Lexer) Position() int {
	return l.position
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespaceAndComments()

	var tok token.Token
	tokPos := token.Position{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.ASSIGN, l.ch, tokPos)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.PLUS_PLUS, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.PLUS, l.ch, tokPos)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MINUS_MINUS, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.MINUS, l.ch, tokPos)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.BANG, l.ch, tokPos)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, tokPos)
	case '/':
		tok = newToken(token.SLASH, l.ch, tokPos)
	case '%':
		tok = newToken(token.MODULO, l.ch, tokPos)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.LT, l.ch, tokPos)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.GT, l.ch, tokPos)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokPos)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokPos)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Position: tokPos}
		} else {
			tok = newToken(token.COLON, l.ch, tokPos)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, tokPos)
	case ',':
		tok = newToken(token.COMMA, l.ch, tokPos)
	case '(':
		tok = newToken(token.LPAREN, l.ch, tokPos)
	case ')':
		tok = newToken(token.RPAREN, l.ch, tokPos)
	case '{':
		tok = newToken(token.LBRACE, l.ch, tokPos)
	case '}':
		tok = newToken(token.RBRACE, l.ch, tokPos)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, tokPos)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, tokPos)
	case '"':
		str, err := l.readString()
		if err != nil {
			tok = token.Token{Type: token.ILLEGAL, Literal: str, Position: tokPos}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Position: tokPos}
		}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Position = tokPos
	default:
		if isLetter(l.ch) {
			tok.Position = tokPos
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Position = tokPos
			lit, tokType := l.readNumber()
			tok.Literal = lit
			tok.Type = tokType
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokPos)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, pos token.Position) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Position: pos}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		// Skip whitespace
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}

		// Check for comments
		if l.ch == '/' {
			if l.peekChar() == '/' {
				// Single line comment
				for l.ch != '\n' && l.ch != 0 {
					l.readChar()
				}
				continue
			} else if l.peekChar() == '*' {
				// Multi-line comment
				l.readChar() // skip '/'
				l.readChar() // skip '*'
				for {
					if l.ch == 0 {
						break
					}
					if l.ch == '\n' {
						l.line++
						l.column = 0
					}
					if l.ch == '*' && l.peekChar() == '/' {
						l.readChar() // skip '*'
						l.readChar() // skip '/'
						break
					}
					l.readChar()
				}
				continue
			}
		}

		break
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	pos := l.position
	isFloat := false
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	tokType := token.TokenType(token.INT)
	if isFloat {
		tokType = token.TokenType(token.FLOAT)
	}
	return l.input[pos:l.position], tokType
}

func (l *Lexer) readString() (string, error) {
	var buf bytes.Buffer
	l.readChar() // skip opening quote
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case 'r':
				buf.WriteByte('\r')
			case '"':
				buf.WriteByte('"')
			case '\\':
				buf.WriteByte('\\')
			default:
				buf.WriteByte('\\')
				buf.WriteByte(l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			buf.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // advance past closing quote
	}
	return buf.String(), nil
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
