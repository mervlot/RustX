package token

// TokenType represents the type of a lexical token.
type TokenType string

// Position tracks file coordinates for diagnostic reporting.
type Position struct {
	Line   int
	Column int
}

// Token represents a lexical token in RustX.
type Token struct {
	Type     TokenType
	Literal  string
	Position Position
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers & Literals
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	ASSIGN       = "="
	COLON_ASSIGN = ":="
	PLUS         = "+"
	MINUS        = "-"
	BANG         = "!"
	ASTERISK     = "*"
	SLASH        = "/"
	PERCENT      = "%"
	PLUS_PLUS    = "++"
	MINUS_MINUS  = "--"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LT_EQ  = "<="
	GT_EQ  = ">="
	AND    = "&&"
	OR     = "||"

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
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	WHILE    = "WHILE"
	FOR      = "FOR"
	RETURN   = "RETURN"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"return": RETURN,
}

// LookupIdent checks whether an identifier string is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
