package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

type Node interface {
	GetToken() token.Token
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) GetToken() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetToken()
	}
	return token.Token{Type: token.EOF, Literal: ""}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// TypeAnnotation (e.g. "int", "string")
type TypeAnnotation struct {
	Token token.Token
	Name  string
}

func (t *TypeAnnotation) Node()             {}
func (t *TypeAnnotation) GetToken() token.Token { return t.Token }
func (t *TypeAnnotation) String() string       { return t.Name }

// VarStatement: var x int = 5; or var x = 5;
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  *TypeAnnotation // optional
	Value Expression
}

func (v *VarStatement) statementNode()       {}
func (v *VarStatement) GetToken() token.Token { return v.Token }
func (v *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(v.Token.Literal + " " + v.Name.String())
	if v.Type != nil {
		out.WriteString(" " + v.Type.String())
	}
	if v.Value != nil {
		out.WriteString(" = " + v.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ShortVarStatement: x := 5
type ShortVarStatement struct {
	Token token.Token // token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (s *ShortVarStatement) statementNode()       {}
func (s *ShortVarStatement) GetToken() token.Token { return s.Token }
func (s *ShortVarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.Name.String() + " := " + s.Value.String() + ";")
	return out.String()
}

// AssignStatement: x = 5 or arr[0] = 5
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name  Expression  // Identifier or IndexExpression
	Value Expression
}

func (a *AssignStatement) statementNode()       {}
func (a *AssignStatement) GetToken() token.Token { return a.Token }
func (a *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(a.Name.String() + " = " + a.Value.String() + ";")
	return out.String()
}

// ReturnStatement: return 5;
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) GetToken() token.Token { return r.Token }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.Token.Literal + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token // first token of expression
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) GetToken() token.Token { return e.Token }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// BlockStatement: { ... }
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) GetToken() token.Token { return b.Token }
func (b *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range b.Statements {
		out.WriteString(s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// IntegerLiteral
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) GetToken() token.Token { return b.Token }
func (b *Boolean) String() string       { return b.Token.Literal }

// ArrayLiteral: [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) GetToken() token.Token { return a.Token }
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression: arr[0]
type IndexExpression struct {
	Token token.Token // '['
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) GetToken() token.Token { return i.Token }
func (i *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")
	return out.String()
}

// PrefixExpression: -5, !true
type PrefixExpression struct {
	Token    token.Token // token.MINUS, token.BANG
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) GetToken() token.Token { return p.Token }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression: 5 + 5
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) GetToken() token.Token { return i.Token }
func (i *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression: i++ or i--
type PostfixExpression struct {
	Token    token.Token // token.INC or token.DEC
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// FunctionParameter
type FunctionParameter struct {
	Name *Identifier
	Type *TypeAnnotation // optional
}

// FunctionLiteral: func add(a int, b int) int { return a + b } or anonymous
type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Name       *Identifier // optional named function
	Parameters []*FunctionParameter
	ReturnType *TypeAnnotation // optional
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) GetToken() token.Token { return f.Token }
func (f *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		pStr := p.Name.String()
		if p.Type != nil {
			pStr += " " + p.Type.String()
		}
		params = append(params, pStr)
	}
	out.WriteString(f.Token.Literal)
	if f.Name != nil {
		out.WriteString(" " + f.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	if f.ReturnType != nil {
		out.WriteString(f.ReturnType.String() + " ")
	}
	out.WriteString(f.Body.String())
	return out.String()
}

// CallExpression: add(2, 3)
type CallExpression struct {
	Token     token.Token // '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) GetToken() token.Token { return c.Token }
func (c *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// IfExpression: if condition { ... } else { ... }
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) GetToken() token.Token { return ie.Token }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

// WhileStatement: while condition { ... }
type WhileStatement struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode()       {}
func (w *WhileStatement) GetToken() token.Token { return w.Token }
func (w *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(w.Condition.String())
	out.WriteString(" ")
	out.WriteString(w.Body.String())
	return out.String()
}

// ForStatement: for i := 0; i < 10; i++ { ... } or for condition { ... }
type ForStatement struct {
	Token     token.Token // token.FOR
	Init      Statement   // e.g., VarStatement, ShortVarStatement, ExpressionStatement, or nil
	Condition Expression  // optional
	Post      Statement   // e.g., ExpressionStatement (i++), AssignStatement, or nil
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) GetToken() token.Token { return fs.Token }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fs.Init != nil || fs.Post != nil {
		if fs.Init != nil {
			out.WriteString(fs.Init.String())
		}
		out.WriteString("; ")
		if fs.Condition != nil {
			out.WriteString(fs.Condition.String())
		}
		out.WriteString("; ")
		if fs.Post != nil {
			out.WriteString(fs.Post.String())
		}
	} else if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}
