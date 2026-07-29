# RustX Starter

RustX Starter is a minimal, clean, and extensible compiled/interpreted programming language inspired by the elegance of Rust and the simplicity of Go. It focuses on speed, readability, simplicity, and future compatibility with safety models (such as borrow checkers and ownership lifetimes).

## Features

- **Variables**: `var` keyword and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Operators**: `+`, `-`, `*`, `/`, comparison (`==`, `!=`, `>`, `<`, `>=`, `<=`), and logical (`&&`, `||`).
- **Control Flow**: `if`, `else`, `while`, and Go-style `for` loops.
- **Functions**: Parameters, optional type annotations, named/anonymous functions, and recursive support.
- **REPL**: Multiline support with interactive variable/function preservation.
- **Standard Library**: Math (`abs`, `pow`, `sin`, `cos`, `tan`), String operations (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), Time (`sleep`, `now`), and Filesystem basic operations (`read_file`, `write_file`).

---

## Folder Structure

```
/lexer      - Tokenizes character input streams
/parser     - Top-down Recursive Descent + Pratt parser
/ast        - Extensible Abstract Syntax Tree node hierarchy
/object     - Object system and runtime environments
/evaluator  - Tree-walking interpreter evaluating AST to Object nodes
/repl       - Interactive multi-line REPL
/token      - Lexical token structures
/errors     - Standard Syntax and Runtime Diagnostics
```

---

## Getting Started

### Installation / Compilation

To build the executable, run:

```bash
go build -o rustx main.go
```

### Running Scripts

```bash
./rustx run samples/factorial.rx
```

### Direct File Execution (Shorthand)

```bash
./rustx samples/factorial.rx
```

### Launching the REPL

```bash
./rustx repl
```

---

## Future Extensibility

The parser and interpreter architecture is designed to support modern constructs without breaking the core parser model. Future additions include:
- Ownership & Borrow Checker
- Traits & Structs
- Generics
- Lightweight Channels/Concurrency
