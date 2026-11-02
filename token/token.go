package token

type TokenType string

type Position struct {
	Line   int
	Column int
}

type Token struct {
	Type     TokenType
	Literal  string
	Position Position
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	FLOAT  = "FLOAT"  // 123.456
	STRING = "STRING" // "foobar"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	MODULO   = "%"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="
	LTE    = "<="
	GTE    = ">="

	AND = "&&"
	OR  = "||"

	PLUS_PLUS   = "++"
	MINUS_MINUS = "--"

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
	FUNCTION = "FUNCTION"
	VAR      = "VAR"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"

	// Type keywords (optional typing)
	TYPE_INT    = "TYPE_INT"
	TYPE_FLOAT  = "TYPE_FLOAT"
	TYPE_BOOL   = "TYPE_BOOL"
	TYPE_STRING = "TYPE_STRING"
	TYPE_VOID   = "TYPE_VOID"
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
	"int":    TYPE_INT,
	"float":  TYPE_FLOAT,
	"bool":   TYPE_BOOL,
	"string": TYPE_STRING,
	"void":   TYPE_VOID,
}

// LookupIdent checks if the identifier is a keyword.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
