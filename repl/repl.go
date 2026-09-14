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

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var buffer strings.Builder
	depth := 0

	for {
		if depth == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if depth == 0 && trimmed == "exit" {
			return
		}

		buffer.WriteString(line + "\n")
		depth = countIndentationDepth(buffer.String())

		if depth <= 0 {
			input := buffer.String()
			buffer.Reset()
			depth = 0

			l := lexer.New(input)
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
		}
	}
}

func countIndentationDepth(input string) int {
	depth := 0
	inString := false
	var stringQuote byte

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '\\' {
				i++ // skip escaped char
			} else if ch == stringQuote {
				inString = false
			}
			continue
		}

		if ch == '"' {
			inString = true
			stringQuote = ch
			continue
		}

		if ch == '{' || ch == '(' {
			depth++
		} else if ch == '}' || ch == ')' {
			depth--
		}
	}
	return depth
}

func printParserErrors(out io.Writer, errs []string) {
	for _, msg := range errs {
		fmt.Fprintln(out, msg)
	}
}
