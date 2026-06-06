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

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	fmt.Fprintf(out, "RustX Starter v0.1\n")
	fmt.Fprintf(out, "Type 'exit' to quit\n")

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" {
			return
		}

		// Basic multiline support: if line ends with '{' or '(', keep reading
		for strings.HasSuffix(strings.TrimSpace(line), "{") || strings.HasSuffix(strings.TrimSpace(line), "(") || countOpenBraces(line) > 0 {
			fmt.Fprintf(out, "... ")
			if !scanner.Scan() {
				break
			}
			nextLine := scanner.Text()
			line += "\n" + nextLine
			if nextLine == "" && countOpenBraces(line) == 0 {
				break
			}
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil && evaluated != evaluator.NULL {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func countOpenBraces(s string) int {
	count := 0
	for _, char := range s {
		if char == '{' || char == '(' {
			count++
		} else if char == '}' || char == ')' {
			count--
		}
	}
	return count
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "SyntaxError:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
