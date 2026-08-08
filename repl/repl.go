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

	fmt.Fprintln(out, "RustX Starter REPL v0.1")
	fmt.Fprintln(out, "Type 'exit' or press Ctrl+D to exit.")
	fmt.Fprint(out, PROMPT)

	var buffer []string
	openBrackets := 0

	for scanner.Scan() {
		line := scanner.Text()

		if len(buffer) == 0 && (strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "exit()") {
			break
		}

		buffer = append(buffer, line)
		openBrackets = calculateBalance(strings.Join(buffer, "\n"))

		if openBrackets > 0 {
			fmt.Fprint(out, CONT_PROMPT)
			continue
		}

		fullInput := strings.Join(buffer, "\n")
		buffer = nil // reset

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
			fmt.Fprintln(out, evaluated.Inspect())
		}

		fmt.Fprint(out, PROMPT)
	}
}

func calculateBalance(text string) int {
	balance := 0
	inString := false
	var stringChar rune

	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		char := runes[i]
		if inString {
			if char == '\\' && i+1 < len(runes) {
				i++
			} else if char == stringChar {
				inString = false
			}
			continue
		}

		if char == '"' {
			inString = true
			stringChar = char
			continue
		}

		if char == '{' || char == '(' || char == '[' {
			balance++
		} else if char == '}' || char == ')' || char == ']' {
			balance--
		}
	}
	return balance
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintln(out, msg)
	}
}
