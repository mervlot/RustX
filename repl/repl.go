package repl

import (
	"bufio"
	"fmt"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"io"
	"strings"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "quit" {
			break
		}

		// Basic multiline support: if line ends with '{' or '(' or has unmatched braces, keep reading
		for !isBalanced(line) {
			fmt.Fprintf(out, "... ")
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
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func isBalanced(input string) bool {
	lbraces := strings.Count(input, "{")
	rbraces := strings.Count(input, "}")
	lparens := strings.Count(input, "(")
	rparens := strings.Count(input, ")")
	return lbraces <= rbraces && lparens <= rparens
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
