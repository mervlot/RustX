# RustX Starter

RustX Starter is a minimal, clean, scalable programming language inspired by Rust and Go, designed for simplicity, performance, readability, and future extensibility.

---

## Language Features

- **Blended Syntax**: Go/Rust style clean syntax (`func`, `var`, `:=`, optional semicolons).
- **Control Flow**: Unparenthesized `if`, `else`, `while`, and `for` loops.
- **Data Types**: `int`, `float`, `bool`, `string`, and `Array`.
- **Functions**: Named and anonymous functions with optional parameter/return type annotations and full recursion support.
- **Lexical Scoping**: Proper local block scoping with variable update mechanics across nested scopes.
- **Built-in Standard Library**: Printing (`println`, `print`), Strings (`len`, `contains`, `split`, `join`, `to_upper`, `to_lower`), Math (`abs`, `pow`, `sin`, `cos`, `tan`), Time (`sleep`, `now`), and Filesystem (`read_file`, `write_file`).
- **Interactive REPL**: Multiline code evaluation with persistent environment state and friendly error diagnostics.

---

## Directory Architecture

```
/token       - Token types, literals, and keyword lookups
/errors      - Standardized human-readable SyntaxError and RuntimeError diagnostics
/lexer       - Tokenizer with line/column tracking, comments (// and /* */), and escape sequences
/ast         - Position-aware Abstract Syntax Tree nodes
/parser      - Pratt / Recursive Descent parser
/object      - Runtime object system and Environment lexical scoping
/evaluator   - Runtime AST evaluator and built-in standard library
/repl        - Multiline interactive REPL prompt
/samples     - Example .rx programs
```

> **Note on Future Compatibility:** The architecture separates tokenization, AST construction, object representations, and evaluation environments. Memory safety systems (ownership, borrow checker, lifetimes) can be integrated into the AST / static analysis pass in future revisions without rewriting the parser or evaluator.

---

## Installation & Build

### Prerequisites
- Go 1.20+

### Build Executable
```bash
go build -o rustx main.go
```

---

## Usage

### 1. Execute a Script
```bash
./rustx run samples/hello.rx
# Or shorthand execution:
./rustx samples/factorial.rx
```

### 2. Interactive REPL
```bash
./rustx repl
# Or simply:
./rustx
```

Example REPL session:
```
RustX Starter v0.1
>>> x := 10
>>> func double(n int) int { return n * 2; }
>>> println(double(x))
20
```

### 3. Syntax Verification / Build
```bash
./rustx build samples/control_flow.rx
```

### 4. Version
```bash
./rustx version
```

---

## Example Program

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

---

## Testing

Run all unit tests across the codebase:
```bash
go test ./...
```
