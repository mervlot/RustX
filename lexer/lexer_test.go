package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken_Basic(t *testing.T) {
	input := `
var age = 16
name := "Muhammad"
var pi = 3.14

func add(a int, b int) int {
    return a + b
}

if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

for i := 0; i < 10; i++ {
    println(i)
}

while age > 0 {
    age--
}

// this is a single line comment
/*
  this is a
  multiline comment
*/
x := 10 % 3
b := true && !false || (5 == 5) != (2 <= 3)
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.INT, "16"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Muhammad"},

		{token.VAR, "var"},
		{token.IDENT, "pi"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3.14"},

		{token.FUNCTION, "func"},
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

		{token.IF, "if"},
		{token.IDENT, "age"},
		{token.GTE, ">="},
		{token.INT, "18"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "Adult"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.STRING, "Minor"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},

		{token.FOR, "for"},
		{token.IDENT, "i"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.PLUS_PLUS, "++"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.RBRACE, "}"},

		{token.WHILE, "while"},
		{token.IDENT, "age"},
		{token.GT, ">"},
		{token.INT, "0"},
		{token.LBRACE, "{"},
		{token.IDENT, "age"},
		{token.MINUS_MINUS, "--"},
		{token.RBRACE, "}"},

		{token.IDENT, "x"},
		{token.COLON_ASSIGN, ":="},
		{token.INT, "10"},
		{token.PERCENT, "%"},
		{token.INT, "3"},

		{token.IDENT, "b"},
		{token.COLON_ASSIGN, ":="},
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.BANG, "!"},
		{token.FALSE, "false"},
		{token.OR, "||"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.NOT_EQ, "!="},
		{token.LPAREN, "("},
		{token.INT, "2"},
		{token.LTE, "<="},
		{token.INT, "3"},
		{token.RPAREN, ")"},

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

func TestNextToken_EscapeStrings(t *testing.T) {
	input := `"hello\nworld\t\"quoted\"\\"`
	l := New(input)
	tok := l.NextToken()
	expected := "hello\nworld\t\"quoted\"\\"
	if tok.Literal != expected {
		t.Fatalf("escape sequence failed. expected=%q, got=%q", expected, tok.Literal)
	}
}
