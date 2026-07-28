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
	fmt.Fprint(out, PROMPT)

	var buffer []string
	braceCount := 0
	parenCount := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Track braces and parentheses to decide if we should continue multiline reading
		for i := 0; i < len(line); i++ {
			char := line[i]
			if char == '{' {
				braceCount++
			} else if char == '}' {
				braceCount--
			} else if char == '(' {
				parenCount++
			} else if char == ')' {
				parenCount--
			}
		}

		buffer = append(buffer, line)

		if braceCount > 0 || parenCount > 0 {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		// Evaluate accumulated buffer once all blocks are closed
		fullInput := strings.Join(buffer, "\n")
		buffer = nil // Reset buffer

		l := lexer.New(fullInput)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}

		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
