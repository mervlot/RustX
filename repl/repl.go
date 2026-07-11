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

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type code and press Enter to evaluate.\n")

	for {
		fmt.Fprintf(out, PROMPT)
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "quit" {
			break
		}

		// Basic multiline support: if line ends with '{' or '(', wait for more
		for strings.HasSuffix(strings.TrimSpace(line), "{") || strings.HasSuffix(strings.TrimSpace(line), "(") {
			fmt.Fprintf(out, MULTI_PROMPT)
			if !scanner.Scan() {
				break
			}
			line += "\n" + scanner.Text()
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
			if evaluated.Type() != object.NULL_OBJ {
				io.WriteString(out, evaluated.Inspect())
				io.WriteString(out, "\n")
			}
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
