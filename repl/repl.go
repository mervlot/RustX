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
const CONTINUE_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type 'exit' to quit\n")

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			break
		}

		fullInput := line
		// Basic multiline support: if the line ends with { or (, or has unbalanced braces
		for shouldContinue(fullInput) {
			fmt.Fprintf(out, CONTINUE_PROMPT)
			scanned := scanner.Scan()
			if !scanned {
				break
			}
			fullInput += "\n" + scanner.Text()
		}

		l := lexer.New(fullInput)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func shouldContinue(input string) bool {
	// Simple brace counting
	openBraces := strings.Count(input, "{")
	closeBraces := strings.Count(input, "}")
	if openBraces > closeBraces {
		return true
	}

	openParens := strings.Count(input, "(")
	closeParens := strings.Count(input, ")")
	if openParens > closeParens {
		return true
	}

	return false
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
