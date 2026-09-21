# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It balances simplicity, speed, readability, and modern syntax aesthetics while maintaining an extensible architecture for future system-level capabilities.

> **Note:** Ownership, borrow checking, and lifetimes are intentionally omitted in this foundational phase. The system architecture (Lexer, Recursive Descent / Pratt Parser, AST, and Evaluator) is explicitly designed so that static analysis, ownership semantics, generics, and native/bytecode backends can be incrementally added without architectural refactoring.

---

## Language Features

- **Blended Go & Rust Syntax**: Uses `func` for function definitions, `var` and `:=` for variable bindings, curly brace blocks, and optional semicolons.
- **Variables & Mutability**: `var` declarations and `:=` shorthand assignments.
- **Data Types**: `int`, `float`, `bool`, `string`, and `array`.
- **Control Flow**: Unparenthesized/parenthesized `if` / `else if` / `else`, `while`, and 3-part `for` loops (e.g., `for i := 0; i < 10; i++`).
- **Operators**: Standard arithmetic (`+`, `-`, `*`, `/`), logical (`&&`, `||`, `!`), postfix increment/decrement (`++`, `--`), and equality/comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`).
- **Functions**: Named and anonymous functions, recursion support, type annotations, and return values.
- **Human-Readable Diagnostics**: Detailed `SyntaxError` and `RuntimeError` reports including line and column positions without uncaught panics.
- **Standard Library Built-ins**:
  - **Printing**: `print`, `println`
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`

---

## Directory & Package Structure

```
.
├── ast/            # AST node definitions and interface implementations
├── token/          # Token specifications and source position tracking
├── errors/         # SyntaxError and RuntimeError diagnostic formatting
├── lexer/          # Lexical scanner for tokens and comments
├── parser/         # Pratt parser & recursive descent statement parser
├── object/         # Runtime representation of types and environment scopes
├── evaluator/      # Evaluator engine and standard library built-ins
├── repl/           # Multiline interactive REPL
├── samples/        # Demonstrative .rx script files
├── main.go         # RustX CLI application entry point
├── go.mod          # Go module specification
└── README.md       # Documentation
```

---

## Installation & Building

### Prerequisites

- Go 1.18 or higher installed on your system.

### Build from Source

To compile the `rustx` executable binary:

```bash
go build -o rustx main.go
```

To run all unit tests across the project:

```bash
go test ./...
```

---

## Usage & CLI Commands

### 1. Execute a Script

Run a `.rx` script using the `run` command or shorthand file argument:

```bash
./rustx run samples/hello.rx
# or
./rustx samples/hello.rx
```

### 2. Interactive REPL

Launch the interactive REPL prompt:

```bash
./rustx repl
# or simply
./rustx
```

Example session:

```text
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> func greet(name string) {
...     println("Hello,", name)
... }
>>> greet("RustX")
Hello, RustX
```

### 3. Syntax Validation / Build Check

Validate syntax and parse correctness of a `.rx` file without executing runtime side-effects:

```bash
./rustx build samples/factorial.rx
```

### 4. Version Information

```bash
./rustx version
```

---

## Syntax & Examples

### Hello World

```rustx
func main() {
    println("Hello, RustX Starter!")
}
```

### Recursion (Factorial)

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

### For Loop & Variable Mutations

```rustx
sum := 0
for i := 1; i <= 10; i++ {
    sum = sum + i
}
println("Sum 1..10:", sum)
```

---

## Future Extensibility Roadmap

1. **Ownership & Borrow Checker**: Introduce static pass in parser/AST before evaluation to check lifetimes and references (`&` / `&mut`).
2. **Structs & Enums**: Extend AST and runtime objects for custom product and sum data structures.
3. **Traits / Interfaces**: Support method dispatch and type contracts.
4. **Bytecode VM or Native Compilation**: Transition evaluator to a bytecode compiler/VM backend or LLVM output.
