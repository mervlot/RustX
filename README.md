# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable compiled/interpreted programming language designed as a hybrid inspired by **Rust** and **Go**. It balances simplicity, execution speed, readability, and future extensibility.

> **Note**: RustX Starter is designed as a foundation phase. Ownership, borrow checking, lifetimes, structs, enums, traits, and bytecode compilation are architected to be added in future versions without rewriting the core engine.

---

## Language Features

- **Blended Syntax**: Combines Go's concise `func`, `var`, and `:=` shorthand declaration with Rust-like modern control structures and expressions.
- **Dynamic & Static Typing Hybrid**: Supports optional type annotations in signatures (`func add(a int, b int) int`) while offering straightforward runtime evaluations.
- **Flexible Control Flow**: `if`, `else if`, `else`, `while`, and 3-part or 1-part `for` loops with optional parentheses around conditions.
- **Rich Operator Set**: Full support for arithmetic (`+`, `-`, `*`, `/`, `%`), comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`, `!`), and post-increment/decrement (`++`, `--`).
- **Standard Library Built-ins**:
  - **Printing**: `println()`, `print()`
  - **Strings**: `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math**: `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time**: `now()`, `sleep()`
  - **Filesystem**: `read_file()`, `write_file()`
- **Interactive REPL**: Features multiline input detection, line history, syntax error reporting, and persistent variable/function state across lines.
- **Human-Readable Errors**: Formatted `SyntaxError` and `RuntimeError` diagnostics with explicit line and column coordinates.

---

## Modular Architecture

The repository is organized into clean, isolated Go packages:

```
├── ast/          # Abstract Syntax Tree definitions and string representations
├── errors/       # Standardized human-readable SyntaxError and RuntimeError diagnostics
├── evaluator/    # AST interpreter, standard library built-ins, and execution scope logic
├── lexer/        # Tokenizer supporting comments, string escapes, numbers, and operators
├── object/       # Runtime objects (Integer, Float, String, Boolean, Function, Array, Environment)
├── parser/       # Recursive descent / Pratt parser with precedence rules
├── repl/         # Interactive REPL with multiline continuation tracking
├── samples/      # Example .rx programs demonstrating language constructs
├── token/        # Token definitions and keyword mapping
├── main.go       # CLI entry point for run, build, repl, and version commands
└── go.mod        # Go module configuration
```

---

## Installation & Build

Ensure you have [Go](https://golang.org/) (v1.22+) installed on your machine.

### Build Executable
```bash
go build -o rustx main.go
```

### Run Tests
Execute the full unit test suite:
```bash
go test ./...
```

---

## CLI Usage

### 1. Direct Execution (`run`)
Execute a `.rx` file directly:
```bash
./rustx run samples/hello.rx
```

Shorthand execution is also supported:
```bash
./rustx samples/factorial.rx
```

### 2. Interactive REPL (`repl`)
Start the interactive multiline prompt:
```bash
./rustx repl
```
Or simply:
```bash
./rustx
```

**REPL Session Example:**
```
RustX Starter v0.1.0 Interactive REPL
Type 'exit' or 'quit' to exit.
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(10, 20)
30
```

### 3. Syntax Verification (`build`)
Perform a dry-run parse check without executing:
```bash
./rustx build samples/factorial.rx
```

### 4. Version Check (`version`)
```bash
./rustx version
```

---

## Syntax Examples

### Example 1: Hello World
```rustx
func main() {
    println("Hello from RustX Starter!")
}
```

### Example 2: Recursive Factorial
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    result := factorial(5)
    println("Factorial of 5 is:", result)
}
```

### Example 3: Control Flow & Variables
```rustx
func main() {
    var count = 0
    name := "RustX"

    if name == "RustX" {
        println("Language initialized!")
    }

    for i := 1; i <= 5; i++ {
        println("Count:", i)
    }

    while count < 3 {
        count++
    }
}
```

---

## Future Compatibility & Roadmap

RustX Starter is designed to evolve into a full systems programming language with:
- [ ] Memory safety system (ownership model, borrow checker, explicit lifetimes)
- [ ] Structs, Enums, and Pattern Matching
- [ ] Traits & Interfaces
- [ ] Concurrency primitive models (channels, goroutines / async-await)
- [ ] Bytecode VM compiler and Native JIT backend

---

## License

MIT License.
