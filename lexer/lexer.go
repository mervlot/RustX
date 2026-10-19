package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
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

	tokLine := l.line
	tokCol := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.ASSIGN, l.ch, tokLine, tokCol)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.COLON, l.ch, tokLine, tokCol)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.PLUS, l.ch, tokLine, tokCol)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.MINUS, l.ch, tokLine, tokCol)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.BANG, l.ch, tokLine, tokCol)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch, tokLine, tokCol)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, tokLine, tokCol)
	case '%':
		tok = newToken(token.MODULO, l.ch, tokLine, tokCol)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.LT, l.ch, tokLine, tokCol)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.GT, l.ch, tokLine, tokCol)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokLine, tokCol)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: tokLine, Column: tokCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokLine, tokCol)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, tokLine, tokCol)
	case ',':
		tok = newToken(token.COMMA, l.ch, tokLine, tokCol)
	case '(':
		tok = newToken(token.LPAREN, l.ch, tokLine, tokCol)
	case ')':
		tok = newToken(token.RPAREN, l.ch, tokLine, tokCol)
	case '{':
		tok = newToken(token.LBRACE, l.ch, tokLine, tokCol)
	case '}':
		tok = newToken(token.RBRACE, l.ch, tokLine, tokCol)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, tokLine, tokCol)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, tokLine, tokCol)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = tokLine
		tok.Column = tokCol
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = tokLine
		tok.Column = tokCol
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = tokLine
			tok.Column = tokCol
			return tok
		} else if isDigit(l.ch) {
			tok.Type, tok.Literal = l.readNumber()
			tok.Line = tokLine
			tok.Column = tokCol
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokLine, tokCol)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}

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

func (l *Lexer) skipSingleLineComment() {
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
	l.readChar() // skip '/'
	l.readChar() // skip '*'

	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
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

func (l *Lexer) readNumber() (token.TokenType, string) {
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

	tokType := token.TokenType(token.INT)
	if isFloat {
		tokType = token.TokenType(token.FLOAT)
	}
	return tokType, l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var out []byte
	l.readChar() // consume opening "

	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}

		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				out = append(out, '\n')
			case 't':
				out = append(out, '\t')
			case 'r':
				out = append(out, '\r')
			case '"':
				out = append(out, '"')
			case '\\':
				out = append(out, '\\')
			default:
				out = append(out, '\\', l.ch)
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
		l.readChar() // consume closing "
	}

	return string(out)
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
