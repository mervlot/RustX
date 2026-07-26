package parser

import (
	"fmt"
	"strconv"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/errors"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/token"
)

const (
	_ int = iota
	LOWEST
	ASSIGN      // = and :=
	OR          // ||
	AND         // &&
	EQUALS      // == or !=
	LESSGREATER // > or < or >= or <=
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -x or !x
	CALL        // myFunction(X), postfix ++ or --
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: ASSIGN,
	token.OR:           OR,
	token.AND:          AND,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LTE:          LESSGREATER,
	token.GTE:          LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.ASTERISK:     PRODUCT,
	token.SLASH:        PRODUCT,
	token.LPAREN:       CALL,
	token.INC:          CALL,
	token.DEC:          CALL,
	token.LBRACKET:     INDEX,
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

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
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
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.COLON_ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) newSyntaxError(msg string, line, col int) {
	p.errors = append(p.errors, &errors.SyntaxError{Message: msg, Line: line, Column: col})
}

func (p *Parser) peekError(t token.TokenType) {
	p.newSyntaxError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type), p.peekToken.Line, p.peekToken.Column)
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if pref, ok := precedences[p.peekToken.Type]; ok {
		return pref
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if pref, ok := precedences[p.curToken.Type]; ok {
		return pref
	}
	return LOWEST
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
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

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = p.curToken.Literal
	}

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

	if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken() // consume WHILE
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	tok := p.curToken

	if p.peekTokenIs(token.LBRACE) {
		p.nextToken() // consume FOR
		body := p.parseBlockStatement()
		return &ast.ForStatement{
			Token: tok,
			Body:  body,
		}
	}

	p.nextToken() // consume FOR to start parsing loop header

	var init ast.Expression
	var cond ast.Expression
	var post ast.Expression

	if p.curTokenIs(token.VAR) {
		initVar := p.parseVarStatement()
		init = initVar

		if !p.curTokenIs(token.SEMICOLON) {
			if !p.expectPeek(token.SEMICOLON) {
				return nil
			}
		}
		p.nextToken() // move past ';' to condition

		if !p.curTokenIs(token.SEMICOLON) {
			cond = p.parseExpression(LOWEST)
		}
		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}
		p.nextToken() // move past ';' to post

		if !p.curTokenIs(token.LBRACE) {
			post = p.parseExpression(LOWEST)
		}
	} else {
		expr := p.parseExpression(LOWEST)

		if p.peekTokenIs(token.SEMICOLON) {
			init = expr
			p.nextToken() // move curToken to ';' (peekToken becomes start of condition)
			p.nextToken() // move curToken to start of condition

			if !p.curTokenIs(token.SEMICOLON) {
				cond = p.parseExpression(LOWEST)
			}
			if !p.expectPeek(token.SEMICOLON) {
				return nil
			}
			p.nextToken() // move past ';' to post

			if !p.curTokenIs(token.LBRACE) {
				post = p.parseExpression(LOWEST)
			}
		} else {
			cond = expr
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &ast.ForStatement{
		Token:     tok,
		Init:      init,
		Condition: cond,
		Post:      post,
		Body:      body,
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken() // consume '{'

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.newSyntaxError(fmt.Sprintf("no prefix parsing function for %s found", p.curToken.Type), p.curToken.Line, p.curToken.Column)
		return nil
	}
	leftExp := prefix()

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
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.newSyntaxError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal), p.curToken.Line, p.curToken.Column)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.newSyntaxError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal), p.curToken.Line, p.curToken.Column)
		return nil
	}
	return &ast.FloatLiteral{Token: p.curToken, Value: val}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
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

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	tok := p.curToken
	op := p.curToken.Literal
	p.nextToken()
	right := p.parseExpression(ASSIGN)
	return &ast.AssignmentExpression{
		Token:    tok,
		Left:     left,
		Operator: op,
		Value:    right,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	al := &ast.ArrayLiteral{Token: p.curToken}
	al.Elements = p.parseExpressionList(token.RBRACKET)
	return al
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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	ie := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	ie.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return ie
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	params := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.Parameter{Token: p.curToken, Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := &ast.Parameter{Token: p.curToken, Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			param.Type = p.curToken.Literal
		}
		params = append(params, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

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

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken() // consume 'if'
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			elseIfExpr := p.parseIfExpression()
			if elseIfExpr != nil {
				stmt := &ast.ExpressionStatement{
					Token:      elseIfExpr.GetToken(),
					Expression: elseIfExpr,
				}
				expression.Alternative = &ast.BlockStatement{
					Token:      elseIfExpr.GetToken(),
					Statements: []ast.Statement{stmt},
				}
			}
		} else if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			expression.Alternative = p.parseBlockStatement()
		} else {
			p.newSyntaxError("expected '{' or 'if' after 'else'", p.peekToken.Line, p.peekToken.Column)
			return nil
		}
	}

	return expression
}
