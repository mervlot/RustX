package token

// TokenType represents the type of a lexical token.
type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + Literals
	IDENT  TokenType = "IDENT"  // add, x, foo
	INT    TokenType = "INT"    // 123
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

	EQ  TokenType = "=="
	NOT_EQ TokenType = "!="
	LT  TokenType = "<"
	GT  TokenType = ">"
	LTE TokenType = "<="
	GTE TokenType = ">="

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
	FUNCTION TokenType = "FUNCTION" // func
	VAR      TokenType = "VAR"      // var
	TRUE     TokenType = "TRUE"     // true
	FALSE    TokenType = "FALSE"    // false
	IF       TokenType = "IF"       // if
	ELSE     TokenType = "ELSE"     // else
	RETURN   TokenType = "RETURN"   // return
	WHILE    TokenType = "WHILE"    // while
	FOR      TokenType = "FOR"      // for
)

// Token represents a lexical token in RustX.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
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

// LookupIdent checks if an identifier is a reserved keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
