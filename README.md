# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable scripting + systems hybrid programming language inspired by **Rust** and **Go**. It focuses on speed, simplicity, readability, and modular language architecture to serve as a strong foundation for future extensibility (such as adding ownership systems, borrow checkers, async support, modules, and native code compilation).

---

## Features

- **Blended Syntax**: Merges Rust and Go idioms naturally (`func`, `var`, `:=`, optional semicolons, curly brace blocks).
- **Core Data Types**: Integers (`int`), Floats (`float`), Booleans (`bool`), Strings (`string`), and Arrays (`Array`).
- **Control Flow**: `if` / `else if` / `else` conditionals, `while` loops, and flexible 3-part or single-condition `for` loops (`for i := 0; i < 10; i++ { ... }`).
- **Functions & Recursion**: Named and anonymous functions with optional parameter and return type annotations.
- **Built-in Standard Library**:
  - **Output**: `print`, `println`
  - **Strings**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`
  - **Time**: `now`, `sleep`
  - **Filesystem**: `read_file`, `write_file`
- **Human-Readable Diagnostics**: Clean position-aware `SyntaxError` and `RuntimeError` reporting with line and column numbers.
- **Interactive REPL**: Multiline prompt support with persistent evaluation state and balanced brace tracking.

---

## Directory Structure

```text
.
├── token/          # Token type definitions and keyword maps
├── errors/         # Diagnostic error structures (SyntaxError, RuntimeError)
├── lexer/          # Lexical analyzer / Tokenizer
├── ast/            # Abstract Syntax Tree node definitions
├── parser/         # Pratt expression parser & recursive descent statement parser
├── object/         # Runtime object representation and Environment scoping
├── evaluator/      # AST evaluator and Standard Library built-in functions
├── repl/           # Multiline interactive REPL prompt
├── samples/        # Sample RustX program scripts (.rx)
├── main.go         # Main CLI entry point
└── go.mod          # Go module definition
```

---

## Getting Started

### Prerequisites
- Go 1.20+ installed.

### Building the Compiler / Interpreter CLI

From the root directory, build the binary:

```bash
go build -o rustx main.go
```

---

## Usage

### 1. Run a Script
Run a `.rx` script directly:

```bash
./rustx run samples/01_hello.rx
```

Or pass the file path directly:

```bash
./rustx samples/02_factorial.rx
```

### 2. Interactive REPL
Launch the REPL by running `rustx` with no arguments or using `rustx repl`:

```bash
./rustx
```

Example session:
```text
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func greet(name) {
...     println("Hello " + name)
... }
>>> greet("RustX")
Hello RustX
```

### 3. Syntax Verification / Build Mode
Validate syntax without executing runtime code:

```bash
./rustx build samples/03_variables_and_loops.rx
```

---

## Sample Code

### Example 1: Hello World
```rustx
func main() {
    println("Hello World from RustX!")
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

func main() {
    println(factorial(5)) // Outputs 120
}
```

### Example 3: Loops and Variables
```rustx
func main() {
    var total = 0
    name := "RustX"

    for i := 0; i < 5; i++ {
        total = total + i
    }

    if name == "RustX" {
        println("Total sum:", total)
    }
}
```

---

## Architectural Extensibility

RustX Starter is designed to serve as a extensible compiler foundation:
1. **Position-Aware AST**: All AST nodes implement `GetToken()`, enabling precise position diagnostics for type checkers and borrow checkers.
2. **Modular Scoping**: The `Environment` structure cleanly separates variable declarations (`Set`) and updates (`Update`), facilitating lexical analysis or static scoping passes.
3. **Pratt Parser**: Simple precedence extension allowing for new operator additions (e.g. `?`, `await`, `->`).

---

## License

MIT
