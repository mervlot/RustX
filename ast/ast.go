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
	GetToken() token.Token
}

// Statement represents statements that do not evaluate to a value
type Statement interface {
	Node
	statementNode()
}

// Expression represents constructs that evaluate to a value
type Expression interface {
	Node
	expressionNode()
}

// Program is the root AST node representing a whole script/file
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

// --- STATEMENTS ---

// VarStatement represents `var name [type] = value` or `var name [type]`
type VarStatement struct {
	Tok   token.Token // token.VAR
	Name  *Identifier
	Type  string // optional type annotation e.g. "int"
	Value Expression
}

func (vs *VarStatement) statementNode()       {}
func (vs *VarStatement) TokenLiteral() string { return vs.Tok.Literal }
func (vs *VarStatement) GetToken() token.Token { return vs.Tok }
func (vs *VarStatement) String() string {
	var out bytes.Buffer
	out.WriteString("var ")
	out.WriteString(vs.Name.String())
	if vs.Type != "" {
		out.WriteString(" " + vs.Type)
	}
	if vs.Value != nil {
		out.WriteString(" = ")
		out.WriteString(vs.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// AssignStatement represents `name = value` or `name := value`
type AssignStatement struct {
	Tok      token.Token // token.ASSIGN or token.COLON_ASSIGN
	Name     *Identifier
	IsShorthand bool
	Value    Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Tok.Literal }
func (as *AssignStatement) GetToken() token.Token { return as.Tok }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	if as.IsShorthand {
		out.WriteString(" := ")
	} else {
		out.WriteString(" = ")
	}
	if as.Value != nil {
		out.WriteString(as.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// ReturnStatement represents `return expr`
type ReturnStatement struct {
	Tok         token.Token // token.RETURN
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Tok.Literal }
func (rs *ReturnStatement) GetToken() token.Token { return rs.Tok }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.Tok.Literal + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatement represents standalone expressions (e.g., function calls, assignments as expressions)
type ExpressionStatement struct {
	Tok        token.Token // First token of expression
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

// BlockStatement represents `{ statements... }`
type BlockStatement struct {
	Tok        token.Token // token.LBRACE
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Tok.Literal }
func (bs *BlockStatement) GetToken() token.Token { return bs.Tok }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString("  " + s.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// IfExpression / Statement represents `if condition { ... } else { ... }`
type IfExpression struct {
	Tok         token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement // *BlockStatement or *IfExpression (else if)
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) statementNode()       {}
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

// WhileStatement represents `while condition { ... }`
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

// ForStatement represents `for init; condition; post { ... }` or `for condition { ... }`
type ForStatement struct {
	Tok       token.Token // token.FOR
	Init      Statement   // VarStatement or AssignStatement or ExpressionStatement
	Condition Expression
	Post      Statement   // AssignStatement or ExpressionStatement
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

// --- EXPRESSIONS ---

// Identifier represents variable/function names
type Identifier struct {
	Tok   token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Tok.Literal }
func (i *Identifier) GetToken() token.Token { return i.Tok }
func (i *Identifier) String() string       { return i.Value }

// Parameter represents function parameter `name type`
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

// IntegerLiteral represents integer values
type IntegerLiteral struct {
	Tok   token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Tok.Literal }
func (il *IntegerLiteral) GetToken() token.Token { return il.Tok }
func (il *IntegerLiteral) String() string       { return il.Tok.Literal }

// FloatLiteral represents float values
type FloatLiteral struct {
	Tok   token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FloatLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FloatLiteral) String() string       { return fl.Tok.Literal }

// StringLiteral represents string values
type StringLiteral struct {
	Tok   token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Tok.Literal }
func (sl *StringLiteral) GetToken() token.Token { return sl.Tok }
func (sl *StringLiteral) String() string       { return `"` + sl.Value + `"` }

// BooleanLiteral represents true / false
type BooleanLiteral struct {
	Tok   token.Token
	Value bool
}

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Tok.Literal }
func (b *BooleanLiteral) GetToken() token.Token { return b.Tok }
func (b *BooleanLiteral) String() string       { return b.Tok.Literal }

// ArrayLiteral represents `[elem1, elem2]`
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

// PrefixExpression represents `!expr` or `-expr`
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

// InfixExpression represents `left op right`
type InfixExpression struct {
	Tok      token.Token
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

// PostfixExpression represents `i++` or `i--`
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

// FunctionLiteral represents `func name(params) returnType { body }`
type FunctionLiteral struct {
	Tok        token.Token // token.FUNC
	Name       string      // optional named function name
	Parameters []*Parameter
	ReturnType string
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) statementNode()       {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Tok.Literal }
func (fl *FunctionLiteral) GetToken() token.Token { return fl.Tok }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("func ")
	if fl.Name != "" {
		out.WriteString(fl.Name + " ")
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

// CallExpression represents `fn(arg1, arg2)`
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

// IndexExpression represents `arr[index]` or `str[index]`
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
