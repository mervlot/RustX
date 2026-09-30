# RustX Starter

RustX Starter is a minimal, clean, and scalable programming language inspired by Rust and Go, focusing on simplicity, speed, readability, and future extensibility.

## Features

- **Blended Syntax:** Blends Rust and Go syntax cleanly (`func`, `var`, `:=`, optional type annotations, curly brace blocks, optional semicolons).
- **Core Types:** Integers, Floats, Booleans, Strings, and Arrays.
- **Control Flow:** `if`/`else`, `while`, and `for` loops (both 3-part C/Go style loops and single-condition loops).
- **Functions:** Parameters, optional return types, recursive call support, closures, and first-class functions.
- **Built-in Standard Library:**
  - **Printing:** `print()`, `println()`
  - **Strings:** `len()`, `contains()`, `split()`, `join()`, `to_upper()`, `to_lower()`
  - **Math:** `abs()`, `pow()`, `sin()`, `cos()`, `tan()`
  - **Time:** `now()`, `sleep()`
  - **Filesystem:** `read_file()`, `write_file()`
- **Interactive REPL:** Multiline input tracking, statement history, state preservation between commands.
- **Rich Error Diagnostics:** Precise line and column location reporting for syntax and runtime errors without crashing the process.

---

## Architecture Overview

```text
/
├── token/       Token definitions and keyword mapping
├── errors/      Structured SyntaxError and RuntimeError diagnostics
├── lexer/       Tokenizer scanning input stream with line & column tracking
├── ast/         Extensible AST node definitions implementing Node interface
├── parser/      Recursive descent Pratt parser for statements & expressions
├── object/      Runtime object representation and enclosed environment scoping
├── evaluator/   Interpreter engine and built-in standard library functions
├── repl/        Interactive shell prompt supporting multiline input
├── samples/     Sample RustX Starter (.rx) scripts
└── main.go      CLI entrypoint (run, repl, build, version)
```

The system is designed with a clean, modular architecture. Future language extensions (such as ownership, borrow checker, concurrency, generics, and native/bytecode compilation) can easily attach static analysis passes or bytecode generation stages on top of the AST without requiring parser or lexer redesigns.

---

## Getting Started

### Prerequisites

- Go 1.20+ installed on your system.

### Build Binary

To build the `rustx` CLI binary:

```bash
go build -o rustx main.go
```

### Running Tests

Run all unit tests across lexer, parser, and evaluator packages:

```bash
go test ./...
```

---

## CLI Usage

### 1. Execute Script Direct

```bash
./rustx run samples/01_hello.rx
```

Or shorthand execution:

```bash
./rustx samples/01_hello.rx
```

### 2. Interactive REPL

Launch the interactive REPL shell:

```bash
./rustx repl
```

or simply run `./rustx` with no arguments.

Example REPL session:

```text
RustX Starter REPL v0.1.0
Type 'exit' or 'quit' to exit.
>>> x := 10
>>> y := 20
>>> println(x + y)
30
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println(add(15, 25))
40
```

### 3. Build / Syntactic Validation

To perform lexical and syntactic validation on a `.rx` file without executing it:

```bash
./rustx build samples/02_factorial.rx
```

### 4. Version Check

```bash
./rustx version
```

---

## Sample Code

### 1. Hello World (`samples/01_hello.rx`)

```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### 2. Recursion & Math (`samples/02_factorial.rx`)

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    res := factorial(5)
    println("Factorial of 5 is:", res)
}
```

### 3. Feature Showcase (`samples/03_showcase.rx`)

```rustx
func main() {
    println("=== RustX Starter Feature Showcase ===")

    var language = "RustX Starter"
    version := 0.1
    isAwesome := true

    println("Language:", language, "Version:", version, "Awesome:", isAwesome)

    sum := 0
    for i := 1; i <= 5; i++ {
        sum = sum + i
    }
    println("Sum 1..5 =", sum)

    greeting := "  hello rustx  "
    clean := to_upper(greeting)
    println("Upper greeting:", clean)

    val := abs(-42.5)
    println("Absolute value of -42.5:", val)

    arr := [10, 20, 30]
    println("Array len:", len(arr), "First el:", arr[0])
}
```
