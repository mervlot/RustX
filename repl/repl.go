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

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type 'exit' to quit.\n")

	var buffer []string

	for {
		if len(buffer) == 0 {
			fmt.Fprintf(out, PROMPT)
		} else {
			fmt.Fprintf(out, "... ")
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" && len(buffer) == 0 {
			return
		}

		buffer = append(buffer, line)
		fullInput := strings.Join(buffer, "\n")

		if !isBalanced(fullInput) {
			continue
		}

		l := lexer.New(fullInput)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			buffer = nil
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
		buffer = nil
	}
}

func isBalanced(input string) bool {
	openBraces := strings.Count(input, "{")
	closeBraces := strings.Count(input, "}")
	openParens := strings.Count(input, "(")
	closeParens := strings.Count(input, ")")
	return openBraces <= closeBraces && openParens <= closeParens
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
