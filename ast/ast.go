package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

// Node represents an AST node.
type Node interface {
	GetToken() token.Token
	String() string
}

// Statement represents a statement node.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every RustX AST.
type Program struct {
	Statements []Statement
}

func (p *Program) GetToken() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetToken()
	}
	return token.Token{}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier represents an identifier expression.
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// VarStatement represents `var x = expr` or `var x int = expr` or `var x int`
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  string // optional type annotation (e.g., "int", "string")
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.Token.Literal + " " + vs.Name.String())
	if vs.Type != "" {
		out.WriteString(" " + vs.Type)
	}
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// AssignStatement represents `x = expr` or shorthand `x := expr`
type AssignStatement struct {
	Token token.Token // token.ASSIGN or token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) GetToken() token.Token { return as.Token }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String() + " " + as.Token.Literal + " " + as.Value.String() + ";")
	return out.String()
}

// ReturnStatement represents `return expr`
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Token.Literal + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents a standalone expression evaluated as a statement
type ExpressionStatement struct {
	Token      token.Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) GetToken() token.Token { return es.Token }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents `{ stmt1; stmt2; }`
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// IntegerLiteral represents an integer constant
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a floating point constant
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string literal constant
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Boolean represents `true` or `false`
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) GetToken() token.Token { return b.Token }
func (b *Boolean) String() string       { return b.Token.Literal }

// PrefixExpression represents `!expr` or `-expr`
type PrefixExpression struct {
	Token    token.Token // token.BANG, token.MINUS
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + pe.Operator + pe.Right.String() + ")")
	return out.String()
}

// PostfixExpression represents `x++` or `x--`
type PostfixExpression struct {
	Token    token.Token // token.INC, token.DEC
	Left     Expression  // Identifier expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + pe.Left.String() + pe.Operator + ")")
	return out.String()
}

// InfixExpression represents binary operations `a + b`, `a == b`, etc.
type InfixExpression struct {
	Token    token.Token // operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) GetToken() token.Token { return ie.Token }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")")
	return out.String()
}

// Parameter represents a function parameter with an optional type annotation
type Parameter struct {
	Name *Identifier
	Type string
}

func (p *Parameter) String() string {
	if p.Type != "" {
		return p.Name.String() + " " + p.Type
	}
	return p.Name.String()
}

// FunctionLiteral represents `func name(a int, b int) int { ... }` or anonymous `func(a int) { ... }`
type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Name       string      // optional function name (if named func declaration)
	Parameters []*Parameter
	ReturnType string          // optional return type annotation
	Body       *BlockStatement
}

func (fl *FunctionLiteral) statementNode()       {}
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.Token.Literal)
	if fl.Name != "" {
		out.WriteString(" " + fl.Name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	if fl.ReturnType != "" {
		out.WriteString(fl.ReturnType + " ")
	}
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents `add(1, 2)` or `println("hello")`
type CallExpression struct {
	Token     token.Token // '(' token or function name token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// IfExpression represents `if condition { ... } else { ... }`
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) statementNode()       {}
func (ie *IfExpression) GetToken() token.Token { return ie.Token }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if " + ie.Condition.String() + " " + ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else " + ie.Alternative.String())
	}
	return out.String()
}

// WhileStatement represents `while condition { ... }`
type WhileStatement struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) GetToken() token.Token { return ws.Token }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while " + ws.Condition.String() + " " + ws.Body.String())
	return out.String()
}

// ForStatement represents both 3-part for loop (`for init; cond; post { ... }`) and single condition for loop (`for cond { ... }`)
type ForStatement struct {
	Token     token.Token // token.FOR
	Init      Statement   // optional: e.g. `i := 0` or `var i = 0`
	Condition Expression  // condition
	Post      Statement   // optional: e.g. `i++` or `i = i + 1`
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) GetToken() token.Token { return fs.Token }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fs.Init != nil {
		out.WriteString(fs.Init.String() + "; " + fs.Condition.String() + "; " + fs.Post.String())
	} else if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	out.WriteString(" " + fs.Body.String())
	return out.String()
}

// ArrayLiteral represents `[1, 2, 3]`
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elems := []string{}
	for _, el := range al.Elements {
		elems = append(elems, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression represents `arr[0]`
type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ie.Left.String() + "[" + ie.Index.String() + "])")
	return out.String()
}
