package parser

import (
	"fmt"
	"rustx/ast"
	"rustx/errors"
	"rustx/lexer"
	"rustx/token"
	"strconv"
)

// Precedence levels for operators
const (
	_ int = iota
	LOWEST
	ASSIGN_PREC // =
	LOR         // ||
	LAND        // &&
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * / %
	PREFIX      // -X or !X
	CALL        // myFunction(X) or X++ or X--
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:       ASSIGN_PREC,
	token.COLON_ASSIGN: ASSIGN_PREC,
	token.OR:           LOR,
	token.AND:          LAND,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LTE:          LESSGREATER,
	token.GTE:          LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.MODULO:       PRODUCT,
	token.LPAREN:       CALL,
	token.PLUS_PLUS:    CALL,
	token.MINUS_MINUS:  CALL,
	token.LBRACKET:     INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses tokens into an AST.
type Parser struct {
	l      *lexer.Lexer
	errors []*errors.SyntaxError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates a new Parser instance.
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseBlockExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
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

	// Read two tokens, so curToken and peekToken are both set
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

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// Errors returns the accumulated syntax errors.
func (p *Parser) Errors() []*errors.SyntaxError {
	return p.errors
}

func (p *Parser) addError(msg string, tok token.Token) {
	p.errors = append(p.errors, &errors.SyntaxError{
		Message: msg,
		Line:    tok.Line,
		Column:  tok.Column,
	})
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s ('%s') instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	p.addError(msg, p.peekToken)
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
	case token.SEMICOLON:
		return nil
	default:
		return p.parseExpressionOrAssignStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: var x int = 5
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // consume '='
		p.nextToken() // advance to start of expression
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

	if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionOrAssignStatement() ast.Statement {
	// Look ahead for := shorthand assignment: ident := expr
	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
		stmt := &ast.AssignStatement{
			Token:    p.peekToken,
			Left:     &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
			IsDefine: true,
		}
		p.nextToken() // curToken is ':='
		p.nextToken() // curToken is start of expr
		stmt.Value = p.parseExpression(LOWEST)

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}

	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}

	// Check if this is an assignment: e.g. x = 5 or arr[0] = 5
	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // curToken is '='
		stmt := &ast.AssignStatement{
			Token:    p.curToken,
			Left:     expr,
			IsDefine: false,
		}
		p.nextToken() // curToken is start of expr
		stmt.Value = p.parseExpression(LOWEST)

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
		p.addError(fmt.Sprintf("No prefix parse function for '%s' found", p.curToken.Literal), p.curToken)
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
	lit := &ast.IntegerLiteral{Token: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = val
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("Could not parse %q as float", p.curToken.Literal), p.curToken)
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

func (p *Parser) parseBlockExpression() ast.Expression {
	// Block parsed as an expression statement if encountered
	return nil
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
		p.nextToken() // curToken is 'else'

		if p.peekTokenIs(token.IF) {
			p.nextToken() // curToken is 'if'
			// Else if: wrap nested if in a block statement
			nestedIf := p.parseIfExpression()
			expression.Alternative = &ast.BlockStatement{
				Token:      p.curToken,
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
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	p.nextToken()

	// Distinguish between `for condition { ... }` and `for init; cond; post { ... }`
	if p.isThreePartForLoop() {
		// Three-part for loop
		if !p.curTokenIs(token.SEMICOLON) {
			expression.Init = p.parseStatement()
		}
		// If parseStatement consumed the semicolon, curToken is the next token.
		// If curToken is still SEMICOLON, advance past it.
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// Condition
		if !p.curTokenIs(token.SEMICOLON) {
			expression.Condition = p.parseExpression(LOWEST)
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// Post statement
		p.nextToken()
		if !p.curTokenIs(token.LBRACE) {
			expression.Post = p.parseStatement()
		}

		if !p.curTokenIs(token.LBRACE) && !p.expectPeek(token.LBRACE) {
			return nil
		}
	} else {
		// Single condition for loop: for cond { ... }
		expression.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
	}

	expression.Body = p.parseBlockStatement()

	return expression
}

func (p *Parser) isThreePartForLoop() bool {
	// Look ahead using lexer copy to detect semicolon before '{'
	subLexer := lexer.New(p.l.Input()[p.l.Position():])
	tok := p.curToken
	if tok.Type == token.SEMICOLON {
		return true
	}
	for {
		t := subLexer.NextToken()
		if t.Type == token.SEMICOLON {
			return true
		}
		if t.Type == token.LBRACE || t.Type == token.EOF {
			return false
		}
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Function could be named: func add(a int, b int) int { ... }
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type: func (...) int { ... }
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.ReturnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.FunctionParameter {
	params := []*ast.FunctionParameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.FunctionParameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	// Optional type annotation on parameter: e.g. `a int`
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
	params = append(params, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		param = &ast.FunctionParameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			param.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}
		params = append(params, param)
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
