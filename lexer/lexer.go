package lexer

import (
	"github.com/username/rustx/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under inspection
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

	l.skipWhitespaceAndComments()

	pos := token.Position{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal, Pos: pos}
		} else {
			tok = newToken(token.ASSIGN, l.ch, pos)
		}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok = token.Token{Type: token.INCREMENT, Literal: "++", Pos: pos}
		} else {
			tok = newToken(token.PLUS, l.ch, pos)
		}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			tok = token.Token{Type: token.DECREMENT, Literal: "--", Pos: pos}
		} else {
			tok = newToken(token.MINUS, l.ch, pos)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal, Pos: pos}
		} else {
			tok = newToken(token.BANG, l.ch, pos)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch, pos)
	case '/':
		tok = newToken(token.SLASH, l.ch, pos)
	case '%':
		tok = newToken(token.MODULO, l.ch, pos)
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: "<=", Pos: pos}
		} else {
			tok = newToken(token.LT, l.ch, pos)
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: ">=", Pos: pos}
		} else {
			tok = newToken(token.GT, l.ch, pos)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = token.Token{Type: token.AND, Literal: "&&", Pos: pos}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, pos)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = token.Token{Type: token.OR, Literal: "||", Pos: pos}
		} else {
			tok = newToken(token.ILLEGAL, l.ch, pos)
		}
	case ':':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.COLON_ASSIGN, Literal: ":=", Pos: pos}
		} else {
			tok = newToken(token.COLON, l.ch, pos)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, pos)
	case ',':
		tok = newToken(token.COMMA, l.ch, pos)
	case '(':
		tok = newToken(token.LPAREN, l.ch, pos)
	case ')':
		tok = newToken(token.RPAREN, l.ch, pos)
	case '{':
		tok = newToken(token.LBRACE, l.ch, pos)
	case '}':
		tok = newToken(token.RBRACE, l.ch, pos)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, pos)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, pos)
	case '"':
		str, err := l.readString()
		if err != "" {
			tok = token.Token{Type: token.ILLEGAL, Literal: err, Pos: pos}
		} else {
			tok = token.Token{Type: token.STRING, Literal: str, Pos: pos}
		}
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.Pos = pos
	default:
		if isLetter(l.ch) {
			tok.Pos = pos
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Pos = pos
			lit, tokType := l.readNumber()
			tok.Literal = lit
			tok.Type = tokType
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, pos)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte, pos token.Position) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch), Pos: pos}
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
			if l.ch == '\n' {
				l.line++
				l.column = 0
			}
			l.readChar()
		}

		// Handle comments
		if l.ch == '/' {
			if l.peekChar() == '/' {
				// Single-line comment
				l.readChar()
				l.readChar()
				for l.ch != '\n' && l.ch != 0 {
					l.readChar()
				}
				continue
			} else if l.peekChar() == '*' {
				// Multi-line comment
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
						l.readChar()
						l.readChar()
						break
					}
					l.readChar()
				}
				continue
			}
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

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
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
	return l.input[position:l.position], tokType
}

func (l *Lexer) readString() (string, string) {
	l.readChar() // skip opening "
	var result []rune
	for {
		if l.ch == 0 || l.ch == '\n' {
			return "", "Unterminated string literal"
		}
		if l.ch == '"' {
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
				result = append(result, rune(l.ch))
			}
		} else {
			result = append(result, rune(l.ch))
		}
		l.readChar()
	}
	l.readChar() // skip closing "
	return string(result), ""
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
