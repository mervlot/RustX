# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable foundation for a modern programming language inspired by **Rust** and **Go**. It balances Go's syntax simplicity and fast interpreter ergonomics with Rust's explicit declarations and future-proof design.

> **Note:** The current release is the **Foundation Phase** (v0.1). Advanced memory safety features such as ownership, borrow checking, and lifetimes are intentionally omitted at this layer so the architecture can serve as an accessible runtime platform. All AST nodes and scope management structures are engineered to accommodate borrow checks and bytecode targets without requiring a total compiler rewrite.

---

## Key Features

- **Blended Syntax:** Function declarations using `func`, mutable variable definitions via `var` and `:=` shorthand, curly brace blocks, optional semicolons.
- **Dynamic Primitive Types & Strong Runtime Evaluation:** Built-in integers (`int`), floating-point numbers (`float`), booleans (`bool`), and UTF-8 strings (`string`).
- **Control Flow Constructs:** `if` / `else if` / `else` conditional branching, `while` loops, and three-part Go-style `for` loops (`for i := 0; i < N; i++`).
- **First-Class & Recursive Functions:** Top-level named functions, anonymous function expressions, parameter and return type annotations, recursive function calls, and automatic `main()` function execution.
- **Human-Readable Diagnostics:** Structured syntax and runtime errors complete with 1-based line and column information.
- **Interactive Multiline REPL:** Shell tracking unclosed `{}` and `()` pairs for seamless multiline snippet execution while preserving state across inputs.
- **Standard Library:** Built-in support for string operations (`to_upper`, `to_lower`, `contains`, `split`, `join`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time utilities (`now`, `sleep`), and file I/O (`read_file`, `write_file`).

---

## Directory Structure

```text
rustx/
├── ast/         # Abstract Syntax Tree interface and node definitions
├── errors/      # Structured SyntaxError and RuntimeError diagnostics
├── evaluator/   # Runtime AST evaluator and standard library built-ins
├── lexer/       # Tokenizer supporting comments, escape sequences, and tracking
├── object/      # Runtime objects, type signatures, and Environment scopes
├── parser/      # Pratt recursive-descent parser
├── repl/        # Interactive multiline REPL shell
├── samples/     # Example .rx programs
├── token/       # Token definitions and keyword lookup maps
├── go.mod       # Go module file
├── main.go      # CLI entry point
└── README.md    # Documentation
```

---

## Installation & Building

Ensure you have **Go 1.21+** installed on your system.

### Build Executable
```bash
go build -o rustx main.go
```

### Run Tests
```bash
go test ./...
```

---

## CLI Usage

RustX Starter provides a full-featured CLI:

### 1. Execute a `.rx` Script
```bash
./rustx run samples/factorial.rx
# OR via shorthand
./rustx samples/hello.rx
```

### 2. Launch Interactive REPL
```bash
./rustx repl
# OR simply
./rustx
```

### 3. Syntax Validation (`build`)
Check a `.rx` source file for syntax/parsing errors without running it:
```bash
./rustx build samples/control_flow.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## Syntax & Code Examples

### Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Recursion & Math (`samples/factorial.rx`)
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

### Control Flow & Variables (`samples/control_flow.rx`)
```rustx
func main() {
    name := "RustX"

    if name == "RustX" {
        println("Language initialized:", name)
    }

    var count = 0
    while count < 3 {
        println("while loop count:", count)
        count++
    }

    for i := 0; i < 3; i++ {
        println("for loop i:", i)
    }
}
```

---

## Architecture & Future Compatibility

RustX Starter's modular architecture is designed for future extension:

1. **Borrow Checker & Ownership Integration:** The `ast.Node` interfaces carry complete token position metadata (`GetToken()`), allowing a future static analysis phase to analyze lifetimes and variable moves prior to evaluation.
2. **Bytecode Compilation:** The clean AST separation allows replacing or supplementing `evaluator/` with a bytecode compiler and stack-based VM.
3. **Structs, Enums, Traits & Generics:** AST node interfaces (`Statement` and `Expression`) allow seamless addition of struct declarations, enum matching, and trait interfaces.

---

## License

Distributed under the MIT License.
