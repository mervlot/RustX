package token

// TokenType represents the type of a lexical token.
type TokenType string

// Source location information
type Position struct {
	Line int
	Col  int
}

// Token represents a lexical token with its literal and source position.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 12345
	FLOAT  TokenType = "FLOAT"  // 123.45
	STRING TokenType = "STRING" // "hello world"

	// Operators
	ASSIGN       TokenType = "="
	PLUS         TokenType = "+"
	MINUS        TokenType = "-"
	BANG         TokenType = "!"
	ASTERISK     TokenType = "*"
	SLASH        TokenType = "/"
	PERCENT      TokenType = "%"
	PLUS_PLUS    TokenType = "++"
	MINUS_MINUS  TokenType = "--"
	COLON_ASSIGN TokenType = ":="

	LT     TokenType = "<"
	GT     TokenType = ">"
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	AND TokenType = "&&"
	OR  TokenType = "||"

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
	FUNC   TokenType = "FUNC"
	VAR    TokenType = "VAR"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	RETURN TokenType = "RETURN"
	WHILE  TokenType = "WHILE"
	FOR    TokenType = "FOR"
)

var keywords = map[string]TokenType{
	"func":   FUNC,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"for":    FOR,
}

// LookupIdent checks whether an identifier is a reserved keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
