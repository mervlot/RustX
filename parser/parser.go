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
	COLON_ASSIGN // :=
	ASSIGN       // =
	EQUALS       // == or !=
	LESSGREATER  // > or < or >= or <=
	LOGICAL_OR   // ||
	LOGICAL_AND  // &&
	SUM          // + or -
	PRODUCT      // * or / or %
	PREFIX       // -X or !X
	CALL         // myFunction(X), X++, X--
	INDEX        // array[index]
)

var precedences = map[token.TokenType]int{
	token.OR:           LOGICAL_OR,
	token.AND:          LOGICAL_AND,
	token.EQ:           EQUALS,
	token.NOT_EQ:        EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LTE:          LESSGREATER,
	token.GTE:          LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.PERCENT:      PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     INDEX,
	token.PLUS_PLUS:    CALL,
	token.MINUS_MINUS:  CALL,
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: COLON_ASSIGN,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses tokens from a Lexer into an AST.
type Parser struct {
	l      *lexer.Lexer
	errors []*errors.SyntaxError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates a new Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []*errors.SyntaxError{},
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
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerInfix(token.MINUS_MINUS, p.parsePostfixExpression)

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

func (p *Parser) Errors() []*errors.SyntaxError {
	return p.errors
}

func (p *Parser) appendError(msg string, line, col int) {
	p.errors = append(p.errors, &errors.SyntaxError{
		Message: msg,
		Line:    line,
		Column:  col,
	})
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s ('%s') instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	p.appendError(msg, p.peekToken.Line, p.peekToken.Column)
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

// ParseProgram parses the entire input program.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

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
		// Check for := short variable assignment (e.g. x := 10)
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

	// Check if next token is a type annotation (e.g. var x int = 5)
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = p.curToken.Literal
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // consume '='
		p.nextToken() // move to expression
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseShortVarStatement() *ast.ShortVarStatement {
	stmt := &ast.ShortVarStatement{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	p.nextToken() // move to ':='
	stmt.Token = p.curToken

	p.nextToken() // move to expression
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

func (p *Parser) parseExpressionStatement() ast.Statement {
	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}

	// Check for assignment: target = value
	if p.peekTokenIs(token.ASSIGN) {
		assignTok := p.peekToken
		p.nextToken() // consume '='
		p.nextToken() // start of value expr
		val := p.parseExpression(LOWEST)

		stmt := &ast.AssignStatement{
			Token: assignTok,
			Name:  expr,
			Value: val,
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
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
		p.appendError(fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type),
			p.curToken.Line, p.curToken.Column)
		return nil
	}
	leftExp := prefix()

	// Stop if we hit '{' which starts block statements in control flow
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
	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.appendError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal),
			p.curToken.Line, p.curToken.Column)
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}
	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.appendError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal),
			p.curToken.Line, p.curToken.Column)
		return nil
	}
	lit.Value = val
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expr.Right = p.parseExpression(PREFIX)
	return expr
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
	return expr
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
	expr := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	expr.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return expr
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

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
	expr := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	// Optional parentheses for condition
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		expr.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	} else {
		expr.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // consume 'else'

		if p.peekTokenIs(token.IF) {
			// else if ...
			p.nextToken()
			elseIfExpr := p.parseIfExpression()
			expr.Alternative = &ast.BlockStatement{
				Token: expr.Consequence.Token,
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token:      elseIfExpr.GetToken(),
						Expression: elseIfExpr,
					},
				},
			}
		} else if p.expectPeek(token.LBRACE) {
			expr.Alternative = p.parseBlockStatement()
		} else {
			return nil
		}
	}

	return expr
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expr := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()
	if p.curTokenIs(token.LPAREN) {
		p.nextToken()
		expr.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
	} else {
		expr.Condition = p.parseExpression(LOWEST)
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()
	return expr
}

func (p *Parser) parseForExpression() ast.Expression {
	forTok := p.curToken
	p.nextToken()

	// Check if this is a 3-part for loop (contains semicolons before '{') or a while-style condition loop
	hasSemicolons := p.scanForSemicolonsBeforeBrace()

	if !hasSemicolons {
		// condition-only for loop: for condition { body }
		var condition ast.Expression
		if !p.curTokenIs(token.LBRACE) {
			if p.curTokenIs(token.LPAREN) {
				p.nextToken()
				condition = p.parseExpression(LOWEST)
				if !p.expectPeek(token.RPAREN) {
					return nil
				}
			} else {
				condition = p.parseExpression(LOWEST)
			}
		}

		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		body := p.parseBlockStatement()

		return &ast.ForExpression{
			Token:     forTok,
			Condition: condition,
			Body:      body,
		}
	}

	// 3-part for loop: for init; condition; post { body }
	var initStmt ast.Statement
	if !p.curTokenIs(token.SEMICOLON) {
		initStmt = p.parseStatement()
	}
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken() // move past first ';'
	}

	var condition ast.Expression
	if !p.curTokenIs(token.SEMICOLON) {
		condition = p.parseExpression(LOWEST)
		p.nextToken()
	}
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken() // move past second ';'
	}

	var postStmt ast.Statement
	if !p.curTokenIs(token.LBRACE) {
		postStmt = p.parseStatement()
	}

	if !p.curTokenIs(token.LBRACE) {
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
	}

	body := p.parseBlockStatement()

	return &ast.ForExpression{
		Token:     forTok,
		Init:      initStmt,
		Condition: condition,
		Post:      postStmt,
		Body:      body,
	}
}

// scanForSemicolonsBeforeBrace looks ahead to check if the for loop header has semicolons.
func (p *Parser) scanForSemicolonsBeforeBrace() bool {
	l := lexer.New(p.l.Input()[p.l.Position():])
	tok := p.curToken
	for tok.Type != token.LBRACE && tok.Type != token.EOF {
		if tok.Type == token.SEMICOLON {
			return true
		}
		tok = l.NextToken()
	}
	return false
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Check if this is a named function (e.g. func add(a, b) { ... })
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type annotation (e.g. func add(a int, b int) int { ... })
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
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		pParam := &ast.FunctionParam{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			pParam.Type = p.curToken.Literal
		}
		params = append(params, pParam)
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
