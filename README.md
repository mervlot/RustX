# RustX Starter

RustX Starter is a clean, minimal, extensible scripting + systems hybrid programming language inspired by Rust and Go syntax. It is designed for fast execution, readability, and future compatibility with memory management systems such as ownership and borrow checking.

---

## 🚀 Features

- **Blended Syntax**: Combines the ergonomics of Go (`var`, `:=`, `func`, concise loop headers) with Rust-inspired control flow and structural design.
- **Dynamic & Safe Primitive Types**: `int`, `float`, `bool`, `string`, `array`.
- **Functions & Recursion**: First-class function definitions, parameters, return statements, and recursive execution.
- **Control Flow**: `if`, `else`, `while`, and flexible Go-style `for` loops (`for init; condition; post {}` or `for condition {}`).
- **Standard Library Built-ins**:
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`
  - **Console**: `print`, `println`
- **Interactive REPL**: Multiline block support with brace and parenthesis tracking.
- **Human-Readable Error Diagnostics**: Clear line and column position reporting for syntax and runtime errors.

---

## 🏗️ Architecture Overview

The codebase is structured modularly in Go:

```
├── token/         # Token definitions and keyword lookup
├── lexer/         # Tokenizer supporting comments and string escapes
├── errors/        # SyntaxError and RuntimeError diagnostics with line/col tracking
├── ast/           # Extensible Abstract Syntax Tree node definitions
├── parser/        # Recursive descent & Pratt expression parser
├── object/        # Internal Object representation and scoped Environment
├── evaluator/     # AST evaluation engine and standard library built-ins
├── repl/          # Interactive REPL with multiline tracking
├── samples/       # Sample .rx source code files
├── main.go        # CLI binary entrypoint
└── README.md
```

---

## 🛠️ Installation & Building

### Prerequisites
- Go 1.20 or newer

### Build Executable
```bash
go build -o rustx main.go
```

---

## 📖 CLI Usage

### Run a Script
```bash
./rustx run samples/01_hello.rx
# Or shorthand execution:
./rustx samples/02_factorial.rx
```

### Interactive REPL Prompt
```bash
./rustx repl
```
Example session:
```text
RustX Starter v0.1
>>> x := 10
>>> y := 20
>>> println("Sum =", x + y)
Sum = 30
```

### Validate Syntax (Build Check)
```bash
./rustx build samples/03_features.rx
```

---

## 📝 Syntax Examples

```rustx
// Function with type annotations and recursion
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    println("Factorial of 5:", factorial(5))

    for i := 0; i < 5; i++ {
        println("Counter:", i)
    }
}
```

---

## 🔮 Future Extensibility

RustX Starter is designed as a modular foundation. The AST and environment layout allow seamless addition of:
- Ownership & borrow checker analysis passes
- Structs, Enums, and Pattern Matching
- Generics and Traits/Interfaces
- Async/Await & Concurrency primitives
- Bytecode compilation / VM or LLVM backend target
