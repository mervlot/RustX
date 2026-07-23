package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // current reading position (after current char)
	ch           byte // current char under examination
	line         int  // line of current char (l.ch)
	column       int  // column of current char (l.ch)
	nextLine     int  // line of next char
	nextColumn   int  // column of next char
}

func New(input string) *Lexer {
	l := &Lexer{
		input:      input,
		line:       1,
		column:     1,
		nextLine:   1,
		nextColumn: 1,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.line = l.nextLine
	l.column = l.nextColumn

	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.nextLine++
		l.nextColumn = 1
	} else {
		l.nextColumn++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
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
			l.readChar() // read '*'
			l.readChar() // read '/'
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readString() (string, int, int) {
	line := l.line
	column := l.column
	l.readChar() // skip '"'

	var runes []rune
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip '\'
			switch l.ch {
			case 'n':
				runes = append(runes, '\n')
			case 't':
				runes = append(runes, '\t')
			case 'r':
				runes = append(runes, '\r')
			case '"':
				runes = append(runes, '"')
			case '\\':
				runes = append(runes, '\\')
			default:
				runes = append(runes, '\\')
				runes = append(runes, rune(l.ch))
			}
		} else {
			runes = append(runes, rune(l.ch))
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // advance past closing quote
	}
	return string(runes), line, column
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() (string, int, int) {
	line := l.line
	column := l.column
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position], line, column
}

func (l *Lexer) readNumber() (token.TokenType, string, int, int) {
	line := l.line
	column := l.column
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

	var tokType token.TokenType = token.INT
	if isFloat {
		tokType = token.FLOAT
	}
	return tokType, l.input[position:l.position], line, column
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaceAndComments()

	line := l.line
	column := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: line, Column: column}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: line, Column: column}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: line, Column: column}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: line, Column: column}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: line, Column: column}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: line, Column: column}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: line, Column: column}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: line, Column: column}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: line, Column: column}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: line, Column: column}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: line, Column: column}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: line, Column: column}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: line, Column: column}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: line, Column: column}
	case '"':
		strLit, sLine, sCol := l.readString()
		return token.Token{Type: token.STRING, Literal: strLit, Line: sLine, Column: sCol}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: line, Column: column}
	default:
		if isLetter(l.ch) {
			ident, iLine, iCol := l.readIdentifier()
			return token.Token{Type: token.LookupIdent(ident), Literal: ident, Line: iLine, Column: iCol}
		} else if isDigit(l.ch) {
			tokType, num, nLine, nCol := l.readNumber()
			return token.Token{Type: tokType, Literal: num, Line: nLine, Column: nCol}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: column}
		}
	}

	l.readChar()
	return tok
}
