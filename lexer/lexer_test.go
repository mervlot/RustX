package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `
	// A single-line comment
	/* A multi-line
	   comment */
	var age = 16;
	name := "Muhammad\n"; // shorthand assign and escape string
	func add(a int, b int) int {
		return a + b;
	}
	12;
	3.14;
	if age <= 20 && age >= 10 || name != "RustX" {
		age++;
	}
	while true {
		break;
	}
	for i := 0; i < 10; i++ {
		println(i);
	}
	[1, 2];
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.INT, "16"},
		{token.SEMICOLON, ";"},

		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Muhammad\n"},
		{token.SEMICOLON, ";"},

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
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "12"},
		{token.SEMICOLON, ";"},
		{token.FLOAT, "3.14"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.IDENT, "age"},
		{token.LTE, "<="},
		{token.INT, "20"},
		{token.AND, "&&"},
		{token.IDENT, "age"},
		{token.GTE, ">="},
		{token.INT, "10"},
		{token.OR, "||"},
		{token.IDENT, "name"},
		{token.NOT_EQ, "!="},
		{token.STRING, "RustX"},
		{token.LBRACE, "{"},
		{token.IDENT, "age"},
		{token.INC, "++"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},
		{token.IDENT, "break"},
		{token.SEMICOLON, ";"},
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
		{token.INC, "++"},
		{token.LBRACE, "{"},
		{token.IDENT, "println"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

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

func TestAdjacentSlashComment(t *testing.T) {
	// Division / * shouldn't be parsed as comment start if they are separated by whitespace.
	// But / followed immediately by * is a multi-line comment start.
	input := `a / * b`
	l := New(input)

	tok := l.NextToken()
	if tok.Type != token.IDENT || tok.Literal != "a" {
		t.Fatalf("expected a, got %s", tok.Literal)
	}
	tok = l.NextToken()
	if tok.Type != token.SLASH {
		t.Fatalf("expected /, got %s", tok.Literal)
	}
	tok = l.NextToken()
	if tok.Type != token.ASTERISK {
		t.Fatalf("expected *, got %s", tok.Literal)
	}
	tok = l.NextToken()
	if tok.Type != token.IDENT || tok.Literal != "b" {
		t.Fatalf("expected b, got %s", tok.Literal)
	}
}
