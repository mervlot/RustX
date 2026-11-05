package lexer

import (
	"rustx/token"
	"strings"
)

// Lexer transforms source code into a stream of tokens.
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

// Input returns the raw input string.
func (l *Lexer) Input() string {
	return l.input
}

// Position returns the current character position.
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

	curLine := l.line
	curCol := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ASSIGN, l.ch, curLine, curCol)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUS_PLUS, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.PLUS, l.ch, curLine, curCol)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MINUS_MINUS, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.MINUS, l.ch, curLine, curCol)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.BANG, l.ch, curLine, curCol)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch, curLine, curCol)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, curLine, curCol)
	case '%':
		tok = newToken(token.PERCENT, l.ch, curLine, curCol)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LE, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.LT, l.ch, curLine, curCol)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GE, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.GT, l.ch, curLine, curCol)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, curLine, curCol)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, curLine, curCol)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.COLON, l.ch, curLine, curCol)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, curLine, curCol)
	case ',':
		tok = newToken(token.COMMA, l.ch, curLine, curCol)
	case '(':
		tok = newToken(token.LPAREN, l.ch, curLine, curCol)
	case ')':
		tok = newToken(token.RPAREN, l.ch, curLine, curCol)
	case '{':
		tok = newToken(token.LBRACE, l.ch, curLine, curCol)
	case '}':
		tok = newToken(token.RBRACE, l.ch, curLine, curCol)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, curLine, curCol)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, curLine, curCol)
	case '"':
		str, ok := l.readString()
		if !ok {
			return token.Token{Type: token.ILLEGAL, Literal: "Unterminated string literal", Line: curLine, Column: curCol}
		}
		// Return directly since readString already advanced past the closing quote
		return token.Token{Type: token.STRING, Literal: str, Line: curLine, Column: curCol}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = curLine
		tok.Column = curCol
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tokType := token.LookupIdent(lit)
			return token.Token{Type: tokType, Literal: lit, Line: curLine, Column: curCol}
		} else if isDigit(l.ch) {
			numLit, numType := l.readNumber()
			return token.Token{Type: numType, Literal: numLit, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, curLine, curCol)
		}
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

		// Check for comments
		if l.ch == '/' && l.peekChar() == '/' {
			// Single-line comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		} else if l.ch == '/' && l.peekChar() == '*' {
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

		break
	}
}

func (l *Lexer) readString() (string, bool) {
	var sb strings.Builder
	l.readChar() // skip opening quote

	for l.ch != '"' {
		if l.ch == 0 {
			return "", false
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			sb.WriteByte(l.ch)
		}
		l.readChar()
	}
	l.readChar() // skip closing quote
	return sb.String(), true
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
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

	lit := l.input[position:l.position]
	if isFloat {
		return lit, token.FLOAT
	}
	return lit, token.INT
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, line int, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}
