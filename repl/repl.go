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
	fmt.Fprintln(out, "Type some code or 'exit' to quit.")
	fmt.Fprint(out, PROMPT)

	var buffer []string

	for scanner.Scan() {
		line := scanner.Text()

		if len(buffer) == 0 && (strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit") {
			break
		}

		buffer = append(buffer, line)
		fullInput := strings.Join(buffer, "\n")

		if !isInputComplete(fullInput) {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		// Input is complete, evaluate it!
		buffer = nil // reset buffer

		l := lexer.New(fullInput, "<repl>")
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		fmt.Fprint(out, PROMPT)
	}
}

func isInputComplete(input string) bool {
	inString := false
	var stack []rune
	for i := 0; i < len(input); i++ {
		char := rune(input[i])
		if char == '"' {
			if i == 0 || input[i-1] != '\\' {
				inString = !inString
			}
		}
		if inString {
			continue
		}
		switch char {
		case '{', '(', '[':
			stack = append(stack, char)
		case '}':
			if len(stack) == 0 || stack[len(stack)-1] != '{' {
				return true // let the parser handle unmatched braces
			}
			stack = stack[:len(stack)-1]
		case ')':
			if len(stack) == 0 || stack[len(stack)-1] != '(' {
				return true
			}
			stack = stack[:len(stack)-1]
		case ']':
			if len(stack) == 0 || stack[len(stack)-1] != '[' {
				return true
			}
			stack = stack[:len(stack)-1]
		}
	}
	return len(stack) == 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
