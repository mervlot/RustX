package parser

import (
	"fmt"
	"rustx/ast"
	"rustx/errors"
	"rustx/lexer"
	"rustx/token"
	"strconv"
)

// Precedence levels for Pratt parsing
const (
	_ int = iota
	LOWEST
	ASSIGN      // = or :=
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALS      // ==, !=
	LESSGREATER // >, <, >=, <=
	SUM         // +, -
	PRODUCT     // *, /, %
	PREFIX      // -X, !X
	CALL        // myFunction(X), X++, X--
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: ASSIGN,
	token.OR:           LOGICAL_OR,
	token.AND:          LOGICAL_AND,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LT_EQ:        LESSGREATER,
	token.GT_EQ:        LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.PERCENT:      PRODUCT,
	token.PLUS_PLUS:    CALL,
	token.MINUS_MINUS:  CALL,
	token.LPAREN:       CALL,
	token.LBRACKET:     INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

// Parser parses tokens from Lexer into an AST.
type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

// New creates a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
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
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.PERCENT, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)

	// Read two tokens so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfix(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors returns all parsing errors encountered.
func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s ('%s') instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	p.errors = append(p.errors, errors.NewSyntaxError(p.peekToken.Position, "%s", msg))
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
	p.peekError(t)
	return false
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

// ParseProgram parses the entire program.
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
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		// Check for shorthand assignment `x := ...`
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
			return p.parseShortVarStatement()
		}
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: var x int = 10; or var x = 10;
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = p.curToken.Literal
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // move to '='
		p.nextToken() // move to expression start
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseShortVarStatement() *ast.ShortVarStatement {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // move to ':='
	colonToken := p.curToken
	p.nextToken() // move to value expression start

	stmt := &ast.ShortVarStatement{
		Token: colonToken,
		Name:  ident,
		Value: p.parseExpression(LOWEST),
	}

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

func (p *Parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}

	// Check if this is an assignment: `left = expr`
	if p.peekTokenIs(token.ASSIGN) {
		assignTok := p.peekToken
		p.nextToken() // move to '='
		p.nextToken() // move to value
		val := p.parseExpression(LOWEST)

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		return &ast.AssignStatement{
			Token: assignTok,
			Left:  expr,
			Value: val,
		}
	}

	stmt := &ast.ExpressionStatement{Token: p.curToken, Expression: expr}

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
		// Stop if peek token is '{' (e.g., if condition { or for condition {)
		if p.peekTokenIs(token.LBRACE) {
			break
		}

		// Check postfix first
		if postfix, ok := p.postfixParseFns[p.peekToken.Type]; ok {
			p.nextToken()
			leftExp = postfix(leftExp)
			continue
		}

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, errors.NewSyntaxError(p.curToken.Position, "%s", msg))
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, errors.NewSyntaxError(p.curToken.Position, "%s", msg))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, errors.NewSyntaxError(p.curToken.Position, "%s", msg))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
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
		Left:     left,
		Operator: p.curToken.Literal,
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

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
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
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()

	hasParen := false
	if p.curTokenIs(token.LPAREN) {
		hasParen = true
		p.nextToken()
	}

	expression.Condition = p.parseExpression(LOWEST)

	if hasParen {
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // move to 'else'

		if p.peekTokenIs(token.IF) {
			p.nextToken() // move to 'if'
			nestedIf := p.parseIfExpression()
			expression.Alternative = &ast.BlockStatement{
				Token:      nestedIf.GetToken(),
				Statements: []ast.Statement{&ast.ExpressionStatement{Token: nestedIf.GetToken(), Expression: nestedIf}},
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()

	hasParen := false
	if p.curTokenIs(token.LPAREN) {
		hasParen = true
		p.nextToken()
	}

	expression.Condition = p.parseExpression(LOWEST)

	if hasParen {
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	p.nextToken()

	hasParen := false
	if p.curTokenIs(token.LPAREN) {
		hasParen = true
		p.nextToken()
	}

	// Check if this is a 3-part for loop or condition-only loop
	// Look ahead to see if there is a semicolon before an opening brace '{'
	isThreePart := p.isForThreePart()

	if isThreePart {
		// Part 1: Init statement
		if !p.curTokenIs(token.SEMICOLON) {
			if p.curTokenIs(token.VAR) {
				expression.Init = p.parseVarStatement()
			} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
				expression.Init = p.parseShortVarStatement()
			} else {
				expression.Init = p.parseExpressionStatement()
			}
		} else {
			p.nextToken() // empty init
		}

		// Ensure we are past the first semicolon
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// Part 2: Condition expression
		if !p.curTokenIs(token.SEMICOLON) {
			expression.Condition = p.parseExpression(LOWEST)
		}

		if !p.curTokenIs(token.SEMICOLON) {
			if !p.expectPeek(token.SEMICOLON) {
				return nil
			}
		}

		p.nextToken() // past semicolon

		// Part 3: Post statement
		if !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.LBRACE) {
			postExpr := p.parseExpression(LOWEST)
			if postExpr != nil {
				expression.Post = &ast.ExpressionStatement{Token: postExpr.GetToken(), Expression: postExpr}
			}
			// advance past postExpr
			if p.peekTokenIs(token.RPAREN) || p.peekTokenIs(token.LBRACE) {
				p.nextToken()
			}
		}
	} else {
		// Condition-only loop: `for condition { ... }` or infinite loop `for { ... }`
		if !p.curTokenIs(token.LBRACE) {
			expression.Condition = p.parseExpression(LOWEST)
		}
	}

	if hasParen && p.curTokenIs(token.RPAREN) {
		p.nextToken()
	}

	if !p.curTokenIs(token.LBRACE) {
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) isForThreePart() bool {
	// Use a temporary lexer scanning ahead from current position
	tempLex := lexer.New(p.l.Input()[p.l.Position():])
	for {
		tok := tempLex.NextToken()
		if tok.Type == token.SEMICOLON {
			return true
		}
		if tok.Type == token.LBRACE || tok.Type == token.EOF {
			return false
		}
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Optional function name (e.g. `func add(...)` vs anonymous `func(...)`)
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type: `func add(a int, b int) int {`
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

func (p *Parser) parseFunctionParameters() []*ast.FunctionParam {
	params := []*ast.FunctionParam{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.FunctionParam{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	// Check for optional parameter type
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // past ','
		p.nextToken() // to next ident
		nextParam := &ast.FunctionParam{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			nextParam.Type = p.curToken.Literal
		}
		params = append(params, nextParam)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}
