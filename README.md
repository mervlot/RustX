# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the syntactic elegance of Rust and Go with simplicity, speed, readability, and future extensibility.

---

## 🚀 Key Features

- **Blended Syntax**: Uses `func` for functions, `var` for mutable variables, and `:=` for shorthand variable declaration.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**:
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Relational: `==`, `!=`, `<`, `>`, `<=`, `>=`
  - Logical: `&&`, `||`, `!`
  - Postfix: `++`, `--`
- **Control Flow**: `if` / `else if` / `else`, `while`, and Go-style `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library**:
  - Printing: `println()`, `print()`
  - String manipulation: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math operations: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time utilities: `now`, `sleep`
  - Filesystem access: `read_file`, `write_file`
- **Human-Readable Diagnostics**: Detailed line and column error messages for both syntax and runtime errors.
- **Interactive REPL**: Interactive prompt with multiline input support and persistent state.

---

## 🏗 Modular Architecture

```text
rustx/
├── token/      # Lexical tokens and type definitions
├── errors/     # Standardized SyntaxError and RuntimeError diagnostics
├── lexer/      # Tokenizer for source code with comment and string escape support
├── ast/        # Extensible Abstract Syntax Tree node definitions
├── parser/     # Recursive descent Pratt parser with optional semicolon support
├── object/     # Runtime environment and object system
├── evaluator/  # AST Evaluator and Built-in Standard Library functions
├── repl/       # Interactive REPL with multiline support
├── samples/    # Sample programs demonstrating language syntax and capabilities
└── main.go     # CLI entry point (run, build, repl, version)
```

---

## 🛠 Building and Running

### Prerequisites
- [Go 1.22+](https://go.dev/)

### Compilation
To compile the `rustx` executable:
```bash
go build -o rustx main.go
```

### Running Tests
To run all unit tests across all packages:
```bash
go test ./...
```

---

## 💻 CLI Usage

### Run a Script
```bash
./rustx run samples/hello.rx
```
Or shorthand:
```bash
./rustx samples/hello.rx
```

### Interactive REPL Prompt
```bash
./rustx repl
```
Or simply:
```bash
./rustx
```

### Validate / Build Syntax
```bash
./rustx build samples/factorial.rx
```

### Display Version
```bash
./rustx version
```

---

## 📝 Code Examples

### 1. Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### 2. Recursion & Functions (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### 3. Loops and Variables (`samples/loops_and_vars.rx`)
```rustx
var count = 0
name := "RustX"

if name == "RustX" {
    println("Language initialized")
}

for i := 1; i <= 5; i++ {
    println("Count:", i)
}
```

### 4. Standard Library (`samples/stdlib_demo.rx`)
```rustx
msg := "rustx starter programming language"
upperMsg := to_upper(msg)
println("Upper:", upperMsg)

words := split(msg, " ")
println("Word count:", len(words))

println("Absolute value of -42:", abs(-42))
println("2^10 =", pow(2, 10))
```

---

## 🔮 Future Extensibility

RustX Starter is built from the ground up to allow future language extensions without requiring core architectural redesigns:
- **Ownership & Borrow Checker**: Variable and value structures can easily incorporate owner pointers and reference scopes.
- **Generics & Type Checker**: AST and Pratt parser nodes support optional type parameters and signatures.
- **Structs, Enums, & Traits**: Clean node hierarchy allows adding new declaration types and interface checking.
- **Bytecode VM / Native Compilation**: The current AST and evaluator pipeline can target a bytecode compiler/VM or LLVM backend.
