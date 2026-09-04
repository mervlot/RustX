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
const CONT_PROMPT = "... "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer strings.Builder

	for {
		if buffer.Len() == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONT_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" && buffer.Len() == 0 {
			return
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		code := buffer.String()
		if isCompleteCode(code) {
			l := lexer.New(code)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				buffer.Reset()
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}

			buffer.Reset()
		}
	}
}

func isCompleteCode(code string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(code); i++ {
		ch := code[i]
		if ch == '"' && (i == 0 || code[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
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

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintln(out, err.Error())
	}
}
