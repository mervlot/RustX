package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
let ten := 10.5
// single line comment
/* multiline
   comment */
func add(x int, y int) int {
	return x + y
}
if (five == 5 && ten != 5) || five <= 10 {
	five++
	ten--
}
name := "Hello\nWorld"
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.IDENT, "let"},
		{token.IDENT, "ten"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},
		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "ten"},
		{token.NOT_EQ, "!="},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.OR, "||"},
		{token.IDENT, "five"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.IDENT, "five"},
		{token.INC, "++"},
		{token.IDENT, "ten"},
		{token.DEC, "--"},
		{token.RBRACE, "}"},
		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Hello\nWorld"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
