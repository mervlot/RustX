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
	INC          = "++"
	DEC          = "--"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="
	LTE    = "<="
	GTE    = ">="
	AND    = "&&"
	OR     = "||"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	LBRACKET  = "["
	RBRACKET  = "]"

	// Keywords
	FUNC   = "FUNC"
	VAR    = "VAR"
	IF     = "IF"
	ELSE   = "ELSE"
	WHILE  = "WHILE"
	FOR    = "FOR"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NULL   = "NULL"

	// Type identifiers (optional for function signatures)
	TYPE_INT    = "INT_TYPE"
	TYPE_FLOAT  = "FLOAT_TYPE"
	TYPE_BOOL   = "BOOL_TYPE"
	TYPE_STRING = "STRING_TYPE"
)

var keywords = map[string]TokenType{
	"func":   FUNC,
	"var":    VAR,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
	"for":    FOR,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"int":    TYPE_INT,
	"float":  TYPE_FLOAT,
	"bool":   TYPE_BOOL,
	"string": TYPE_STRING,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
