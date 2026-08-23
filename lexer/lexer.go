package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under inspection
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

	tokLine := l.line
	tokColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.ASSIGN, l.ch, tokLine, tokColumn)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.PLUS, l.ch, tokLine, tokColumn)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.MINUS, l.ch, tokLine, tokColumn)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.BANG, l.ch, tokLine, tokColumn)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, tokLine, tokColumn)
	case '/':
		tok = newToken(token.SLASH, l.ch, tokLine, tokColumn)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.LT, l.ch, tokLine, tokColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.GT, l.ch, tokLine, tokColumn)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokLine, tokColumn)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokLine, tokColumn)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: tokLine, Column: tokColumn}
		} else {
			tok = newToken(token.COLON, l.ch, tokLine, tokColumn)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, tokLine, tokColumn)
	case ',':
		tok = newToken(token.COMMA, l.ch, tokLine, tokColumn)
	case '(':
		tok = newToken(token.LPAREN, l.ch, tokLine, tokColumn)
	case ')':
		tok = newToken(token.RPAREN, l.ch, tokLine, tokColumn)
	case '{':
		tok = newToken(token.LBRACE, l.ch, tokLine, tokColumn)
	case '}':
		tok = newToken(token.RBRACE, l.ch, tokLine, tokColumn)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, tokLine, tokColumn)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, tokLine, tokColumn)
	case '"':
		strLit := l.readString()
		tok = token.Token{Type: token.STRING, Literal: strLit, Line: tokLine, Column: tokColumn}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = tokLine
		tok.Column = tokColumn
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = tokLine
			tok.Column = tokColumn
			return tok
		} else if isDigit(l.ch) {
			tokType, tokLit := l.readNumber()
			tok.Type = tokType
			tok.Literal = tokLit
			tok.Line = tokLine
			tok.Column = tokColumn
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tokLine, tokColumn)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
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
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // consume '*'
			l.readChar() // consume '/'
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
	tokType := token.TokenType(token.INT)
	for isDigit(l.ch) {
		l.readChar()
	}
	if l.ch == '.' && isDigit(l.peekChar()) {
		tokType = token.FLOAT
		l.readChar() // consume '.'
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return tokType, l.input[position:l.position]
}

func (l *Lexer) readString() string {
	var result []byte
	l.readChar() // consume opening quote
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			default:
				result = append(result, '\\', l.ch)
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
		l.readChar() // consume closing quote
	}
	return string(result)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
