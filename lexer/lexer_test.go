package lexer

import (
	"testing"
	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
ten := 10.5

// This is a single line comment
/* This is a
   multiline comment */

func add(x int, y int) int {
	return x + y
}

result := add(five, ten)
! - / * 5
5 < 10 > 5

if 5 < 10 {
	return true
} else {
	return false
}

10 == 10
10 != 9
10 <= 10
10 >= 9
a && b || c
i++
i--
"hello world"
"hello \"world\""
[1, 2]
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},

		{token.IDENT, "ten"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "10.5"},

		{token.FUNC, "func"},
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

		{token.IDENT, "result"},
		{token.COLON_ASSIGN, ":="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},

		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},

		{token.INT, "10"},
		{token.LTE, "<="},
		{token.INT, "10"},

		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "9"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "c"},

		{token.IDENT, "i"},
		{token.INC, "++"},

		{token.IDENT, "i"},
		{token.DEC, "--"},

		{token.STRING, "hello world"},
		{token.STRING, "hello \"world\""},

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
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q)",
				i, tt.expectedType, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestDetailedTokens(t *testing.T) {
	input := `x := 5 + 3.14
// comment
/* multi
line */
if x >= 5 && x <= 10 || false {
	println("yes")
}
`
	l := New(input)
	expected := []struct {
		tokType token.TokenType
		literal string
	}{
		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "5"},
		{token.PLUS, "+"},
		{token.FLOAT, "3.14"},
		{token.IF, "if"},
		{token.IDENT, "x"},
		{token.GTE, ">="},
		{token.INT, "5"},
		{token.AND, "&&"},
		{token.IDENT, "x"},
		{token.LTE, "<="},
		{token.INT, "10"},
		{token.OR, "||"},
		{token.FALSE, "false"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "yes"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	for i, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.tokType {
			t.Fatalf("test[%d] failed - type expected=%q, got=%q", i, tt.tokType, tok.Type)
		}
		if tok.Literal != tt.literal {
			t.Fatalf("test[%d] failed - literal expected=%q, got=%q", i, tt.literal, tok.Literal)
		}
	}
}
