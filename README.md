# RustX Starter

RustX Starter is a lightweight scripting + systems hybrid programming language designed as a clean, scalable foundation inspired by Rust and Go. It combines clean syntax (`func`, `var`, `:=`, optional semicolons, curly braces) with speed, readability, and future extensibility.

---

## 🚀 Features

- **Variables & Mutability**: Declare variables with `var` or `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`), comparisons (`==`, `!=`, `<`, `>`, `<=`, `>=`), logical (`&&`, `||`), postfix (`++`, `--`).
- **Control Flow**: `if` / `else`, `while` loops, and `for` loops (both 3-part Go-style and single-condition).
- **Functions**: Named and anonymous function definitions, recursive call support, parameters, and optional return type annotations.
- **Built-in Standard Library**:
  - Printing: `print`, `println`
  - Strings: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`
  - Math: `abs`, `pow`, `sin`, `cos`, `tan`
  - Time: `now`, `sleep`
  - Filesystem: `read_file`, `write_file`
- **Comments**: Single-line (`//`) and multi-line (`/* */`).
- **Human-Readable Diagnostics**: Precise line and column location reporting for syntax and runtime errors without panicking.

---

## 📁 Architecture & Folder Structure

```
.
├── token/        # Lexical tokens and position tracking
├── errors/       # SyntaxError and RuntimeError diagnostic reporting
├── lexer/        # Tokenizer supporting single/multi-line comments & string escape decoding
├── ast/          # Abstract Syntax Tree nodes designed for future extensions
├── parser/       # Recursive descent / Pratt parser with precedence rules
├── object/       # Runtime object definitions and lexically scoped environments
├── evaluator/    # AST evaluation engine & standard library builtins
├── repl/         # Interactive prompt with multiline block detection
├── samples/      # Example .rx programs
├── main.go       # CLI entry point
├── go.mod        # Module configuration
└── README.md     # Documentation
```

---

## 🛠️ Installation & Building

Ensure you have **Go 1.20+** installed on your system.

```bash
# Clone repository and build binary
go build -o rustx main.go
```

---

## 💻 CLI Usage

```bash
# Run a script file
./rustx run samples/hello.rx
# Shorthand script execution
./rustx samples/hello.rx

# Start the interactive REPL prompt
./rustx repl
# Or run without arguments
./rustx

# Parse and validate syntax without executing
./rustx build samples/factorial.rx

# View version
./rustx version
```

---

## 📝 Example Code

### Hello World (`samples/hello.rx`)
```rustx
func main() {
    println("Hello World from RustX Starter!")
}
```

### Recursion & Math (`samples/factorial.rx`)
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

func main() {
    result := factorial(5)
    println("Factorial of 5 is:", result)
}
```

---

## 🧪 Testing

Run the unit test suite across all packages:

```bash
go test ./...
```
