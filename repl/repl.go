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
	fmt.Fprintf(out, "Type 'exit()' to quit\n")

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		// Handle multi-line input (simplistic: check for open braces)
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "(") {
			for {
				fmt.Fprintf(out, "... ")
				if !scanner.Scan() {
					break
				}
				nextLine := scanner.Text()
				line += "\n" + nextLine
				if !strings.HasSuffix(nextLine, "{") && !strings.HasSuffix(nextLine, "(") && (strings.HasSuffix(nextLine, "}") || strings.HasSuffix(nextLine, ")") || nextLine == "") {
					// This is a very basic heuristic for multi-line
					if countBraces(line) == 0 {
						break
					}
				}
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
		if evaluated != nil && evaluated.Type() != object.NULL_OBJ {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func countBraces(s string) int {
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
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
