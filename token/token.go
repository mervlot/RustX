package token

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 12345
	FLOAT  TokenType = "FLOAT"  // 123.45
	STRING TokenType = "STRING" // "foo"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%"

	INC TokenType = "++"
	DEC TokenType = "--"

	COLON_ASSIGN TokenType = ":="

	LT  TokenType = "<"
	GT  TokenType = ">"
	LTE TokenType = "<="
	GTE TokenType = ">="
	EQ  TokenType = "=="
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

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

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

// LookupIdent checks if an identifier is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
