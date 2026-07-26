package ast

import (
	"bytes"
	"github.com/username/rustx/token"
)

// The base Node interface
type Node interface {
	GetToken() token.Token
	String() string
}

// All statements implement this
type Statement interface {
	Node
	statementNode()
}

// All expressions implement this
type Expression interface {
	Node
	expressionNode()
}

// Program node is the root of the AST
type Program struct {
	Statements []Statement
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

// VarStatement represents: var x = value
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  string      // optional type (e.g. "int")
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) expressionNode()      {}
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString(vs.Token.Literal + " " + vs.Name.String())
	if vs.Type != "" {
		out.WriteString(" " + vs.Type)
	}
	if vs.Value != nil {
		out.WriteString(" = " + vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents: return expression
type ReturnStatement struct {
	Token       token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Token.Literal + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents an expression evaluated on its own
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) GetToken() token.Token { return es.Token }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement represents a sequence of statements inside curly braces
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

// Identifier represents variable names
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral represents int values
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents float values
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// BooleanLiteral represents true/false
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) GetToken() token.Token { return bl.Token }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// StringLiteral represents a string value
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// ArrayLiteral represents [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
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

// IndexExpression represents arr[idx]
type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
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

// PrefixExpression represents prefix operators like ! or -
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents binary operators
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
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

// PostfixExpression represents post-increment/decrement like x++ or x--
type PostfixExpression struct {
	Token    token.Token // ++ or --
	Left     Expression  // Must be an Identifier or IndexExpression (usually Identifier)
	Operator string
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// AssignmentExpression represents x = value or x := value
type AssignmentExpression struct {
	Token    token.Token // = or :=
	Left     Expression  // Usually Identifier or IndexExpression
	Operator string      // "=" or ":="
	Value    Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) GetToken() token.Token { return ae.Token }
func (ae *AssignmentExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ae.Left.String())
	out.WriteString(" " + ae.Operator + " ")
	out.WriteString(ae.Value.String())
	out.WriteString(")")
	return out.String()
}

// Parameter represents a function parameter with optional type
type Parameter struct {
	Token token.Token // token.IDENT
	Name  *Identifier
	Type  string // e.g. "int" (optional)
}

func (p *Parameter) GetToken() token.Token { return p.Token }
func (p *Parameter) String() string {
	if p.Type != "" {
		return p.Name.String() + " " + p.Type
	}
	return p.Name.String()
}

// FunctionLiteral represents: func name(a, b) int { ... }
type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Name       *Identifier // Optional (nil for anonymous function literals)
	Parameters []*Parameter
	ReturnType string // Optional return type (e.g. "int")
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(fl.Token.Literal)
	if fl.Name != nil {
		out.WriteString(" " + fl.Name.String())
	}
	out.WriteString("(")
	for i, p := range fl.Parameters {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(p.String())
	}
	out.WriteString(") ")
	if fl.ReturnType != "" {
		out.WriteString(fl.ReturnType + " ")
	}
	out.WriteString(fl.Body.String())
	return out.String()
}

// CallExpression represents: add(1, 2)
type CallExpression struct {
	Token     token.Token // token.LPAREN
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}
	out.WriteString(")")
	return out.String()
}

// IfExpression represents: if cond { ... } else { ... }
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // Optional
}

func (ie *IfExpression) expressionNode()      {}
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

// WhileStatement represents: while condition { ... }
type WhileStatement struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) GetToken() token.Token { return ws.Token }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Body.String())
	return out.String()
}

// ForStatement represents: for init; cond; post { ... } OR for cond { ... } OR for { ... }
type ForStatement struct {
	Token     token.Token // token.FOR
	Init      Expression  // assignment / var statement / nil
	Condition Expression  // Expression or nil
	Post      Expression  // assignment / postfix / nil
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) GetToken() token.Token { return fs.Token }
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
