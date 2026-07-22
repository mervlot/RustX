# RustX Starter

A minimal, clean, and scalable programming language inspired by Rust and Go. Focused on simplicity, speed, readability, and future extensibility.

Initially implemented as a lightweight scripting and tree-walking interpreter.

## Syntax Style
- `func` for functions
- `var` for mutable variable declarations
- `:=` shorthand assignment
- Optional semicolons
- Block scoping and curly braces
- Unparenthesized condition clauses and loop headers

## Architecture & Module Layout
The interpreter and standard runtime library are fully designed using clean Go:
- `/errors`: Custom error types for position-aware error diagnostics reporting.
- `/token`: Lexical tokens.
- `/lexer`: Lexical analyzer.
- `/ast`: Position-aware AST nodes.
- `/parser`: Pratt-parsing based recursive descent parser.
- `/object`: Object runtime system & Environment scope tracker.
- `/evaluator`: Tree-walking AST evaluator with standard library built-ins.
- `/repl`: Interactive REPL with multiline inputs balancing.
- `main.go`: RustX CLI.

## Building and Running

### Requirements
- Go 1.18+

### Build
To build the `rustx` executable binary, run:
```bash
go build -o rustx main.go
```

### Run Scripts
Execute `.rx` files directly or with `run` command:
```bash
./rustx samples/hello.rx
# OR (using run command)
./rustx run samples/hello.rx
```

### Interactive REPL Prompt
Start interactive REPL session:
```bash
./rustx repl
```

### Display Version
```bash
./rustx version
```

### Tests
To run the standard Go tests:
```bash
go test ./...
```
