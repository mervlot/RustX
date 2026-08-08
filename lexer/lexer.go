package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line (1-based)
	column       int  // current column (1-based)
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
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
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) Input() string {
	return l.input
}

func (l *Lexer) Position() int {
	return l.position
}

func (l *Lexer) Clone() *Lexer {
	return &Lexer{
		input:        l.input,
		position:     l.position,
		readPosition: l.readPosition,
		ch:           l.ch,
		line:         l.line,
		column:       l.column,
	}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		l.skipWhitespace()
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipSingleLineComment()
			continue
		}
		if l.ch == '/' && l.peekChar() == '*' {
			l.skipMultiLineComment()
			continue
		}
		break
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipSingleLineComment() {
	// Skip '//'
	l.readChar()
	l.readChar()
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	// Skip '/*'
	l.readChar()
	l.readChar()
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // consume '*'
			l.readChar() // consume '/'
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readString() (string, int, int) {
	tokLine := l.line
	tokColumn := l.column

	// Skip opening double quote
	l.readChar()

	var out []byte
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '"' {
			l.readChar() // advance past closing quote
			break
		}
		if l.ch == '\\' {
			l.readChar() // consume '\'
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
			l.readChar()
		} else {
			out = append(out, l.ch)
			l.readChar()
		}
	}
	return string(out), tokLine, tokColumn
}

func (l *Lexer) readNumber() (token.TokenType, string, int, int) {
	tokLine := l.line
	tokColumn := l.column
	startPos := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
		return token.FLOAT, l.input[startPos:l.position], tokLine, tokColumn
	}

	return token.INT, l.input[startPos:l.position], tokLine, tokColumn
}

func (l *Lexer) readIdentifier() (string, int, int) {
	tokLine := l.line
	tokColumn := l.column
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position], tokLine, tokColumn
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaceAndComments()

	tokLine := l.line
	tokColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
	case '"':
		lit, line, col := l.readString()
		return token.Token{Type: token.STRING, Literal: lit, Line: line, Column: col}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: tokLine, Column: tokColumn}
	default:
		if isLetter(l.ch) {
			lit, line, col := l.readIdentifier()
			return token.Token{Type: token.LookupIdent(lit), Literal: lit, Line: line, Column: col}
		} else if isDigit(l.ch) {
			typ, lit, line, col := l.readNumber()
			return token.Token{Type: typ, Literal: lit, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tokLine, Column: tokColumn}
		}
	}

	l.readChar()
	return tok
}
