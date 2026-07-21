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

	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '"':
		str, err := l.readString()
		if err != nil {
			tok = token.Token{Type: token.ILLEGAL, Literal: str, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Line: startLine, Column: startColumn}
			return tok
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startColumn}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			tok = token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startColumn}
			return tok
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			var tokType token.TokenType = token.INT
			if isFloat {
				tokType = token.FLOAT
			}
			tok = token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startColumn}
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startColumn}
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

		// Single line comment
		if l.ch == '/' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Multiline comment
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar()
			l.readChar()
			for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
				if l.ch == '\n' {
					l.line++
					l.column = 0
				}
				l.readChar()
			}
			if l.ch == '*' {
				l.readChar() // eat '*'
				l.readChar() // eat '/'
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

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() (string, bool) {
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
	return l.input[position:l.position], isFloat
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readString() (string, error) {
	l.readChar() // consume '"'
	var result []byte
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			next := l.peekChar()
			switch next {
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
		l.readChar() // advance past closing quote as per memory/NextToken specs
	}
	return string(result), nil
}
