package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken_Basic(t *testing.T) {
	input := `var x = 5
y := 10.5
func add(a int, b int) int {
	return a + b
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.IDENT, "y"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},
		{token.FUNC, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.IDENT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.IDENT, "int"},
		{token.RPAREN, ")"},
		{token.IDENT, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.RBRACE, "}"},
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

func TestNextToken_OperatorsAndComments(t *testing.T) {
	input := `
// Single line comment
var a = 1 + 2 * 3 / 4 % 5;
/* Multi
   line
   comment */
if a == 1 && a != 2 || a < 3 && a <= 4 || a > 0 && a >= 0 {
    a++
    a--
}
str := "hello\nworld\""
[1, 2]
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "a"},
		{token.ASSIGN, "="},
		{token.INT, "1"},
		{token.PLUS, "+"},
		{token.INT, "2"},
		{token.ASTERISK, "*"},
		{token.INT, "3"},
		{token.SLASH, "/"},
		{token.INT, "4"},
		{token.MODULO, "%"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.INT, "1"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.INT, "2"},
		{token.OR, "||"},
		{token.IDENT, "a"},
		{token.LT, "<"},
		{token.INT, "3"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.LTE, "<="},
		{token.INT, "4"},
		{token.OR, "||"},
		{token.IDENT, "a"},
		{token.GT, ">"},
		{token.INT, "0"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.GTE, ">="},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "a"},
		{token.PLUS_PLUS, "++"},
		{token.IDENT, "a"},
		{token.MINUS_MINUS, "--"},
		{token.RBRACE, "}"},
		{token.IDENT, "str"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "hello\nworld\""},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q at line %d col %d",
				i, tt.expectedType, tok.Type, tok.Line, tok.Column)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
