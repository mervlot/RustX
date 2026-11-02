package ast

import (
	"bytes"
	"rustx/token"
	"strings"
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

// Program is the root AST node.
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

// Identifier represents an identifier expression.
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) GetToken() token.Token { return i.Token }
func (i *Identifier) String() string       { return i.Value }

// TypeAnnotation represents optional type annotations (e.g. `int`, `string`, `func(...)`).
type TypeAnnotation struct {
	Token token.Token
	Name  string
}

func (t *TypeAnnotation) String() string { return t.Name }

// VarStatement represents `var x = 5` or `var x int = 5`.
type VarStatement struct {
	Token token.Token // token.VAR
	Name  *Identifier
	Type  *TypeAnnotation // optional
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Token }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString("var ")
	out.WriteString(vs.Name.String())
	if vs.Type != nil {
		out.WriteString(" " + vs.Type.String())
	}
	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	return out.String()
}

// AssignStatement represents assignment `x = 5` or `x := 5` or `arr[0] = 5`.
type AssignStatement struct {
	Token    token.Token // token.ASSIGN or token.COLON_ASSIGN
	Target   Expression  // Identifier or IndexExpression
	Value    Expression
	IsDefine bool // true if :=
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) GetToken() token.Token { return as.Token }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Target.String())
	if as.IsDefine {
		out.WriteString(" := ")
	} else {
		out.WriteString(" = ")
	}
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	return out.String()
}

// ReturnStatement represents `return expr`.
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
	return out.String()
}

// ExpressionStatement represents an expression wrapped as a statement.
type ExpressionStatement struct {
	Token      token.Token // first token of expression
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

// BlockStatement represents a sequence of statements inside `{ ... }`.
type BlockStatement struct {
	Token      token.Token // '{'
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) GetToken() token.Token { return bs.Token }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// IntegerLiteral represents an integer literal.
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Token }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral represents a floating-point literal.
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Token }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral represents a string literal.
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Token }
func (sl *StringLiteral) String() string       { return `"` + sl.Value + `"` }

// BooleanLiteral represents true / false.
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) GetToken() token.Token { return bl.Token }
func (bl *BooleanLiteral) String() string       { return bl.Token.Literal }

// PrefixExpression represents `!x`, `-x`, etc.
type PrefixExpression struct {
	Token    token.Token // '!', '-'
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

// InfixExpression represents `a + b`, `x == y`, etc.
type InfixExpression struct {
	Token    token.Token // operator token
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

// PostfixExpression represents `x++` or `x--`.
type PostfixExpression struct {
	Token    token.Token // '++', '--'
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

// IfExpression represents `if condition { consequence } else { alternative }`.
type IfExpression struct {
	Token       token.Token // 'if'
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // optional
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

// WhileStatement represents `while condition { body }`.
type WhileStatement struct {
	Token     token.Token // 'while'
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

// ForStatement represents `for init; condition; post { body }` or `for condition { body }`.
type ForStatement struct {
	Token     token.Token // 'for'
	Init      Statement   // can be nil, or VarStatement / AssignStatement / ExpressionStatement
	Condition Expression  // can be nil
	Post      Statement   // can be nil, or ExpressionStatement / AssignStatement
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
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
		out.WriteString(" ")
	} else if fs.Condition != nil {
		out.WriteString(fs.Condition.String())
		out.WriteString(" ")
	}
	out.WriteString(fs.Body.String())
	return out.String()
}

// FunctionParameter represents a parameter with an optional type.
type FunctionParameter struct {
	Name *Identifier
	Type *TypeAnnotation
}

func (fp *FunctionParameter) String() string {
	if fp.Type != nil {
		return fp.Name.String() + " " + fp.Type.String()
	}
	return fp.Name.String()
}

// FunctionLiteral represents `func name?(params) returnType? { body }`.
type FunctionLiteral struct {
	Token      token.Token // 'func'
	Name       string      // optional name if named func
	Parameters []*FunctionParameter
	ReturnType *TypeAnnotation // optional
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) statementNode()       {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Token }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	var params []string
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	if fl.Name != "" {
		out.WriteString(" " + fl.Name)
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

// CallExpression represents `funcName(arg1, arg2)`.
type CallExpression struct {
	Token     token.Token // '('
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

// IndexExpression represents `arr[idx]`.
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
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// ArrayLiteral represents `[1, 2, 3]`.
type ArrayLiteral struct {
	Token    token.Token // '['
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
