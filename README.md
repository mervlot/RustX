# RustX Starter

![RustX](assets/rustx-logo.png)

RustX Starter is a clean, scalable programming language inspired by the strengths of **Rust** and **Go**. It blends high readability, modern syntax, and expressive control flow into a lightweight scripting + systems hybrid language foundation.

> **Design Note:** The foundation is designed to allow future additions of ownership models, borrow checkers, lifetimes, async/await, structs, traits, and bytecode compilation without requiring architecture rewrites.

---

## Features

- **Blended Rust + Go Syntax:**
  - `func` for function declarations
  - `var` declarations and `:=` shorthand assignment
  - Optional semicolons
  - Optional parentheses for `if`, `while`, and `for` conditions
  - Post-increment and decrement (`++`, `--`)
  - Curly-brace scoped blocks
- **Data Types:**
  - `int` (64-bit integers)
  - `float` (64-bit floating point)
  - `bool` (`true`, `false`)
  - `string` (with escape sequence support: `\n`, `\t`, `\r`, `\"`, `\\`)
  - `array` (list literals `[1, 2, 3]` and indexing `arr[0]`)
- **Operators:**
  - Arithmetic: `+`, `-`, `*`, `/`, `%`
  - Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=`
  - Logical: `&&`, `||`, `!`
  - Assignment: `=`, `:=`
- **Control Flow:**
  - `if`, `else if`, `else`
  - `while` loops
  - `for` loops (both 3-part `for i := 0; i < 10; i++ { ... }` and single-condition `for cond { ... }`)
- **Functions:**
  - First-class function support
  - Optional type annotations in signatures (`func add(a int, b int) int`)
  - Recursion and closures
  - Automatic `func main()` entrypoint invocation when running `.rx` files
- **Comments:**
  - Single-line `// ...`
  - Multi-line `/* ... */`
- **Standard Library Built-ins:**
  - **I/O & Printing:** `println(...)`, `print(...)`
  - **Strings:** `len(s)`, `str_contains(s, sub)`, `str_split(s, sep)`, `str_upper(s)`, `str_lower(s)`, `str_trim(s)`, `str_replace(s, old, new)`
  - **Math:** `math_abs(x)`, `math_sqrt(x)`, `math_pow(b, exp)`, `math_min(a, b)`, `math_max(a, b)`, `math_floor(x)`, `math_ceil(x)`
  - **Time:** `time_now()`, `time_millis()`, `time_sleep(ms)`
  - **Filesystem:** `read_file(path)`, `write_file(path, content)`, `file_exists(path)`
- **Human-Readable Diagnostics:**
  - Clear `SyntaxError` and `RuntimeError` messages with line and column numbers.

---

## Project Architecture

The project is structured into modular packages to ensure maintainability and future extensibility:

```
├── token/        # Token type definitions, token constants, and keyword lookups
├── errors/       # Standardized SyntaxError and RuntimeError definitions with position tracking
├── lexer/        # Lexical analyzer / tokenizer handling literals, operators, strings, comments
├── ast/          # Abstract Syntax Tree node definitions (statements, expressions, literals)
├── parser/       # Pratt / recursive descent parser generating AST from token stream
├── object/       # Runtime object system and lexical environment scoping
├── evaluator/    # AST evaluator / interpreter and standard library built-ins
├── repl/         # Interactive REPL supporting multiline input and balance checking
├── samples/      # Example programs demonstrating RustX features
└── main.go       # CLI router supporting run, repl, build, and version commands
```

### Future Extensibility Roadmap
- **Ownership & Borrow Checker:** AST expressions and object representations maintain source tokens and distinct binding identifiers, making it straightforward to attach static lifetime and ownership analysis passes.
- **Type Checking:** Parser preserves parameter and return type annotations. A dedicated type checking pass can be inserted before evaluation.
- **Bytecode VM / Compiler:** The AST structure decouples parsing from execution, allowing a bytecode compiler or LLVM backend to replace the tree-walk evaluator.

---

## Installation & Build Instructions

### Prerequisites
- Go 1.20+ (tested on Go 1.24)

### Building from Source

To compile the `rustx` CLI executable:

```bash
go build -o rustx main.go
```

To run all tests:

```bash
go test ./...
```

---

## CLI Usage

### Interactive REPL
Launch the REPL by running `rustx` or `rustx repl`:

```bash
$ ./rustx
RustX Starter v0.1
Type RustX expressions or statements. Press Ctrl+D (or Ctrl+C) to exit.
>>> x := 5
>>> y := 10
>>> println(x + y)
15
```

The REPL supports multiline statements:
```bash
>>> func double(n int) int {
...     return n * 2
... }
>>> double(21)
42
```

### Running Scripts
Execute any `.rx` source file:

```bash
./rustx run samples/01_hello.rx
# Or using shorthand:
./rustx samples/01_hello.rx
```

### Validating / Checking Syntax
Verify syntax without executing:

```bash
./rustx build samples/02_factorial.rx
```

### Version Information
```bash
./rustx version
```

---

## Language Syntax & Examples

### Hello World (`samples/01_hello.rx`)
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Recursion & Functions (`samples/02_factorial.rx`)
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

### Control Flow (`samples/03_control_flow.rx`)
```rustx
func main() {
    name := "RustX"

    if name == "RustX" {
        println("Language initialized successfully!")
    }

    println("Counting with for loop:")
    for i := 0; i < 5; i++ {
        println("  i =", i)
    }

    println("Counting with while loop:")
    count := 3
    while count > 0 {
        println("  count =", count)
        count--
    }
}
```

### Standard Library (`samples/04_stdlib.rx`)
```rustx
func main() {
    // Strings
    clean := str_trim("   RustX   ")
    println("Upper:", str_upper(clean))

    // Math
    println("Sqrt(25):", math_sqrt(25))
    println("2^10:", math_pow(2, 10))

    // Time & Filesystem
    println("Timestamp:", time_now())
    write_file("output.txt", "RustX works!")
    if file_exists("output.txt") {
        println("File content:", read_file("output.txt"))
    }
}
```
