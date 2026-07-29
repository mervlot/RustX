package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // line tracking (1-based)
	column       int  // column tracking (1-based)
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

	line := l.line
	col := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: line, Column: col}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: line, Column: col}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: line, Column: col}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: line, Column: col}
		}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: line, Column: col}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: line, Column: col}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LT_EQ, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: line, Column: col}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GT_EQ, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: line, Column: col}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: col}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: col}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: line, Column: col}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: line, Column: col}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: line, Column: col}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: line, Column: col}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: line, Column: col}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: line, Column: col}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: line, Column: col}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: line, Column: col}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: line, Column: col}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: line, Column: col}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = col
		// Since readString advances past the closing quote, we return immediately to avoid call to readChar()
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		tok.Line = line
		tok.Column = col
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = line
			tok.Column = col
			return tok
		} else if isDigit(l.ch) {
			lit, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = lit
			tok.Line = line
			tok.Column = col
			return tok
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Column: col}
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
			l.readChar() // consume first '/'
			l.readChar() // consume second '/'
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Multi-line comment
		if l.ch == '/' && l.peekChar() == '*' {
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
	var out []byte
	l.readChar() // consume opening '"'
	for l.ch != '"' && l.ch != 0 {
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
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			out = append(out, l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // consume closing '"'
	}
	return string(out)
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
		l.readChar() // consume '.'
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
