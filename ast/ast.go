package ast

import (
	"bytes"

	"github.com/username/rustx/token"
)

// Node represents any node in the AST.
type Node interface {
	TokenLiteral() string
	String() string
	GetToken() token.Token
}

// Statement is an AST node that represents a statement.
type Statement interface {
	Node
	statementNode()
}

// Expression is an AST node that represents an expression.
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

// ExpressionStatement represents a statement that consists of a single expression.
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

// VarStatement represents 'var name = value' or 'var name Type = value'.
type VarStatement struct {
	Tok   token.Token // The token.VAR token
	Name  *Identifier
	Type  string // optional type annotation (e.g., "int", "string")
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Tok.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Tok }
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

// ReturnStatement represents 'return expression'.
type ReturnStatement struct {
	Tok         token.Token // The token.RETURN token
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

// BlockStatement represents '{ statements }'.
type BlockStatement struct {
	Tok        token.Token // The token.LBRACE token
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

// Identifier represents variable names, function names, etc.
type Identifier struct {
	Tok   token.Token // The token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Tok.Literal }
func (i *Identifier) GetToken() token.Token { return i.Tok }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents integer values.
type IntegerLiteral struct {
	Tok   token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Tok.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Tok }
func (il *IntegerLiteral) String() string       { return il.Tok.Literal }

// FloatLiteral represents floating point values.
type FloatLiteral struct {
	Tok   token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FloatLiteral) String() string       { return fl.Tok.Literal }

// Boolean represents true or false.
type Boolean struct {
	Tok   token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Tok.Literal }
func (b *Boolean) GetToken() token.Token { return b.Tok }
func (b *Boolean) String() string       { return b.Tok.Literal }

// StringLiteral represents string values.
type StringLiteral struct {
	Tok   token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Tok.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Tok }
func (sl *StringLiteral) String() string       { return sl.Tok.Literal }

// PrefixExpression represents operations like -x, !x.
type PrefixExpression struct {
	Tok      token.Token // The prefix token, e.g. ! or -
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

// InfixExpression represents binary operations like a + b, x == y.
type InfixExpression struct {
	Tok      token.Token // The operator token, e.g. +
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

// PostfixExpression represents operations like x++ or y--.
type PostfixExpression struct {
	Tok      token.Token // The ++ or -- token
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

// ShortAssignExpression represents 'name := value'.
type ShortAssignExpression struct {
	Tok   token.Token // The := token
	Name  *Identifier
	Value Expression
}

func (sae *ShortAssignExpression) expressionNode()      {}
func (sae *ShortAssignExpression) TokenLiteral() string { return sae.Tok.Literal }
func (sae *ShortAssignExpression) GetToken() token.Token { return sae.Tok }
func (sae *ShortAssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString(sae.Name.String())
	out.WriteString(" := ")
	out.WriteString(sae.Value.String())
	return out.String()
}

// AssignmentExpression represents 'name = value'.
type AssignmentExpression struct {
	Tok   token.Token // The = token
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Tok.Literal }
func (ae *AssignmentExpression) GetToken() token.Token { return ae.Tok }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())
	return out.String()
}

// Parameter represents a function parameter with an optional type annotation.
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

// FunctionLiteral represents 'func name(a int, b int) int { ... }' or anonymous functions.
type FunctionLiteral struct {
	Tok        token.Token // The 'func' token
	Name       string      // optional function name (if named func declaration)
	Parameters []*Parameter
	ReturnType string
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("func")
	if fl.Name != "" {
		out.WriteString(" " + fl.Name)
	}
	out.WriteString("(")
	for i, p := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(")")
	if fl.ReturnType != "" {
		out.WriteString(" " + fl.ReturnType)
	}
	out.WriteString(" ")
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents 'function(arg1, arg2)'.
type CallExpression struct {
	Tok       token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Tok.Literal }
func (ce *CallExpression) GetToken() token.Token { return ce.Tok }
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

// IfExpression represents 'if condition { ... } else { ... }'.
type IfExpression struct {
	Tok         token.Token // The 'if' token
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

// WhileStatement represents 'while condition { ... }'.
type WhileStatement struct {
	Tok       token.Token // The 'while' token
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

// ForStatement represents 'for init; condition; post { body }' or 'for condition { body }'.
type ForStatement struct {
	Tok       token.Token // The 'for' token
	Init      Statement   // optional: var i = 0 or i := 0 or i = 0
	Condition Expression  // optional
	Post      Expression  // optional: i++ or i = i + 1
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
