package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
)

const PROMPT = ">>> "
const MULTI_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1 Interactive REPL")
	fmt.Fprintln(out, "Type statements or expressions to evaluate. Ctrl+C or Exit to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer string
	inMultiLine := false

	for scanner.Scan() {
		line := scanner.Text()

		if buffer == "" {
			buffer = line
		} else {
			buffer += "\n" + line
		}

		// Count brace balance and paren balance
		braceCount, parenCount := countBalance(buffer)
		if braceCount > 0 || parenCount > 0 {
			inMultiLine = true
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		inMultiLine = false
		l := lexer.New(buffer)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			buffer = ""
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		buffer = ""
		fmt.Fprint(out, PROMPT)
	}

	if inMultiLine {
		fmt.Fprintln(out)
	}
}

func countBalance(s string) (braceCount int, parenCount int) {
	inString := false
	var stringQuote byte
	escaped := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && inString {
			escaped = true
			continue
		}

		if (ch == '"' || ch == '\'') && !inString {
			inString = true
			stringQuote = ch
			continue
		} else if inString && ch == stringQuote {
			inString = false
			continue
		}

		if !inString {
			switch ch {
			case '{':
				braceCount++
			case '}':
				if braceCount > 0 {
					braceCount--
				}
			case '(':
				parenCount++
			case ')':
				if parenCount > 0 {
					parenCount--
				}
			}
		}
	}
	return braceCount, parenCount
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
