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

// Identifier represents variable/function names or type annotations.
type Identifier struct {
	Tok   token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Tok.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) GetToken() token.Token { return i.Tok }

// VarStatement represents `var x = 5` or `var x int = 5`
type VarStatement struct {
	Tok   token.Token // token.VAR
	Name  *Identifier
	Type  *Identifier // optional type annotation
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Tok.Literal }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString("var " + vs.Name.String())
	if vs.Type != nil {
		out.WriteString(" " + vs.Type.String())
	}
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
func (vs *VarStatement) GetToken() token.Token { return vs.Tok }

// ReturnStatement represents `return expr`
type ReturnStatement struct {
	Tok         token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Tok.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Tok.Literal + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
func (rs *ReturnStatement) GetToken() token.Token { return rs.Tok }

// ExpressionStatement allows standalone expressions to act as statements.
type ExpressionStatement struct {
	Tok        token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Tok.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) GetToken() token.Token { return es.Tok }

// BlockStatement represents `{ statements... }`
type BlockStatement struct {
	Tok        token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Tok.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}
func (bs *BlockStatement) GetToken() token.Token { return bs.Tok }

// IntegerLiteral represents integer numbers.
type IntegerLiteral struct {
	Tok   token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Tok.Literal }
func (il *IntegerLiteral) String() string       { return il.Tok.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Tok }

// FloatLiteral represents floating point numbers.
type FloatLiteral struct {
	Tok   token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FloatLiteral) String() string       { return fl.Tok.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Tok }

// StringLiteral represents string values.
type StringLiteral struct {
	Tok   token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Tok.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }
func (sl *StringLiteral) GetToken() token.Token { return sl.Tok }

// BooleanLiteral represents `true` or `false`.
type BooleanLiteral struct {
	Tok   token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Tok.Literal }
func (b *BooleanLiteral) String() string       { return b.Tok.Literal }
func (b *BooleanLiteral) GetToken() token.Token { return b.Tok }

// PrefixExpression represents `-x` or `!x`
type PrefixExpression struct {
	Tok      token.Token // operator token, e.g. ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Tok.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
func (pe *PrefixExpression) GetToken() token.Token { return pe.Tok }

// InfixExpression represents binary operations like `a + b`, `a == b`.
type InfixExpression struct {
	Tok      token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode()      {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Tok.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}
func (oe *InfixExpression) GetToken() token.Token { return oe.Tok }

// PostfixExpression represents postfix operations like `i++` or `j--`.
type PostfixExpression struct {
	Tok      token.Token // ++ or --
	Left     Expression
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Tok.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}
func (pe *PostfixExpression) GetToken() token.Token { return pe.Tok }

// AssignExpression represents re-assignment: `x = 10`
type AssignExpression struct {
	Tok   token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Tok.Literal }
func (ae *AssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.String() + " = " + ae.Value.String())
	return out.String()
}
func (ae *AssignExpression) GetToken() token.Token { return ae.Tok }

// ColonAssignExpression represents shorthand declaration & assignment: `x := 5`
type ColonAssignExpression struct {
	Tok   token.Token // token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (cae *ColonAssignExpression) expressionNode()      {}
func (cae *ColonAssignExpression) TokenLiteral() string { return cae.Tok.Literal }
func (cae *ColonAssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString(cae.Name.String() + " := " + cae.Value.String())
	return out.String()
}
func (cae *ColonAssignExpression) GetToken() token.Token { return cae.Tok }

// IfExpression represents `if condition { consequence } else { alternative }`
type IfExpression struct {
	Tok         token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Tok.Literal }
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
func (ie *IfExpression) GetToken() token.Token { return ie.Tok }

// WhileStatement represents `while condition { body }`
type WhileStatement struct {
	Tok       token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Tok.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}
func (ws *WhileStatement) GetToken() token.Token { return ws.Tok }

// ForStatement represents `for init; condition; post { body }` or `for condition { body }`
type ForStatement struct {
	Tok       token.Token // token.FOR
	Init      Statement   // optional init statement (e.g. VarStatement, AssignExpression, ColonAssignExpression)
	Condition Expression  // condition
	Post      Statement   // optional post statement (e.g. PostfixExpression, AssignExpression)
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Tok.Literal }
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
func (fs *ForStatement) GetToken() token.Token { return fs.Tok }

// Parameter represents a parameter in a function definition: `a int` or `a`
type Parameter struct {
	Tok  token.Token // Parameter identifier token
	Name *Identifier
	Type *Identifier // optional type annotation
}

func (p *Parameter) TokenLiteral() string { return p.Tok.Literal }
func (p *Parameter) String() string {
	if p.Type != nil {
		return p.Name.String() + " " + p.Type.String()
	}
	return p.Name.String()
}
func (p *Parameter) GetToken() token.Token { return p.Tok }

// FunctionLiteral represents `func name(params) ReturnType { body }` or anonymous `func(params) { body }`
type FunctionLiteral struct {
	Tok        token.Token // token.FUNCTION
	Name       *Identifier // optional function name
	Parameters []*Parameter
	ReturnType *Identifier     // optional return type annotation
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) statementNode()       {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("func ")
	if fl.Name != nil {
		out.WriteString(fl.Name.String())
	}
	out.WriteString("(")
	for i, p := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(") ")
	if fl.ReturnType != nil {
		out.WriteString(fl.ReturnType.String() + " ")
	}
	out.WriteString(fl.Body.String())
	return out.String()
}
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Tok }

// CallExpression represents `fn(arg1, arg2)`
type CallExpression struct {
	Tok       token.Token // '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Tok.Literal }
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
func (ce *CallExpression) GetToken() token.Token { return ce.Tok }

// ArrayLiteral represents `[1, 2, 3]`
type ArrayLiteral struct {
	Tok      token.Token // '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Tok.Literal }
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
func (al *ArrayLiteral) GetToken() token.Token { return al.Tok }

// IndexExpression represents `arr[index]`
type IndexExpression struct {
	Tok   token.Token // '[' token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Tok.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
func (ie *IndexExpression) GetToken() token.Token { return ie.Tok }
