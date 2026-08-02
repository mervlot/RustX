package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

// Node is the base interface for all AST nodes
type Node interface {
	TokenLiteral() string
	String() string
	GetToken() token.Token // position-aware diagnostics
}

// Statement is an AST node representing an action/statement
type Statement interface {
	Node
	statementNode()
}

// Expression is an AST node representing a value producing operation
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of our AST
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

// VarStatement represents a 'var name = value' statement
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
	Type  string // optional type annotation (ignored during evaluation but parsed)
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
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
func (vs *VarStatement) GetToken() token.Token { return vs.Token }

// ShorthandAssignStatement represents 'name := value'
type ShorthandAssignStatement struct {
	Token token.Token // token.COLON_ASSIGN
	Name  *Identifier
	Value Expression
}

func (sas *ShorthandAssignStatement) statementNode()       {}
func (sas *ShorthandAssignStatement) TokenLiteral() string { return sas.Token.Literal }
func (sas *ShorthandAssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(sas.Name.String() + " := " + sas.Value.String() + ";")
	return out.String()
}
func (sas *ShorthandAssignStatement) GetToken() token.Token { return sas.Token }

// AssignStatement represents 'name = value'
type AssignStatement struct {
	Token token.Token // token.ASSIGN
	Name  *Identifier
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String() + " = " + as.Value.String() + ";")
	return out.String()
}
func (as *AssignStatement) GetToken() token.Token { return as.Token }

// ReturnStatement represents 'return expression'
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}
func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }

// ExpressionStatement wrapper for expressions inside statements
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) GetToken() token.Token { return es.Token }

// BlockStatement is a series of statements inside braces
type BlockStatement struct {
	Token      token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }

// Identifier AST Node
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) GetToken() token.Token { return i.Token }

// IntegerLiteral AST Node
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }

// FloatLiteral AST Node
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }

// BooleanLiteral AST Node
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }
func (bl *BooleanLiteral) GetToken() token.Token { return bl.Token }

// StringLiteral AST Node
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }

// ArrayLiteral AST Node e.g., [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token // token.LBRACKET
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
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
func (al *ArrayLiteral) GetToken() token.Token { return al.Token }

// IndexExpression AST Node e.g., arr[0]
type IndexExpression struct {
	Token token.Token // token.LBRACKET
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}
func (ie *IndexExpression) GetToken() token.Token { return ie.Token }

// PrefixExpression e.g., -5, !true
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }

// PostfixExpression e.g., x++, y--
type PostfixExpression struct {
	Token    token.Token // token.INCREMENT or token.DECREMENT
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode()      {}
func (pe *PostfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}
func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }

// InfixExpression e.g., x + y, x && y
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
func (ie *InfixExpression) GetToken() token.Token { return ie.Token }

// Parameter represents a parameter with name and optional type
type Parameter struct {
	Name *Identifier
	Type string // optional e.g., "int", "string"
}

// FunctionLiteral represents 'func(a, b) { ... }' or named 'func foo(a, b) { ... }'
type FunctionLiteral struct {
	Token      token.Token // token.FUNC
	Name       *Identifier // optional function name for recursive / named registration
	Parameters []*Parameter
	ReturnType string // optional return type annotation
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		s := p.Name.String()
		if p.Type != "" {
			s += " " + p.Type
		}
		params = append(params, s)
	}

	out.WriteString(fl.TokenLiteral())
	if fl.Name != nil {
		out.WriteString(" " + fl.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fl.ReturnType != "" {
		out.WriteString(" " + fl.ReturnType)
	}
	out.WriteString(" ")
	out.WriteString(fl.Body.String())

	return out.String()
}
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }

// CallExpression represents function calls e.g., add(1, 2)
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
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
func (ce *CallExpression) GetToken() token.Token { return ce.Token }

// IfExpression represents 'if condition { ... } else { ... }'
type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}
func (ie *IfExpression) GetToken() token.Token { return ie.Token }

// WhileExpression represents 'while condition { ... }'
type WhileExpression struct {
	Token     token.Token // token.WHILE
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) expressionNode()      {}
func (we *WhileExpression) TokenLiteral() string { return we.Token.Literal }
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())
	return out.String()
}
func (we *WhileExpression) GetToken() token.Token { return we.Token }

// ForExpression represents 'for init; condition; post { ... }'
type ForExpression struct {
	Token     token.Token // token.FOR
	Init      Statement   // e.g., var i = 0 or i := 0 or i = 0 or nil
	Condition Expression  // e.g., i < 10 or nil
	Post      Statement   // e.g., i++ or i = i + 1 or nil
	Body      *BlockStatement
}

func (fe *ForExpression) expressionNode()      {}
func (fe *ForExpression) TokenLiteral() string { return fe.Token.Literal }
func (fe *ForExpression) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	if fe.Init != nil {
		out.WriteString(fe.Init.String() + " ")
	} else {
		out.WriteString("; ")
	}
	if fe.Condition != nil {
		out.WriteString(fe.Condition.String() + " ")
	}
	out.WriteString("; ")
	if fe.Post != nil {
		out.WriteString(fe.Post.String())
	}
	out.WriteString(" ")
	out.WriteString(fe.Body.String())
	return out.String()
}
func (fe *ForExpression) GetToken() token.Token { return fe.Token }
