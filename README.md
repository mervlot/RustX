<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

# RustX Starter

**RustX Starter** is a minimal, clean, and scalable interpreted programming language inspired by **Rust** and **Go**. It combines the elegant syntax of Rust with the simplicity and speed of Go.

Designed as a lightweight scripting and systems hybrid foundation, RustX Starter features an extensible architecture prepared for future systems-level extensions (such as ownership models, generics, concurrency, and native bytecode compilation).

---

## 🚀 Core Features

- **Blended Syntax:** `func` for functions, `var` for mutable variables, `:=` shorthand assignment, curly brace blocks, optional semicolons.
- **Data Types:** `int`, `float`, `bool`, `string`, arrays, and functions.
- **Operators:** Arithmetic (`+`, `-`, `*`, `/`, `%`), Comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), Logical (`&&`, `||`), and Postfix (`++`, `--`).
- **Control Flow:** `if` / `else if` / `else`, `while` loops, and `for` loops (both 3-part header and single condition).
- **Function Support:** First-class functions, parameter type annotations, return type annotations, recursion, and auto-invoking `main()`.
- **Error Handling:** Standardized `SyntaxError` and `RuntimeError` diagnostic reporting with accurate line and column numbers.
- **Interactive REPL:** Multiline input support, command history, and persistent memory across commands.
- **Standard Library:** Built-in functions for printing (`print`, `println`), strings (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`sleep`, `now`), and file I/O (`read_file`, `write_file`).

---

## 📁 Architecture & Directory Structure

```
├── ast/         # Abstract Syntax Tree node definitions
├── errors/      # Structured SyntaxError and RuntimeError diagnostics
├── evaluator/   # Tree-walking evaluator & standard library implementation
├── lexer/       # Lexical analyzer / Tokenizer
├── object/      # Runtime object system and Environment scopes
├── parser/      # Recursive descent Pratt parser
├── repl/        # Interactive REPL console with multiline support
├── samples/     # Example programs (.rx files)
├── token/       # Token type definitions and keyword mappings
├── main.go      # CLI entry point
└── go.mod       # Go module declaration
```

---

## 🛠️ Build & Installation

### Requirements
- **Go**: Version 1.22+ installed.

### Build
To compile the `rustx` executable binary from the root directory:

```bash
go build -o rustx main.go
```

To run all package unit tests:

```bash
go test ./...
```

---

## 💻 CLI Usage

```bash
# Start the interactive REPL
./rustx

# Run a RustX script file
./rustx run samples/01_hello.rx

# Direct script execution shortcut
./rustx samples/02_factorial.rx

# Validate script syntax without executing
./rustx build samples/03_variables_and_control.rx

# View version information
./rustx version
```

---

## 📝 Syntax Examples

### 1. Hello World

```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### 2. Variables & Operators

```rustx
func main() {
    var name string = "RustX"
    version := 0.1

    if name == "RustX" && version >= 0.1 {
        println("Initialized", name, "v", version)
    }
}
```

### 3. Recursive Functions

```rustx
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

### 4. For Loops & Post-increment

```rustx
func main() {
    sum := 0
    for i := 1; i <= 10; i++ {
        sum = sum + i
    }
    println("Sum of 1..10:", sum)
}
```

### 5. Standard Library Usage

```rustx
func main() {
    println("Uppercase:", to_upper("rustx"))
    println("Power 2^10:", pow(2, 10))

    write_file("output.txt", "RustX Stdlib!")
    data := read_file("output.txt")
    println("Read file:", data)
}
```

---

## 🔮 Future Roadmap

RustX Starter is designed to serve as a solid foundation for future additions:
- Ownership & Borrow Checker
- Async / Await & Concurrency Primitives
- Structs, Enums, and Pattern Matching
- Generics & Traits
- Native Bytecode VM or LLVM Compiler Backend
