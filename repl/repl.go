package repl

import (
	"bufio"
	"fmt"
	"io"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/parser"
	"github.com/username/rustx/evaluator"
	"github.com/username/rustx/object"
	"strings"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == "exit" || line == "quit" {
			return
		}

		// Simple multiline support: if line ends with { or (, or has unbalanced ones, continue reading
		for (strings.Count(line, "{") > strings.Count(line, "}")) || (strings.Count(line, "(") > strings.Count(line, ")")) {
			fmt.Fprint(out, "... ")
			if !scanner.Scan() {
				break
			}
			line += "\n" + scanner.Text()
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

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, msg+"\n")
	}
}
