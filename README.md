# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a lightweight, compiled/interpreted programming language designed by blending the elegance of Rust and the simplicity of Go. It features dynamic scripting behavior with strong, clean syntactic conventions.

> **Note on Architecture:** Advanced memory systems such as borrow checking, ownership models, and lifetimes are intentionally omitted in this foundation phase. The AST and interpreter architecture are modularly built so these features can be integrated seamlessly in future releases.

---

## Language Features

- **Blended Syntax:** `func` for function declarations, `var` for variables, `:=` shorthand assignments, optional semicolons, curly-brace blocks.
- **Data Types:** `int`, `float`, `bool`, `string`, `array`.
- **Control Flow:** `if / else if / else`, `while`, `for` (both 3-part `for init; cond; post` and single condition loops).
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`), Logical (`&&`, `||`, `!`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Postfix (`++`, `--`).
- **Standard Library Built-ins:**
  - **Printing:** `print()`, `println()`
  - **Strings:** `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math:** `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time:** `now()`, `sleep()`
  - **Filesystem:** `read_file()`, `write_file()`
- **Human-Readable Error Diagnostics:** Line and column position reporting for `SyntaxError` and `RuntimeError`.

---

## Directory Architecture

```
rustx/
├── token/        # Token definitions and keyword map
├── errors/       # Standardized error reporting types
├── lexer/        # Tokenizer with position tracking and comment skipping
├── ast/          # Abstract Syntax Tree node definitions
├── parser/       # Pratt parser with precedence parsing
├── object/       # Runtime object model and scope environment
├── evaluator/    # Interpreter evaluator and standard library built-ins
├── repl/         # Multiline interactive REPL prompt
├── samples/      # Example programs (.rx)
├── main.go       # CLI application entrypoint
└── go.mod        # Go module specification
```

---

## Getting Started

### Prerequisites
- Go 1.20+ installed.

### Building
Compile the executable binary:
```bash
go build -o rustx main.go
```

### Running Programs
Execute a `.rx` script directly or via the `run` command:
```bash
./rustx run samples/hello.rx
# OR shorthand:
./rustx samples/hello.rx
```

### Interactive REPL
Launch the interactive REPL prompt:
```bash
./rustx repl
# OR simply:
./rustx
```
Inside the REPL:
```
RustX Starter v0.1
Type 'exit' or 'quit' to exit.
>>> name := "RustX"
>>> println("Hello", name)
Hello RustX
```

### Validating Syntax
Check for syntax errors without executing:
```bash
./rustx build samples/factorial.rx
```

### Running Unit Tests
Run the entire unit test suite:
```bash
go test ./...
```

---

## Code Examples

### 1. Hello World
```rustx
func main() {
    println("Hello World")
}
```

### 2. Factorial & Recursion
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

### 3. Loops and Standard Library
```rustx
func main() {
    name := "RustX"
    if name == "RustX" {
        println("Language initialized")
    }

    var sum = 0
    for i := 0; i < 10; i++ {
        sum = sum + i
    }
    println("Sum:", sum)
}
```

---

## Future Extensibility Roadmap

The clean separation of AST nodes, Lexer tokens, Environment scoping, and Evaluation logic allows for straightforward future extensions:
1. **Ownership & Borrow Checker:** Adding static analysis passes on AST nodes before evaluation.
2. **Bytecode VM & Native Compilation:** Replacing or supplementing the AST evaluator with a compiler emitting bytecode or native LLVM IR.
3. **Structs, Enums, & Traits:** Expanding parser expression nodes and runtime objects.
4. **Modules & Package Manager:** Enhancing standard library resolution and AST import nodes.
