package lexer

import (
	"strings"

	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	col          int  // current column number
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

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INCREMENT, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DECREMENT, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			col := l.col
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: l.line, Col: col}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: l.line, Col: l.col}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: l.line, Col: l.col}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: l.line, Col: l.col}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: l.line, Col: l.col}
	case '"':
		tok.Type = token.STRING
		tok.Line = l.line
		tok.Col = l.col
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		tok.Line = l.line
		tok.Col = l.col
	default:
		if isLetter(l.ch) {
			tok.Line = l.line
			tok.Col = l.col
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = l.line
			tok.Col = l.col
			literal, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = literal
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: l.line, Col: l.col}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		// First skip normal whitespace
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.col = 0
			}
			l.readChar()
		}

		// Check for single line comments
		if l.ch == '/' && l.peekChar() == '/' {
			// Read until newline or EOF
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Check for multi line comments
		if l.ch == '/' && l.peekChar() == '*' {
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
					l.col = 0
				}
				l.readChar()
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

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	isFloat := false
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			// If it's a double dot or followed by something else, let's be careful.
			// But for now, we just flag it as float.
			isFloat = true
		}
		l.readChar()
	}
	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	l.readChar() // eat the opening double quote
	var sb strings.Builder
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // eat the backslash
			switch l.ch {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			default:
				sb.WriteByte('\\')
				sb.WriteByte(l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.col = 0
			}
			sb.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // eat the closing double quote
	}
	return sb.String()
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
