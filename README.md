# RustX Starter

A minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility. This repository contains the foundation-level implementation in Go (interpreter, parser, lexer, REPL, and CLI commands).

## Core Language Design & Syntax Style
RustX Starter blends elements of Rust and Go syntax:
- `func` for function definitions.
- `var` for mutable variable declarations.
- `:=` shorthand block scope variable assignment.
- Optional semicolons.
- Curly braces for scoping.
- Optional parentheses for control flow conditions (e.g. `if condition { ... }`).

## Features
1. **Variables & Scope**: Supports recursive outer-scope variable resolution and updates.
2. **Operators**: Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical AND (`&&`) and OR (`||`) with short-circuiting, postfix operators (`++`, `--`).
3. **Data Types**: `int`, `float`, `bool`, `string`, and basic `Array` types.
4. **Functions**: Local block scope parameters, return values, and recursion.
5. **Control Flow**: `if`, `else`, `while`, and 3-part Go-style `for` loops (as well as infinite and single-condition loops).
6. **Standard Library**: Core built-ins including printing (`print`/`println`), trigonometry/math (`abs`, `pow`, `sin`, `cos`, `tan`), string manipulation (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), time manipulation (`sleep`, `now`), and basic safe file I/O (`read_file`, `write_file`).

## Architecture Details
The project is split into clean, decoupled, and easily extensible Go packages:
- `/errors`: Standardized diagnostic reporting for `SyntaxError` and `RuntimeError`.
- `/token`: Keywords, delimiters, operators, and location tracking (line/column).
- `/lexer`: A high-performance, single/multi-line-comment-skipping lexical scanner.
- `/ast`: Extensible Abstract Syntax Tree node representations.
- `/parser`: A powerful hybrid Recursive Descent + Pratt Parser.
- `/object`: Object runtime and scoped variable evaluation context.
- `/evaluator`: Tree-walking interpreter with short-circuiting, division-by-zero protection, and standard built-ins.
- `/repl`: Multiline balanced-brace tracking interactive shell.

## Build and Run Instructions

### Build CLI
To compile the `rustx` compiler:
```bash
go build -o rustx main.go
```

### Run Scripts
Execute `.rx` scripts directly or using the run command:
```bash
./rustx samples/factorial.rx
./rustx run samples/factorial.rx
```

### Start Interactive REPL
Start the interactive command line environment:
```bash
./rustx repl
# OR simply
./rustx
```

### Run All Unit Tests
Run the entire Go test suite:
```bash
go test ./... -v
```
