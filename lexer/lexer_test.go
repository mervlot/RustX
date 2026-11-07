package lexer

import (
	"rustx/token"
	"testing"
)

func TestNextToken_Basic(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
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

func TestNextToken_CodeSnippet(t *testing.T) {
	input := `
var age = 16
name := "Muhammad\n"
pi := 3.14

// This is a single line comment
/* This is a
   multiline comment */

func add(a int, b int) int {
    if a == b {
        return a + b
    } else if a != b && a >= 0 || a <= 10 {
        return a - b * 2 / 1 % 3
    }
}

for i := 0; i < 10; i++ {
    i--
}

while true {
    return false
}
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		// var age = 16
		{token.VAR, "var"},
		{token.IDENT, "age"},
		{token.ASSIGN, "="},
		{token.INT, "16"},

		// name := "Muhammad\n"
		{token.IDENT, "name"},
		{token.COLON_ASSIGN, ":="},
		{token.STRING, "Muhammad\n"},

		// pi := 3.14
		{token.IDENT, "pi"},
		{token.COLON_ASSIGN, ":="},
		{token.FLOAT, "3.14"},

		// func add(a int, b int) int {
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

		// if a == b {
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.EQ, "=="},
		{token.IDENT, "b"},
		{token.LBRACE, "{"},

		// return a + b
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},

		// } else if a != b && a >= 0 || a <= 10 {
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.IF, "if"},
		{token.IDENT, "a"},
		{token.NOT_EQ, "!="},
		{token.IDENT, "b"},
		{token.AND, "&&"},
		{token.IDENT, "a"},
		{token.GT_EQ, ">="},
		{token.INT, "0"},
		{token.OR, "||"},
		{token.IDENT, "a"},
		{token.LT_EQ, "<="},
		{token.INT, "10"},
		{token.LBRACE, "{"},

		// return a - b * 2 / 1 % 3
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.MINUS, "-"},
		{token.IDENT, "b"},
		{token.ASTERISK, "*"},
		{token.INT, "2"},
		{token.SLASH, "/"},
		{token.INT, "1"},
		{token.MODULO, "%"},
		{token.INT, "3"},

		// }
		{token.RBRACE, "}"},

		// }
		{token.RBRACE, "}"},

		// for i := 0; i < 10; i++ {
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

		// i--
		{token.IDENT, "i"},
		{token.MINUS_MINUS, "--"},

		// }
		{token.RBRACE, "}"},

		// while true {
		{token.WHILE, "while"},
		{token.TRUE, "true"},
		{token.LBRACE, "{"},

		// return false
		{token.RETURN, "return"},
		{token.FALSE, "false"},

		// }
		{token.RBRACE, "}"},

		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q (literal=%q, line=%d, col=%d)",
				i, tt.expectedType, tok.Type, tok.Literal, tok.Line, tok.Column)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLineAndColumnTracking(t *testing.T) {
	input := "var x = 1\nvar y = 2"
	l := New(input)

	tok1 := l.NextToken()
	if tok1.Line != 1 || tok1.Column != 1 {
		t.Errorf("expected line 1, col 1, got %d, %d", tok1.Line, tok1.Column)
	}

	// Read until second line
	for {
		tok := l.NextToken()
		if tok.Line == 2 && tok.Type == token.VAR {
			if tok.Column != 1 {
				t.Errorf("expected line 2, col 1 for second var, got col %d", tok.Column)
			}
			break
		}
		if tok.Type == token.EOF {
			t.Fatalf("EOF before line 2 var")
		}
	}
}
