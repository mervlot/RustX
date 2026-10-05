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
	Tok   token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Tok.Literal }
func (i *Identifier) GetToken() token.Token { return i.Tok }
func (i *Identifier) String() string       { return i.Value }

// TypeAnnotation (e.g. int, float, string, bool)
type TypeAnnotation struct {
	Tok   token.Token
	Value string
}

func (t *TypeAnnotation) Node()               {}
func (t *TypeAnnotation) TokenLiteral() string { return t.Tok.Literal }
func (t *TypeAnnotation) GetToken() token.Token { return t.Tok }
func (t *TypeAnnotation) String() string       { return t.Value }

// VarStatement: var x int = 5 or var x = 5
type VarStatement struct {
	Tok   token.Token // the token.VAR token
	Name  *Identifier
	Type  *TypeAnnotation // optional
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Tok.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Tok }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.TokenLiteral() + " " + vs.Name.String())
	if vs.Type != nil {
		out.WriteString(" " + vs.Type.String())
	}
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ShortVarStatement: x := 5
type ShortVarStatement struct {
	Tok   token.Token // the token.COLON_ASSIGN token
	Name  *Identifier
	Value Expression
}

func (sv *ShortVarStatement) statementNode()       {}
func (sv *ShortVarStatement) TokenLiteral() string { return sv.Tok.Literal }
func (sv *ShortVarStatement) GetToken() token.Token { return sv.Tok }
func (sv *ShortVarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(sv.Name.String() + " := " + sv.Value.String() + ";")
	return out.String()
}

// AssignStatement: x = 5 or arr[0] = 5
type AssignStatement struct {
	Tok   token.Token // token.ASSIGN
	Name  Expression  // *Identifier or *IndexExpression
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Tok.Literal }
func (as *AssignStatement) GetToken() token.Token { return as.Tok }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String() + " = " + as.Value.String() + ";")
	return out.String()
}

// ReturnStatement: return expr;
type ReturnStatement struct {
	Tok         token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Tok.Literal }
func (rs *ReturnStatement) GetToken() token.Token { return rs.Tok }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement
type ExpressionStatement struct {
	Tok        token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Tok.Literal }
func (es *ExpressionStatement) GetToken() token.Token { return es.Tok }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement
type BlockStatement struct {
	Tok        token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Tok.Literal }
func (bs *BlockStatement) GetToken() token.Token { return bs.Tok }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String() + " ")
	}
	out.WriteString("}")
	return out.String()
}

// Literals
type IntegerLiteral struct {
	Tok   token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Tok.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Tok }
func (il *IntegerLiteral) String() string       { return il.Tok.Literal }

type FloatLiteral struct {
	Tok   token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FloatLiteral) String() string       { return fl.Tok.Literal }

type BooleanLiteral struct {
	Tok   token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Tok.Literal }
func (b *BooleanLiteral) GetToken() token.Token { return b.Tok }
func (b *BooleanLiteral) String() string       { return b.Tok.Literal }

type StringLiteral struct {
	Tok   token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Tok.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Tok }
func (sl *StringLiteral) String() string       { return sl.Tok.Literal }

type ArrayLiteral struct {
	Tok      token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Tok.Literal }
func (al *ArrayLiteral) GetToken() token.Token { return al.Tok }
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

type IndexExpression struct {
	Tok   token.Token // token.LBRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Tok.Literal }
func (ie *IndexExpression) GetToken() token.Token { return ie.Tok }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

type PrefixExpression struct {
	Tok      token.Token // token.BANG, token.MINUS
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Tok.Literal }
func (pe *PrefixExpression) GetToken() token.Token { return pe.Tok }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Tok      token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Tok.Literal }
func (oe *InfixExpression) GetToken() token.Token { return oe.Tok }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}

type PostfixExpression struct {
	Tok      token.Token // token.INC or token.DEC
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Tok.Literal }
func (pe *PostfixExpression) GetToken() token.Token { return pe.Tok }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

type FunctionParameter struct {
	Name *Identifier
	Type *TypeAnnotation // optional
}

func (fp *FunctionParameter) String() string {
	str := fp.Name.String()
	if fp.Type != nil {
		str += " " + fp.Type.String()
	}
	return str
}

type FunctionLiteral struct {
	Tok        token.Token // token.FUNC
	Name       *Identifier // optional named function: func add(...)
	Parameters []*FunctionParameter
	ReturnType *TypeAnnotation // optional
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	if fl.Name != nil {
		out.WriteString(" " + fl.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	if fl.ReturnType != nil {
		out.WriteString(fl.ReturnType.String() + " ")
	}
	out.WriteString(fl.Body.String())
	return out.String()
}

type CallExpression struct {
	Tok       token.Token // token.LPAREN
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Tok.Literal }
func (ce *CallExpression) GetToken() token.Token { return ce.Tok }
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

type IfExpression struct {
	Tok         token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // optional, or could be another IfExpression statement/block
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Tok.Literal }
func (ie *IfExpression) GetToken() token.Token { return ie.Tok }
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

type WhileExpression struct {
	Tok       token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Tok.Literal }
func (we *WhileExpression) GetToken() token.Token { return we.Tok }
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())
	return out.String()
}

type ForExpression struct {
	Tok       token.Token // token.FOR
	Init      Statement   // e.g. i := 0 or var i = 0
	Condition Expression  // e.g. i < 10
	Post      Statement   // e.g. i++ or i = i + 1
	Body      *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Tok.Literal }
func (fe *ForExpression) GetToken() token.Token { return fe.Tok }
func (fe *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fe.Init != nil {
		out.WriteString(fe.Init.String() + " ")
	}
	if fe.Condition != nil {
		out.WriteString(fe.Condition.String() + "; ")
	}
	if fe.Post != nil {
		out.WriteString(fe.Post.String() + " ")
	}
	out.WriteString(fe.Body.String())
	return out.String()
}
