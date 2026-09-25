package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

// The base Node interface
type Node interface {
	TokenLiteral() string
	String() string
	GetToken() token.Token
}

// Statement interface
type Statement interface {
	Node
	statementNode()
}

// Expression interface
type Expression interface {
	Node
	expressionNode()
}

// Program node is the root node of every AST
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
	Tok   token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Tok.Literal }
func (i *Identifier) GetToken() token.Token { return i.Tok }
func (i *Identifier) String() string       { return i.Value }

// TypeAnnotation represents optional type annotations like "int", "string", etc.
type TypeAnnotation struct {
	Tok   token.Token
	Value string
}

func (t *TypeAnnotation) NodeLiteral() string  { return t.Tok.Literal }
func (t *TypeAnnotation) GetToken() token.Token { return t.Tok }
func (t *TypeAnnotation) String() string       { return t.Value }

// VarStatement represents `var x = expr;` or `var x int = expr;`
type VarStatement struct {
	Tok   token.Token // token.VAR
	Name  *Identifier
	Type  *TypeAnnotation
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
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ColonAssignStatement represents `x := expr;`
type ColonAssignStatement struct {
	Tok   token.Token // token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (ca *ColonAssignStatement) statementNode()       {}
func (ca *ColonAssignStatement) TokenLiteral() string { return ca.Tok.Literal }
func (ca *ColonAssignStatement) GetToken() token.Token { return ca.Tok }
func (ca *ColonAssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ca.Name.String() + " := " + ca.Value.String() + ";")
	return out.String()
}

// AssignStatement represents `x = expr;`
type AssignStatement struct {
	Tok   token.Token // token.ASSIGN
	Name  *Identifier
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

// ReturnStatement represents `return expr;`
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

// ExpressionStatement represents a standalone expression as a statement
type ExpressionStatement struct {
	Tok        token.Token // first token of expression
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

// BlockStatement represents `{ stmt1; stmt2; }`
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
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

// IntegerLiteral
type IntegerLiteral struct {
	Tok   token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Tok.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Tok }
func (il *IntegerLiteral) String() string       { return il.Tok.Literal }

// FloatLiteral
type FloatLiteral struct {
	Tok   token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FloatLiteral) String() string       { return fl.Tok.Literal }

// StringLiteral
type StringLiteral struct {
	Tok   token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Tok.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Tok }
func (sl *StringLiteral) String() string       { return sl.Tok.Literal }

// BooleanLiteral
type BooleanLiteral struct {
	Tok   token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Tok.Literal }
func (bl *BooleanLiteral) GetToken() token.Token { return bl.Tok }
func (bl *BooleanLiteral) String() string       { return bl.Tok.Literal }

// ArrayLiteral represents `[1, 2, 3]`
type ArrayLiteral struct {
	Tok      token.Token // token.LBRACK
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

// IndexExpression represents `arr[index]`
type IndexExpression struct {
	Tok   token.Token // token.LBRACK
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

// PrefixExpression represents `-x` or `!x`
type PrefixExpression struct {
	Tok      token.Token // prefix token, e.g. - or !
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

// PostfixExpression represents `x++` or `x--`
type PostfixExpression struct {
	Tok      token.Token // token.INC or token.DEC
	Operator string
	Left     Expression
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

// InfixExpression represents `a + b`, `a == b`, etc.
type InfixExpression struct {
	Tok      token.Token // operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Tok.Literal }
func (ie *InfixExpression) GetToken() token.Token { return ie.Tok }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// IfExpression represents `if cond { ... } else { ... }`
type IfExpression struct {
	Tok         token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
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

// WhileStatement represents `while cond { ... }`
type WhileStatement struct {
	Tok       token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Tok.Literal }
func (ws *WhileStatement) GetToken() token.Token { return ws.Tok }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement represents `for init; cond; post { ... }` or `for cond { ... }`
type ForStatement struct {
	Tok       token.Token // token.FOR
	Init      Statement   // optional init, e.g., VarStatement, ColonAssignStatement, AssignStatement
	Condition Expression  // optional cond
	Post      Statement   // optional post, e.g., PostfixExpression as ExpressionStatement or AssignStatement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Tok.Literal }
func (fs *ForStatement) GetToken() token.Token { return fs.Tok }
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

// FunctionParameter represents a parameter with name and optional type annotation
type FunctionParameter struct {
	Name *Identifier
	Type *TypeAnnotation
}

// FunctionLiteral represents `func name(a int, b int) int { ... }` or anonymous `func(a int) { ... }`
type FunctionLiteral struct {
	Tok        token.Token // token.FUNC
	Name       *Identifier // optional named function
	Parameters []*FunctionParameter
	ReturnType *TypeAnnotation
	Body       *BlockStatement
}

func (fl *FunctionLiteral) statementNode()       {}
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		paramStr := p.Name.String()
		if p.Type != nil {
			paramStr += " " + p.Type.String()
		}
		params = append(params, paramStr)
	}

	out.WriteString(fl.TokenLiteral() + " ")
	if fl.Name != nil {
		out.WriteString(fl.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fl.ReturnType != nil {
		out.WriteString(" " + fl.ReturnType.String())
	}
	out.WriteString(" ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents `add(1, 2)`
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
