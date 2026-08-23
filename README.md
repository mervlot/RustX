# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It blends the syntax and features of both languages to prioritize simplicity, performance, readability, and future extensibility.

> **Note**: This foundation phase behaves as a lightweight scripting + systems hybrid language. Memory safety features such as ownership, lifetimes, and a borrow checker are intentionally deferred for future architectural phases, but the codebase architecture is modular and ready for these systems without requiring a rewrite.

---

## Language Highlights & Syntax Blend

- **`func`** for function definitions.
- **`var`** keyword and **`:=`** shorthand for variable declaration/assignment.
- **Optional Semicolons**: Code naturally reads clean without strict semicolon enforcement.
- **Strong & Simple Typing**: Support for `int`, `float`, `bool`, `string`, and `array`. Optional type annotations in function parameters and return signatures.
- **Control Flow**: Curly-brace blocks for `if`, `else`, `while`, and 3-part `for` loops (`for i := 0; i < 10; i++`).
- **Comments**: Single-line (`//`) and multi-line (`/* ... */`).
- **Standard Library**: Built-in functions for printing (`print`, `println`), strings (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), math (`abs`, `pow`, `sin`, `cos`, `tan`), time (`now`, `sleep`), and filesystem ops (`read_file`, `write_file`).

---

## Directory Architecture

```
rustx/
├── ast/         # Abstract Syntax Tree nodes implementing Node/Statement/Expression interfaces
├── errors/      # Custom SyntaxError and RuntimeError diagnostics with line & column info
├── evaluator/   # Tree-walking interpreter & standard library builtins
├── lexer/       # Lexical analyzer / tokenizer with position tracking and string escape handling
├── object/      # Runtime object model & lexical Environment scoping
├── parser/      # Pratt/Recursive Descent Parser for operator precedence and syntax rules
├── repl/        # Multiline REPL prompt preserving state across commands
├── samples/     # Sample RustX (.rx) programs
├── token/       # Token types, literals, keywords, and mapping
├── go.mod       # Go module declaration
├── main.go      # CLI entry point (run, build, repl, version)
└── README.md    # Language specification and manual
```

---

## Getting Started

### Building RustX

Ensure you have **Go 1.18+** installed. Clone the repository and build the binary:

```bash
go build -o rustx main.go
```

### Running Tests

Run unit tests across all packages:

```bash
go test ./...
```

---

## CLI Usage

### 1. Execute a Script
Run `.rx` files using either `rustx run <file>` or shorthand `rustx <file>`:

```bash
./rustx run samples/hello.rx
# Output: Hello, RustX Starter!

./rustx samples/factorial.rx
# Output: Factorial of 5 is: 120
```

### 2. Interactive REPL
Launch the REPL by invoking `rustx` with no arguments or `rustx repl`:

```bash
./rustx
```

```
RustX Starter v0.1 REPL
Type 'exit' or press Ctrl+C to quit.
>>> x := 5
>>> y := 10
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println("Sum:", add(x, y))
Sum: 15
```

### 3. Syntax Validation (`build`)
Validate syntax without executing runtime code:

```bash
./rustx build samples/std_demo.rx
# Output: Build successful for samples/std_demo.rx (syntax validated)
```

### 4. Version Info
Check the language version:

```bash
./rustx version
# Output: RustX Starter v0.1
```

---

## Code Examples

### Example 1: Function Recursion & Entry Point (`main`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    val := factorial(5)
    println("Factorial of 5 is:", val)
}
```

### Example 2: Loops & Variable Scoping
```rustx
func main() {
    println("For Loop:")
    for i := 0; i < 5; i++ {
        println("i =", i)
    }

    println("While Loop:")
    count := 3
    while count > 0 {
        println("count =", count)
        count--
    }
}
```

---

## Future Extensibility

The parser and AST designs allow seamlessly plugging in:
- Ownership & Borrow Checker analysis passes before evaluation/codegen.
- Structs, Enums, and Pattern Matching.
- Modules / Package Import systems.
- Concurrency constructs (goroutines/channels or async/await).
- Bytecode compilation or Native LLVM codegen.
