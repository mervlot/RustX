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

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		// readString already advanced past the ending double quote. Return immediately.
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = literal
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tok.Line, Column: tok.Column}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
			continue
		}

		if l.ch == '/' && l.peekChar() == '/' {
			// Single-line comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		if l.ch == '/' && l.peekChar() == '*' {
			// Multi-line comment
			l.readChar() // eat '/'
			l.readChar() // eat '*'
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
			continue
		}

		break
	}
}

func (l *Lexer) readString() string {
	l.readChar() // step over the opening double quote
	var result []byte
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			peek := l.peekChar()
			switch peek {
			case 'n':
				result = append(result, '\n')
				l.readChar()
			case 't':
				result = append(result, '\t')
				l.readChar()
			case 'r':
				result = append(result, '\r')
				l.readChar()
			case '"':
				result = append(result, '"')
				l.readChar()
			case '\\':
				result = append(result, '\\')
				l.readChar()
			default:
				result = append(result, '\\')
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			result = append(result, l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // step over the closing double quote
	}
	return string(result)
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

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
