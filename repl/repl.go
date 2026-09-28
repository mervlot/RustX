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

	fmt.Fprintln(out, "RustX Starter v0.1 REPL")
	fmt.Fprintln(out, "Type 'exit' or Ctrl+D to quit.")
	fmt.Fprint(out, PROMPT)

	var inputBuffer strings.Builder
	openBraces := 0
	openParens := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "exit" && inputBuffer.Len() == 0 {
			break
		}

		inputBuffer.WriteString(line)
		inputBuffer.WriteString("\n")

		// Count brace and paren depth for multiline support
		inString := false
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c == '"' && (i == 0 || line[i-1] != '\\') {
				inString = !inString
			}
			if !inString {
				if c == '{' {
					openBraces++
				} else if c == '}' {
					openBraces--
				} else if c == '(' {
					openParens++
				} else if c == ')' {
					openParens--
				}
			}
		}

		if openBraces > 0 || openParens > 0 {
			fmt.Fprint(out, MULTI_PROMPT)
			continue
		}

		fullInput := inputBuffer.String()
		inputBuffer.Reset()
		openBraces = 0
		openParens = 0

		if strings.TrimSpace(fullInput) == "" {
			fmt.Fprint(out, PROMPT)
			continue
		}

		l := lexer.New(fullInput)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Fprintln(out, err.Error())
			}
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
