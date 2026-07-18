package lexer

import (
	"testing"

	"github.com/username/rustx/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
let ten = 10; // note: let is fine as an identifier, var/:= is used for vars
func add(x, y) {
    return x + y;
}

result := add(five, ten);
!- / * 5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

10 == 10;
10 != 9;
"foobar"
"foo bar"
[1, 2];
/* multi line
   comment */
x++
y--
5 <= 10 && 10 >= 5 || true;
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedLine    int
		expectedCol     int
	}{
		{token.VAR, "var", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		{token.SEMICOLON, ";", 1, 13},

		{token.IDENT, "let", 2, 1},
		{token.IDENT, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "10", 2, 11},
		{token.SEMICOLON, ";", 2, 13},

		{token.FUNC, "func", 3, 1},
		{token.IDENT, "add", 3, 6},
		{token.LPAREN, "(", 3, 9},
		{token.IDENT, "x", 3, 10},
		{token.COMMA, ",", 3, 11},
		{token.IDENT, "y", 3, 13},
		{token.RPAREN, ")", 3, 14},
		{token.LBRACE, "{", 3, 16},

		{token.RETURN, "return", 4, 5},
		{token.IDENT, "x", 4, 12},
		{token.PLUS, "+", 4, 14},
		{token.IDENT, "y", 4, 16},
		{token.SEMICOLON, ";", 4, 17},

		{token.RBRACE, "}", 5, 1},

		{token.IDENT, "result", 7, 1},
		{token.COLON_ASSIGN, ":=", 7, 8},
		{token.IDENT, "add", 7, 11},
		{token.LPAREN, "(", 7, 14},
		{token.IDENT, "five", 7, 15},
		{token.COMMA, ",", 7, 19},
		{token.IDENT, "ten", 7, 21},
		{token.RPAREN, ")", 7, 24},
		{token.SEMICOLON, ";", 7, 25},

		{token.BANG, "!", 8, 1},
		{token.MINUS, "-", 8, 2},
		{token.SLASH, "/", 8, 4},
		{token.ASTERISK, "*", 8, 6},
		{token.INT, "5", 8, 8},
		{token.SEMICOLON, ";", 8, 9},

		{token.INT, "5", 9, 1},
		{token.LT, "<", 9, 3},
		{token.INT, "10", 9, 5},
		{token.GT, ">", 9, 8},
		{token.INT, "5", 9, 10},
		{token.SEMICOLON, ";", 9, 11},

		{token.IF, "if", 11, 1},
		{token.LPAREN, "(", 11, 4},
		{token.INT, "5", 11, 5},
		{token.LT, "<", 11, 7},
		{token.INT, "10", 11, 9},
		{token.RPAREN, ")", 11, 11},
		{token.LBRACE, "{", 11, 13},

		{token.RETURN, "return", 12, 5},
		{token.TRUE, "true", 12, 12},
		{token.SEMICOLON, ";", 12, 16},

		{token.RBRACE, "}", 13, 1},
		{token.ELSE, "else", 13, 3},
		{token.LBRACE, "{", 13, 8},

		{token.RETURN, "return", 14, 5},
		{token.FALSE, "false", 14, 12},
		{token.SEMICOLON, ";", 14, 17},

		{token.RBRACE, "}", 15, 1},

		{token.INT, "10", 17, 1},
		{token.EQ, "==", 17, 4},
		{token.INT, "10", 17, 7},
		{token.SEMICOLON, ";", 17, 9},

		{token.INT, "10", 18, 1},
		{token.NOT_EQ, "!=", 18, 4},
		{token.INT, "9", 18, 7},
		{token.SEMICOLON, ";", 18, 8},

		{token.STRING, "foobar", 19, 1},
		{token.STRING, "foo bar", 20, 1},

		{token.LBRACKET, "[", 21, 1},
		{token.INT, "1", 21, 2},
		{token.COMMA, ",", 21, 3},
		{token.INT, "2", 21, 5},
		{token.RBRACKET, "]", 21, 6},
		{token.SEMICOLON, ";", 21, 7},

		{token.IDENT, "x", 24, 1},
		{token.INC, "++", 24, 2},

		{token.IDENT, "y", 25, 1},
		{token.DEC, "--", 25, 2},

		{token.INT, "5", 26, 1},
		{token.LTE, "<=", 26, 3},
		{token.INT, "10", 26, 6},
		{token.AND, "&&", 26, 9},
		{token.INT, "10", 26, 12},
		{token.GTE, ">=", 26, 15},
		{token.INT, "5", 26, 18},
		{token.OR, "||", 26, 20},
		{token.TRUE, "true", 26, 23},
		{token.SEMICOLON, ";", 26, 27},

		{token.EOF, "", 27, 1},
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

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}

		if tok.Col != tt.expectedCol {
			t.Fatalf("tests[%d] - column wrong for literal %q. expected=%d, got=%d",
				i, tt.expectedLiteral, tt.expectedCol, tok.Col)
		}
	}
}
