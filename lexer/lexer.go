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
	col          int  // current column number (1-based)
}

func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, col: 0}
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
	l.col++
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

	line := l.line
	col := l.col

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch), Line: line, Col: col}
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.PLUS_PLUS, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.PLUS, Literal: string(l.ch), Line: line, Col: col}
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.MINUS_MINUS, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.MINUS, Literal: string(l.ch), Line: line, Col: col}
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.BANG, Literal: string(l.ch), Line: line, Col: col}
		}
	case '/':
		tok = token.Token{Type: token.SLASH, Literal: string(l.ch), Line: line, Col: col}
	case '*':
		tok = token.Token{Type: token.ASTERISK, Literal: string(l.ch), Line: line, Col: col}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.LTE, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.LT, Literal: string(l.ch), Line: line, Col: col}
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.GTE, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.GT, Literal: string(l.ch), Line: line, Col: col}
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.AND, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Col: col}
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.OR, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Col: col}
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: line, Col: col}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: line, Col: col}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: line, Col: col}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch), Line: line, Col: col}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch), Line: line, Col: col}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch), Line: line, Col: col}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch), Line: line, Col: col}
	case '[':
		tok = token.Token{Type: token.LBRACKET, Literal: string(l.ch), Line: line, Col: col}
	case ']':
		tok = token.Token{Type: token.RBRACKET, Literal: string(l.ch), Line: line, Col: col}
	case '"':
		str, err := l.readString()
		if err != nil {
			tok = token.Token{Type: token.ILLEGAL, Literal: err.Error(), Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Line: line, Col: col}
		}
		return tok
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: line, Col: col}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := token.LookupIdent(literal)
			return token.Token{Type: tokType, Literal: literal, Line: line, Col: col}
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			var tokType token.TokenType = token.INT
			if isFloat {
				tokType = token.FLOAT
			}
			return token.Token{Type: tokType, Literal: literal, Line: line, Col: col}
		} else {
			tok = token.Token{Type: token.ILLEGAL, Literal: string(l.ch), Line: line, Col: col}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
			if l.ch == '\n' {
				l.line++
				l.col = 0
			}
			l.readChar()
			continue
		}

		// Check for single-line or multi-line comments
		if l.ch == '/' && l.peekChar() == '/' {
			// Single-line comment
			l.readChar()
			l.readChar()
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		if l.ch == '/' && l.peekChar() == '*' {
			// Multi-line comment
			l.readChar()
			l.readChar()
			for {
				if l.ch == 0 {
					break
				}
				if l.ch == '\n' {
					l.line++
					l.col = 0
				}
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar()
					l.readChar()
					break
				}
				l.readChar()
			}
			continue
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

func (l *Lexer) readNumber() (string, bool) {
	position := l.position
	isFloat := false
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			// make sure it's not a second decimal point
			if isFloat {
				break
			}
			// make sure it's followed by a digit (e.g. not a method call or range)
			next := l.peekChar()
			if !isDigit(next) {
				break
			}
			isFloat = true
		}
		l.readChar()
	}
	return l.input[position:l.position], isFloat
}

func (l *Lexer) readString() (string, error) {
	l.readChar() // Skip starting double quote
	var result []byte

	for l.ch != '"' && l.ch != 0 {
		if l.ch == '\\' {
			// Handle escape sequence
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
				result = append(result, '\\')
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	l.readChar() // Skip closing double quote

	return string(result), nil
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
