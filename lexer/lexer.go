package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number (1-based, points to current char)
	file         string
}

func New(input string, file string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
		file:   file,
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

	tok.Line = l.line
	tok.Column = l.column
	tok.File = l.file

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ASSIGN
			tok.Literal = string(l.ch)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok.Type = token.INC
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.PLUS
			tok.Literal = string(l.ch)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok.Type = token.DEC
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.MINUS
			tok.Literal = string(l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.NOT_EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.BANG
			tok.Literal = string(l.ch)
		}
	case '*':
		tok.Type = token.ASTERISK
		tok.Literal = string(l.ch)
	case '/':
		tok.Type = token.SLASH
		tok.Literal = string(l.ch)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.LT_EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.LT
			tok.Literal = string(l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.GT_EQ
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.GT
			tok.Literal = string(l.ch)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok.Type = token.AND
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok.Type = token.OR
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Type = token.COLON_ASSIGN
			tok.Literal = string(ch) + string(l.ch)
		} else {
			tok.Type = token.COLON
			tok.Literal = string(l.ch)
		}
	case ';':
		tok.Type = token.SEMICOLON
		tok.Literal = string(l.ch)
	case ',':
		tok.Type = token.COMMA
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = token.LPAREN
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = token.RPAREN
		tok.Literal = string(l.ch)
	case '{':
		tok.Type = token.LBRACE
		tok.Literal = string(l.ch)
	case '}':
		tok.Type = token.RBRACE
		tok.Literal = string(l.ch)
	case '[':
		tok.Type = token.LBRACKET
		tok.Literal = string(l.ch)
	case ']':
		tok.Type = token.RBRACKET
		tok.Literal = string(l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		return tok
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			lit, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = lit
			return tok
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	}

	l.readChar()
	return tok
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
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) skipSingleLineComment() {
	// Skip "//"
	l.readChar()
	l.readChar()
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	// Skip "/*"
	l.readChar()
	l.readChar()
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // eat '*'
			l.readChar() // eat '/'
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	isFloat := false
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		isFloat = true
		l.readChar() // eat '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
	}
	lit := l.input[position:l.position]
	if l.ch == '"' {
		l.readChar() // advance past closing quote
	}
	return lit
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
