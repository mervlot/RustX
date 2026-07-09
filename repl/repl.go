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
	fmt.Fprintf(out, "Type 'exit' to quit\n")

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		if line == "" {
			continue
		}

		// Basic multiline support: if line ends with '{' or '(', keep reading
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "(") {
			for {
				fmt.Fprintf(out, "... ")
				if !scanner.Scan() {
					break
				}
				nextLine := scanner.Text()
				line += "\n" + nextLine
				if strings.HasSuffix(nextLine, "}") || strings.HasSuffix(nextLine, ")") {
					// This is very naive, but okay for a starter.
					// A better way would be counting braces/parens.
					break
				}
			}
		}

		l := lexer.New(line)
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

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
