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
	col          int  // current column number (1-based index)
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
		col:   0,
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
	l.col++
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

	// Capture the start of the current token's column
	colStart := l.col

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: l.line, Col: colStart}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: l.line, Col: colStart}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: l.line, Col: colStart}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: l.line, Col: colStart}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: l.line, Col: colStart}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: l.line, Col: colStart}
	case '"':
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Col = colStart
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		tok.Line = l.line
		tok.Col = colStart
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Col = colStart
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = l.line
			tok.Col = colStart
			literal, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = literal
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.line, Col: colStart}
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
				l.col = 0
			}
			l.readChar()
			continue
		}

		// Single line comment
		if l.ch == '/' && l.peekChar() == '/' {
			l.readChar() // eat '/'
			l.readChar() // eat '/'
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Multi line comment
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar() // eat '/'
			l.readChar() // eat '*'
			for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
				if l.ch == '\n' {
					l.line++
					l.col = 0
				}
				l.readChar()
			}
			if l.ch != 0 {
				l.readChar() // eat '*'
				l.readChar() // eat '/'
			}
			continue
		}

		break
	}
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
			l.col = 0
		}
	}
	return l.input[position:l.position]
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
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			isFloat = true
		}
		l.readChar()
	}
	return l.input[position:l.position], isFloat
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
