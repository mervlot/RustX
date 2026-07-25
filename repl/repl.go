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
	fmt.Fprint(out, PROMPT)

	var buffer []string
	openCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		openCount += countOpen(line)
		buffer = append(buffer, line)

		if openCount <= 0 {
			fullInput := strings.Join(buffer, "\n")
			buffer = nil
			openCount = 0

			l := lexer.New(fullInput)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) > 0 {
				printParserErrors(out, p.Errors())
				fmt.Fprint(out, PROMPT)
				continue
			}

			evaluated := evaluator.Eval(program, env)
			if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
				if errObj, ok := evaluated.(*object.Error); ok {
					fmt.Fprintf(out, "RuntimeError:\n%s\n", errObj.Message)
				} else {
					fmt.Fprintln(out, evaluated.Inspect())
				}
			}
			fmt.Fprint(out, PROMPT)
		} else {
			fmt.Fprint(out, MULTI_PROMPT)
		}
	}
}

func countOpen(line string) int {
	open := 0
	inString := false
	var stringChar rune
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if inString {
			if r == '\\' && i+1 < len(runes) {
				i++
			} else if r == rune(stringChar) {
				inString = false
			}
		} else {
			if r == '"' {
				inString = true
				stringChar = '"'
			} else if r == '/' && i+1 < len(runes) && runes[i+1] == '/' {
				break
			} else if r == '{' || r == '(' || r == '[' {
				open++
			} else if r == '}' || r == ')' || r == ']' {
				open--
			}
		}
	}
	return open
}

func printParserErrors(out io.Writer, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(out, "%s\n", err.Error())
	}
}
