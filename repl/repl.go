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

	fmt.Fprintln(out, "RustX Starter v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+C to quit.")

	var inputBuffer strings.Builder
	depth := 0

	for {
		if depth == 0 {
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, CONTINUATION_PROMPT)
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
		depth = countUnclosedBrackets(inputBuffer.String())

		if depth > 0 {
			continue
		}

		code := inputBuffer.String()
		inputBuffer.Reset()

		l := lexer.New(code)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			if evaluated.Type() == object.ERROR_OBJ {
				fmt.Fprintln(out, evaluated.Inspect())
			} else {
				fmt.Fprintln(out, evaluated.Inspect())
			}
		}
	}
}

func countUnclosedBrackets(input string) int {
	depth := 0
	inString := false
	var stringChar byte

	for i := 0; i < len(input); i++ {
		ch := input[i]
		if inString {
			if ch == '\\' {
				i++ // skip next char
			} else if ch == stringChar {
				inString = false
			}
		} else {
			if ch == '"' || ch == '\'' {
				inString = true
				stringChar = ch
			} else if ch == '{' || ch == '(' || ch == '[' {
				depth++
			} else if ch == '}' || ch == ')' || ch == ']' {
				depth--
			}
		}
	}

	if depth < 0 {
		return 0
	}
	return depth
}

func printParserErrors(out io.Writer, errs []*errors.SyntaxError) {
	for _, msg := range errs {
		fmt.Fprintln(out, msg.Error())
	}
}
