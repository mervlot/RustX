# RustX Starter

RustX Starter is a foundation-level programming language inspired by Rust and Go, focused on simplicity, speed, readability, and future extensibility.

## Features

- **Variables**: `var` keyword and `:=` shorthand.
- **Data Types**: `int`, `float`, `bool`, `string`.
- **Operators**: `+`, `-`, `*`, `/`, `==`, `!=`, `>`, `<`, `>=`, `<=`, `&&`, `||`, `++`.
- **Functions**: Parameters, return types, and recursion support.
- **Control Flow**: `if`, `else`, `while`, `for` loops.
- **Printing**: `println()`, `print()`.
- **Comments**: `//` single line, `/* */` multiline.
- **Built-in Functions**: `len`, `abs`, `pow`, `sin`, `cos`, `tan`, `to_upper`, `to_lower`, `contains`, `split`, `join`, `now`, `sleep`, `read_file`, `write_file`.

## Project Structure

- `ast/`: Abstract Syntax Tree node definitions.
- `evaluator/`: Tree-walking interpreter and built-in functions.
- `lexer/`: Lexical analyzer.
- `object/`: Runtime object system and environment.
- `parser/`: Recursive descent + Pratt parser.
- `repl/`: Interactive REPL.
- `token/`: Token definitions.
- `samples/`: Example programs in RustX.

## Usage

### Run a script
```bash
go run main.go run samples/hello.rx
```

### Interactive REPL
```bash
go run main.go repl
```

### Version
```bash
go run main.go version
```

## Sample Code

### Named Function and Recursion
```rustx
func factorial(n) {
    if (n <= 1) {
        return 1
    }
    return n * factorial(n - 1)
}

println(factorial(5))
```

### Loops and Variables
```rustx
var total = 0
for (i := 0; i <= 10; i++) {
    total = total + i
}
println(total)
```

## Architecture

RustX uses a modular architecture:
1. **Lexer**: Tokenizes the input string.
2. **Parser**: Builds an AST from tokens using Recursive Descent for statements and Pratt Parsing for expressions.
3. **Evaluator**: Traverses the AST and executes it using an object-based runtime system.
4. **Environment**: Handles variable scoping and persistence.

The AST and Object systems are designed to be extensible, allowing for future additions like an ownership system, borrow checker, and native compilation.
