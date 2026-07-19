package ast

import (
	"bytes"
	"strings"

	"github.com/username/rustx/token"
)

// The base Node interface
type Node interface {
	GetToken() token.Token
	String() string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of every AST
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

// --- Statements ---

// VarStatement represents mutable variable declarations: e.g. "var x = 5" or with types "var x int = 5"
// Shorthand assignment also can be modeled similarly or separately, we'll model it as a separate expression or statement as needed, but standardizing:
// "var x = 5;" => VarStatement
type VarStatement struct {
	Token token.Token // the token.VAR token
	Name  *Identifier
	Value Expression
	Type  string // optional type annotation, e.g. "int" or "float"
}

func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) statementNode()        {}
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

// ReturnStatement represents "return <expression>;"
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) GetToken() token.Token { return rs.Token }
func (rs *ReturnStatement) statementNode()        {}
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Token.Literal + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement wrapper to use expressions where statements are expected
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) GetToken() token.Token { return es.Token }
func (es *ExpressionStatement) statementNode()        {}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// BlockStatement is a sequence of statements wrapped in braces
type BlockStatement struct {
	Token      token.Token // the token.LBRACE token
	Statements []Statement
}

func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) statementNode()        {}
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

// --- Expressions ---

// Identifier represents variable names, function names, etc.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) expressionNode()        {}
func (i *Identifier) String() string         { return i.Value }

// IntegerLiteral represents an integer constant
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) expressionNode()        {}
func (il *IntegerLiteral) String() string         { return il.Token.Literal }

// FloatLiteral represents a floating-point constant
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) expressionNode()        {}
func (fl *FloatLiteral) String() string         { return fl.Token.Literal }

// Boolean represents true / false
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) GetToken() token.Token { return b.Token }
func (b *Boolean) expressionNode()        {}
func (b *Boolean) String() string         { return b.Token.Literal }

// StringLiteral represents string constants
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) expressionNode()        {}
func (sl *StringLiteral) String() string         { return sl.Token.Literal }

// ArrayLiteral represents array list constants, e.g. [1, 2, 3]
type ArrayLiteral struct {
	Token    token.Token // the token.LBRACKET token
	Elements []Expression
}

func (al *ArrayLiteral) GetToken() token.Token { return al.Token }
func (al *ArrayLiteral) expressionNode()        {}
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

// IndexExpression represents indexing, e.g. "array[index]" or "string[index]"
type IndexExpression struct {
	Token token.Token // the [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) GetToken() token.Token { return ie.Token }
func (ie *IndexExpression) expressionNode()        {}
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// PrefixExpression represents prefix operators: e.g. "-5", "!true"
type PrefixExpression struct {
	Token    token.Token // prefix token, e.g. ! or -
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) GetToken() token.Token { return pe.Token }
func (pe *PrefixExpression) expressionNode()        {}
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression represents infix operators: e.g. "5 + 5", "x > y"
type InfixExpression struct {
	Token    token.Token // infix token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) GetToken() token.Token { return ie.Token }
func (ie *InfixExpression) expressionNode()        {}
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// PostfixExpression represents postfix operators: e.g. "i++", "i--"
type PostfixExpression struct {
	Token    token.Token // e.g. ++ or --
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) GetToken() token.Token { return pe.Token }
func (pe *PostfixExpression) expressionNode()        {}
func (pe *PostfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")
	return out.String()
}

// AssignExpression represents variable updates or shorthand bindings
// This includes:
// "x = 5" (ASSIGN)
// "x := 5" (COLON_ASSIGN)
type AssignExpression struct {
	Token    token.Token // = or :=
	Name     *Identifier
	Value    Expression
	Operator string // = or :=
}

func (ae *AssignExpression) GetToken() token.Token { return ae.Token }
func (ae *AssignExpression) expressionNode()        {}
func (ae *AssignExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ae.Name.String())
	out.WriteString(" " + ae.Operator + " ")
	out.WriteString(ae.Value.String())
	out.WriteString(")")
	return out.String()
}

// IfExpression representing: if <condition> { <consequence> } else { <alternative> }
type IfExpression struct {
	Token       token.Token // the token.IF token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) GetToken() token.Token { return ie.Token }
func (ie *IfExpression) expressionNode()        {}
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

// FunctionLiteral represents function literals: "func(x int, y int) int { ... }"
// Support optional parameter/return types to be ignored or kept.
type Parameter struct {
	Name *Identifier
	Type string // e.g. "int" (optional)
}

type FunctionLiteral struct {
	Token      token.Token // the token.FUNC token
	Parameters []*Parameter
	ReturnType string // optional, e.g. "int"
	Body       *BlockStatement
	Name       string // optional, for named functions so they can reference themselves recursively
}

func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) expressionNode()        {}
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		paramStr := p.Name.String()
		if p.Type != "" {
			paramStr += " " + p.Type
		}
		params = append(params, paramStr)
	}
	out.WriteString(fl.Token.Literal)
	if fl.Name != "" {
		out.WriteString(" " + fl.Name)
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

// CallExpression represents function calls: "add(2, 3)"
type CallExpression struct {
	Token     token.Token // the token.LPAREN token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) GetToken() token.Token { return ce.Token }
func (ce *CallExpression) expressionNode()        {}
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

// WhileExpression represents while loops: "while <condition> { ... }"
type WhileExpression struct {
	Token     token.Token // the token.WHILE token
	Condition Expression
	Body      *BlockStatement
}

func (we *WhileExpression) GetToken() token.Token { return we.Token }
func (we *WhileExpression) expressionNode()        {}
func (we *WhileExpression) String() string {
	var out bytes.Buffer
	out.WriteString("while ")
	out.WriteString(we.Condition.String())
	out.WriteString(" ")
	out.WriteString(we.Body.String())
	return out.String()
}

// ForExpression represents three-part or simpler Go-style for loops:
// "for Init; Condition; Post { Body }" or "for Condition { Body }" (similar to while)
type ForExpression struct {
	Token     token.Token // the token.FOR token
	Init      Statement   // e.g. i := 0 or var i = 0 (optional)
	Condition Expression  // e.g. i < 10 (optional)
	Post      Statement   // e.g. i++ or i = i + 1 (optional)
	Body      *BlockStatement
}

func (fe *ForExpression) GetToken() token.Token { return fe.Token }
func (fe *ForExpression) expressionNode()        {}
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
		out.WriteString(fe.Post.String())
	}
	out.WriteString(" ")
	out.WriteString(fe.Body.String())
	return out.String()
}
