package lexer

import (
	"rustx/token"
	"strings"
)

// Lexer transforms source code into a stream of tokens.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number (1-based)
	column       int  // current column number (1-based)
}

// New creates a new Lexer instance.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// Input returns the full source input string.
func (l *Lexer) Input() string {
	return l.input
}

// Position returns the current index position in the input string.
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

// NextToken scans and returns the next token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	for {
		l.skipWhitespace()
		// Check for single line or multi line comments
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

	startLine := l.line
	startCol := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.ASSIGN, l.ch, startLine, startCol)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.PLUS, l.ch, startLine, startCol)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.MINUS, l.ch, startLine, startCol)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.BANG, l.ch, startLine, startCol)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, startLine, startCol)
	case '/':
		tok = newToken(token.SLASH, l.ch, startLine, startCol)
	case '%':
		tok = newToken(token.MODULO, l.ch, startLine, startCol)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.LT, l.ch, startLine, startCol)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GE, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.GT, l.ch, startLine, startCol)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startCol)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startCol)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Line: startLine, Column: startCol}
		} else {
			tok = newToken(token.COLON, l.ch, startLine, startCol)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, startLine, startCol)
	case ',':
		tok = newToken(token.COMMA, l.ch, startLine, startCol)
	case '(':
		tok = newToken(token.LPAREN, l.ch, startLine, startCol)
	case ')':
		tok = newToken(token.RPAREN, l.ch, startLine, startCol)
	case '{':
		tok = newToken(token.LBRACE, l.ch, startLine, startCol)
	case '}':
		tok = newToken(token.RBRACE, l.ch, startLine, startCol)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, startLine, startCol)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, startLine, startCol)
	case '"':
		str, ok := l.readString()
		if !ok {
			return token.Token{Type: token.ILLEGAL, Literal: str, Line: startLine, Column: startCol}
		}
		return token.Token{Type: token.STRING, Literal: str, Line: startLine, Column: startCol}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = startLine
		tok.Column = startCol
	default:
		if isLetter(l.ch) {
			lit := l.readIdentifier()
			return token.Token{
				Type:    token.LookupIdent(lit),
				Literal: lit,
				Line:    startLine,
				Column:  startCol,
			}
		} else if isDigit(l.ch) {
			numLit, tokType := l.readNumber()
			return token.Token{
				Type:    tokType,
				Literal: numLit,
				Line:    startLine,
				Column:  startCol,
			}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, startLine, startCol)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, line int, col int) token.Token {
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
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	// Skip the initial '/*'
	l.readChar()
	l.readChar()

	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // eat '*'
			l.readChar() // eat '/'
			break
		}
		l.readChar()
	}
}

func (l *Lexer) readString() (string, bool) {
	var sb strings.Builder
	for {
		l.readChar()
		if l.ch == 0 {
			return sb.String(), false // Unterminated string
		}
		if l.ch == '"' {
			break
		}
		if l.ch == '\\' {
			l.readChar()
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
			case '0':
				sb.WriteByte(0)
			default:
				sb.WriteByte('\\')
				sb.WriteByte(l.ch)
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			sb.WriteByte(l.ch)
		}
	}
	l.readChar() // Advance past closing quote
	return sb.String(), true
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	startPos := l.position
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

	return l.input[startPos:l.position], tokType
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
