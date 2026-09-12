package ast

import (
	"bytes"

	"github.com/username/rustx/token"
)

type Node interface {
	TokenLiteral() string
	String() string
	GetToken() token.Token
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

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
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

// Identifier
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// Parameter
type Parameter struct {
	Token          token.Token // token.IDENT
	Name           *Identifier
	TypeAnnotation string
}

func (p *Parameter) Node()                {}
func (p *Parameter) TokenLiteral() string { return p.Token.Literal }
func (p *Parameter) GetToken() token.Token { return p.Token }
func (p *Parameter) String() string {
	if p.TypeAnnotation != "" {
		return p.Name.String() + " " + p.TypeAnnotation
	}
	return p.Name.String()
}

// VarStatement (e.g., var age int = 16 or var age = 16)
type VarStatement struct {
	Token          token.Token // token.VAR
	Name           *Identifier
	TypeAnnotation string
	Value          Expression
}

func (v *VarStatement) statementNode()       {}
func (v *VarStatement) TokenLiteral() string { return v.Token.Literal }
func (v *VarStatement) GetToken() token.Token { return v.Token }
func (v *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(v.TokenLiteral() + " " + v.Name.String())
	if v.TypeAnnotation != "" {
		out.WriteString(" " + v.TypeAnnotation)
	}
	if v.Value != nil {
		out.WriteString(" = " + v.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ShortVarStatement (e.g., name := "Muhammad")
type ShortVarStatement struct {
	Token token.Token // token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (s *ShortVarStatement) statementNode()       {}
func (s *ShortVarStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ShortVarStatement) GetToken() token.Token { return s.Token }
func (s *ShortVarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.Name.String() + " := " + s.Value.String() + ";")
	return out.String()
}

// AssignStatement (e.g., name = "Muhammad")
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

func (a *AssignStatement) statementNode()       {}
func (a *AssignStatement) TokenLiteral() string { return a.Token.Literal }
func (a *AssignStatement) GetToken() token.Token { return a.Token }
func (a *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(a.Name.String() + " = " + a.Value.String() + ";")
	return out.String()
}

// ReturnStatement
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) GetToken() token.Token { return r.Token }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(r.TokenLiteral() + " ")
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
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) GetToken() token.Token { return e.Token }
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

// BlockStatement
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
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

// FunctionStatement / FunctionLiteral
type FunctionStatement struct {
	Token      token.Token // token.FUNC
	Name       *Identifier // optional for anonymous functions if any
	Parameters []*Parameter
	ReturnType string
	Body       *BlockStatement
}

func (f *FunctionStatement) statementNode()       {}
func (f *FunctionStatement) expressionNode()      {}
func (f *FunctionStatement) TokenLiteral() string { return f.Token.Literal }
func (f *FunctionStatement) GetToken() token.Token { return f.Token }
func (f *FunctionStatement) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	if f.Name != nil {
		out.WriteString(f.Name.String())
	}
	out.WriteString("(")
	for i, p := range f.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(")")
	if f.ReturnType != "" {
		out.WriteString(" " + f.ReturnType)
	}
	out.WriteString(" ")
	if f.Body != nil {
		out.WriteString(f.Body.String())
	}
	return out.String()
}

// IfExpression / IfStatement
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // *BlockStatement or *IfExpression
}

func (i *IfExpression) statementNode()       {}
func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) GetToken() token.Token { return i.Token }
func (i *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())
	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}
	return out.String()
}

// WhileStatement
type WhileStatement struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileStatement) statementNode()       {}
func (w *WhileStatement) TokenLiteral() string { return w.Token.Literal }
func (w *WhileStatement) GetToken() token.Token { return w.Token }
func (w *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(w.Condition.String())
	out.WriteString(" ")
	out.WriteString(w.Body.String())
	return out.String()
}

// ForStatement (for init; cond; post { ... })
type ForStatement struct {
	Token     token.Token // token.FOR
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (f *ForStatement) statementNode()       {}
func (f *ForStatement) TokenLiteral() string { return f.Token.Literal }
func (f *ForStatement) GetToken() token.Token { return f.Token }
func (f *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if f.Init != nil {
		out.WriteString(f.Init.String() + " ")
	}
	if f.Condition != nil {
		out.WriteString(f.Condition.String() + "; ")
	}
	if f.Post != nil {
		out.WriteString(f.Post.String() + " ")
	}
	out.WriteString(f.Body.String())
	return out.String()
}

// IntegerLiteral
type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral
type FloatLiteral struct {
	Token token.Token // token.FLOAT
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral
type StringLiteral struct {
	Token token.Token // token.STRING
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// Boolean
type Boolean struct {
	Token token.Token // token.TRUE / token.FALSE
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) GetToken() token.Token { return b.Token }
func (b *Boolean) String() string       { return b.Token.Literal }

// PrefixExpression
type PrefixExpression struct {
	Token    token.Token // token.BANG, token.MINUS
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression
type InfixExpression struct {
	Token    token.Token // operator token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) GetToken() token.Token { return ie.Token }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression (e.g. i++, i--)
type PostfixExpression struct {
	Token    token.Token // token.INCREMENT or token.DECREMENT
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) statementNode()       {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// CallExpression
type CallExpression struct {
	Token     token.Token // token.LPAREN
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, a := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(a.String())
	}
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, el := range al.Elements {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(el.String())
	}
	out.WriteString("]")
	return out.String()
}

// IndexExpression
type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
