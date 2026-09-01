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
const CONTINUATION_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder
	depthBraces := 0
	depthParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		// Track block/paren balance ignoring characters inside quotes
		inString := false
		for i := 0; i < len(line); i++ {
			ch := line[i]
			if ch == '"' && (i == 0 || line[i-1] != '\\') {
				inString = !inString
			}
			if !inString {
				if ch == '{' {
					depthBraces++
				} else if ch == '}' && depthBraces > 0 {
					depthBraces--
				} else if ch == '(' {
					depthParens++
				} else if ch == ')' && depthParens > 0 {
					depthParens--
				}
			}
		}

		if depthBraces > 0 || depthParens > 0 {
			fmt.Fprint(out, CONTINUATION_PROMPT)
			continue
		}

		input := buffer.String()
		buffer.Reset()
		depthBraces = 0
		depthParens = 0

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			fmt.Fprint(out, PROMPT)
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			fmt.Fprintln(out, evaluated.Inspect())
		}

		fmt.Fprint(out, PROMPT)
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
