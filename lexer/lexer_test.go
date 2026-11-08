package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextTokenBasic(t *testing.T) {
	input := `var five = 5
ten := 10.5
func add(x int, y int) int {
    return x + y
}
result := add(five, ten)
!- / * 5;
5 < 10 > 5
if 5 < 10 {
    return true
} else {
    return false
}
10 == 10
10 != 9
x <= 10
y >= 20
a && b || c
i++
i--
"hello world"
"hello \"world\"\n"
// comment line
/* multiline
comment */
5 % 2
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
		{token.SEMICOLON, ";"},

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

		{token.IDENT, "x"},
		{token.LE, "<="},
		{token.INT, "10"},

		{token.IDENT, "y"},
		{token.GE, ">="},
		{token.INT, "20"},

		{token.IDENT, "a"},
		{token.AND, "&&"},
		{token.IDENT, "b"},
		{token.OR, "||"},
		{token.IDENT, "c"},

		{token.IDENT, "i"},
		{token.INCREMENT, "++"},

		{token.IDENT, "i"},
		{token.DECREMENT, "--"},

		{token.STRING, "hello world"},
		{token.STRING, "hello \"world\"\n"},

		{token.INT, "5"},
		{token.MODULO, "%"},
		{token.INT, "2"},

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

func TestUnterminatedString(t *testing.T) {
	input := `"unterminated string`
	l := New(input)
	tok := l.NextToken()
	if tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL token, got %q", tok.Type)
	}
}
