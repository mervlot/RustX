# RustX Starter

A minimal, clean, scalable, and position-aware programming language inspired by Go and Rust. RustX Starter balances safety, simplicity, speed, and readability. It initially acts as a lightweight hybrid of a compiled CLI with a tree-walking interpreter, designed cleanly so advanced systems (like borrow checking, ownership, lifetime, or a bytecode VM) can be added in the future without rewriting the compiler front-end.

---

## 🚀 Key Features

1. **Clean Syntax**: Merges the clean syntax of Go (e.g. shorthand assignments `:=`, for-loop structure, curly-brace blocks) with the elegant expression-oriented model of Rust. Semicolons are entirely optional!
2. **Variable Bindings & Shadowing**: Supports explicit mutable bindings with `var`, shorthand automatic assignments with `:=`, and proper lexical block scoping/shadowing.
3. **Robust Control Flow**: Includes unparenthesized headers for `if`/`else` conditionals, `while` loops, and three-part Go-style `for` loops (`for i := 0; i < 10; i++`).
4. **First-Class Functions & Recursion**: Supports anonymous lambdas, first-class functions, closures, and named recursive function calls.
5. **Basic Standard Library**:
   - **Printing**: `print`, `println` (space-separated, position-aware).
   - **Math**: `abs`, `pow`, `sin`, `cos`, `tan`.
   - **Strings & Arrays**: `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`, string and array 0-based indexing (`[index]`), array element assignment.
   - **Time**: `now` (millisecond timestamp), `sleep`.
   - **Filesystem**: `read_file`, `write_file`.
6. **Detailed Error Diagnostics**: Precise position-aware `SyntaxError` and `RuntimeError` reporting tracking line and column numbers.
7. **Interactive Multiline REPL**: Supports multiline commands by balancing brackets `{}` and parentheses `()`, while preserving scope state across inputs.

---

## 📂 Project Architecture

The workspace is highly modular and strictly structured to prevent circular dependencies and promote future extensibility:

- **`/token`**: Holds structural token definitions and token lookup helpers.
- **`/errors`**: Houses custom human-readable `SyntaxError` and `RuntimeError` diagnostic definitions.
- **`/lexer`**: Scans input string characters into a stream of structured tokens, preserving exact line and column locations, skipping whitespace, and decoding escape sequences within strings.
- **`/parser`**: Implements a highly efficient Recursive Descent parser combined with a Pratt parser for mathematical expressions, ensuring correct operator precedence (including AND > OR). Handles optional semicolons and unparenthesized loop headers.
- **`/ast`**: Defines clear, strongly-typed syntax tree nodes, each implementing the `GetToken()` interface for location tracking.
- **`/object`**: Houses the runtime representation of RustX values (integers, floats, booleans, strings, arrays, functions, errors) and environment scopes (with Set for local bindings and recursive Update for nested assignment scope traversal).
- **`/evaluator`**: Interprets the AST recursively, manages short-circuiting logical operations, checks for division-by-zero, compares primitives by value, and registers built-in library functions.
- **`/repl`**: Provides the interactive Read-Eval-Print-Loop with multiline buffer tracking.
- **`/samples`**: Standard code programs showcasing various capabilities of the language.

---

## 🛠️ Build and Test Instructions

### Prerequisites
- Go 1.18 or higher

### Running All Unit Tests
Run the comprehensive test suite directly from the project root:
```bash
go test ./...
```

### Compiling the Binary
To compile the standalone `rustx` executable binary, run:
```bash
go build -o rustx main.go
```

---

## 💻 CLI Usage

The compiled `rustx` executable provides the following primary commands:

### 1. Interactive REPL (interactive mode)
To launch the multiline REPL:
```bash
./rustx
```
or
```bash
./rustx repl
```

**Example session:**
```
RustX Starter v0.1
>>> func fib(n) {
...     if n <= 1 { return n; }
...     return fib(n - 1) + fib(n - 2);
... }
>>> println("fib(10) =", fib(10))
fib(10) = 55
```

### 2. Run a Script File
To execute an `.rx` file directly:
```bash
./rustx run samples/hello.rx
```

### 3. Shorthand Execution
You can skip the `run` command by passing the `.rx` file as the first argument:
```bash
./rustx samples/factorial.rx
```

### 4. Syntax Verification (Build Check)
To parse and validate the syntax of an `.rx` file without executing it (exits with code 1 if syntax errors exist):
```bash
./rustx build samples/loops_and_builtins.rx
```

### 5. Check Language Version
To print the current version:
```bash
./rustx version
```

---

## 📜 Example Syntax & Sample Code

Take a look at `/samples` for complete, runnable program examples.

### Example 1: Standard Greeting (`samples/hello.rx`)
```go
func main() {
    println("Hello World from RustX Starter!")
}

main()
```

### Example 2: Recursive Factorial (`samples/factorial.rx`)
```go
func factorial(n int) int {
    if n <= 1 {
        return 1
    }

    return n * factorial(n - 1)
}

println("Factorial of 5 is:", factorial(5))
```

### Example 3: Scoping, Loops, and Built-ins (`samples/loops_and_builtins.rx`)
```go
// This program showcases loops and standard library capabilities of RustX.

/* Let's start with a standard
   three-part for-loop */
println("Three-part loop:")
for i := 1; i <= 5; i++ {
    println("Iteration:", i)
}

/* Now let's try a single-condition loop */
println("While-style loop:")
count := 3
for count > 0 {
    println("Count:", count)
    count--
}

/* Shadowing and Scoping */
var x = "outer"
{
    x := "inner"
    println("Nested scope:", x) // prints "inner"
}
println("Outer scope:", x) // prints "outer"

/* Built-ins math and strings */
msg := "rustx is awesome"
upperMsg := to_upper(msg)
println("Uppercase message:", upperMsg)

b := pow(2, 10)
println("2 raised to 10 is:", b)

arr := ["RustX", "is", "fast"]
joined := join(arr, " ")
println("Joined array:", joined)
```
