package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // line number (1-based)
	column       int  // column number (1-based)
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
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

	l.skipWhitespace()

	// In nextToken, we handle potentially multiple comments back-to-back.
	// Comments should be skipped, and we check whitespace again.
	for l.ch == '/' {
		if l.peekChar() == '/' {
			l.skipSingleLineComment()
			l.skipWhitespace()
			continue
		} else if l.peekChar() == '*' {
			l.skipMultiLineComment()
			l.skipWhitespace()
			continue
		}
		break
	}

	startLine := l.line
	startColumn := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ',':
		tok = token.Token{Type: token.COMMA, Line: startLine, Column: startColumn}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
		tok.Line = startLine
		tok.Column = startColumn
		// Note: readString already consumed the trailing '"', and we don't want to call readChar here
		// as we might skip the very next non-whitespace char.
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
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
			literal, isFloat := l.readNumber()
			tok.Literal = literal
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
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
	if l.ch == '\n' {
		l.line++
		l.column = 0
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	// Skip '/*'
	l.readChar()
	l.readChar()
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

func (l *Lexer) readString() string {
	l.readChar() // Skip opening quote
	var runes []rune
	for {
		if l.ch == '"' || l.ch == 0 {
			break
		}
		if l.ch == '\\' {
			l.readChar() // consume '\\'
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
				// Fallback
				runes = append(runes, '\\')
				runes = append(runes, rune(l.ch))
			}
		} else {
			runes = append(runes, rune(l.ch))
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // consume closing quote
	}
	return string(runes)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
