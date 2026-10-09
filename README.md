# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It combines the elegant, lightweight syntax of Go with modern programming features inspired by Rust, designed as an extensible foundation for future system features (such as borrow checking, ownership models, modules, and native compilation).

---

## 🚀 Features

- **Blazing Fast & Lightweight**: Clean AST-based interpreter written in Go.
- **Hybrid Scripting & Systems Syntax**:
  - `func` for function declarations with optional type annotations.
  - `var` for variable definitions.
  - `:=` shorthand assignment operator.
  - `++` and `--` post-increment/decrement operators.
  - Optional semicolons and unparenthesized control flow headers (`if`, `while`, `for`).
- **Rich Control Flow**:
  - `if` / `else if` / `else` conditional logic.
  - `while` loops.
  - C/Go-style 3-part `for` loops (`for i := 0; i < 10; i++ { ... }`).
- **Core Standard Library**:
  - **I/O**: `print`, `println`.
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`.
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`.
  - **Filesystem**: `read_file`, `write_file`.
- **Interactive Multiline REPL**: Preserves variables and functions across prompts with intelligent multiline input tracking.
- **Human-Readable Error Diagnostics**: Positional error reporting with line and column tracking (`SyntaxError` and `RuntimeError`).

---

## 📁 Directory Structure

```
.
├── ast/        # Abstract Syntax Tree node definitions
├── errors/     # Standardized SyntaxError and RuntimeError diagnostics
├── evaluator/  # AST Evaluator and Standard Library built-ins
├── lexer/      # Lexical analyzer / Tokenizer
├── object/     # Object runtime environment and scoping
├── parser/     # Pratt recursive descent parser
├── repl/       # Interactive REPL prompt implementation
├── samples/    # Sample RustX Starter (.rx) scripts
├── token/      # Token type definitions
├── main.go     # CLI entrypoint
└── go.mod      # Go module configuration
```

---

## 🛠️ Build & Installation

### Prerequisites
- Go 1.20 or later installed.

### Building RustX Binary
```bash
go build -o rustx main.go
```

---

## 🎮 CLI Usage

### 1. Run a Script
```bash
./rustx run samples/hello.rx
```
or use shorthand direct execution:
```bash
./rustx samples/hello.rx
```

### 2. Start Interactive REPL
```bash
./rustx repl
```
or run without arguments:
```bash
./rustx
```

Example REPL Session:
```
RustX Starter v0.1 REPL
Type 'exit' or Ctrl+C to quit.
>>> func add(a int, b int) int {
...     return a + b
... }
>>> x := 10
>>> y := 20
>>> println("Sum is:", add(x, y))
Sum is: 30
```

### 3. Syntax Verification (Build)
```bash
./rustx build samples/factorial.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## 📄 Example Programs

### Factorial Recursion (`samples/factorial.rx`)
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

### Standard Library Demo (`samples/stdlib_demo.rx`)
```rustx
func main() {
    name := "RustX Starter"
    println("Original name:", name)
    println("Upper name:", to_upper(name))
    println("Lower name:", to_lower(name))

    numbers := [1, 2, 3, 4, 5]
    println("Array length:", len(numbers))

    sum := 0
    for i := 0; i < len(numbers); i++ {
        sum = sum + numbers[i]
    }
    println("Sum of array elements:", sum)
}
```

---

## 🧪 Running Tests

Execute all unit and integration tests:
```bash
go test ./...
```

---

## 🔮 Future Architecture Compatibility

RustX Starter is designed as a clean foundation. Future releases plan to introduce:
- Memory ownership and borrow checking without breaking the interpreter pipeline.
- Structs, Enums, and Trait interfaces.
- Generics and Package Management.
- Bytecode VM / Native LLVM Compiler backend.
