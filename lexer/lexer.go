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

func (l *Lexer) Line() int {
	return l.line
}

func (l *Lexer) Column() int {
	return l.column
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

		// Check for single-line // or multi-line /* comments
		if l.ch == '/' && l.peekChar() == '/' {
			l.skipSingleLineComment()
			continue
		} else if l.ch == '/' && l.peekChar() == '*' {
			l.skipMultiLineComment()
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
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.ASSIGN, l.ch, startLine, startColumn)
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
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.BANG, l.ch, startLine, startColumn)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, startLine, startColumn)
	case '/':
		tok = newToken(token.SLASH, l.ch, startLine, startColumn)
	case '%':
		tok = newToken(token.PERCENT, l.ch, startLine, startColumn)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LT_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.LT, l.ch, startLine, startColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GT_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
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
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startColumn}
		} else {
			tok = newToken(token.COLON, l.ch, startLine, startColumn)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, startLine, startColumn)
	case ',':
		tok = newToken(token.COMMA, l.ch, startLine, startColumn)
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
		strLiteral := l.readString()
		// Return directly so we don't call l.readChar() again at end of NextToken
		return token.Token{
			Type:    token.STRING,
			Literal: strLiteral,
			Line:    startLine,
			Column:  startColumn,
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startColumn}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			return token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startColumn}
		} else if isDigit(l.ch) {
			tokType, literal := l.readNumber()
			return token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startColumn}
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

	for isDigit(l.ch) || (l.ch == '.' && !isFloat && isDigit(l.peekChar())) {
		if l.ch == '.' {
			isFloat = true
		}
		l.readChar()
	}

	literal := l.input[position:l.position]
	if isFloat {
		return token.FLOAT, literal
	}
	return token.INT, literal
}

func (l *Lexer) readString() string {
	l.readChar() // Skip opening quote
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
		l.readChar() // Advance past closing double quote
	}

	return string(result)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
