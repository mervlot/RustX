package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

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

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type commands or 'exit' to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "exit" {
			break
		}

		buffer = append(buffer, line)
		fullInput := strings.Join(buffer, "\n")

		// Track balanced braces and parentheses to decide on multiline evaluation
		if !isBalanced(fullInput) {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			buffer = nil
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		buffer = nil
		fmt.Fprint(out, PROMPT)
	}
}

func isBalanced(s string) bool {
	braces := 0
	parens := 0
	brackets := 0
	inString := false

	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' {
			if i > 0 && s[i-1] == '\\' {
				// escaped quote, ignore
			} else {
				inString = !inString
			}
		}
		if inString {
			continue
		}

		// Skip comments
		if ch == '/' && i+1 < len(s) && s[i+1] == '/' {
			break // rest of line is a comment
		}

		switch ch {
		case '{':
			braces++
		case '}':
			braces--
		case '(':
			parens++
		case ')':
			parens--
		case '[':
			brackets++
		case ']':
			brackets--
		}
	}

	return braces <= 0 && parens <= 0 && brackets <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
