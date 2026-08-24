package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, x, y
	INT    = "INT"    // 123
	FLOAT  = "FLOAT"  // 3.14
	STRING = "STRING" // "hello"

	// Operators
	ASSIGN       = "="
	COLON_ASSIGN = ":="
	PLUS         = "+"
	MINUS        = "-"
	BANG         = "!"
	ASTERISK     = "*"
	SLASH        = "/"
	INC          = "++"
	DEC          = "--"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="
	LTE    = "<="
	GTE    = ">="

	AND = "&&"
	OR  = "||"

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

	INT_TYPE    = "INT_TYPE"
	FLOAT_TYPE  = "FLOAT_TYPE"
	BOOL_TYPE   = "BOOL_TYPE"
	STRING_TYPE = "STRING_TYPE"
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
	"int":    INT_TYPE,
	"float":  FLOAT_TYPE,
	"bool":   BOOL_TYPE,
	"string": STRING_TYPE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
