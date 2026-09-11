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
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	FLOAT  TokenType = "FLOAT"
	STRING TokenType = "STRING"

	// Operators
	ASSIGN       TokenType = "="
	COLON_ASSIGN TokenType = ":="
	PLUS         TokenType = "+"
	MINUS        TokenType = "-"
	BANG         TokenType = "!"
	ASTERISK     TokenType = "*"
	SLASH        TokenType = "/"
	INC          TokenType = "++"
	DEC          TokenType = "--"

	// Comparisons
	LT     TokenType = "<"
	GT     TokenType = ">"
	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LTE    TokenType = "<="
	GTE    TokenType = ">="

	// Logical
	AND TokenType = "&&"
	OR  TokenType = "||"

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	COLON     TokenType = ":"

	LPAREN TokenType = "("
	RPAREN TokenType = ")"
	LBRACE TokenType = "{"
	RBRACE TokenType = "}"
	LBRACK TokenType = "["
	RBRACK TokenType = "]"

	// Keywords
	FUNC   TokenType = "FUNC"
	VAR    TokenType = "VAR"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	WHILE  TokenType = "WHILE"
	FOR    TokenType = "FOR"
	RETURN TokenType = "RETURN"
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
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
