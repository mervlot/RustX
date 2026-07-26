package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // line number (1-based)
	column       int  // column number (1-based)
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
		l.column = 0
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

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() bool {
	if l.ch == '/' && l.peekChar() == '/' {
		// single line comment
		l.readChar() // consume first '/'
		l.readChar() // consume second '/'
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		return true
	}
	if l.ch == '/' && l.peekChar() == '*' {
		// multi-line comment
		l.readChar() // consume '/'
		l.readChar() // consume '*'
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
		return true
	}
	return false
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		l.skipWhitespace()
		if (l.ch == '/' && l.peekChar() == '/') || (l.ch == '/' && l.peekChar() == '*') {
			l.skipComment()
		} else {
			break
		}
	}
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	startPos := l.position
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
	tokType := token.INT
	if isFloat {
		tokType = token.FLOAT
	}
	return l.input[startPos:l.position], tokType
}

func (l *Lexer) readString() (string, token.TokenType) {
	l.readChar() // consume opening quote

	var chars []byte
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			next := l.peekChar()
			switch next {
			case 'n':
				chars = append(chars, '\n')
				l.readChar() // consume '\\'
				l.readChar() // consume 'n'
			case 't':
				chars = append(chars, '\t')
				l.readChar()
				l.readChar()
			case 'r':
				chars = append(chars, '\r')
				l.readChar()
				l.readChar()
			case '"':
				chars = append(chars, '"')
				l.readChar()
				l.readChar()
			case '\\':
				chars = append(chars, '\\')
				l.readChar()
				l.readChar()
			default:
				chars = append(chars, l.ch)
				l.readChar()
			}
		} else {
			chars = append(chars, l.ch)
			l.readChar()
		}
	}

	var tokType token.TokenType = token.STRING
	if l.ch == '"' {
		l.readChar() // consume closing quote
	} else {
		tokType = token.ILLEGAL
	}
	return string(chars), tokType
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaceAndComments()

	startLine := l.line
	startCol := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: startLine, Column: startCol}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: startLine, Column: startCol}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: startLine, Column: startCol}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: startLine, Column: startCol}
	case '"':
		str, tokType := l.readString()
		tok = token.Token{Type: tokType, Literal: str, Line: startLine, Column: startCol}
		return tok
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startCol}
	default:
		if isLetter(l.ch) {
			startPos := l.position
			for isLetter(l.ch) || isDigit(l.ch) {
				l.readChar()
			}
			literal := l.input[startPos:l.position]
			tok = token.Token{Type: token.LookupIdent(literal), Literal: literal, Line: startLine, Column: startCol}
			return tok
		} else if isDigit(l.ch) {
			literal, tokType := l.readNumber()
			tok = token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startCol}
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startCol}
		}
	}

	l.readChar()
	return tok
}
