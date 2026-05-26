package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 1343456
	FLOAT  TokenType = "FLOAT"  // 123.45
	STRING TokenType = "STRING" // "foobar"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	INC      TokenType = "++"

	LT TokenType = "<"
	GT TokenType = ">"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LTE    TokenType = "<="
	GTE    TokenType = ">="

	AND TokenType = "&&"
	OR  TokenType = "||"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	VAR      TokenType = "VAR"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	RETURN   TokenType = "RETURN"
	WHILE    TokenType = "WHILE"
	FOR      TokenType = "FOR"

	// Shorthand
	DECLARE TokenType = ":="
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
