package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under inspection
	line         int
	column       int
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
			tok = newTok(token.ASSIGN, l.ch, startLine, startColumn)
		}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: "++", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.PLUS, l.ch, startLine, startColumn)
		}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: "--", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.MINUS, l.ch, startLine, startColumn)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.BANG, l.ch, startLine, startColumn)
		}
	case '/':
		tok = newTok(token.SLASH, l.ch, startLine, startColumn)
	case '*':
		tok = newTok(token.ASTERISK, l.ch, startLine, startColumn)
	case '%':
		tok = newTok(token.MOD, l.ch, startLine, startColumn)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: "<=", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.LT, l.ch, startLine, startColumn)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: ">=", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.GT, l.ch, startLine, startColumn)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: "&&", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: "||", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	case ':':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: ":=", Line: startLine, Column: startColumn}
		} else {
			tok = newTok(token.COLON, l.ch, startLine, startColumn)
		}
	case ';':
		tok = newTok(token.SEMICOLON, l.ch, startLine, startColumn)
	case ',':
		tok = newTok(token.COMMA, l.ch, startLine, startColumn)
	case '(':
		tok = newTok(token.LPAREN, l.ch, startLine, startColumn)
	case ')':
		tok = newTok(token.RPAREN, l.ch, startLine, startColumn)
	case '{':
		tok = newTok(token.LBRACE, l.ch, startLine, startColumn)
	case '}':
		tok = newTok(token.RBRACE, l.ch, startLine, startColumn)
	case '[':
		tok = newTok(token.LBRACKET, l.ch, startLine, startColumn)
	case ']':
		tok = newTok(token.RBRACKET, l.ch, startLine, startColumn)
	case '"':
		strLit := l.readString()
		tok = token.Token{Type: token.STRING, Literal: strLit, Line: startLine, Column: startColumn}
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
			lit, tokType := l.readNumber()
			tok.Literal = lit
			tok.Type = tokType
			tok.Line = startLine
			tok.Column = startColumn
			return tok
		} else {
			tok = newTok(token.ILLEGAL, l.ch, startLine, startColumn)
		}
	}

	l.readChar()
	return tok
}

func newTok(tokenType token.TokenType, ch byte, line, column int) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
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
	l.readChar() // consume '/'
	l.readChar() // consume '*'
	for l.ch != 0 {
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar()
			l.readChar()
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

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	isFloat := false
	for isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar())) {
		if l.ch == '.' {
			isFloat = true
		}
		l.readChar()
	}
	tokType := token.TokenType(token.INT)
	if isFloat {
		tokType = token.FLOAT
	}
	return l.input[position:l.position], tokType
}

func (l *Lexer) readString() string {
	var out []rune
	l.readChar() // pass initial opening quote
	for l.ch != '"' && l.ch != 0 {
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
				out = append(out, rune(l.ch))
			}
		} else {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			out = append(out, rune(l.ch))
		}
		l.readChar()
	}
	l.readChar() // consume closing quote
	return string(out)
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
