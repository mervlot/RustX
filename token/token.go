package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 134345
	FLOAT  = "FLOAT"  // 12.34
	STRING = "STRING" // "foobar"

	// Operators
	ASSIGN       = "="
	COLON_ASSIGN = ":="
	PLUS         = "+"
	MINUS        = "-"
	BANG         = "!"
	ASTERISK     = "*"
	SLASH        = "/"

	LT    = "<"
	GT    = ">"
	LT_EQ = "<="
	GT_EQ = ">="
	EQ    = "=="
	NOT_EQ = "!="

	AND = "&&"
	OR  = "||"

	INC = "++"
	DEC = "--"

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
	FUNC   = "FUNC"
	VAR    = "VAR"
	IF     = "IF"
	ELSE   = "ELSE"
	WHILE  = "WHILE"
	FOR    = "FOR"
	RETURN = "RETURN"

	// Types
	INT_TYPE    = "int"
	FLOAT_TYPE  = "float"
	BOOL_TYPE   = "bool"
	STRING_TYPE = "string"

	// Booleans
	TRUE  = "TRUE"
	FALSE = "FALSE"
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
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"return": RETURN,
	"int":    INT_TYPE,
	"float":  FLOAT_TYPE,
	"bool":   BOOL_TYPE,
	"string": STRING_TYPE,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
