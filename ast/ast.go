package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

type Node interface {
	GetToken() token.Token
	TokenLiteral() string
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
	return token.Token{}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// Identifier expression
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) TokenLiteral() string  { return i.Token.Literal }
func (i *Identifier) String() string        { return i.Value }

// IntegerLiteral
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) TokenLiteral() string  { return il.Token.Literal }
func (il *IntegerLiteral) String() string        { return il.Token.Literal }

// FloatLiteral
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) TokenLiteral() string  { return fl.Token.Literal }
func (fl *FloatLiteral) String() string        { return fl.Token.Literal }

// StringLiteral
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) TokenLiteral() string  { return sl.Token.Literal }
func (sl *StringLiteral) String() string        { return sl.Token.Literal }

// Boolean
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) GetToken() token.Token { return b.Token }
func (b *Boolean) TokenLiteral() string  { return b.Token.Literal }
func (b *Boolean) String() string        { return b.Token.Literal }

// ArrayLiteral
type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) TokenLiteral() string  { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression
type IndexExpression struct {
	Token token.Token // '['
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) TokenLiteral() string  { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// PrefixExpression
type PrefixExpression struct {
	Token    token.Token // prefix token, e.g. ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) TokenLiteral() string  { return pe.Token.Literal }
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
	Token    token.Token // operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) GetToken() token.Token { return oe.Token }
func (oe *InfixExpression) TokenLiteral() string  { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression (for ++ and --)
type PostfixExpression struct {
	Token    token.Token // ++ or --
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) TokenLiteral() string  { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// AssignmentExpression (for = and :=)
type AssignmentExpression struct {
	Token token.Token // = or :=
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) GetToken() token.Token { return ae.Token }
func (ae *AssignmentExpression) TokenLiteral() string  { return ae.Token.Literal }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.String())
	out.WriteString(" " + ae.Token.Literal + " ")
	out.WriteString(ae.Value.String())
	return out.String()
}

// CallExpression
type CallExpression struct {
	Token     token.Token // '('
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) TokenLiteral() string  { return ce.Token.Literal }
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

// Parameter definition in function
type Parameter struct {
	Name *Identifier
	Type string // e.g., "int", "string"
}

// FunctionLiteral
type FunctionLiteral struct {
	Token      token.Token // 'func'
	Name       string      // optional function name (e.g. func add(...) ...)
	Parameters []*Parameter
	ReturnType string      // optional return type annotation
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) TokenLiteral() string  { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		if p.Type != "" {
			params = append(params, p.Name.String()+" "+p.Type)
		} else {
			params = append(params, p.Name.String())
		}
	}
	out.WriteString(fl.TokenLiteral())
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

// IfExpression
type IfExpression struct {
	Token       token.Token // 'if'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) GetToken() token.Token { return ie.Token }
func (ie *IfExpression) TokenLiteral() string  { return ie.Token.Literal }
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

// WhileStatement
type WhileStatement struct {
	Token     token.Token // 'while'
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) GetToken() token.Token { return ws.Token }
func (ws *WhileStatement) TokenLiteral() string  { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement (for init; condition; post { body })
type ForStatement struct {
	Token     token.Token // 'for'
	Init      Statement
	Condition Expression
	Post      Statement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) GetToken() token.Token { return fs.Token }
func (fs *ForStatement) TokenLiteral() string  { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
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
	out.WriteString(" ")
	out.WriteString(fs.Body.String())
	return out.String()
}

// VarStatement (var x int = 5)
type VarStatement struct {
	Token token.Token // 'var'
	Name  *Identifier
	Type  string // optional type annotation
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) TokenLiteral() string  { return vs.Token.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	if vs.Type != "" {
		out.WriteString(" " + vs.Type)
	}
	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	return out.String()
}

// ReturnStatement
type ReturnStatement struct {
	Token       token.Token // 'return'
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }
func (rs *ReturnStatement) TokenLiteral() string  { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}

// ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) GetToken() token.Token { return es.Token }
func (es *ExpressionStatement) TokenLiteral() string  { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement
type BlockStatement struct {
	Token      token.Token // '{'
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) TokenLiteral() string  { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String() + "; ")
	}
	out.WriteString("}")
	return out.String()
}
