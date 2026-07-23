# RustX Starter

RustX Starter is a minimal, clean, and extensible programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility (such as adding ownership and borrowing systems without rewriting the entire compiler/interpreter).

RustX behaves like a lightweight scripting and systems hybrid language, with syntax elements naturally blended from both Rust and Go.

## Features

- **Variables**: `var` declarations and shorthand assignments (`:=`). Mutable by default.
- **Strong Types**: `int`, `float`, `bool`, and `string`. Semicolons are optional.
- **Control Flow**: Unparenthesized `if`, `else`, `while`, and Go-style `for` loops.
- **Functions**: Parameter lists, return types, recursive named functions, and anonymous functions.
- **Arrays & Indexing**: Zero-based index access and index-based updates/assignments for Arrays.
- **Built-in Functions**: Math operations (`abs`, `pow`, `sin`, `cos`, `tan`), String manipulation (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), Filesystem basics (`read_file`, `write_file`), and Time (`sleep`, `now`).
- **Interactive REPL**: A multiline-aware REPL prompt displaying diagnostics cleanly.
- **Modular Architecture**: Clean separation between Lexer, Parser (Pratt Operator Precedence), AST nodes, Object Runtime, and Evaluator.

---

## Folder Structure

- `/errors`: Custom position-aware syntax and runtime diagnostics.
- `/token`: Language syntax tokens and location definitions.
- `/lexer`: Stream-based tokenizer with line and column tracking, multi-line comment support, and string escape sequences.
- `/ast`: Extensible tree nodes tracking source position.
- `/parser`: Recursive descent Pratt precedence parser.
- `/object`: Object model, type declarations, and Environment with nested scopes.
- `/evaluator`: Tree-walking interpreter with short-circuiting, zero-division checking, and the standard library.
- `/repl`: Multiline interactive prompt.
- `/samples`: Sample programs demonstrating features (`factorial.rx`, `variables.rx`, `builtins.rx`).

---

## Build & Installation

Prerequisites: [Go](https://go.dev/doc/install) 1.24+

```bash
# Clone the repository and navigate to the directory
cd rustx

# Run the unit tests
go test ./...

# Build the binary
go build -o rustx main.go
```

---

## CLI Usage

### Running Scripts
Execute any `.rx` file directly:
```bash
./rustx run samples/factorial.rx
```
Or use the shorthand command:
```bash
./rustx samples/factorial.rx
```

### Interactive REPL
Launch the interactive shell:
```bash
./rustx repl
# Or simply
./rustx
```
Inside the REPL, blocks remain balanced automatically before execution:
```rustx
>>> x := 5
>>> for i := 0; i < x; i++ {
...     println("Iteration", i)
... }
Iteration 0
Iteration 1
Iteration 2
Iteration 3
Iteration 4
```

### Show Version
```bash
./rustx version
```

---

## Error Handling

Diagnostics show line and column details with human-readable descriptions:

```rustx
SyntaxError:
expected next token to be LBRACE, got INT instead
at line 4 column 12
```

---

## Future Extensibility

The parser and AST have been designed cleanly so that subsequent versions can layer on:
1. Lifetimes and Ownership systems
2. Borrow checker
3. Package modules and import systems
4. Concurrency via goroutines/channels or async/await syntax.
