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
			l.readChar()
			l.readChar()
			break
		}
		l.readChar()
	}
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
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ASSIGN, l.ch, startLine, startColumn)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.COLON, l.ch, startLine, startColumn)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.PLUS, l.ch, startLine, startColumn)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.MINUS, l.ch, startLine, startColumn)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.BANG, l.ch, startLine, startColumn)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, startLine, startColumn)
	case '/':
		tok = newToken(token.SLASH, l.ch, startLine, startColumn)
	case '%':
		tok = newToken(token.MODULO, l.ch, startLine, startColumn)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.LT, l.ch, startLine, startColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.GT, l.ch, startLine, startColumn)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	case ',':
		tok = newToken(token.COMMA, l.ch, startLine, startColumn)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, startLine, startColumn)
	case '(':
		tok = newToken(token.LPAREN, l.ch, startLine, startColumn)
	case ')':
		tok = newToken(token.RPAREN, l.ch, startLine, startColumn)
	case '{':
		tok = newToken(token.LBRACE, l.ch, startLine, startColumn)
	case '}':
		tok = newToken(token.RBRACE, l.ch, startLine, startColumn)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, startLine, startColumn)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, startLine, startColumn)
	case '"':
		str, err := l.readString()
		if err != "" {
			tok = token.Token{Type: token.ILLEGAL, Literal: err, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Line: startLine, Column: startColumn}
		}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = startLine
		tok.Column = startColumn
	default:
		if isLetter(l.ch) {
			tok.Line = startLine
			tok.Column = startColumn
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = startLine
			tok.Column = startColumn
			lit, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = lit
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
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

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
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

func (l *Lexer) readString() (string, string) {
	// Skip initial quote
	l.readChar()
	var runes []rune
	for {
		if l.ch == '"' {
			// Read past the closing quote
			l.readChar()
			break
		}
		if l.ch == 0 {
			return "", "Unterminated string literal"
		}
		if l.ch == '\\' {
			l.readChar()
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
				runes = append(runes, '\\')
				runes = append(runes, rune(l.ch))
			}
		} else {
			runes = append(runes, rune(l.ch))
		}
		l.readChar()
	}
	return string(runes), ""
}
