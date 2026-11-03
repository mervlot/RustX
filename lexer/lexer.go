package lexer

import (
	"bytes"
	"rustx/token"
)

// Lexer converts source code into a stream of tokens.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number (1-based)
	column       int  // current column number (1-based)
}

// New creates a new Lexer instance.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// Input returns the raw source input string.
func (l *Lexer) Input() string {
	return l.input
}

// Position returns the current byte offset.
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

// NextToken scans and returns the next token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaceAndComments()

	pos := token.Position{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Position: pos}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Position: pos}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.PLUS_PLUS, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Position: pos}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.MINUS_MINUS, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Position: pos}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Position: pos}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Position: pos}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Position: pos}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Position: pos}
	case '%':
		tok = token.Token{Type: token.PERCENT, Literal: string(l.ch), Position: pos}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Position: pos}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Position: pos}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Position: pos}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Position: pos}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Position: pos}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Position: pos}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Position: pos}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Position: pos}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Position: pos}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Position: pos}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Position: pos}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Position: pos}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Position: pos}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Position: pos}
	case '"':
		str := l.readString()
		return token.Token{Type: token.STRING, Literal: str, Position: pos}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Position: pos}
	default:
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			return token.Token{Type: token.LookupIdent(ident), Literal: ident, Position: pos}
		} else if isDigit(l.ch) {
			num, tokType := l.readNumber()
			return token.Token{Type: tokType, Literal: num, Position: pos}
		}
		tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Position: pos}
	}

	l.readChar()
	return tok
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

		// Single line comment: //
		if l.ch == '/' && l.peekChar() == '/' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Multi-line comment: /* ... */
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar() // consume '/'
			l.readChar() // consume '*'
			for {
				if l.ch == 0 {
					break
				}
				if l.ch == '\n' {
					l.line++
					l.column = 0
				}
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
					break
				}
				l.readChar()
			}
			continue
		}

		break
	}
}

func (l *Lexer) readString() string {
	var out bytes.Buffer
	l.readChar() // skip opening quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out.WriteByte('\n')
			case 't':
				out.WriteByte('\t')
			case 'r':
				out.WriteByte('\r')
			case '\\':
				out.WriteByte('\\')
			case '"':
				out.WriteByte('"')
			default:
				out.WriteByte('\\')
				out.WriteByte(l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			out.WriteByte(l.ch)
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // advance past closing quote
	}

	return out.String()
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	startPos := l.position
	var tokType token.TokenType = token.INT

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		tokType = token.FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}

	return l.input[startPos:l.position], tokType
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
