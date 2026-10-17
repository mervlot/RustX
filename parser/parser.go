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
	ASSIGN      // = or :=
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALS      // == or !=
	LESSGREATER // > or < or >= or <=
	SUM         // + or -
	PRODUCT     // * or / or %
	PREFIX      // -X or !X
	CALL        // myFunction(X) or X++ or X--
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.DECL_ASSIGN: ASSIGN,
	token.ASSIGN:      ASSIGN,
	token.OR:          LOGICAL_OR,
	token.AND:         LOGICAL_AND,
	token.EQ:          EQUALS,
	token.NOT_EQ:      EQUALS,
	token.LT:          LESSGREATER,
	token.GT:          LESSGREATER,
	token.LTE:         LESSGREATER,
	token.GTE:         LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.SLASH:       PRODUCT,
	token.ASTERISK:    PRODUCT,
	token.PERCENT:     PRODUCT,
	token.INC:         CALL,
	token.DEC:         CALL,
	token.LPAREN:      CALL,
	token.LBRACKET:    INDEX,
}

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) (*Parser, error) {
	p := &Parser{
		l:              l,
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	err := p.nextToken()
	if err != nil {
		return nil, err
	}
	err = p.nextToken()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Parser) nextToken() error {
	p.curToken = p.peekToken
	var err error
	p.peekToken, err = p.l.NextToken()
	return err
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.FUNCTION:
		if p.peekToken.Type == token.IDENT {
			// Named top-level or statement function
			return p.parseFunctionStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() (ast.Statement, error) {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected variable name after 'var'")
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: var x int = 5
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		stmt.Type = p.curToken.Literal
	}

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
		p.nextToken()

		val, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.Value = val
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseReturnStatement() (ast.Statement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		val, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt.ReturnValue = val
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseExpressionStatement() (ast.Statement, error) {
	startToken := p.curToken
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	// Check if this expression is part of a shorthand assignment: name := expr
	if p.peekToken.Type == token.DECL_ASSIGN {
		ident, ok := expr.(*ast.Identifier)
		if !ok {
			return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Left side of ':=' must be an identifier")
		}
		p.nextToken() // move to ':='
		assignToken := p.curToken
		p.nextToken() // move to right-hand side
		val, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		stmt := &ast.ShortVarStatement{
			Token: assignToken,
			Name:  ident,
			Value: val,
		}
		if p.peekToken.Type == token.SEMICOLON {
			p.nextToken()
		}
		return stmt, nil
	}

	// Check if this expression is standard assignment: name = expr
	if p.peekToken.Type == token.ASSIGN {
		if ident, ok := expr.(*ast.Identifier); ok {
			p.nextToken() // move to '='
			assignToken := p.curToken
			p.nextToken() // move to RHS
			val, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			stmt := &ast.AssignStatement{
				Token: assignToken,
				Name:  ident,
				Value: val,
			}
			if p.peekToken.Type == token.SEMICOLON {
				p.nextToken()
			}
			return stmt, nil
		} else if idxExpr, ok := expr.(*ast.IndexExpression); ok {
			p.nextToken() // move to '='
			assignToken := p.curToken
			p.nextToken() // move to RHS
			val, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			stmt := &ast.IndexAssignStatement{
				Token: assignToken,
				Left:  idxExpr,
				Value: val,
			}
			if p.peekToken.Type == token.SEMICOLON {
				p.nextToken()
			}
			return stmt, nil
		}
	}

	stmt := &ast.ExpressionStatement{Token: startToken, Expression: expr}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt, nil
}

func (p *Parser) parseWhileStatement() (ast.Statement, error) {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	cond, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	stmt.Condition = cond

	if !p.expectPeek(token.LBRACE) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected '{' after while condition")
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	return stmt, nil
}

func (p *Parser) parseForStatement() (ast.Statement, error) {
	stmt := &ast.ForStatement{Token: p.curToken}

	// Check if this is a 3-part for loop (contains ';') or single condition for loop
	// We scan ahead in input using a temporary lexer to check for ';' before '{'
	isThreePart := p.isThreePartForHeader()

	p.nextToken()

	if isThreePart {
		// Part 1: Init statement
		if p.curToken.Type != token.SEMICOLON {
			initStmt, err := p.parseForHeaderStatement()
			if err != nil {
				return nil, err
			}
			stmt.Init = initStmt
			if p.peekToken.Type == token.SEMICOLON {
				p.nextToken()
			}
		}

		if !p.curTokenIs(token.SEMICOLON) {
			return nil, errors.NewSyntaxError(p.curToken.Line, p.curToken.Column, "Expected ';' after for loop init")
		}
		p.nextToken()

		// Part 2: Condition
		if p.curToken.Type != token.SEMICOLON {
			cond, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			stmt.Condition = cond
		}

		if !p.expectPeek(token.SEMICOLON) {
			return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected ';' after for loop condition")
		}
		p.nextToken()

		// Part 3: Post statement
		if p.curToken.Type != token.LBRACE {
			postStmt, err := p.parseForHeaderStatement()
			if err != nil {
				return nil, err
			}
			stmt.Post = postStmt
			if p.peekToken.Type == token.SEMICOLON {
				p.nextToken()
			}
			p.nextToken()
		}
	} else {
		// Single condition 'for condition { ... }' or infinite loop 'for { ... }'
		if p.curToken.Type != token.LBRACE {
			cond, err := p.parseExpression(LOWEST)
			if err != nil {
				return nil, err
			}
			stmt.Condition = cond
			p.nextToken()
		}
	}

	if !p.curTokenIs(token.LBRACE) {
		return nil, errors.NewSyntaxError(p.curToken.Line, p.curToken.Column, "Expected '{' to start for loop body")
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	stmt.Body = body

	return stmt, nil
}

func (p *Parser) isThreePartForHeader() bool {
	// Create a temporary lexer scanning ahead from the current position
	tempLexer := lexer.New(p.l.Input()[p.l.Position():])
	braceDepth := 0

	for {
		tok, err := tempLexer.NextToken()
		if err != nil || tok.Type == token.EOF {
			break
		}
		if tok.Type == token.LBRACE {
			if braceDepth == 0 {
				return false
			}
			braceDepth++
		} else if tok.Type == token.RBRACE {
			braceDepth--
		} else if tok.Type == token.SEMICOLON && braceDepth == 0 {
			return true
		}
	}
	return false
}

func (p *Parser) parseForHeaderStatement() (ast.Statement, error) {
	if p.curTokenIs(token.VAR) {
		return p.parseVarStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseFunctionStatement() (ast.Statement, error) {
	fnTok := p.curToken
	p.nextToken() // move to fn name

	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.LPAREN) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected '(' after function name")
	}

	params, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}

	var returnType string
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		returnType = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected '{' before function body")
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return &ast.FunctionLiteral{
		Token:      fnTok,
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}, nil
}

func (p *Parser) parseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil, errors.NewSyntaxError(p.curToken.Line, p.curToken.Column, "No prefix parse function for %s found", p.curToken.Literal)
	}

	leftExp, err := prefix()
	if err != nil {
		return nil, err
	}

	for p.peekToken.Type != token.SEMICOLON && p.peekToken.Type != token.LBRACE && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp, nil
		}

		p.nextToken()

		leftExp, err = infix(leftExp)
		if err != nil {
			return nil, err
		}
	}

	return leftExp, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		return nil, errors.NewSyntaxError(p.curToken.Line, p.curToken.Column, "Could not parse %q as integer", p.curToken.Literal)
	}

	lit.Value = value
	return lit, nil
}

func (p *Parser) parseFloatLiteral() (ast.Expression, error) {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		return nil, errors.NewSyntaxError(p.curToken.Line, p.curToken.Column, "Could not parse %q as float", p.curToken.Literal)
	}

	lit.Value = value
	return lit, nil
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}

	expression.Right = right
	return expression, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()

	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}

	expression.Right = right
	return expression, nil
}

func (p *Parser) parsePostfixExpression(left ast.Expression) (ast.Expression, error) {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if !p.expectPeek(token.RPAREN) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected ')' after expression")
	}

	return exp, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	cond, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	expression.Condition = cond

	if !p.expectPeek(token.LBRACE) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected '{' after if condition")
	}

	consequence, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	expression.Consequence = consequence

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			// else if block
			p.nextToken()
			elseIfExpr, err := p.parseIfExpression()
			if err != nil {
				return nil, err
			}
			stmt := &ast.ExpressionStatement{
				Token:      elseIfExpr.GetToken(),
				Expression: elseIfExpr,
			}
			expression.Alternative = &ast.BlockStatement{
				Token:      elseIfExpr.GetToken(),
				Statements: []ast.Statement{stmt},
			}
		} else if p.expectPeek(token.LBRACE) {
			alt, err := p.parseBlockStatement()
			if err != nil {
				return nil, err
			}
			expression.Alternative = alt
		}
	}

	return expression, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected '(' after func")
	}

	params, err := p.parseFunctionParameters()
	if err != nil {
		return nil, err
	}
	lit.Parameters = params

	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		lit.ReturnType = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected '{' before function body")
	}

	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	lit.Body = body

	return lit, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.Parameter, error) {
	identifiers := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers, nil
	}

	p.nextToken()

	param := &ast.Parameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	identifiers = append(identifiers, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		pParam := &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekToken.Type == token.IDENT {
			p.nextToken()
			pParam.Type = p.curToken.Literal
		}
		identifiers = append(identifiers, pParam)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected ')' after function parameters")
	}

	return identifiers, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block, nil
}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	args, err := p.parseExpressionList(token.RPAREN)
	if err != nil {
		return nil, err
	}
	exp.Arguments = args
	return exp, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	array := &ast.ArrayLiteral{Token: p.curToken}
	elements, err := p.parseExpressionList(token.RBRACKET)
	if err != nil {
		return nil, err
	}
	array.Elements = elements
	return array, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	index, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	exp.Index = index

	if !p.expectPeek(token.RBRACKET) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, "Expected ']' after index expression")
	}

	return exp, nil
}

func (p *Parser) parseExpressionList(end token.TokenType) ([]ast.Expression, error) {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list, nil
	}

	p.nextToken()
	expr, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	list = append(list, expr)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		expr, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
	}

	if !p.expectPeek(end) {
		return nil, errors.NewSyntaxError(p.peekToken.Line, p.peekToken.Column, fmt.Sprintf("Expected %q at end of expression list", end))
	}

	return list, nil
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
	}
	return false
}
