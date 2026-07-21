# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility. This project implements a fully working Go-based tree-walking interpreter (evaluator), recursive-descent parser with Pratt expression parsing, lexer, REPL, and a complete command line interface.

## Architecture

The project has been written following clean and modular compiler architecture guidelines:

```
├── main.go             # CLI Entry Point
├── token/              # Token types, definitions, and positional info
├── errors/             # Position-aware Syntax & Runtime custom errors
├── lexer/              # Lexical scanner & tokenizer (supports strings, escapes, comments, etc.)
├── ast/                # Extensible Abstract Syntax Tree (AST) node definitions
├── parser/             # Recursive-descent & Pratt Expression parser
├── object/             # Object-based runtime and Lexical Environments with block scoping
├── evaluator/          # Tree-walking interpreter and standard library built-ins
├── repl/               # Interactive REPL with multiline balance checking
└── samples/            # Example RustX programs
```

## Features

1. **Variables**: Mutability by default, declared with `var` or `:=` shorthand assignment.
2. **Data Types**: `int`, `float`, `bool`, `string`, `Array`.
3. **Operators**: Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), and logical (`&&`, `||`) with short-circuiting.
4. **Control Flow**: Unparenthesized `if-else` blocks, `while` loops, and Go-style three-part `for` loops.
5. **Functions**: Named and anonymous functions, recursive call support, parameters, and optional type annotations in signatures.
6. **Comments**: Support for single-line (`//`) and multi-line (`/* */`) comments.
7. **Basic Standard Library**:
   - Printing: `print()`, `println()` (space-separated with newline).
   - Math: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`.
   - Strings: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`.
   - Time: `sleep()`, `now()`.
   - Filesystem: `read_file()`, `write_file()`.

## Build and Installation

To build the executable, run:
```bash
go build -o rustx main.go
```

## Usage

### Run a Script
To run a `.rx` file:
```bash
./rustx run samples/hello.rx
```

### Shorthand Execution
You can also execute `.rx` files directly without the `run` command:
```bash
./rustx samples/hello.rx
```

### Start interactive REPL
To start the interactive prompt:
```bash
./rustx repl
```
Or simply run `./rustx` with no arguments.

The REPL supports multiline buffers by automatically detecting unclosed parenthesis `(` and braces `{`.

### Show Version
```bash
./rustx version
```

## Running Tests

To run all package unit tests:
```bash
go test ./...
```
