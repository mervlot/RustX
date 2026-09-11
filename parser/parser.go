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
	ASSIGN       // = or :=
	OR           // ||
	AND          // &&
	EQUALS       // == or !=
	LESSGREATER  // > or < or >= or <=
	SUM          // + or -
	PRODUCT      // * or /
	PREFIX       // -x or !x
	CALL         // myFunction(X), i++, i--
	INDEX        // array[index]
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
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.INC:          CALL,
	token.DEC:          CALL,
	token.LPAREN:       CALL,
	token.LBRACK:       INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []*errors.SyntaxError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

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
	p.registerPrefix(token.LBRACE, p.parseBlockExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACK, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
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
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACK, p.parseIndexExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []*errors.SyntaxError {
	return p.errors
}

func (p *Parser) addError(msg string, tok token.Token) {
	p.errors = append(p.errors, errors.NewSyntaxError(msg, tok.Line, tok.Column))
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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
	case token.FUNC:
		// Named function definition at top level or statement level
		if p.peekToken.Type == token.IDENT {
			return p.parseFunctionLiteralStatement()
		}
		return p.parseExpressionStatement()
	default:
		// Check for shorthand assignment `x := expr`
		if p.curToken.Type == token.IDENT && p.peekToken.Type == token.COLON_ASSIGN {
			return p.parseShortVarStatement()
		}
		// Check for reassignment `x = expr`
		if p.curToken.Type == token.IDENT && p.peekToken.Type == token.ASSIGN {
			return p.parseAssignmentStatement()
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

	// Optional type annotation: `var x int = 5` or `var x int`
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		stmt.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if p.peekToken.Type == token.ASSIGN {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseShortVarStatement() *ast.ShortVarStatement {
	stmt := &ast.ShortVarStatement{}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // curToken is now :=
	stmt.Token = p.curToken

	p.nextToken() // move to expression start
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // curToken is now =
	stmt.Token = p.curToken

	p.nextToken() // move to expression start
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if p.curToken.Type != token.SEMICOLON && p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Token: p.curToken}

	p.nextToken()

	// Check if this is a 3-part for loop (contains semicolons) or single condition loop
	// We scan ahead in input to check if there is a semicolon before `{`
	if p.isThreePartForLoopHeader() {
		// 1. Init (optional)
		if p.curToken.Type != token.SEMICOLON {
			if p.curToken.Type == token.VAR {
				stmt.Init = p.parseVarStatement()
			} else if p.curToken.Type == token.IDENT && p.peekToken.Type == token.COLON_ASSIGN {
				stmt.Init = p.parseShortVarStatement()
			} else if p.curToken.Type == token.IDENT && p.peekToken.Type == token.ASSIGN {
				stmt.Init = p.parseAssignmentStatement()
			} else {
				stmt.Init = p.parseExpressionStatement()
			}
		}

		if p.curToken.Type == token.SEMICOLON {
			// passed init semicolon
		} else if p.peekToken.Type == token.SEMICOLON {
			p.nextToken()
		}

		p.nextToken() // move to condition

		// 2. Condition (optional)
		if p.curToken.Type != token.SEMICOLON {
			stmt.Condition = p.parseExpression(LOWEST)
		}

		if p.curToken.Type == token.SEMICOLON {
			// passed condition semicolon
		} else if p.peekToken.Type == token.SEMICOLON {
			p.nextToken()
		}

		p.nextToken() // move to post statement

		// 3. Post (optional)
		if p.curToken.Type != token.LBRACE {
			if p.curToken.Type == token.IDENT && p.peekToken.Type == token.ASSIGN {
				stmt.Post = p.parseAssignmentStatement()
			} else {
				stmt.Post = p.parseExpressionStatement()
			}
			p.nextToken()
		}
	} else {
		// Single condition for loop e.g. `for i < 10 { ... }` or `for { ... }`
		if p.curToken.Type != token.LBRACE {
			stmt.Condition = p.parseExpression(LOWEST)
			p.nextToken()
		}
	}

	if p.curToken.Type != token.LBRACE {
		p.addError(fmt.Sprintf("Expected '{' for loop body, got %s", p.curToken.Literal), p.curToken)
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) isThreePartForLoopHeader() bool {
	// Look ahead from current position to find if there is a semicolon before `{`
	lexerCopy := lexer.New(p.l.Input()[p.l.Position():])
	tok := lexerCopy.NextToken()
	for tok.Type != token.EOF && tok.Type != token.LBRACE {
		if tok.Type == token.SEMICOLON {
			return true
		}
		tok = lexerCopy.NextToken()
	}
	return false
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for p.curToken.Type != token.RBRACE && p.curToken.Type != token.EOF {
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
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		// If peek token is LBRACE, stop unless infix parse handles it
		if p.peekToken.Type == token.LBRACE {
			break
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("Could not parse %q as float", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curToken.Type == token.TRUE}
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
		Left:     left,
		Operator: p.curToken.Literal,
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

func (p *Parser) parseBlockExpression() ast.Expression {
	return p.parseBlockStatement()
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if p.peekToken.Type == token.IF {
			p.nextToken()
			expression.Alternative = &ast.BlockStatement{
				Token:      p.curToken,
				Statements: []ast.Statement{&ast.ExpressionStatement{Expression: p.parseIfExpression()}},
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseFunctionLiteralStatement() ast.Statement {
	fn := p.parseFunctionLiteral().(*ast.FunctionLiteral)
	return &ast.ExpressionStatement{
		Token:      fn.Token,
		Expression: fn,
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type: `func add(a int, b int) int { ... }`
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		lit.ReturnType = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.FunctionParam {
	params := []*ast.FunctionParam{}

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return params
	}

	p.nextToken()

	param := &ast.FunctionParam{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}

	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		param.Type = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}
	params = append(params, param)

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()

		param := &ast.FunctionParam{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}

		if p.peekToken.Type == token.IDENT {
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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACK) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACK)
	return array
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekToken.Type == end {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekToken.Type == token.COMMA {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
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

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.addError(msg, p.peekToken)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.addError(msg, p.curToken)
}
