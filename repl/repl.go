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
		trimmed := strings.TrimSpace(line)

		if trimmed == "exit" {
			break
		}

		buffer = append(buffer, line)
		fullInput := strings.Join(buffer, "\n")

		if !isBalanced(fullInput) {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		// Evaluate the complete buffer
		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			buffer = nil
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintf(out, "%s\n", evaluated.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}

		buffer = nil
		fmt.Fprint(out, PROMPT)
	}
}

func isBalanced(input string) bool {
	var braces []rune
	var parens []rune
	var brackets []rune
	inString := false
	escaped := false

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		if inString {
			if escaped {
				escaped = false
			} else if r == '\\' {
				escaped = true
			} else if r == '"' {
				inString = false
			}
			continue
		}

		if r == '"' {
			inString = true
			continue
		}

		// Ignore characters inside comments
		if r == '/' && i+1 < len(runes) && runes[i+1] == '/' {
			// Single line comment, skip until newline
			for i < len(runes) && runes[i] != '\n' {
				i++
			}
			continue
		}

		if r == '/' && i+1 < len(runes) && runes[i+1] == '*' {
			// Multiline comment, skip until */
			i += 2
			for i+1 < len(runes) && !(runes[i] == '*' && runes[i+1] == '/') {
				i++
			}
			i++
			continue
		}

		switch r {
		case '{':
			braces = append(braces, r)
		case '}':
			if len(braces) == 0 {
				return true // Let parser handle syntax errors
			}
			braces = braces[:len(braces)-1]
		case '(':
			parens = append(parens, r)
		case ')':
			if len(parens) == 0 {
				return true
			}
			parens = parens[:len(parens)-1]
		case '[':
			brackets = append(brackets, r)
		case ']':
			if len(brackets) == 0 {
				return true
			}
			brackets = brackets[:len(brackets)-1]
		}
	}

	return len(braces) == 0 && len(parens) == 0 && len(brackets) == 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
