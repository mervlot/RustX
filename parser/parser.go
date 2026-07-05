package parser

import (
	"fmt"
	"strconv"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // = or :=
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	POSTFIX     // x++
)

var precedences = map[token.TokenType]int{
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.LTE:          LESSGREATER,
	token.GT:           LESSGREATER,
	token.GTE:          LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     INDEX,
	token.INCREMENT:    POSTFIX,
	token.DECREMENT:    POSTFIX,
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: ASSIGN,
	token.AND:          EQUALS,
	token.OR:           EQUALS,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.COLON_ASSIGN, p.parseAssignExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("SyntaxError: expected next token to be %s, got %s instead at line %d column %d",
		t, p.peekToken.Type, p.peekToken.Line, p.peekToken.Column)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("SyntaxError: no prefix parse function for %s found at line %d column %d",
		t, p.curToken.Line, p.curToken.Column)
	p.errors = append(p.errors, msg)
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
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("RuntimeError: could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("RuntimeError: could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
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

	// Rust/Go style: parentheses are optional.
	// But our Pratt parser might need a bit of help if they are not there.
	// Actually, p.parseExpression(LOWEST) will stop at LBRACE if we configure it.

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
			// Nested if-else if
			elseIf := p.parseIfExpression()
			// Wrap it in a BlockStatement
			expression.Alternative = &ast.BlockStatement{
				Token:      token.Token{Type: token.LBRACE, Literal: "{"},
				Statements: []ast.Statement{&ast.ExpressionStatement{Expression: elseIf}},
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

	// Check if it's a named function: func add(x, y) { ... }
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type: func add(x, y) int { ... }
	if p.peekTokenIs(token.INT_TYPE) || p.peekTokenIs(token.FLOAT_TYPE) || p.peekTokenIs(token.BOOL_TYPE) || p.peekTokenIs(token.STRING_TYPE) || p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.ReturnType = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	// Optional type annotation: func add(x int, y int)
	if p.peekTokenIs(token.INT_TYPE) || p.peekTokenIs(token.FLOAT_TYPE) || p.peekTokenIs(token.BOOL_TYPE) || p.peekTokenIs(token.STRING_TYPE) || p.peekTokenIs(token.IDENT) {
		p.nextToken()
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)

		if p.peekTokenIs(token.INT_TYPE) || p.peekTokenIs(token.FLOAT_TYPE) || p.peekTokenIs(token.BOOL_TYPE) || p.peekTokenIs(token.STRING_TYPE) || p.peekTokenIs(token.IDENT) {
			p.nextToken()
		}
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

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseWhileExpression() ast.Expression {
	exp := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseForExpression() ast.Expression {
	exp := &ast.ForExpression{Token: p.curToken}

	// for i := 0; i < 10; i++ { ... }
	if !p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		if !p.curTokenIs(token.SEMICOLON) {
			exp.Init = p.parseStatement()
		}

		if !p.curTokenIs(token.SEMICOLON) {
			p.expectPeek(token.SEMICOLON)
		}

		p.nextToken()
		if !p.curTokenIs(token.SEMICOLON) {
			exp.Condition = p.parseExpression(LOWEST)
		}

		if !p.curTokenIs(token.SEMICOLON) {
			p.expectPeek(token.SEMICOLON)
		}

		p.nextToken()
		if !p.curTokenIs(token.LBRACE) {
			exp.Post = p.parseStatement()
		}
	}

	if !p.curTokenIs(token.LBRACE) {
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	exp := &ast.AssignExpression{
		Token: p.curToken,
	}

	if ident, ok := left.(*ast.Identifier); ok {
		exp.Name = ident
	} else {
		msg := fmt.Sprintf("SyntaxError: expected identifier on left side of assignment at line %d column %d",
			p.curToken.Line, p.curToken.Column)
		p.errors = append(p.errors, msg)
		return nil
	}

	p.nextToken()
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
