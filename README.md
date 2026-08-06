# RustX Starter

A minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility. Implemented in Go.

## Core Language Design
- **Syntax**: Merges patterns of Go (e.g., shorthand assignments `:=`, post-increment `++`, optional semicolons) and Rust (`func` keyword, block structures).
- **Execution**: Direct interpreter execution, syntactic validation checking, and an interactive REPL with multiline support.
- **Errors**: Non-crashing, precise position-aware (line and column) lexical, syntax, and runtime diagnostic reporting.

## Directory Structure
- `token/`: Token types, keywords, and structural representations.
- `lexer/`: Word tokenization and code scanners.
- `parser/`: Recursive Descent and Pratt parser with customized operators and optional type annotations support.
- `ast/`: Extensible Abstract Syntax Tree representation of statements, expressions, and literals.
- `object/`: Runtime object-based value representation and scoped environment hierarchies.
- `evaluator/`: Scoping, logical expressions evaluation, loops, function literals execution, and builtins registry.
- `repl/`: Interactive CLI with multiline balance checks and persistent scoped state.
- `samples/`: Core demo files of RustX.

## Getting Started

### Building the Compiler/Interpreter
To compile the `rustx` executable from source, run:
```bash
go build -o rustx main.go
```

### Running Scripts
To run a `.rx` file directly:
```bash
./rustx samples/hello.rx
```
or via the `run` command:
```bash
./rustx run samples/factorial.rx
```

### Checking Syntax (Build Mode)
To check a `.rx` file for lexical or syntax errors without evaluating it:
```bash
./rustx build samples/fibonacci.rx
```

### Running the REPL
Start the interactive session:
```bash
./rustx repl
```
or simply:
```bash
./rustx
```

### Running Tests
To run the full suite of unit tests:
```bash
go test ./...
```
