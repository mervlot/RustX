package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

// Node represents any node in the AST.
type Node interface {
	TokenLiteral() string
	String() string
	GetToken() token.Token
}

// Statement represents a statement node.
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node.
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST.
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

// VarStatement represents: var name [type] = value;
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  string // optional type annotation, e.g. "int"
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

// ExpressionStatement represents a statement consisting of a single expression.
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

// ReturnStatement represents a return statement.
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral())
	if rs.ReturnValue != nil {
		out.WriteString(" " + rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// BlockStatement represents a series of statements inside braces.
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String() + " ")
	}
	out.WriteString("}")
	return out.String()
}

// Identifier represents an identifier.
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents an integer.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a float.
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// BooleanLiteral represents a boolean.
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) GetToken() token.Token { return bl.Token }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// StringLiteral represents a string.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// ArrayLiteral represents an array.
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
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
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// IndexExpression represents indexing of an array or string.
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

// PrefixExpression represents a prefix operator expression.
type PrefixExpression struct {
	Token    token.Token // prefix token, e.g. ! or -
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

// InfixExpression represents an infix operator expression.
type InfixExpression struct {
	Token    token.Token // operator token, e.g. +
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

// AssignExpression represents: Left = Right
type AssignExpression struct {
	Token token.Token // token.ASSIGN
	Left  Expression  // e.g. Identifier or IndexExpression
	Value Expression
}

func (ae *AssignExpression) expressionNode()      {}
func (ae *AssignExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignExpression) GetToken() token.Token { return ae.Token }
func (ae *AssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(" = ")
	out.WriteString(ae.Value.String())
	out.WriteString(")")
	return out.String()
}

// ColonAssignExpression represents: Left := Right
type ColonAssignExpression struct {
	Token token.Token // token.COLON_ASSIGN
	Left  Expression  // e.g. Identifier
	Value Expression
}

func (cae *ColonAssignExpression) expressionNode()      {}
func (cae *ColonAssignExpression) TokenLiteral() string { return cae.Token.Literal }
func (cae *ColonAssignExpression) GetToken() token.Token { return cae.Token }
func (cae *ColonAssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(cae.Left.String())
	out.WriteString(" := ")
	out.WriteString(cae.Value.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression represents: Left++ or Left--
type PostfixExpression struct {
	Token    token.Token // ++ or --
	Left     Expression  // e.g. Identifier
	Operator string      // ++ or --
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

// IfExpression represents an if / else statement or expression.
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

// WhileExpression represents: while condition { body }
type WhileExpression struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) GetToken() token.Token { return we.Token }
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())
	return out.String()
}

// ForExpression represents Go-style for: for init; condition; post { body }
type ForExpression struct {
	Token     token.Token // token.FOR
	Init      Statement   // can be empty (nil)
	Condition Expression  // can be empty (nil)
	Post      Statement   // can be empty (nil)
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
	} else {
		out.WriteString("; ")
	}
	if fe.Condition != nil {
		out.WriteString(fe.Condition.String() + "; ")
	} else {
		out.WriteString("; ")
	}
	if fe.Post != nil {
		out.WriteString(fe.Post.String() + " ")
	} else {
		out.WriteString(" ")
	}
	out.WriteString(fe.Body.String())
	return out.String()
}

// Parameter represents a function parameter: name [type]
type Parameter struct {
	Token token.Token // token.IDENT
	Name  string
	Type  string // optional type, e.g. "int"
}

func (p *Parameter) TokenLiteral() string { return p.Token.Literal }
func (p *Parameter) GetToken() token.Token { return p.Token }
func (p *Parameter) String() string {
	if p.Type != "" {
		return p.Name + " " + p.Type
	}
	return p.Name
}

// FunctionLiteral represents: func [name](params) [type] { body }
type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Name       string      // optional name (for named function literal / recursion)
	Parameters []*Parameter
	ReturnType string // optional, e.g. "int"
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral() + " ")
	if fl.Name != "" {
		out.WriteString(fl.Name)
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

// CallExpression represents: function(arguments)
type CallExpression struct {
	Token     token.Token // token.LPAREN
	Function  Expression  // e.g. Identifier or FunctionLiteral
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
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}
