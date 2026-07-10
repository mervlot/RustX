package repl

import (
	"bufio"
	"fmt"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/object"
	"github.com/username/rustx/parser"
	"io"
	"strings"
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

		// Basic multiline support: if line ends with {, keep reading until }
		if strings.HasSuffix(strings.TrimSpace(line), "{") || strings.HasSuffix(strings.TrimSpace(line), "(") {
			line = readMultiline(scanner, out, line)
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func readMultiline(scanner *bufio.Scanner, out io.Writer, firstLine string) string {
	var lines []string
	lines = append(lines, firstLine)

	openBraces := strings.Count(firstLine, "{") - strings.Count(firstLine, "}")
	openParens := strings.Count(firstLine, "(") - strings.Count(firstLine, ")")

	for openBraces > 0 || openParens > 0 {
		fmt.Fprintf(out, "... ")
		scanned := scanner.Scan()
		if !scanned {
			break
		}
		line := scanner.Text()
		lines = append(lines, line)
		openBraces += strings.Count(line, "{") - strings.Count(line, "}")
		openParens += strings.Count(line, "(") - strings.Count(line, ")")
	}

	return strings.Join(lines, "\n")
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
