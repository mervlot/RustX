package parser

import (
	"fmt"
	"strconv"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/errors"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/token"
)

// Precedences
const (
	_ int = iota
	LOWEST
	ASSIGN      // = or :=
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALS      // == or !=
	LESS_GREATER // > or < or >= or <=
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // myFunction(X), array[index], post++ / post--
)

var precedences = map[token.Type]int{
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: ASSIGN,
	token.OR:           LOGICAL_OR,
	token.AND:          LOGICAL_AND,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESS_GREATER,
	token.GT:           LESS_GREATER,
	token.LTE:          LESS_GREATER,
	token.GTE:          LESS_GREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     CALL,
	token.INCREMENT:    CALL,
	token.DECREMENT:    CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) newSyntaxError(msg string, tok token.Token) *errors.SyntaxError {
	return &errors.SyntaxError{
		Message: msg,
		Line:    tok.Line,
		Column:  tok.Column,
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, p.newSyntaxError(msg, p.peekToken))
}

func (p *Parser) registerPrefix(tokenType token.Type, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.Type, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Skip extra semicolons
	for p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	if p.curTokenIs(token.EOF) {
		return nil
	}

	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatementOrAssignment()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: var name type = value or var name = value
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = p.curToken.Literal
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // move to =
		p.nextToken() // move to expression start
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatementOrAssignment() ast.Statement {
	// We might have an assignment 'x = expr' or 'x := expr'
	// Or a plain expression statement.
	// Let's parse an expression first, and then check if the next token is an assignment operator,
	// OR if we can parse it as a standalone expression statement.

	// Check if this looks like 'ident := expression' or 'ident = expression' before parsing.
	// But actually, we can also use Pratt parser directly! If ASSIGN precedence is handled correctly,
	// an assignment parses as an infix expression. Let's do that or parse specifically.
	// Let's parse expression first:
	expr := p.parseExpression(LOWEST)

	if expr == nil {
		return nil
	}

	// Check if the parsed expression is an identifier and is followed by '=' or ':='
	if ident, ok := expr.(*ast.Identifier); ok {
		if p.peekTokenIs(token.ASSIGN) {
			p.nextToken() // move to =
			tok := p.curToken
			p.nextToken() // move to expr start
			val := p.parseExpression(LOWEST)
			stmt := &ast.AssignStatement{Token: tok, Name: ident, Value: val}
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}
			return stmt
		} else if p.peekTokenIs(token.COLON_ASSIGN) {
			p.nextToken() // move to :=
			tok := p.curToken
			p.nextToken() // move to expr start
			val := p.parseExpression(LOWEST)
			stmt := &ast.ShorthandAssignStatement{Token: tok, Name: ident, Value: val}
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}
			return stmt
		}
	}

	stmt := &ast.ExpressionStatement{Token: expr.GetToken(), Expression: expr}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, p.newSyntaxError(fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type), p.curToken))
		return nil
	}
	leftExp := prefix()

	// To support unparenthesized control flow, stop parsing if peekToken is '{'
	for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.LBRACE) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, p.newSyntaxError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal), p.curToken))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, p.newSyntaxError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal), p.curToken))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	al := &ast.ArrayLiteral{Token: p.curToken}
	al.Elements = p.parseExpressionList(token.RBRACKET)
	return al
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var list []ast.Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next expression
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			// Nested if
			nestedIf := p.parseIfExpression()
			if nestedIf != nil {
				expression.Alternative = &ast.BlockStatement{
					Token:      nestedIf.GetToken(),
					Statements: []ast.Statement{&ast.ExpressionStatement{Token: nestedIf.GetToken(), Expression: nestedIf}},
				}
			}
		} else {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}
	p.nextToken() // move past 'for'

	// Semicolons are required for standard 3-part for loop, e.g., for init; cond; post { ... }
	// We might also have `for cond { ... }` which is basically `while`, but let's handle the 3-part loop first.
	// Actually, Go supports `for cond { ... }` or `for init; cond; post { ... }` or `for { ... }`.
	// Let's implement full Go-style 3-part loop `for init; condition; post {}`
	// Check if there is NO loop header (infinite loop): `for {`
	if p.curTokenIs(token.LBRACE) {
		expression.Body = p.parseBlockStatement()
		return expression
	}

	// Try parsing init statement. But we don't know yet if it's `for cond {` or `for init; cond; post {`.
	// Let's parse an expression/statement first, but handle it carefully.
	// If the current statement is followed by a semicolon, it's definitely a 3-part loop!
	// E.g., `for i := 0; i < 10; i++ {`
	// Let's look ahead or parse a statement. Let's parse a statement first if possible.
	// Actually, we can check if there's a semicolon in the future before a brace.
	// A simple scanner search in the lexer input can tell us. Or we can parse an expression and look if we see a semicolon.
	// In RustX/Go style:
	// Let's parse the first part. It can be a VarStatement, ShorthandAssignStatement, AssignStatement, or Expression.
	var init ast.Statement
	if p.curTokenIs(token.VAR) {
		init = p.parseVarStatement()
	} else {
		// Try parsing expression or assignment
		init = p.parseExpressionStatementOrAssignment()
	}

	if p.curTokenIs(token.SEMICOLON) {
		// It was an init statement! So this is a 3-part loop.
		expression.Init = init
		p.nextToken() // consume semicolon

		if !p.curTokenIs(token.SEMICOLON) {
			expression.Condition = p.parseExpression(LOWEST)
		}
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}
		p.nextToken() // move past semicolon

		if !p.curTokenIs(token.LBRACE) {
			// Parse post statement, can be a PostfixExpression (like x++) or an AssignStatement
			expression.Post = p.parseExpressionStatementOrAssignment()
		}
	} else {
		// No semicolon, so it was actually the condition of a single-condition loop `for cond {`
		// Cast init to Expression if possible
		if exprStmt, ok := init.(*ast.ExpressionStatement); ok {
			expression.Condition = exprStmt.Expression
		} else {
			p.errors = append(p.errors, p.newSyntaxError("invalid condition in for loop", p.curToken))
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Optional function name for recursion: func factorial(...)
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type: func add(...) int {
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.ReturnType = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	var identifiers []*ast.Parameter

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	// Parse parameter name
	param := &ast.Parameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	// Optional parameter type: name type
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	identifiers = append(identifiers, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next identifier
		param = &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			param.Type = p.curToken.Literal
		}
		identifiers = append(identifiers, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}
