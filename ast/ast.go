package ast

import (
	"bytes"
	"strings"

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

// Program is the root node of every AST.
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
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// VarStatement represents `var x int = 5;`
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  string // optional type annotation e.g. "int"
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " " + vs.Name.String())
	if vs.Type != "" {
		out.WriteString(" " + vs.Type)
	}
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// AssignStatement represents `x := 5;` or `x = 5;`
type AssignStatement struct {
	Token       token.Token // token.COLON_ASSIGN or token.ASSIGN
	Name        *Identifier
	IsShorthand bool
	Value       Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) expressionNode()      {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) GetToken() token.Token { return as.Token }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	if as.IsShorthand {
		out.WriteString(" := ")
	} else {
		out.WriteString(" = ")
	}
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents `return expr;`
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement wrapper for expressions as statements
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
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
func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// Parameter for functions
type Parameter struct {
	Name *Identifier
	Type string // optional
}

func (p *Parameter) String() string {
	if p.Type != "" {
		return p.Name.String() + " " + p.Type
	}
	return p.Name.String()
}

// FunctionStatement represents `func add(a int, b int) int { ... }`
type FunctionStatement struct {
	Token      token.Token // token.FUNC
	Name       *Identifier
	Parameters []*Parameter
	ReturnType string
	Body       *BlockStatement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) expressionNode()      {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) GetToken() token.Token { return fs.Token }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func " + fs.Name.String() + "(" + strings.Join(params, ", ") + ")")
	if fs.ReturnType != "" {
		out.WriteString(" " + fs.ReturnType)
	}
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// IntegerLiteral
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) GetToken() token.Token { return b.Token }
func (b *Boolean) String() string       { return b.Token.Literal }

// StringLiteral
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// PrefixExpression represents `-a` or `!a`
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + pe.Operator + pe.Right.String() + ")")
	return out.String()
}

// InfixExpression represents `a + b`
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) GetToken() token.Token { return ie.Token }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")")
	return out.String()
}

// PostfixExpression represents `i++` or `i--`
type PostfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) String() string {
	return pe.Left.String() + pe.Operator
}

// IfExpression represents `if cond { ... } else { ... }` or `if cond { ... } else if cond2 { ... }`
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // *BlockStatement or *IfExpression / ExpressionStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) statementNode()       {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) GetToken() token.Token { return ie.Token }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if " + ie.Condition.String() + " " + ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(" else " + ie.Alternative.String())
	}
	return out.String()
}

// WhileStatement represents `while cond { ... }`
type WhileStatement struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) GetToken() token.Token { return ws.Token }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while " + ws.Condition.String() + " " + ws.Body.String())
	return out.String()
}

// ForStatement represents `for init; cond; post { body }` or `for cond { body }`
type ForStatement struct {
	Token     token.Token // token.FOR
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) GetToken() token.Token { return fs.Token }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fs.Init != nil {
		out.WriteString(fs.Init.String() + "; ")
	}
	if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
	}
	if fs.Post != nil {
		out.WriteString("; " + fs.Post.String())
	}
	out.WriteString(" " + fs.Body.String())
	return out.String()
}

// CallExpression represents `add(1, 2)`
type CallExpression struct {
	Token     token.Token // '(' token
	Function  Expression  // Identifier or FunctionStatement
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String() + "(" + strings.Join(args, ", ") + ")")
	return out.String()
}

// IndexExpression represents `arr[index]`
type IndexExpression struct {
	Token token.Token // '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) String() string {
	return ie.Left.String() + "[" + ie.Index.String() + "]"
}

// ArrayLiteral represents `[1, 2, 3]`
type ArrayLiteral struct {
	Token    token.Token // '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[" + strings.Join(elements, ", ") + "]")
	return out.String()
}
