package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/username/rustx/errors"
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

	fmt.Fprintln(out, "RustX Starter REPL v0.1.0")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to leave.")
	fmt.Fprintln(out)

	var buffer strings.Builder
	inMultiLine := false

	for {
		if inMultiLine {
			fmt.Fprint(out, CONTINUATION_PROMPT)
		} else {
			fmt.Fprint(out, PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if !inMultiLine && trimmed == "exit" {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		code := buffer.String()

		if isCompleteBlock(code) {
			inMultiLine = false
			buffer.Reset()

			l := lexer.New(code)
			p := parser.New(l)
			program := p.ParseProgram()

			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		} else {
			inMultiLine = true
		}
	}
}

func isCompleteBlock(code string) bool {
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

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, err := range errs {
		fmt.Fprintln(out, err.Error())
	}
}
