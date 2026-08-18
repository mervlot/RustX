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
	ASSIGN   = "="
	COLON_EQ = ":="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	INC      = "++"
	DEC      = "--"

	EQ     = "=="
	NOT_EQ = "!="
	LT     = "<"
	GT     = ">"
	LTE    = "<="
	GTE    = ">="
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
	FUNC   = "FUNC"
	VAR    = "VAR"
	IF     = "IF"
	ELSE   = "ELSE"
	WHILE  = "WHILE"
	FOR    = "FOR"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
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
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
