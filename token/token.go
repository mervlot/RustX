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
	IDENT  TokenType = "IDENT"  // add, foobar, x, y, ...
	INT    TokenType = "INT"    // 123
	FLOAT  TokenType = "FLOAT"  // 3.14
	STRING TokenType = "STRING" // "foo"

	// Operators
	ASSIGN       TokenType = "="
	COLON_ASSIGN TokenType = ":="
	PLUS         TokenType = "+"
	MINUS        TokenType = "-"
	BANG         TokenType = "!"
	ASTERISK     TokenType = "*"
	SLASH        TokenType = "/"

	EQ     TokenType = "=="
	NOT_EQ TokenType = "!="
	LT     TokenType = "<"
	GT     TokenType = ">"
	LT_EQ  TokenType = "<="
	GT_EQ  TokenType = ">="

	AND TokenType = "&&"
	OR  TokenType = "||"

	INC TokenType = "++"
	DEC TokenType = "--"

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
