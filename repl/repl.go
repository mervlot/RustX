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
	fmt.Fprintf(out, "Type exit() or press Ctrl+C to quit\n")

	var buffer []string
	openBraces := 0

	for {
		if openBraces == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, "... ")
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit()" {
			return
		}

		buffer = append(buffer, line)
		openBraces += countBraces(line)

		if openBraces <= 0 {
			fullInput := strings.Join(buffer, "\n")
			l := lexer.New(fullInput)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				buffer = nil
				openBraces = 0
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
			buffer = nil
			openBraces = 0
		}
	}
}

func countBraces(line string) int {
	count := 0
	for _, char := range line {
		if char == '{' || char == '(' {
			count++
		} else if char == '}' || char == ')' {
			count--
		}
	}
	return count
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
