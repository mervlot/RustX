package lexer

import (
	"strings"

	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
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

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	for {
		l.skipWhitespace()

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
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.ASSIGN, l.ch, startLine, startCol)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.COLON, l.ch, startLine, startCol)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.PLUS, l.ch, startLine, startCol)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.MINUS, l.ch, startLine, startCol)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.BANG, l.ch, startLine, startCol)
		}
	case '*':
		tok = l.newToken(token.ASTERISK, l.ch, startLine, startCol)
	case '/':
		tok = l.newToken(token.SLASH, l.ch, startLine, startCol)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LT_EQ, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.LT, l.ch, startLine, startCol)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GT_EQ, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.GT, l.ch, startLine, startCol)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch, startLine, startCol)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch, startLine, startCol)
		}
	case ';':
		tok = l.newToken(token.SEMICOLON, l.ch, startLine, startCol)
	case ',':
		tok = l.newToken(token.COMMA, l.ch, startLine, startCol)
	case '(':
		tok = l.newToken(token.LPAREN, l.ch, startLine, startCol)
	case ')':
		tok = l.newToken(token.RPAREN, l.ch, startLine, startCol)
	case '{':
		tok = l.newToken(token.LBRACE, l.ch, startLine, startCol)
	case '}':
		tok = l.newToken(token.RBRACE, l.ch, startLine, startCol)
	case '[':
		tok = l.newToken(token.LBRACK, l.ch, startLine, startCol)
	case ']':
		tok = l.newToken(token.RBRACK, l.ch, startLine, startCol)
	case '"':
		strLiteral := l.readString()
		tok = token.Token{Type: token.STRING, Literal: strLiteral, Line: startLine, Column: startCol}
		return tok
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startCol}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			return token.Token{Type: tokType, Literal: literal, Line: startLine, Column: startCol}
		} else if isDigit(l.ch) {
			numLiteral, tokType := l.readNumber()
			return token.Token{Type: tokType, Literal: numLiteral, Line: startLine, Column: startCol}
		} else {
			tok = l.newToken(token.ILLEGAL, l.ch, startLine, startCol)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType token.TokenType, ch byte, line, col int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: col}
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
	for isDigit(l.ch) || (l.ch == '.' && !isFloat && isDigit(l.peekChar())) {
		if l.ch == '.' {
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

func (l *Lexer) readString() string {
	var sb strings.Builder
	l.readChar() // skip opening '"'
	for {
		if l.ch == '"' || l.ch == 0 {
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
			default:
				sb.WriteByte('\\')
				sb.WriteByte(l.ch)
			}
		} else {
			sb.WriteByte(l.ch)
		}
		l.readChar()
	}
	if l.ch == '"' {
		l.readChar() // advance past closing quote
	}
	return sb.String()
}

func (l *Lexer) skipSingleLineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipMultiLineComment() {
	l.readChar() // consume '/'
	l.readChar() // consume '*'
	for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
		l.readChar()
	}
	if l.ch == '*' && l.peekChar() == '/' {
		l.readChar() // consume '*'
		l.readChar() // consume '/'
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
