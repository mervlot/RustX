package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // line number
	col          int  // column number (1-based, tracking accurately)
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
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

	// Update line and col
	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
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
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
		return true
	}
	if l.ch == '/' && l.peekChar() == '*' {
		// multiline comment
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
			l.readChar()
		}
		return true
	}
	return false
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		l.skipWhitespace()
		if l.skipComment() {
			continue
		}
		break
	}
}

func (l *Lexer) readString() string {
	var out []byte
	l.readChar() // advance past opening quote
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '"' {
			l.readChar() // advance past the closing '"'
			break
		}
		if l.ch == '\\' {
			next := l.peekChar()
			switch next {
			case 'n':
				out = append(out, '\n')
				l.readChar()
			case 't':
				out = append(out, '\t')
				l.readChar()
			case 'r':
				out = append(out, '\r')
				l.readChar()
			case '"':
				out = append(out, '"')
				l.readChar()
			case '\\':
				out = append(out, '\\')
				l.readChar()
			default:
				out = append(out, '\\')
			}
		} else {
			out = append(out, l.ch)
		}
		l.readChar()
	}
	return string(out)
}

func (l *Lexer) readNumber() (token.TokenType, string) {
	startPos := l.position
	isFloat := false
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' {
		isFloat = true
		l.readChar() // eat '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	literal := l.input[startPos:l.position]
	if isFloat {
		return token.FLOAT, literal
	}
	return token.INT, literal
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaceAndComments()

	tok.Line = l.line
	tok.Col = l.col

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: tok.Line, Col: tok.Col}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
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
			typ, lit := l.readNumber()
			tok.Type = typ
			tok.Literal = lit
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: tok.Line, Col: tok.Col}
		}
	}

	l.readChar()
	return tok
}
