# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="200"/>
</p>

**RustX Starter** is a clean, modular, foundation-level programming language inspired by **Rust** and **Go**. It balances simplicity, speed, readability, and future extensibility.

---

## Features

- **Blended Syntax**: Combines the clarity of Go (`func`, `var`, `:=`) with Rust-inspired syntax patterns.
- **Semicolons Optional**: Flexible statement termination.
- **Unparenthesized Control Flow**: Supports `if`, `while`, and `for` without requiring surrounding parentheses.
- **Strong & Simple Dynamic Runtime**: Primitive types (`int`, `float`, `bool`, `string`), arrays, and high-order functions.
- **Standard Library Included**: Built-in packages for `strings`, `math`, `time`, and basic `filesystem` I/O.
- **Interactive REPL**: Interactive prompt with multiline block tracking and preserved state.
- **Clean Diagnostics**: Human-readable syntax and runtime error diagnostics with line/column reporting.

---

## Directory & Architecture Overview

The project is structured into clean, modular Go packages designed for future expansion (e.g. static typing, ownership/borrow checker, bytecode VM):

```
.
├── token/        # Lexical token definitions and keywords
├── errors/       # Diagnostic error reporting types (SyntaxError, RuntimeError)
├── lexer/        # Tokenizer with line/column tracking & comment processing
├── ast/          # Abstract Syntax Tree node definitions
├── parser/       # Recursive descent + Pratt expression parser
├── object/       # Runtime object system and scoped Environments
├── evaluator/    # AST evaluator, short-circuiting logic, and built-in stdlib
├── repl/         # Interactive multiline REPL interface
├── samples/      # Sample RustX scripts (.rx)
├── main.go       # CLI entry point
└── README.md
```

---

## Language Reference & Syntax

### 1. Variables & Assignment
```rustx
var x int = 10;   // Explicit variable declaration (optional type)
var y = 20;       // Type inference
z := 30;          // Shorthand assignment
z = 40;           // Update existing variable
```

### 2. Postfix Operators
```rustx
count := 0;
count++;
count--;
```

### 3. Control Flow

#### `if` / `else`
```rustx
name := "RustX"

if name == "RustX" {
    println("Language initialized")
} else {
    println("Unknown language")
}
```

#### `while` Loop
```rustx
i := 0
while i < 5 {
    println(i)
    i++
}
```

#### `for` Loop
```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

### 4. Functions & Recursion
```rustx
func add(a int, b int) int {
    return a + b
}

func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

If a zero-argument `main()` function is defined, it will automatically execute when running a script file.

---

## Standard Library Functions

| Function | Category | Description |
|---|---|---|
| `print(...)`, `println(...)` | I/O | Outputs arguments to stdout |
| `len(val)` | General | Returns length of string or array |
| `contains(str, sub)` | Strings | Checks if string contains substring |
| `split(str, sep)` | Strings | Splits string into array by delimiter |
| `join(arr, sep)` | Strings | Joins array of elements with delimiter |
| `to_upper(str)`, `to_lower(str)` | Strings | Converts case |
| `abs(n)`, `pow(b, e)` | Math | Absolute value, exponentiation |
| `sin(x)`, `cos(x)`, `tan(x)` | Math | Trigonometric functions |
| `now()`, `sleep(ms)` | Time | Unix time in ms, pause execution |
| `read_file(path)`, `write_file(path, content)` | Filesystem | Basic file read/write |

---

## Building & Installation

### Requirements
- [Go 1.20+](https://golang.org)

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

```bash
# Execute a script
./rustx run samples/factorial.rx

# Direct file execution shorthand
./rustx samples/hello.rx

# Syntactic validation without execution
./rustx build samples/std_demo.rx

# Interactive REPL prompt
./rustx repl

# Print version
./rustx version
```

---

## Interactive REPL

Run `./rustx` or `./rustx repl` to start:

```
RustX Starter v0.1 Interactive REPL
Type 'exit' or press Ctrl+D to quit.
>>> x := 10
>>> func double(n) { return n * 2; }
>>> double(x)
20
```

---

## Future Compatibility & Extensibility

RustX Starter is built with an AST and Environment architecture specifically structured to accommodate future systems without requiring a rewrite:

- **Type Checker Pass**: Type annotations in function parameters (`a int`) and variable statements are preserved in the AST nodes, enabling a dedicated static type checking pass before evaluation.
- **Ownership & Borrow Checker**: The AST statement and variable node structure supports adding lifetime parameters, borrow flags (`&`, `&mut`), and move semantics checks.
- **Bytecode Compiler & VM**: AST nodes map cleanly to stack-based bytecode instructions for future high-performance execution engines.
