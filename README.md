# RustX Starter

<p align="center">
  <img src="assets/rustx-logo.png" alt="RustX Logo" width="180"/>
</p>

A clean, minimal, and scalable programming language inspired by **Rust** and **Go**, focused on simplicity, speed, readability, and future extensibility.

---

## Overview

RustX Starter blends the expressive ergonomics of Rust and Go into a lightweight scripting + systems hybrid language. It features:

- Familiar, elegant syntax: `func`, `var`, and `:=` shorthand.
- Strong yet flexible primitive types (`int`, `float`, `bool`, `string`, `array`).
- Arithmetic, comparison, and short-circuit logical operators.
- First-class and recursive functions with optional type hints.
- Flexible control flow (`if` / `else`, `while`, and Go-style 3-part `for` loops).
- Built-in standard library utilities for strings, math, time, and filesystem I/O.
- Extensible recursive descent Pratt parser and environment runtime designed to accommodate future ownership models, borrow checkers, async systems, and byte-code/native compilers.
- Interactive multiline REPL with persistent environment.

---

## Architecture

The project is structured in modular, decoupled packages:

```
├── ast/            # Abstract Syntax Tree definitions (extensible nodes with token tracking)
├── errors/         # Human-readable SyntaxError and RuntimeError diagnostics with line & column
├── evaluator/      # Runtime evaluation engine, environment scoping, and built-in standard library
├── lexer/          # Tokenizer with position, line, and column tracking
├── object/         # Runtime object representation and environment scoping (Get/Set/Update)
├── parser/         # Pratt recursive descent parser with operator precedence
├── repl/           # Multiline interactive REPL
├── samples/        # Example RustX source scripts (.rx)
├── token/          # Token definitions and keyword lookup
└── main.go         # CLI entry point (run, repl, build, version)
```

---

## Language Syntax & Features

### 1. Variables & Assignments
Variables are declared using `var` or the shorthand `:=`:

```rustx
var age = 16
var score int = 100
name := "Muhammad"

// Reassignment
age = 17
```

### 2. Data Types
- **int**: 64-bit signed integers (`42`, `-10`)
- **float**: 64-bit IEEE 754 floating point numbers (`3.14`)
- **bool**: `true` or `false`
- **string**: UTF-8 strings with escape sequence support (`"Hello\nWorld"`)
- **array**: Dynamic lists (`[1, 2, 3]`)

### 3. Operators
- **Arithmetic**: `+`, `-`, `*`, `/`, `%`
- **Comparison**: `==`, `!=`, `<`, `>`, `<=`, `>=`
- **Logical**: `&&`, `||`, `!`
- **Postfix**: `++`, `--`

### 4. Control Flow

#### If / Else
Condition parentheses are optional:

```rustx
if score >= 90 {
    println("Grade: A")
} else if score >= 80 {
    println("Grade: B")
} else {
    println("Grade: C")
}
```

#### While Loop
```rustx
count := 0
while count < 5 {
    println(count)
    count++
}
```

#### For Loop
Supports classic 3-part Go-style loop headers:
```rustx
for i := 0; i < 10; i++ {
    println(i)
}
```

### 5. Functions & Recursion
Functions can be declared as top-level or anonymous first-class values:

```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5:", factorial(5))
```

If a zero-argument `func main()` is defined in a script, it is automatically invoked as the program entry point upon script execution.

### 6. Standard Library

- **I/O & Printing**:
  - `print(...)`: Prints space-separated values without a trailing newline.
  - `println(...)`: Prints space-separated values with a newline.
- **Conversion & Length**:
  - `len(item)`: Returns string character length or array element count.
  - `to_string(val)`: Converts values to string.
  - `to_int(val)`: Converts values to integer.
  - `to_float(val)`: Converts values to float.
- **Strings**:
  - `str_len(s)`: String length.
  - `str_contains(s, substr)`: Check substring inclusion.
  - `str_split(s, sep)`: Split string into array of substrings.
  - `str_to_upper(s)`: Uppercase conversion.
  - `str_to_lower(s)`: Lowercase conversion.
- **Math**:
  - `math_abs(x)`: Absolute value.
  - `math_max(a, b)`: Maximum of two values.
  - `math_min(a, b)`: Minimum of two values.
  - `math_pow(base, exp)`: Power exponentiation.
  - `math_sqrt(x)`: Square root.
- **Time**:
  - `time_now()`: Current UNIX timestamp in seconds.
  - `time_now_millis()`: Current UNIX timestamp in milliseconds.
  - `time_sleep(ms)`: Pause execution for `ms` milliseconds.
- **Filesystem**:
  - `read_file(path)`: Read text file contents as string.
  - `write_file(path, content)`: Write string contents to file.

---

## Installation & Building

### Prerequisites
- Go 1.18 or later

### Build Binary
Compile the binary directly from source:

```bash
go build -o rustx main.go
```

---

## Usage & CLI

### 1. Run a script
```bash
./rustx run samples/hello.rx
# Or using shorthand:
./rustx samples/hello.rx
```

### 2. Interactive REPL
Launch the REPL with multiline editing and state persistence:

```bash
./rustx repl
# Or simply:
./rustx
```

Example session:
```
RustX Starter v0.1
Type RustX expressions or statements. Press Ctrl+C or Ctrl+D to exit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func greet(name) {
...     println("Hello", name)
... }
>>> greet("Developer")
Hello Developer
```

### 3. Syntax Verification / Build
Validate file syntax without executing:

```bash
./rustx build samples/factorial.rx
```

### 4. Version Check
```bash
./rustx version
```

---

## Running Tests

Run the complete test suite across all packages:

```bash
go test -v ./...
```

---

## Future Roadmap

RustX Starter is designed so that future systems can be integrated seamlessly without rewriting the core AST or interpreter:

- Memory safety: Static ownership tracking, borrow checker, and lifetime analysis.
- Concurrency: Lightweight green threads and communication channels.
- Object-oriented/Functional primitives: Structs, enums, pattern matching, and traits.
- Packaging: Module import resolution and package manager.
- Performance: Bytecode virtual machine and LLVM/native code generation.
