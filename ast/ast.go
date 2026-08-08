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

// Program root node
type Program struct {
	Statements []Statement
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

func (p *Program) GetToken() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].GetToken()
	}
	return token.Token{}
}

// VarStatement represent: var x = value
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " " + vs.Name.String() + " = ")
	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represent: return value
type ReturnStatement struct {
	Token       token.Token // token.RETURN
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

// ExpressionStatement represent: x + y;
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
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

// Identifier
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

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

// StringLiteral
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return `"` + sl.Value + `"` }

// BooleanLiteral
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) GetToken() token.Token { return bl.Token }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// PrefixExpression
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. ! or -
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
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) GetToken() token.Token { return oe.Token }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression (e.g. i++, i--)
type PostfixExpression struct {
	Token    token.Token // operator token (++, --)
	Left     Expression  // should be identifier
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
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

// AssignExpression (e.g. x = value)
type AssignExpression struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) GetToken() token.Token { return ae.Token }
func (ae *AssignExpression) String() string {
	return ae.Name.String() + " = " + ae.Value.String()
}

// ShortAssignExpression (e.g. x := value)
type ShortAssignExpression struct {
	Token token.Token // token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (sae *ShortAssignExpression) expressionNode()      {}
func (sae *ShortAssignExpression) TokenLiteral() string { return sae.Token.Literal }
func (sae *ShortAssignExpression) GetToken() token.Token { return sae.Token }
func (sae *ShortAssignExpression) String() string {
	return sae.Name.String() + " := " + sae.Value.String()
}

// BlockStatement
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// IfExpression
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) GetToken() token.Token { return ie.Token }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if " + ie.Condition.String() + " " + ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else " + ie.Alternative.String())
	}
	return out.String()
}

// WhileExpression
type WhileExpression struct {
	Token     token.Token // The 'while' token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) GetToken() token.Token { return we.Token }
func (we *WhileExpression) String() string {
	return "while " + we.Condition.String() + " { " + we.Body.String() + " }"
}

// ForExpression
type ForExpression struct {
	Token     token.Token // The 'for' token
	Init      Node        // can be expression (ShortAssignExpression) or Statement (VarStatement) or nil
	Condition Expression
	Post      Expression
	Body      *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) GetToken() token.Token { return fe.Token }
func (fe *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fe.Init != nil {
		out.WriteString(fe.Init.String() + "; ")
	}
	if fe.Condition != nil {
		out.WriteString(fe.Condition.String())
	}
	out.WriteString("; ")
	if fe.Post != nil {
		out.WriteString(fe.Post.String())
	}
	out.WriteString(" { " + fe.Body.String() + " }")
	return out.String()
}

// Parameter definition with optional type
type Parameter struct {
	Name *Identifier
	Type token.Token // e.g., INT_TYPE, FLOAT_TYPE, etc.
}

// FunctionLiteral
type FunctionLiteral struct {
	Token      token.Token // The 'func' token
	Name       *Identifier // optional named function identifier
	Parameters []*Parameter
	ReturnType token.Token // optional return type token
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	if fl.Name != nil {
		out.WriteString(fl.Name.String())
	}
	out.WriteString("(")
	for i, p := range fl.Parameters {
		out.WriteString(p.Name.String())
		if p.Type.Type != "" {
			out.WriteString(" " + p.Type.Literal)
		}
		if i < len(fl.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	if fl.ReturnType.Type != "" {
		out.WriteString(fl.ReturnType.Literal + " ")
	}
	out.WriteString("{ " + fl.Body.String() + " }")
	return out.String()
}

// CallExpression
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String() + "(")
	for i, a := range ce.Arguments {
		out.WriteString(a.String())
		if i < len(ce.Arguments)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	return out.String()
}

// ArrayLiteral
type ArrayLiteral struct {
	Token    token.Token // '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")
	for i, el := range al.Elements {
		out.WriteString(el.String())
		if i < len(al.Elements)-1 {
			out.WriteString(", ")
		}
	}
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
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String() + "[" + ie.Index.String() + "]")
	out.WriteString(")")
	return out.String()
}
