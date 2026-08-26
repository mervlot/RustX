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
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// TypeAnnotation (Optional type annotations, e.g. int, float, string)
type TypeAnnotation struct {
	Token token.Token
	Value string
}

func (t *TypeAnnotation) Node()               {}
func (t *TypeAnnotation) TokenLiteral() string { return t.Token.Literal }
func (t *TypeAnnotation) GetToken() token.Token { return t.Token }
func (t *TypeAnnotation) String() string       { return t.Value }

// VarStatement (var x int = 5 or var x = 5)
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  *TypeAnnotation
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
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

// ShortVarStatement (x := 5)
type ShortVarStatement struct {
	Token token.Token // token.COLON_EQ
	Name  *Identifier
	Value Expression
}

func (svs *ShortVarStatement) statementNode()       {}
func (svs *ShortVarStatement) TokenLiteral() string { return svs.Token.Literal }
func (svs *ShortVarStatement) GetToken() token.Token { return svs.Token }
func (svs *ShortVarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(svs.Name.String() + " := ")
	if svs.Value != nil {
		out.WriteString(svs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// AssignmentStatement (x = 5)
type AssignmentStatement struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

func (as *AssignmentStatement) statementNode()       {}
func (as *AssignmentStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignmentStatement) GetToken() token.Token { return as.Token }
func (as *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String() + " = ")
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement
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

// ExpressionStatement
type ExpressionStatement struct {
	Token      token.Token // first token of the expression
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

// BlockStatement
type BlockStatement struct {
	Token      token.Token // token.LBRACE
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

// Literals
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanLiteral) GetToken() token.Token { return b.Token }
func (b *BooleanLiteral) String() string       { return b.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

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

// Prefix Expression (-x, !x)
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
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// Postfix Expression (x++, x--)
type PostfixExpression struct {
	Token    token.Token // token.INC or token.DEC
	Left     Expression
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

// Infix Expression (a + b, a == b, a && b)
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
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// IfExpression
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
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

// WhileStatement
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
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement (for init; condition; post { body })
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
		out.WriteString(fs.Init.String() + " ")
	}
	if fs.Condition != nil {
		out.WriteString(fs.Condition.String() + "; ")
	}
	if fs.Post != nil {
		out.WriteString(fs.Post.String() + " ")
	}
	out.WriteString(fs.Body.String())
	return out.String()
}

// Parameter definition in function
type Parameter struct {
	Name *Identifier
	Type *TypeAnnotation
}

// FunctionLiteral (func name(a int, b int) int { ... })
type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Name       *Identifier // optional function name for top-level/named funcs
	Parameters []*Parameter
	ReturnType *TypeAnnotation
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		pStr := p.Name.String()
		if p.Type != nil {
			pStr += " " + p.Type.String()
		}
		params = append(params, pStr)
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
	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
