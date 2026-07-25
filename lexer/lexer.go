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
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
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
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INCREMENT, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DECREMENT, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: startLine, Column: startColumn}
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
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: startLine, Column: startColumn}
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
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = startLine
		tok.Column = startColumn
		return tok
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
		tok.Line = startLine
		tok.Column = startColumn
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = startLine
			tok.Column = startColumn
			return tok
		} else if isDigit(l.ch) {
			tok.Literal, tok.Type = l.readNumber()
			tok.Line = startLine
			tok.Column = startColumn
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
	l.readChar()
	l.readChar()
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	l.readChar()
	l.readChar()
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
}

func (l *Lexer) readString() string {
	l.readChar() // skip '"'
	var runes []rune
	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			l.readChar() // skip '\'
			switch l.ch {
			case 'n':
				runes = append(runes, '\n')
			case 't':
				runes = append(runes, '\t')
			case 'r':
				runes = append(runes, '\r')
			case '"':
				runes = append(runes, '"')
			case '\\':
				runes = append(runes, '\\')
			default:
				runes = append(runes, '\\', rune(l.ch))
			}
		} else {
			runes = append(runes, rune(l.ch))
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // advance past the closing double-quote
	}
	return string(runes)
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

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	isFloat := false
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			if isFloat {
				break
			}
			isFloat = true
		}
		l.readChar()
	}
	literal := l.input[position:l.position]
	if isFloat {
		return literal, token.FLOAT
	}
	return literal, token.INT
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
