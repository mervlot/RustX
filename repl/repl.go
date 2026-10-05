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

	fmt.Fprintf(out, "RustX Starter v0.1.0 Interactive REPL\n")
	fmt.Fprintf(out, "Type statements or expressions to evaluate. Type 'exit' to quit.\n\n")

	var buffer strings.Builder
	inMultiLine := false

	for {
		if inMultiLine {
			fmt.Fprintf(out, MULTI_PROMPT)
		} else {
			fmt.Fprintf(out, PROMPT)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if !inMultiLine && strings.TrimSpace(line) == "exit" {
			break
		}

		buffer.WriteString(line)
		buffer.WriteString("\n")

		inputStr := buffer.String()

		if isCompleteBlock(inputStr) {
			l := lexer.New(inputStr)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
				buffer.Reset()
				inMultiLine = false
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				fmt.Fprintf(out, "%s\n", evaluated.Inspect())
			}

			buffer.Reset()
			inMultiLine = false
		} else {
			inMultiLine = true
		}
	}
}

func isCompleteBlock(input string) bool {
	openBraces := 0
	openParens := 0
	inString := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if ch == '"' && (i == 0 || input[i-1] != '\\') {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if ch == '{' {
			openBraces++
		} else if ch == '}' {
			openBraces--
		} else if ch == '(' {
			openParens++
		} else if ch == ')' {
			openParens--
		}
	}

	return openBraces <= 0 && openParens <= 0
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "%s\n", msg)
	}
}
