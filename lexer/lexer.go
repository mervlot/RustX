package lexer

import (
	"bytes"
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number
	column       int  // current column number
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

func (l *Lexer) Input() string {
	return l.input
}

func (l *Lexer) Position() int {
	return l.position
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
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.ASSIGN, l.ch, line, col)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.PLUS, l.ch, line, col)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.MINUS, l.ch, line, col)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.BANG, l.ch, line, col)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, line, col)
	case '/':
		tok = newToken(token.SLASH, l.ch, line, col)
	case '%':
		tok = newToken(token.PERCENT, l.ch, line, col)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.LT, l.ch, line, col)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.GT, l.ch, line, col)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, line, col)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, line, col)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: line, Column: col}
		} else {
			tok = newToken(token.COLON, l.ch, line, col)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, line, col)
	case ',':
		tok = newToken(token.COMMA, l.ch, line, col)
	case '(':
		tok = newToken(token.LPAREN, l.ch, line, col)
	case ')':
		tok = newToken(token.RPAREN, l.ch, line, col)
	case '{':
		tok = newToken(token.LBRACE, l.ch, line, col)
	case '}':
		tok = newToken(token.RBRACE, l.ch, line, col)
	case '[':
		tok = newToken(token.LBRACK, l.ch, line, col)
	case ']':
		tok = newToken(token.RBRACK, l.ch, line, col)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = col
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
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
			tok.Literal = lit
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Line = line
			tok.Column = col
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, line, col)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		l.skipWhitespace()
		if l.ch == '/' {
			if l.peekChar() == '/' {
				l.skipSingleLineComment()
				continue
			} else if l.peekChar() == '*' {
				l.skipMultiLineComment()
				continue
			}
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
	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip '*'
			l.readChar() // skip '/'
			break
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
	for isDigit(l.ch) || (l.ch == '.' && !isFloat && isDigit(l.peekChar())) {
		if l.ch == '.' {
			isFloat = true
		}
		l.readChar()
	}
	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	l.readChar() // Skip opening quote
	var out bytes.Buffer
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out.WriteByte('\n')
			case 't':
				out.WriteByte('\t')
			case 'r':
				out.WriteByte('\r')
			case '"':
				out.WriteByte('"')
			case '\\':
				out.WriteByte('\\')
			default:
				out.WriteByte('\\')
				out.WriteByte(l.ch)
			}
		} else {
			out.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // consume closing quote
	}
	return out.String()
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}
