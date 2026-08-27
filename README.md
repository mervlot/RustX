# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200" />
</p>

**RustX Starter** is a minimal, clean, and scalable interpreted programming language designed as a scripting + systems hybrid. It blends the structural elegance of **Rust** with the concise syntax and simplicity of **Go**.

---

## 🌟 Features

- **Blended Rust/Go Syntax**: Uses `func` for functions, `var` for mutable variable declarations, and `:=` for shorthand declaration and assignment.
- **Optional Semicolons & Optional Parentheses**: Write clean, uncluttered code while maintaining standard brace blocks.
- **Dynamic Scoped Interpreter**: Lexical scope chain supporting local variable bindings, global variables, closures, and recursion.
- **Rich Data Types**: `int`, `float`, `bool`, `string`, `null`, and `array`.
- **Full Control Flow**: `if` / `else if` / `else`, `while` loops, and 3-part Go-style `for` loops (`for i := 0; i < 10; i++`).
- **Standard Library Included**:
  - **Printing**: `print`, `println`
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`
- **Interactive Multiline REPL**: Persistent environment across statements with automatic block brace matching.
- **Diagnostic Error Handling**: Position-aware error reporting providing precise line and column numbers.

---

## 🏗️ Architecture & Modular Design

```
.
├── token/       # Token taxonomy and keyword map
├── errors/      # Structured diagnostic error types (SyntaxError, RuntimeError)
├── lexer/       # Lexical analyzer (supports single & multi-line comments, escapes)
├── ast/         # Extensible Abstract Syntax Tree node definitions
├── parser/      # Pratt recursive descent parser (handles operator precedence)
├── object/      # Runtime object model and environment scoping chain
├── evaluator/   # Core interpreter and standard library builtins
├── repl/        # Interactive multiline REPL interface
├── samples/     # Example RustX (.rx) programs
└── main.go      # CLI entry point
```

---

## 🚀 Getting Started & Build Instructions

### Prerequisites
- [Go 1.20+](https://golang.org/dl/)

### Building the Compiler/Interpreter

From the project root:

```bash
go build -o rustx main.go
```

### Running Tests

Execute unit tests across all packages:

```bash
go test -v ./...
```

---

## 💻 CLI Usage

```bash
# Execute a RustX script
./rustx run samples/hello.rx

# Shorthand execution
./rustx samples/factorial.rx

# Start the interactive REPL prompt
./rustx repl
# OR simply
./rustx

# Perform syntactic validation
./rustx build samples/std_demo.rx

# View version information
./rustx version
```

---

## 📝 Syntax Guide & Code Examples

### Example 1: Hello World
```rustx
func main() {
    println("RustX Starter")
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

println("Factorial of 5 is:", factorial(5))
```

### Example 3: Variable Declarations & Loops
```rustx
name := "RustX"
var age int = 16

if name == "RustX" {
    println("Language initialized for age:", age)
}

for i := 0; i < 5; i++ {
    println("Iteration:", i)
}
```

### Example 4: Standard Library Capabilities
```rustx
text := "rustx starter"
println("Uppercase:", to_upper(text))

words := split(text, " ")
println("Word count:", len(words))

println("Power 2^10:", pow(2, 10))
write_file("out.txt", "Hello from RustX!")
```

---

## 🔮 Future Extensibility Design

RustX Starter is architected from the ground up to allow future extensions without rewriting the frontend parser or interpreter:

1. **Borrow Checker & Ownership Model**: AST nodes retain source location (`token.Token`) and parameter types (`*ast.Parameter`). An AST pass for static ownership analysis can be inserted after parsing and before evaluation.
2. **Native Compilation / Bytecode VM**: The AST nodes mirror standard SSA format capabilities, allowing a future bytecode compiler or LLVM backend targeting native binaries.
3. **Structs, Enums & Traits**: Precedence tables and Pratt parser definitions are structured to easily add `struct`, `enum`, `trait`, and `impl` syntax rules.
