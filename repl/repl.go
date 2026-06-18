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

	var buffer []string
	braceCount := 0
	parenCount := 0

	for {
		if braceCount == 0 && parenCount == 0 {
			fmt.Fprintf(out, PROMPT)
		} else {
			fmt.Fprintf(out, "... ")
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		buffer = append(buffer, line)

		// Simple multiline tracking
		braceCount += strings.Count(line, "{") - strings.Count(line, "}")
		parenCount += strings.Count(line, "(") - strings.Count(line, ")")

		if braceCount <= 0 && parenCount <= 0 {
			fullInput := strings.Join(buffer, "\n")
			l := lexer.New(fullInput)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				buffer = nil
				braceCount = 0
				parenCount = 0
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated != evaluator.NULL {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}

			buffer = nil
			braceCount = 0
			parenCount = 0
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
