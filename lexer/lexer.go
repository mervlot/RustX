package lexer

import (
	"github.com/username/rustx/token"
)

// Lexer processes source code character by character
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int  // current line number (1-based)
	column       int  // current column number (1-based)
}

// New creates a new Lexer instance
func New(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
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
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken scans the next token from the input
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
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Line: startLine, Column: startColumn}
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
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: literal, Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.COLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
		}
	case ';':
		tok = token.Token{Type: token.SEMICOLON, Literal: string(l.ch), Line: startLine, Column: startColumn}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch), Line: startLine, Column: startColumn}
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
		str, err := l.readString()
		if err != nil {
			tok = token.Token{Type: token.ILLEGAL, Literal: err.Error(), Line: startLine, Column: startColumn}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Line: startLine, Column: startColumn}
		}
		return tok
	case 0:
		tok = token.Token{Type: token.EOF, Literal: "", Line: startLine, Column: startColumn}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Line = startLine
			tok.Column = startColumn
			return tok
		} else if isDigit(l.ch) {
			lit, isF := l.readNumber()
			tok.Literal = lit
			if isF {
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

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		if l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
			continue
		}

		if l.ch == '/' && l.peekChar() == '/' {
			// Single-line comment
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		if l.ch == '/' && l.peekChar() == '*' {
			// Multi-line comment
			l.readChar() // consume '/'
			l.readChar() // consume '*'
			for {
				if l.ch == 0 {
					break // unfinished comment handled gracefully or left as EOF
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

func (l *Lexer) readString() (string, error) {
	l.readChar() // step over the opening quote
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
				// Return as is if unknown escape
				result = append(result, '\\')
				result = append(result, l.ch)
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
	// Note: We MUST advance past the closing quote in NextToken if it returned,
	// but since we returned inside case '"', let's do it here.
	l.readChar() // step over closing quote
	return string(result), nil
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
