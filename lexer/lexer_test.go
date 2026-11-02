package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5
ten := 10.5
func add(x int, y int) int {
	return x + y
}

// this is a comment
/* multiline
comment */
if (5 < 10) {
	return true
} else {
	return false
}

10 == 10
10 != 9
1 <= 2
2 >= 1
x && y
x || y
x++
y--
"hello \"world\"\n"
[1, 2]
while
for
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
		{token.FUNCTION, "func"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.TYPE_INT, "int"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.TYPE_INT, "int"},
		{token.RPAREN, ")"},
		{token.TYPE_INT, "int"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
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
		{token.INT, "1"},
		{token.LTE, "<="},
		{token.INT, "2"},
		{token.INT, "2"},
		{token.GTE, ">="},
		{token.INT, "1"},
		{token.IDENT, "x"},
		{token.AND, "&&"},
		{token.IDENT, "y"},
		{token.IDENT, "x"},
		{token.OR, "||"},
		{token.IDENT, "y"},
		{token.IDENT, "x"},
		{token.PLUS_PLUS, "++"},
		{token.IDENT, "y"},
		{token.MINUS_MINUS, "--"},
		{token.STRING, "hello \"world\"\n"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.WHILE, "while"},
		{token.FOR, "for"},
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

func TestPositionTracking(t *testing.T) {
	input := "x :=\n  42"
	l := New(input)

	tok1 := l.NextToken()
	if tok1.Position.Line != 1 || tok1.Position.Column != 1 {
		t.Fatalf("tok1 position wrong: line=%d, col=%d", tok1.Position.Line, tok1.Position.Column)
	}

	tok2 := l.NextToken()
	if tok2.Position.Line != 1 || tok2.Position.Column != 3 {
		t.Fatalf("tok2 position wrong: line=%d, col=%d", tok2.Position.Line, tok2.Position.Column)
	}

	tok3 := l.NextToken()
	if tok3.Position.Line != 2 || tok3.Position.Column != 3 {
		t.Fatalf("tok3 position wrong: line=%d, col=%d", tok3.Position.Line, tok3.Position.Column)
	}
}
