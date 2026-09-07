# RustX Starter

**RustX Starter** is a lightweight, compiled/interpreted scripting + systems hybrid programming language inspired by Rust and Go. It is designed for simplicity, speed, readability, and future extensibility.

> **Note on Architecture**: RustX Starter is intentionally foundation-focused. It does not enforce a borrow checker, ownership model, or lifetimes in this phase, but its AST and object architecture are engineered to support those systems in future iterations without compiler redesigns.

---

## Language Features

- **Blended Syntax**: Combines Go-style declaration shorthand (`:=`), function signatures (`func`), and Rust-style constructs with optional semicolons.
- **Data Types**: `int`, `float`, `bool`, `string`, `Array`, `null`.
- **Variables**: `var` keyword with optional explicit typing, `:=` shorthand assignment, scoped mutation.
- **Control Flow**: `if` / `else if` / `else`, `while` loops, and `for` loops (supports 3-part `for i := 0; i < N; i++`, single-condition `for cond`, and infinite `for`).
- **Functions**: Parameter parsing, optional return types, recursion support, higher-order functions.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), increment/decrement (`++`, `--`), relational (`==`, `!=`, `<`, `>`, `<=`, `>=`), and logical (`&&`, `||`).
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Standard Library Built-ins**:
  - Printing: `print()`, `println()`
  - String Operations: `len()`, `split()`, `join()`, `contains()`, `to_upper()`, `to_lower()`
  - Math Operations: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - System/Time: `sleep()`, `now()`
  - File I/O: `read_file()`, `write_file()`

---

## Folder Architecture

```text
rustx/
├── ast/         # Abstract Syntax Tree nodes with token location metadata
├── errors/      # Human-readable SyntaxError and RuntimeError diagnostics
├── evaluator/   # Interpreter, tree-walk evaluator, and stdlib built-in functions
├── lexer/       # Lexical analyzer / tokenizer with position tracking
├── object/      # Runtime object system and scoped Environment handling
├── parser/      # Recursive descent / Pratt parser
├── repl/        # Interactive multiline REPL shell
├── samples/     # Example programs demonstrating language features
├── token/       # Token definitions and keyword mapping
├── go.mod       # Go module file definition
├── main.go      # CLI entrypoint for rustx execution
└── README.md    # Documentation and user guide
```

---

## Getting Started & Build Instructions

### Prerequisites
- Go 1.20+ installed on your system.

### Build Executable
To compile the `rustx` executable binary:

```bash
go build -o rustx main.go
```

### Running Tests
To run unit test suites across all packages:

```bash
go test ./...
```

---

## CLI Usage

### 1. Execute a `.rx` file
You can execute a file using the `run` command or pass the `.rx` file directly:

```bash
./rustx run samples/factorial.rx
# OR
./rustx samples/hello.rx
```

### 2. Interactive REPL
Start the interactive REPL prompt:

```bash
./rustx repl
# OR simply
./rustx
```

Example session:
```text
RustX Starter v0.1.0 Interactive REPL
Type 'exit' or press Ctrl+C to exit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> exit
```

### 3. Syntax Check / Build
To validate script syntax without running it:

```bash
./rustx build samples/demo.rx
```

---

## Example Programs

### Hello World
```rustx
func main() {
    println("Hello World")
}
```

### Recursive Factorial
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println(factorial(5))
}
```

---

## Future Extensibility

The architecture provides clean entry points for upcoming features:
- **Ownership & Borrowing**: AST nodes and Environment values can wrap ownership lifetimes.
- **Bytecode Virtual Machine**: The current evaluator can be replaced/complemented with a bytecode compiler and VM.
- **Structs & Traits**: Type system AST nodes easily accommodate custom struct types and trait interfaces.
- **Concurrency**: Built-in async/await or Go-style channels can be added to the runtime object system.
