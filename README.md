# RustX Starter

**RustX Starter** is a minimal, clean, and scalable compiled/interpreted hybrid programming language inspired by **Rust** and **Go**. It blends the syntactic elegance and static typing feel of Rust with the simplicity and rapid developer feedback of Go.

---

## 🌟 Architectural Philosophy & Key Principles

1. **Simplicity First**: Clean, uncluttered syntax using `func`, `var`, `:=`, optional semicolons, and natural control flow headers.
2. **Modular Architecture**: Clean separation between Lexer, Pratt Parser, AST, Runtime Objects, Evaluator Environment, and CLI/REPL.
3. **Future Extensibility**: Designed ground-up so advanced systems (borrow checker, ownership model, lifetimes, async/await, structs, enums, traits, bytecode VM) can be introduced in later versions **without rewriting the parser or AST core**.
4. **Human-Readable Diagnostics**: Precise `SyntaxError` and `RuntimeError` line-and-column reporting with zero panics on standard user errors.

---

## 📁 Repository Structure

```
.
├── ast/            # Abstract Syntax Tree nodes & interfaces
├── errors/         # Standardized SyntaxError and RuntimeError diagnostics
├── evaluator/      # AST Evaluator and standard library built-ins (math, string, time, fs)
├── lexer/          # Tokenizer supporting literals, comments, escape sequences, positions
├── object/         # Runtime Object system & Environment with lexical scoping
├── parser/         # Pratt Recursive Descent Parser
├── repl/           # Interactive REPL with multiline balance tracking & persistent scope
├── samples/        # Sample RustX scripts (.rx)
├── token/          # Token types, keywords, and location data
├── go.mod          # Go module specification
├── main.go         # CLI main entry point
└── README.md       # Project documentation
```

---

## 🚀 Getting Started

### Prerequisites
- **Go 1.18+** installed on your system.

### Installation / Building

Build the `rustx` executable binary:

```bash
go build -o rustx main.go
```

---

## 💻 CLI & Usage

### 1. Run a Script
Execute `.rx` scripts using the `run` command or direct file shorthand:

```bash
./rustx run samples/hello.rx
# OR shorthand:
./rustx samples/factorial.rx
```

### 2. Interactive REPL
Start the interactive REPL prompt:

```bash
./rustx repl
# OR run without arguments:
./rustx
```

REPL Session Example:
```text
RustX Starter REPL v0.1.0
Type 'exit' or press Ctrl+C to exit.
>>> x := 10
>>> y := 20
>>> println("Sum is:", x + y)
Sum is: 30
>>> func add(a int, b int) int {
...     return a + b
... }
>>> println(add(15, 25))
40
```

### 3. Syntax Validation (`build`)
Check syntax without evaluating the script:

```bash
./rustx build samples/control_flow.rx
```

### 4. Version Check
```bash
./rustx version
```

---

## 📖 Language Syntax & Features

### 1. Variables & Assignment
```rustx
var x = 10
var y int = 20
name := "RustX"
x = x + 5
x++
```

### 2. Data Types
- **Integer**: `42`, `-7`
- **Float**: `3.14159`, `-0.05`
- **Boolean**: `true`, `false`
- **String**: `"Hello\nWorld"`
- **Array**: `[1, 2, 3, 4]`

### 3. Functions & Recursion
```rustx
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5)) // 120
```

### 4. Control Flow
```rustx
// Conditionals
if age >= 18 {
    println("Adult")
} else {
    println("Minor")
}

// While Loops
var count = 3
while count > 0 {
    println(count)
    count--
}

// C/Go style For Loops
for i := 0; i < 5; i++ {
    println(i)
}
```

### 5. Standard Library Built-ins
- **Printing**: `print(...)`, `println(...)`
- **Strings**: `len(s)`, `contains(s, sub)`, `split(s, sep)`, `join(arr, sep)`, `to_upper(s)`, `to_lower(s)`
- **Math**: `abs(x)`, `pow(base, exp)`, `sin(rad)`, `cos(rad)`, `tan(rad)`
- **Time**: `now()`, `sleep(ms)`
- **Filesystem**: `read_file(path)`, `write_file(path, content)`

---

## 🧪 Testing

Run all unit tests across the repository:

```bash
go test ./...
```

---

## 🔮 Future Extensibility Roadmap
- **Phase 2**: Ownership, Borrow Checker, Lifetimes
- **Phase 3**: Structs, Methods, Traits/Interfaces, Enums, Pattern Matching
- **Phase 4**: Modules, Package Management, Bytecode Compiler & Virtual Machine
