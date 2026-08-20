package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
)

const PROMPT = ">>> "
const CONTINUATION_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer string
	openBraces := 0
	openParens := 0

	for {
		if buffer == "" {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if buffer == "" && line == "exit" {
			return
		}

		buffer += line + "\n"

		// Count brace and paren balances
		for _, ch := range line {
			switch ch {
			case '{':
				openBraces++
			case '}':
				openBraces--
			case '(':
				openParens++
			case ')':
				openParens--
			}
		}

		// Continue reading if unbalanced
		if openBraces > 0 || openParens > 0 {
			continue
		}

		l := lexer.New(buffer)
		p := parser.New(l)

		program := p.ParseProgram()
		buffer = ""
		openBraces = 0
		openParens = 0

		if len(p.Errors()) != 0 {
			printParserErrors(out, p)
			continue
		}

		evaluated := evaluator.EvalProgram(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, p *parser.Parser) {
	for _, err := range p.Errors() {
		fmt.Fprintln(out, err.Error())
	}
}
