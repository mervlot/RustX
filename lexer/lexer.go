package lexer

import (
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

	curLine := l.line
	curCol := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ASSIGN, l.ch, curLine, curCol)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.INC, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.PLUS, l.ch, curLine, curCol)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.DEC, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.MINUS, l.ch, curLine, curCol)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.BANG, l.ch, curLine, curCol)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, curLine, curCol)
	case '/':
		tok = newToken(token.SLASH, l.ch, curLine, curCol)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.LT, l.ch, curLine, curCol)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.GT, l.ch, curLine, curCol)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, curLine, curCol)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, curLine, curCol)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: curLine, Column: curCol}
		} else {
			tok = newToken(token.COLON, l.ch, curLine, curCol)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, curLine, curCol)
	case ',':
		tok = newToken(token.COMMA, l.ch, curLine, curCol)
	case '(':
		tok = newToken(token.LPAREN, l.ch, curLine, curCol)
	case ')':
		tok = newToken(token.RPAREN, l.ch, curLine, curCol)
	case '{':
		tok = newToken(token.LBRACE, l.ch, curLine, curCol)
	case '}':
		tok = newToken(token.RBRACE, l.ch, curLine, curCol)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, curLine, curCol)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, curLine, curCol)
	case '"':
		tok.Line = curLine
		tok.Column = curCol
		tok.Type = token.STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Line = curLine
		tok.Column = curCol
	default:
		if isLetter(l.ch) {
			tok.Line = curLine
			tok.Column = curCol
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Line = curLine
			tok.Column = curCol
			literal, isFloat := l.readNumber()
			if isFloat {
				tok.Type = token.FLOAT
			} else {
				tok.Type = token.INT
			}
			tok.Literal = literal
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, curLine, curCol)
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
	l.readChar() // skip /
	l.readChar() // skip *
	for {
		if l.ch == 0 {
			break
		}
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // skip *
			l.readChar() // skip /
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
	for isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar())) {
		if l.ch == '.' {
			isFloat = true
		}
		l.readChar()
	}
	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() string {
	var result []byte
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
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
	}
	l.readChar() // consume closing quote
	return string(result)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
