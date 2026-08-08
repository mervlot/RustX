# RustX Starter

RustX Starter is a minimal, clean, and highly scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility. It blends Go and Rust syntax naturally and acts as a foundation-level programming language hybrid.

Initially, RustX Starter behaves as a lightweight scripting + systems hybrid language. The design has been architected to allow future enhancements such as a borrow checker, ownership models, lifetimes, concurrency, and native compilation without requiring a rewrite of the core lexing, parsing, and execution systems.

## Key Features

1. **Variables & Mutability**: Mutable variables by default via `var` and shorthand variable assignments `:=`.
2. **Elegant Control Flow**: Support for unparenthesized conditional `if`/`else` structures, `while` loops, and Go-style three-part `for` loops.
3. **Functions**: Supports named/anonymous function signatures, optional parameter and return types (ignored during evaluation but parsed for future extensibility), and full recursion.
4. **Pratt & Recursive Descent Parsing**: Built with modularity in mind. Every expression implements operator precedence and stops cleanly at braces `{` when parsing conditions.
5. **Data Types**: Implements native `int`, `float`, `bool`, and `string` primitives, plus `array` types with zero-based index access.
6. **Detailed Error Handling**: Returns precise, formatted, position-aware `SyntaxError` and `RuntimeError` diagnostics containing exact line and column locations.
7. **Comprehensive Standard Library**: Supports `print`, `println`, `abs`, `pow`, `sin`, `cos`, `tan`, `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`, `sleep`, `now`, `read_file`, and `write_file` out of the box.

---

## Directory Structure

```text
.
├── token/          # Token types, definitions, and keywords
├── lexer/          # Lexer implementation with line & column tracking, strings, escapes, and comments
├── ast/            # Complete abstract syntax tree nodes (position-aware)
├── parser/         # Pratt expression and recursive descent parsing engine
├── object/         # Runtime representation, types, and Environment scoping
├── evaluator/      # Tree-walking interpreter and standard library builtins
├── repl/           # Multiline interactive REPL prompt with block-balance checking
├── errors/         # Custom SyntaxError and RuntimeError diagnostics
├── samples/        # Interactive sample programs (.rx)
├── main.go         # Command-line application entrypoint
└── README.md       # Project documentation
```

---

## Build & Run Instructions

To compile the `rustx` command-line tool, make sure you have Go installed, and execute:

```bash
go build -o rustx main.go
```

### CLI Options

#### 1. Execute an `.rx` Script
```bash
./rustx run samples/hello.rx
```

You can also use the shorthand path-execution (without the `run` command):
```bash
./rustx samples/hello.rx
```

#### 2. Start the REPL Prompt
```bash
./rustx repl
```
Or start it simply by calling the executable directly with no arguments:
```bash
./rustx
```

#### 3. Syntactic Validation (Dry Run)
Check a script for syntactic issues without execution:
```bash
./rustx build samples/loops_builtins.rx
```
This command performs validation and returns an exit code of `1` on error.

#### 4. Print Version
```bash
./rustx version
```

---

## Example Program

```rx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }

    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```
