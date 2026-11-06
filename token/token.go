package token

// TokenType represents the type of a token.
type TokenType string

// Token represents a lexical token with its position in source code.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENT  TokenType = "IDENT"  // add, x, y, ...
	INT    TokenType = "INT"    // 12345
	FLOAT  TokenType = "FLOAT"  // 123.45
	STRING TokenType = "STRING" // "abc"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	PERCENT  TokenType = "%"

	LT TokenType = "<"
	GT TokenType = ">"
	LE TokenType = "<="
	GE TokenType = ">="
	EQ TokenType = "=="
	NOT_EQ TokenType = "!="

	AND TokenType = "&&"
	OR  TokenType = "||"

	PLUS_PLUS   TokenType = "++"
	MINUS_MINUS TokenType = "--"
	COLON_ASSIGN TokenType = ":="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

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
)

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

// LookupIdent checks if an identifier is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
