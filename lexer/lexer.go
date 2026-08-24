package lexer

import (
	"github.com/username/rustx/token"
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
	var tok token.Token

	l.skipWhitespaceAndComments()

	line := l.line
	col := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.ASSIGN, l.ch, line, col)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.COLON, l.ch, line, col)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.PLUS, l.ch, line, col)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.MINUS, l.ch, line, col)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.BANG, l.ch, line, col)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch, line, col)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, line, col)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.LT, l.ch, line, col)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.GT, l.ch, line, col)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, line, col)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, line, col)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, line, col)
	case ',':
		tok = newToken(token.COMMA, l.ch, line, col)
	case '(':
		tok = newToken(token.LPAREN, l.ch, line, col)
	case ')':
		tok = newToken(token.RPAREN, l.ch, line, col)
	case '{':
		tok = newToken(token.LBRACE, l.ch, line, col)
	case '}':
		tok = newToken(token.RBRACE, l.ch, line, col)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, line, col)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, line, col)
	case '"':
		strLit := l.readString()
		tok = token.Token{Type: token.STRING, Literal: strLit, Line: line, Column: col}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = line
		tok.Column = col
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			tokType := token.LookupIdent(lit)
			return token.Token{Type: tokType, Literal: lit, Line: line, Column: col}
		} else if isDigit(l.ch) {
			lit, tokType := l.readNumber()
			return token.Token{Type: tokType, Literal: lit, Line: line, Column: col}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, line, col)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}

		if l.ch == '/' && l.peekChar() == '/' {
			// Single line comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			if l.ch == '\n' {
				l.line++
				l.column = 0
				l.readChar()
			}
			continue
		}

		if l.ch == '/' && l.peekChar() == '*' {
			// Multiline comment
			l.readChar() // consume '/'
			l.readChar() // consume '*'
			for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
				if l.ch == '\n' {
					l.line++
					l.column = 0
				}
				l.readChar()
			}
			if l.ch == '*' && l.peekChar() == '/' {
				l.readChar() // consume '*'
				l.readChar() // consume '/'
			}
			continue
		}

		break
	}
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
	tokType := token.TokenType(token.INT)

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

	return l.input[position:l.position], tokType
}

func (l *Lexer) readString() string {
	var result []rune
	l.readChar() // skip initial quote

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			default:
				result = append(result, rune(l.ch))
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			result = append(result, rune(l.ch))
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // consume closing quote
	}
	return string(result)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}
