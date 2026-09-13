# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

**RustX Starter** is a minimal, clean, and scalable compiled/interpreted hybrid programming language designed as a blend of **Rust** and **Go**. It focuses on simplicity, fast evaluation, syntax readability, and future extensibility.

> **Note:** The current version is a lightweight foundation. Advanced memory safety systems such as borrow checking, ownership models, lifetimes, and static reference graph verification are designed into the AST and interpreter architecture to be integrated in future compiler phases without requiring a top-to-bottom redesign.

---

## Key Features

- **Blended Syntax:** Combines Go's clean syntax (`func`, `:=`, optional semicolons) with Rust's modern features and block structures.
- **Data Types:** Native support for `int`, `float`, `bool`, `string`, and `array`.
- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`
  - Comparison: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`
  - Postfix: `++`, `--`
- **Control Flow:** `if`, `else`, `while`, and Go-style 3-part `for` loops.
- **Functions:** Named and anonymous functions, parameters with type annotations, return types, and full recursion support.
- **Standard Library:** Built-in functions for string manipulation (`to_upper`, `to_lower`, `contains`, `split`, `join`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), filesystem operations (`read_file`, `write_file`), and printing (`print`, `println`).
- **Interactive REPL:** Multiline input support with bracket/brace balance detection and persistent scope between prompt evaluations.
- **Diagnostics:** Human-readable `SyntaxError` and `RuntimeError` messages containing line and column position tracking.

---

## Directory & Package Structure

```
.
├── ast/         # Abstract Syntax Tree node definitions
├── errors/      # Position-aware error types (SyntaxError, RuntimeError)
├── evaluator/   # AST Interpreter and runtime built-in function library
├── lexer/       # Lexical analyzer / Tokenizer
├── object/      # Runtime Object type system and Lexical Environment
├── parser/      # Recursive descent Pratt parser
├── repl/        # Interactive REPL prompt shell
├── samples/     # Example .rx scripts
├── token/       # Token types and keyword mappings
├── go.mod       # Go module file
├── main.go      # CLI entrypoint
└── README.md    # Documentation
```

---

## Getting Started

### Prerequisites

- [Go 1.20+](https://golang.org/dl/) installed on your machine.

### Building RustX Starter

To build the executable from source:

```bash
go build -o rustx main.go
```

---

## Usage & CLI Commands

### 1. Execute a Script
Run `.rx` files directly using the `run` command or shorthand file argument:

```bash
./rustx run samples/01_hello.rx
# OR
./rustx samples/01_hello.rx
```

### 2. Interactive REPL Prompt
Launch the multiline interactive shell:

```bash
./rustx repl
# OR simply
./rustx
```

Example session:
```
RustX Starter v0.1 Prompt
Type 'exit' or press Ctrl+C to quit.
>>> func add(a int, b int) int {
...     return a + b
... }
>>> add(10, 32)
42
>>> exit
```

### 3. Syntax Verification (`build`)
Validate `.rx` files for lexical and syntactical correctness without execution:

```bash
./rustx build samples/02_factorial.rx
```

### 4. Display Version
```bash
./rustx version
```

---

## Code Examples

### Example 1: Hello World (`samples/01_hello.rx`)
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### Example 2: Factorial Recursion (`samples/02_factorial.rx`)
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

### Example 3: Control Flow & Loops (`samples/03_control_flow.rx`)
```rustx
func main() {
    name := "RustX"

    if name == "RustX" {
        println("Language initialized:", name)
    }

    println("Counting using for loop:")
    for i := 0; i < 5; i++ {
        println("i =", i)
    }

    println("Counting using while loop:")
    var count = 3
    while count > 0 {
        println("count =", count)
        count--
    }
}
```

---

## Future Architecture Roadmap

RustX Starter is designed to serve as the foundation for a full systems language. Upcoming extensions will include:
- Ownership system & Borrow checker (reference lifetimes and mutability permissions)
- Module system & multi-file `import` support
- Structs, Enums, and Trait / Interface abstractions
- Bytecode Virtual Machine & Native Compiler backend
