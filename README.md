# RustX Starter

A minimal, clean, and scalable programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

RustX Starter is designed with a clean, modular architecture. It is implemented in Go and executes scripts either directly or through an interactive REPL.

## Core Language Design
* **Functions**: Declared using `func` keyword. Supports optional parameter and return types (ignored during execution, but parsed as AST), recursion, and anonymous functions.
* **Variables**: Declared using `var` keyword (mutable by default) or Go-style shorthand assignment `:=`.
* **Types**: Built-in primitives: `int`, `float`, `bool`, and `string`. Semiautomatic optional type annotations.
* **Control Flow**: `if`, `else`, `while` loops, and Go-style `for` loops (`for init; condition; post { ... }`). Parentheses are optional.
* **Operators**: Arithmetic (`+`, `-`, `*`, `/`), logical with short-circuiting (`&&`, `||`), comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`), and postfix operators (`++`, `--`).
* **Comments**: Supports single-line (`// ...`) and multi-line (`/* ... */`) comments.
* **Standard Library**:
  * `print`, `println`: Console output.
  * `abs`, `pow`, `sin`, `cos`, `tan`: Math.
  * `len`, `contains`, `split`, `join`, `to_upper`, `to_lower`: String processing.
  * `sleep`, `now`: Time.
  * `read_file`, `write_file`: Basic I/O.

## Project Architecture
The project structure is clean, modular, and separates concerns into isolated Go packages:
* `/token`: Token definitions, literal structures, and lookup tables.
* `/lexer`: Hand-written scanner that supports multi-line comments, line and column tracking, and escape characters in string literals.
* `/parser`: Recursive Descent and Pratt parser that supports optional semicolons and unparenthesized block control flow.
* `/ast`: Strongly typed Abstract Syntax Tree (AST) node representations of the language syntax.
* `/object`: Object system and runtime types representing live values, as well as an `Environment` scope mapping with support for recursive parent environment updates.
* `/evaluator`: High-performance tree-walking interpreter that includes division-by-zero validation, logical short-circuiting, and standard library builtins.
* `/repl`: Interactive prompt that supports multi-line inputs by tracking balanced braces and parenthetical closures.
* `/errors`: Custom `SyntaxError` and `RuntimeError` objects that provide precise, position-aware diagnostics (line and column).

## Installation & Build Instructions

Prerequisites: **Go 1.18+**

To build the executable:
```bash
go build -o rustx main.go
```

To run unit tests across all packages:
```bash
go test ./...
```

## How to Run

### 1. Execute a script
```bash
./rustx run samples/hello.rx
```

### 2. Shorthand direct execution
```bash
./rustx samples/hello.rx
```

### 3. Start the Interactive REPL
```bash
./rustx
```
or
```bash
./rustx repl
```

Within the REPL:
```
RustX Starter v0.1
Type your code below. Semicolons are optional.
>>> x := 5
>>> x++
6
>>> if x > 5 {
...     println("x is greater than 5")
... }
x is greater than 5
```

## Sample Programs
Located in `/samples`:
* `hello.rx`: Standard printing and greeting.
* `factorial.rx`: Recursive mathematical factorial function.
* `fibonacci.rx`: Recursive fibonacci computation.
* `loop.rx`: Demonstrating `for` loops, `while` loops, and variable updates.
* `builtins.rx`: Standard library built-ins for math, strings, time, and file system basics.
