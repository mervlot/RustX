package token

import "fmt"

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // main, add, x, y, ...
	INT    = "INT"    // 12345
	FLOAT  = "FLOAT"  // 123.45
	STRING = "STRING" // "hello"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"

	INC = "++"
	DEC = "--"

	LT  = "<"
	GT  = ">"
	EQ  = "=="
	NOT_EQ = "!="
	LTE = "<="
	GTE = ">="

	AND = "&&"
	OR  = "||"

	// Shorthand assignment / declarations
	COLON_ASSIGN = ":="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNC     = "FUNC"
	VAR      = "VAR"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	WHILE    = "WHILE"
	FOR      = "FOR"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%s, %q, %d:%d)", t.Type, t.Literal, t.Line, t.Column)
}

var keywords = map[string]TokenType{
	"func":   FUNC,
	"var":    VAR,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"while":  WHILE,
	"for":    FOR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
