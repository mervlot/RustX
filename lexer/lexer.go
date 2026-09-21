package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
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

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespaceAndComments()

	pos := token.Position{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.ASSIGN, l.ch, pos)
		}
	case '+':
		if l.peekChar() == '+' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.INC, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.PLUS, l.ch, pos)
		}
	case '-':
		if l.peekChar() == '-' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.DEC, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.MINUS, l.ch, pos)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.BANG, l.ch, pos)
		}
	case '*':
		tok = newSingleCharToken(token.ASTERISK, l.ch, pos)
	case '/':
		tok = newSingleCharToken(token.SLASH, l.ch, pos)
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.LT, l.ch, pos)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.GT, l.ch, pos)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.ILLEGAL, l.ch, pos)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.ILLEGAL, l.ch, pos)
		}
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: string(ch) + string(l.ch), Pos: pos}
		} else {
			tok = newSingleCharToken(token.COLON, l.ch, pos)
		}
	case ';':
		tok = newSingleCharToken(token.SEMICOLON, l.ch, pos)
	case ',':
		tok = newSingleCharToken(token.COMMA, l.ch, pos)
	case '(':
		tok = newSingleCharToken(token.LPAREN, l.ch, pos)
	case ')':
		tok = newSingleCharToken(token.RPAREN, l.ch, pos)
	case '{':
		tok = newSingleCharToken(token.LBRACE, l.ch, pos)
	case '}':
		tok = newSingleCharToken(token.RBRACE, l.ch, pos)
	case '[':
		tok = newSingleCharToken(token.LBRACKET, l.ch, pos)
	case ']':
		tok = newSingleCharToken(token.RBRACKET, l.ch, pos)
	case '"':
		strLiteral, err := l.readString()
		if err != nil {
			tok = token.Token{Type: token.ILLEGAL, Literal: strLiteral, Pos: pos}
		} else {
			tok = token.Token{Type: token.STRING, Literal: strLiteral, Pos: pos}
		}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Pos = pos
		return tok
	default:
		if isLetter(l.ch) {
			tok.Pos = pos
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Pos = pos
			tok.Literal, tok.Type = l.readNumber()
			return tok
		} else {
			tok = newSingleCharToken(token.ILLEGAL, l.ch, pos)
		}
	}

	l.readChar()
	return tok
}

func newSingleCharToken(tokenType token.TokenType, ch byte, pos token.Position) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Pos: pos}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		// Skip whitespace
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			l.readChar()
		}

		// Check for single line comment //
		if l.ch == '/' && l.peekChar() == '/' {
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
			continue
		}

		// Check for multi-line comment /* */
		if l.ch == '/' && l.peekChar() == '*' {
			l.readChar() // consume '/'
			l.readChar() // consume '*'
			for {
				if l.ch == 0 {
					break
				}
				if l.ch == '*' && l.peekChar() == '/' {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
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

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	isFloat := false

	for isDigit(l.ch) || (l.ch == '.' && isDigit(l.peekChar()) && !isFloat) {
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

func (l *Lexer) readString() (string, error) {
	var result []byte
	l.readChar() // pass starting double quote

	for {
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
				result = append(result, '\\', l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
		l.readChar()
	}

	if l.ch == '"' {
		l.readChar() // pass closing quote
	}

	return string(result), nil
}
