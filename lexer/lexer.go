package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number (1-based)
	column       int  // current column number (1-based)
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

	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ASSIGN, l.ch, startLine, startColumn)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.PLUS, l.ch, startLine, startColumn)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.MINUS, l.ch, startLine, startColumn)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.BANG, l.ch, startLine, startColumn)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, startLine, startColumn)
	case '/':
		tok = newToken(token.SLASH, l.ch, startLine, startColumn)
	case '%':
		tok = newToken(token.MOD, l.ch, startLine, startColumn)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.LT, l.ch, startLine, startColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.GT, l.ch, startLine, startColumn)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.COLON, l.ch, startLine, startColumn)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, startLine, startColumn)
	case ',':
		tok = newToken(token.COMMA, l.ch, startLine, startColumn)
	case '(':
		tok = newToken(token.LPAREN, l.ch, startLine, startColumn)
	case ')':
		tok = newToken(token.RPAREN, l.ch, startLine, startColumn)
	case '{':
		tok = newToken(token.LBRACE, l.ch, startLine, startColumn)
	case '}':
		tok = newToken(token.RBRACE, l.ch, startLine, startColumn)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, startLine, startColumn)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, startLine, startColumn)
	case '"':
		str, ok := l.readString()
		if !ok {
			tok = token.Token{Type: token.ILLEGAL, Literal: str, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Line: startLine, Column: startColumn}
		}
		return tok
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startColumn}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			return token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startColumn}
		} else if isDigit(l.ch) {
			literal, tokType := l.readNumber()
			return token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		skipped := false
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
			skipped = true
		}

		if l.ch == '/' {
			if l.peekChar() == '/' {
				l.readChar()
				l.readChar()
				for l.ch != '\n' && l.ch != 0 {
					l.readChar()
				}
				skipped = true
				continue
			} else if l.peekChar() == '*' {
				l.readChar()
				l.readChar()
				for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
					if l.ch == '\n' {
						l.line++
						l.column = 0
					}
					l.readChar()
				}
				if l.ch != 0 {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
				}
				skipped = true
				continue
			}
		}

		if !skipped {
			break
		}
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
	tokType := token.INT
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			if tokType == token.FLOAT { // second dot
				break
			}
			tokType = token.FLOAT
		}
		l.readChar()
	}
	return l.input[position:l.position], tokType
}

func (l *Lexer) readString() (string, bool) {
	var out []byte
	l.readChar() // skip opening quote
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out = append(out, '\n')
			case 't':
				out = append(out, '\t')
			case 'r':
				out = append(out, '\r')
			case '"':
				out = append(out, '"')
			case '\\':
				out = append(out, '\\')
			default:
				out = append(out, l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			out = append(out, l.ch)
		}
		l.readChar()
	}

	if l.ch != '"' {
		return "", false
	}
	l.readChar() // skip closing quote
	return string(out), true
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
