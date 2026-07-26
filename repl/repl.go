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
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprint(out, PROMPT)

	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()
		buffer = append(buffer, line)

		fullInput := strings.Join(buffer, "\n")
		if isBalanced(fullInput) {
			l := lexer.New(fullInput)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				evaluated := evaluator.Eval(program, env)
				if evaluated != nil && evaluated != evaluator.NULL {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			buffer = nil
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}
	}
}

func isBalanced(input string) bool {
	inString := false
	var escape bool
	var braces int
	var parens int
	var brackets int

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if escape {
				escape = false
			} else if ch == '\\' {
				escape = true
			} else if ch == '"' {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			escape = false
			continue
		}

		if ch == '/' && i+1 < len(input) && input[i+1] == '/' {
			for i < len(input) && input[i] != '\n' {
				i++
			}
			continue
		}
		if ch == '/' && i+1 < len(input) && input[i+1] == '*' {
			i += 2
			for i+1 < len(input) && !(input[i] == '*' && input[i+1] == '/') {
				i++
			}
			i++
			continue
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
