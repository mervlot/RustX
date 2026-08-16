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

	fmt.Fprintln(out, "RustX Starter v0.1.0")
	fmt.Fprintln(out, "Type 'exit' or Ctrl+C to quit.")

	var inputBuffer strings.Builder
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
		if depth == 0 && strings.TrimSpace(line) == "exit" {
			return
		}

		inputBuffer.WriteString(line + "\n")
		bufStr := inputBuffer.String()

		depth = calculateDepth(bufStr)
		if depth > 0 {
			continue
		}

		// Process complete statement or block
		l := lexer.New(bufStr)
		p := parser.New(l)

		program := p.ParseProgram()
		inputBuffer.Reset()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintf(out, "%s\n", evaluated.Inspect())
			} else {
				fmt.Fprintf(out, "%s\n", evaluated.Inspect())
			}
		}
	}
}

func calculateDepth(input string) int {
	depth := 0
	inString := false
	var stringChar byte

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == stringChar && (i == 0 || input[i-1] != '\\') {
				inString = false
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			inString = true
			stringChar = ch
			continue
		}
		if ch == '{' || ch == '(' {
			depth++
		} else if ch == '}' || ch == ')' {
			depth--
		}
	}

	if depth < 0 {
		return 0
	}
	return depth
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
