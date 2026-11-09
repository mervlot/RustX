package repl

import (
	"bufio"
	"fmt"
	"io"
	"rustx/errors"
	"rustx/evaluator"
	"rustx/lexer"
	"rustx/object"
	"rustx/parser"
	"strings"
)

const PROMPT = ">>> "
const CONT_PROMPT = "... "

// Start launches the interactive REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	io.WriteString(out, "RustX Starter v0.1\n")
	io.WriteString(out, "Type 'exit' or Ctrl+D to quit.\n\n")

	var buffer strings.Builder
	braceDepth := 0
	parenDepth := 0

	for {
		if buffer.Len() == 0 {
			io.WriteString(out, PROMPT)
		} else {
			io.WriteString(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if buffer.Len() == 0 && (trimmed == "exit" || trimmed == "quit") {
			return
		}

		// Update block depths
		for _, ch := range line {
			switch ch {
			case '{':
				braceDepth++
			case '}':
				braceDepth--
			case '(':
				parenDepth++
			case ')':
				parenDepth--
			}
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// If multiline statement is open, keep reading lines
		if braceDepth > 0 || parenDepth > 0 {
			continue
		}

		input := buffer.String()
		buffer.Reset()
		braceDepth = 0
		parenDepth = 0

		if strings.TrimSpace(input) == "" {
			continue
		}

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			io.WriteString(out, evaluated.Inspect()+"\n")
		}
	}
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, msg := range errs {
		fmt.Fprintf(out, "%s\n", msg.Error())
	}
}
