# RustX Starter

![RustX Logo](assets/rustx-logo.png)

RustX Starter is a clean, scalable programming language inspired by **Rust** and **Go**. It blends Go's simplicity, quick learning curve, and lightweight syntax with Rust's expressive elegance and modern design principles.

> **Note on Architecture & Roadmap**: RustX Starter currently operates as a lightweight scripting + systems hybrid language foundation. Advanced memory safety features (borrow checking, lifetimes, ownership) are omitted in this foundation release, but the modular lexer, Pratt parser, AST, and runtime architecture are explicitly structured to allow future integration of static typing checks, ownership analysis, native bytecode, or LLVM backends without requiring compiler rewrites.

---

## Language Highlights & Features

- **Blended Syntax**: Uses `func` for functions, `var` and `:=` shorthand for variable bindings, curly brace blocks, optional semicolons, and strong primitive types (`int`, `float`, `bool`, `string`).
- **Control Flow**: Supports `if` / `else`, `while`, and flexible Go-style 3-part (`for i := 0; i < 10; i++`) or single-condition `for` loops. Parentheses around conditions are optional.
- **Built-in Functions & Standard Library**:
  - **I/O & Printing**: `print()`, `println()`
  - **Strings & Collections**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time**: `now()`, `sleep()`
  - **Filesystem**: `read_file()`, `write_file()`
- **Comments**: Supports single-line (`//`) and multi-line (`/* ... */`) comments.
- **Human-Readable Diagnostics**: Precise error diagnostics formatted as `SyntaxError` or `RuntimeError` with line and column position tracking.
- **Interactive REPL**: Interactive prompt supporting multiline input tracking (unbalanced `{` or `(`) and persistent environment state across commands.

---

## Directory & Package Structure

```
.
├── ast/         # Abstract Syntax Tree interface and node definitions
├── errors/      # Position-aware SyntaxError and RuntimeError diagnostics
├── evaluator/   # Tree-walking interpreter & standard library builtins
├── lexer/       # Lexical scanner & tokenizer with location tracking
├── object/      # Runtime type system (Integer, Float, String, Function, Array, Environment)
├── parser/      # Recursive descent Pratt expression and statement parser
├── repl/        # Interactive REPL prompt with multiline tracking
├── samples/     # Example .rx source scripts demonstrating language features
├── token/       # Token types, definitions, and keyword maps
├── main.go      # CLI entrypoint and script execution runner
└── README.md
```

---

## Getting Started & Build Instructions

### Prerequisites
- [Go](https://golang.org/) 1.18 or higher.

### Building RustX CLI Binary

Compile the executable from the repository root:

```bash
go build -o rustx main.go
```

Verify binary installation:

```bash
./rustx version
# Output: RustX Starter v0.1.0
```

---

## CLI Usage

### Running Scripts
Execute a `.rx` file using the `run` command or direct execution:

```bash
./rustx run samples/hello.rx
# or shorthand direct execution:
./rustx samples/hello.rx
```

If a zero-argument `main()` function is present in the file, `rustx` automatically invokes it as the entry point.

### Interactive REPL
Launch the interactive shell:

```bash
./rustx repl
# or simply run without arguments:
./rustx
```

Example REPL Session:
```text
RustX Starter v0.1 Interactive REPL
Type statements or expressions to evaluate. Ctrl+C or Exit to quit.
>>> x := 10
>>> y := 20
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println("Sum:", add(x, y))
Sum: 30
```

### Syntax Validation (`build`)
Check a `.rx` source file for syntax errors without executing its runtime logic:

```bash
./rustx build samples/features.rx
# Output: Build successful! samples/features.rx contains valid RustX syntax.
```

---

## Code Examples

### 1. Variables & Control Flow
```rust
var x int = 10;
y := 20.5;

if x < 15 {
    println("x is less than 15");
} else {
    println("x is 15 or greater");
}

for i := 0; i < 5; i++ {
    println("Iteration:", i);
}
```

### 2. Recursion & Functions
```rust
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5:", factorial(5))
}
```

### 3. Standard Library Usage
```rust
func main() {
    // Strings & Math
    greeting := "hello world"
    println(to_upper(greeting))
    println("Power 2^10 =", pow(2, 10))

    // Filesystem I/O
    write_file("output.txt", "RustX File Output Sample\n")
    content := read_file("output.txt")
    print("Read content:", content)
}
```

---

## Running Tests

Run all unit test suites across all packages:

```bash
go test -v ./...
```

---

## Future Extensibility

RustX Starter is designed with modular AST structures and Pratt parser operator precedences to seamlessly support upcoming features:
- **Structs & Enums**: Data structure AST node additions without parser re-architecting.
- **Ownership & Borrow Checker**: Static analysis passes over AST nodes before evaluation or bytecode compilation.
- **Generics & Traits/Interfaces**: Type annotations and constraint resolution in the AST/Object layer.
- **Concurrency & Async/Await**: Goroutine/channel runtime wrappers or async task queues.
- **Bytecode VM**: Compiling AST nodes to VM instructions for increased execution performance.
