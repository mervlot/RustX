# RustX Starter

![RustX Logo](assets/rustx-logo.png)

**RustX Starter** is a minimal, clean, and scalable programming language inspired by **Rust** and **Go**. It combines the elegant syntax of Rust and Go with simplicity, speed, readability, and future extensibility.

> **Note:** RustX Starter is designed as a foundational hybrid scripting + systems language. Advanced memory safety features (borrow checker, ownership, lifetimes) are intentionally excluded from this initial release, but the compiler/interpreter architecture is built so those systems can be added seamlessly without a complete redesign.

---

## Language Features

- **File Extension**: `.rx`
- **Variable Declarations**: `var` keyword and `:=` shorthand assignment.
- **Data Types**: `int`, `float`, `bool`, `string`, `array`.
- **Operators**: `+`, `-`, `*`, `/`, `++`, `--`, `==`, `!=`, `<`, `>`, `<=`, `>=`, `&&`, `||`.
- **Functions**: Parameter list, return types, recursive call support, closures.
- **Control Flow**: `if`, `else if`, `else`, `while`, and `for` loops.
- **Printing**: Built-in `print(...)` and `println(...)`.
- **Comments**: `//` single-line and `/* ... */` multi-line.
- **Standard Library**: String manipulation, math utilities, time helpers, and file I/O.
- **Human-Readable Diagnostics**: Clear line and column position reporting on syntax and runtime errors.

---

## System Architecture

The project is structured into clean, modular Go packages:

```text
.
├── token/        # Token type definitions and keyword lookup
├── errors/       # Standardized SyntaxError and RuntimeError diagnostics
├── lexer/        # Lexical analyzer / tokenizer with comment & string escape handling
├── ast/          # Abstract Syntax Tree nodes for statements and expressions
├── parser/       # Pratt recursive descent parser
├── object/       # Runtime object model and environment environment chain
├── evaluator/    # AST interpreter and standard library built-ins
├── repl/         # Interactive REPL prompt with multi-line brace tracking
├── samples/      # Example .rx programs
└── main.go       # CLI interface entry point
```

### Architectural Scalability & Future Extensibility
The codebase is structured to easily integrate future language extensions:
1. **Ownership & Borrow Checker**: The AST nodes and Environment objects maintain token location metadata and reference graphs, enabling a static analysis pass between parsing and evaluation.
2. **Structs, Enums, Traits/Interfaces**: The expression and statement interfaces in `ast/ast.go` allow new AST types to be registered with minimal boilerplate.
3. **Generics & Type Checking**: Function and variable signatures support type annotations (`int`, `float`, `string`) which can be expanded to a full static type checker.
4. **Bytecode VM or Native Compiler**: The AST representation can easily be targeted by a bytecode compiler/VM or LLVM backend.

---

## Installation & Build Instructions

### Prerequisites
- [Go 1.22+](https://golang.org/dl/)

### Build Binary
To compile the `rustx` executable from source:

```bash
go build -o rustx main.go
```

### Run Tests
To run the full suite of unit tests across all packages:

```bash
go test ./...
```

---

## CLI Usage

### 1. Run Script Directly
Execute any `.rx` file using the `run` command or shorthand file execution:

```bash
./rustx run samples/hello.rx
# or
./rustx samples/factorial.rx
```

### 2. Interactive REPL
Start the interactive RustX REPL:

```bash
./rustx repl
# or simply
./rustx
```

Example REPL session:
```text
RustX Starter v0.1
Type 'exit' or press Ctrl+C to quit.
>>> x := 10
>>> y := 20
>>> println("Sum:", x + y)
Sum: 30
>>> func greet(name) {
...     println("Hello,", name)
... }
>>> greet("RustX")
Hello, RustX
```

### 3. Syntax Validation / Build Check
Check a `.rx` script for lexical or syntax errors without running it:

```bash
./rustx build samples/loop.rx
```

### 4. Version Check
Display the current language version:

```bash
./rustx version
```

---

## Syntax Examples

### Function Declaration & Recursion
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### Control Flow & Loops
```rustx
name := "RustX Starter"

if name == "RustX Starter" {
    println("Language initialized successfully!")
}

for i := 0; i < 5; i++ {
    println("i =", i)
}

count := 3
while count > 0 {
    println("count =", count)
    count--
}
```

### Standard Library Built-ins
```rustx
str := "Hello RustX Language"
println("Uppercase:", to_upper(str))
println("Length:", len(str))

arr := [10, 20, 30]
println("First element:", arr[0])

write_file("output.txt", "RustX Starter File I/O")
content := read_file("output.txt")
println("File content:", content)
```

---

## License

MIT License.
