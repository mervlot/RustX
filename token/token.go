package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT  TokenType = "IDENT"  // add, x, y, ...
	INT    TokenType = "INT"    // 123
	FLOAT  TokenType = "FLOAT"  // 3.14
	STRING TokenType = "STRING" // "hello"

	// Operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	BANG     TokenType = "!"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%"

	LT     TokenType = "<"
	GT     TokenType = ">"
	LTE    TokenType = "<="
	GTE    TokenType = ">="
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="

	AND TokenType = "&&"
	OR  TokenType = "||"

	INC TokenType = "++"
	DEC TokenType = "--"

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

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
