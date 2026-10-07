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

	fmt.Fprintln(out, "RustX Starter v0.1.0")
	fmt.Fprintln(out, "Type \"exit\" or press Ctrl+C / Ctrl+D to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "exit" && buffer.Len() == 0 {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		input := buffer.String()

		if isCompleteInput(input) {
			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				for _, err := range p.Errors() {
					fmt.Fprintln(out, err.Error())
				}
				buffer.Reset()
				fmt.Fprint(out, PROMPT)
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil {
				if errObj, ok := evaluated.(*object.Error); ok {
					fmt.Fprintln(out, errObj.Inspect())
				} else if evaluated.Type() != object.NULL_OBJ {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}

			buffer.Reset()
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		}
	}
}

func isCompleteInput(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false
	escaped := false

	for i := 0; i < len(input); i++ {
		ch := input[i]

		if escaped {
			escaped = false
			continue
		}

		if ch == '\\' && inString {
			escaped = true
			continue
		}

		if ch == '"' {
			inString = !inString
			continue
		}

		if !inString {
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
	}

	return openBraces <= 0 && openParens <= 0
}
