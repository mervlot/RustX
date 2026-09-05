package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under inspection
	line         int  // current line number (1-indexed)
	column       int  // current column number (1-indexed)
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

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ASSIGN, l.ch, tok.Line, tok.Column)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.COLON, l.ch, tok.Line, tok.Column)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.PLUS, l.ch, tok.Line, tok.Column)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.MINUS, l.ch, tok.Line, tok.Column)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.BANG, l.ch, tok.Line, tok.Column)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, tok.Line, tok.Column)
	case '/':
		tok = newToken(token.SLASH, l.ch, tok.Line, tok.Column)
	case '%':
		tok = newToken(token.MOD, l.ch, tok.Line, tok.Column)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.LT, l.ch, tok.Line, tok.Column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.GT, l.ch, tok.Line, tok.Column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: tok.Line, Column: tok.Column}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, tok.Line, tok.Column)
	case ',':
		tok = newToken(token.COMMA, l.ch, tok.Line, tok.Column)
	case '(':
		tok = newToken(token.LPAREN, l.ch, tok.Line, tok.Column)
	case ')':
		tok = newToken(token.RPAREN, l.ch, tok.Line, tok.Column)
	case '{':
		tok = newToken(token.LBRACE, l.ch, tok.Line, tok.Column)
	case '}':
		tok = newToken(token.RBRACE, l.ch, tok.Line, tok.Column)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, tok.Line, tok.Column)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, tok.Line, tok.Column)
	case '"':
		strVal := l.readString()
		tok.Type = token.STRING
		tok.Literal = strVal
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			lit, tokType := l.readNumber()
			tok.Literal = lit
			tok.Type = tokType
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, tok.Line, tok.Column)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
			continue
		}
		if l.ch == '/' {
			if l.peekChar() == '/' {
				l.readChar() // consume '/'
				l.readChar() // consume '/'
				for l.ch != '\n' && l.ch != 0 {
					l.readChar()
				}
				continue
			} else if l.peekChar() == '*' {
				l.readChar() // consume '/'
				l.readChar() // consume '*'
				for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
					l.readChar()
				}
				if l.ch == '*' {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
				}
				continue
			}
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

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	isFloat := false
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			if isFloat {
				break
			}
			if !isDigit(l.peekChar()) {
				break
			}
			isFloat = true
		}
		l.readChar()
	}
	lit := l.input[position:l.position]
	if isFloat {
		return lit, token.FLOAT
	}
	return lit, token.INT
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	var result []byte
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
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // advance past closing quote
	}
	return string(result)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
