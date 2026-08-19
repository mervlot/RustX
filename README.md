# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="300"/>
</p>

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the syntax and simplicity of both languages into a lightweight scripting and systems foundation designed for speed, readability, and future extensibility.

---

## Features

- **Blended Rust & Go Syntax**: Uses `func` for functions, `var` and `:=` for variable declaration/assignment, strong primitive types, optional semicolons, and brace-delimited blocks.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`), unary (`!`, `-`), and postfix (`++`, `--`).
- **Control Flow**: Unparenthesized `if`/`else if`/`else`, `while` loops, and Go-style 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Functions**: Parameters, optional type annotations, recursive support, and automatic top-level `main()` entrypoint execution.
- **Built-in Standard Library**:
  - Printing: `print`, `println`
  - String manipulation: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math utilities: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `sleep`, `now`
  - Filesystem: `read_file`, `write_file`
- **Error Diagnostics**: Standardized `SyntaxError` and `RuntimeError` with line and column position tracking without panic/crash.
- **Interactive REPL**: Multiline code block support, persistent environment, clean output formatting.

---

## Architecture Overview

RustX Starter is designed with a clean, modular architecture:

```
/token      - Token definitions, constants, and keywords with position metadata
/errors     - Position-aware SyntaxError and RuntimeError diagnostics
/lexer      - Handwritten lexer with character tracking and escape decoding
/ast        - Extensible AST node interface and structures
/parser     - Recursive descent + Pratt parser for precedence expression parsing
/object     - Runtime object types, environment scoping, and equality logic
/evaluator  - Tree-walking AST evaluator with standard library builtins
/repl       - Interactive multiline REPL shell
/samples    - Demonstrative RustX source files (.rx)
```

> **Note on Memory Safety Architecture**:
> RustX Starter is intentionally designed as a foundation phase without implementing borrow checking, lifetimes, or ownership rules. The AST nodes and environment architecture allow ownership and lifetime systems to be added in future compiler passes without rewriting the core interpreter.

---

## Getting Started

### Prerequisites
- [Go](https://golang.org/) 1.18 or higher.

### Building RustX Starter

Build the standalone binary from the repository root:

```bash
go build -o rustx main.go
```

Or run tests across all modules:

```bash
go test ./...
```

---

## Usage & CLI Commands

### 1. Execute a Script
Execute `.rx` files using the `run` command or shorthand:

```bash
./rustx run samples/hello.rx
# OR shorthand:
./rustx samples/factorial.rx
```

### 2. Interactive REPL
Start the interactive shell:

```bash
./rustx repl
# OR simply:
./rustx
```

Example session:
```
RustX Starter v0.1.0 Interactive REPL
Type 'exit' or press Ctrl+C to quit.
>>> name := "RustX"
>>> println("Hello,", name)
Hello, RustX
>>> func add(a, b) { return a + b }
>>> add(10, 20)
30
```

### 3. Syntax Validation (`build`)
Validate `.rx` source code syntax without executing:

```bash
./rustx build samples/control_flow.rx
```

### 4. Version Info
Display compiler version:

```bash
./rustx version
```

---

## Code Examples

### Factorial (Recursion)
```rust
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5 is:", factorial(5))
}
```

### Control Flow & Loops
```rust
func main() {
    name := "RustX"

    if name == "RustX" {
        println("Language initialized:", name)
    }

    for i := 0; i < 5; i++ {
        println("i =", i)
    }

    var count = 3
    while count > 0 {
        println("count =", count)
        count--
    }
}
```

---

## License

MIT License.
