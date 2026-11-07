package parser

import (
	"fmt"
	"rustx/ast"
	"rustx/errors"
	"rustx/lexer"
	"rustx/token"
	"strconv"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN_PREC // =
	LOR         // ||
	LAND        // &&
	EQUALS      // ==, !=
	LESSGREATER // >, <, >=, <=
	SUM         // +
	PRODUCT     // *, /, %
	PREFIX      // -X, !X
	POSTFIX     // X++, X--
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.COLON_ASSIGN: ASSIGN_PREC,
	token.ASSIGN:       ASSIGN_PREC,
	token.OR:           LOR,
	token.AND:          LAND,
	token.EQ:           EQUALS,
	token.NOT_EQ:        EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LT_EQ:        LESSGREATER,
	token.GT_EQ:        LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.MODULO:       PRODUCT,
	token.PLUS_PLUS:    POSTFIX,
	token.MINUS_MINUS:  POSTFIX,
	token.LPAREN:       CALL,
	token.LBRACKET:     INDEX,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
	postfixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []*errors.SyntaxError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
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
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

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
	p.registerInfix(token.LT_EQ, p.parseInfixExpression)
	p.registerInfix(token.GT_EQ, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	p.postfixParseFns = make(map[token.TokenType]postfixParseFn)
	p.registerPostfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerPostfix(token.MINUS_MINUS, p.parsePostfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.addError(msg, p.peekToken)
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

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

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
	case token.WHILE:
		return p.parseWhileStatement()
	case token.FOR:
		return p.parseForStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
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

	// Check if optional type hint follows, e.g. `var age int = 16`
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.TypeHint = p.curToken.Literal
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // consume '='
		p.nextToken() // move to value
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

	// RustX supports:
	// for i := 0; i < 10; i++ { ... }
	// for cond { ... }
	// We can check if upcoming tokens before '{' contain a semicolon
	p.nextToken()

	// Check if this is a Go-style 3-part for loop
	// Lookahead to see if there is a semicolon before the block starts
	isThreePart := p.lookaheadForSemicolonBeforeBrace()

	if isThreePart {
		// 1. Init
		if !p.curTokenIs(token.SEMICOLON) {
			stmt.Init = p.parseStatementInForHeader()
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// 2. Cond
		if !p.curTokenIs(token.SEMICOLON) {
			stmt.Condition = p.parseExpression(LOWEST)
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken() // consume peek semicolon
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// 3. Post
		if !p.curTokenIs(token.LBRACE) {
			stmt.Post = p.parseStatementInForHeader()
		}

		if !p.curTokenIs(token.LBRACE) {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
		}
	} else {
		// Single condition: for cond { ... }
		stmt.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) lookaheadForSemicolonBeforeBrace() bool {
	pos := p.l.Position()
	if pos < 0 {
		pos = 0
	}
	tempLex := lexer.New(p.l.Input()[pos:])
	// scan tokens from upcoming input
	for {
		tok := tempLex.NextToken()
		if tok.Type == token.EOF || tok.Type == token.LBRACE {
			return false
		}
		if tok.Type == token.SEMICOLON {
			return true
		}
	}
}

func (p *Parser) parseStatementInForHeader() ast.Statement {
	if p.curTokenIs(token.VAR) {
		return p.parseVarStatement()
	}
	return p.parseExpressionOrAssignStatement()
}

func (p *Parser) parseExpressionOrAssignStatement() ast.Statement {
	startTok := p.curToken

	// Check for shorthand assignment: IDENT := EXPR
	if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() // on :=
		colonTok := p.curToken
		p.nextToken() // on value expr
		val := p.parseExpression(LOWEST)
		stmt := &ast.VarStatement{
			Token: colonTok,
			Name:  ident,
			Value: val,
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}

	expr := p.parseExpression(LOWEST)
	if expr == nil {
		return nil
	}

	// Check if this is an assignment: expr = expr
	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // move to =
		assignTok := p.curToken
		p.nextToken() // move to value
		val := p.parseExpression(LOWEST)
		stmt := &ast.AssignStatement{
			Token: assignTok,
			Left:  expr,
			Value: val,
		}
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
		return stmt
	}

	stmt := &ast.ExpressionStatement{Token: startTok, Expression: expr}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// Handle postfix if applicable (e.g. i++)
	if postfix, ok := p.postfixParseFns[p.peekToken.Type]; ok {
		p.nextToken()
		leftExp = postfix(leftExp)
	}

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// Stop infix parsing if next token is opening brace '{' (for unparenthesized if/while/for conditions)
		if p.peekTokenIs(token.LBRACE) {
			break
		}

		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)

		if postfix, ok := p.postfixParseFns[p.peekToken.Type]; ok {
			p.nextToken()
			leftExp = postfix(leftExp)
		}
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.addError(msg, p.curToken)
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.addError(msg, p.curToken)
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
		p.addError(msg, p.curToken)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // on 'else'

		if p.peekTokenIs(token.IF) {
			p.nextToken() // on 'if'
			nestedIf := p.parseIfExpression()
			expression.Alternative = &ast.BlockStatement{
				Token:      p.curToken,
				Statements: []ast.Statement{&ast.ExpressionStatement{Token: p.curToken, Expression: nestedIf}},
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Check if function has a name, e.g. `func add(a int, b int) int { ... }`
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type, e.g. `func add(...) int {`
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

func (p *Parser) parseFunctionParameters() []*ast.FunctionParameter {
	identifiers := []*ast.FunctionParameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	param := &ast.FunctionParameter{
		Token: p.curToken,
		Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	// Check optional type annotation
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.TypeHint = p.curToken.Literal
	}
	identifiers = append(identifiers, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // comma
		p.nextToken() // next param name
		pParam := &ast.FunctionParameter{
			Token: p.curToken,
			Name:  &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			pParam.TypeHint = p.curToken.Literal
		}
		identifiers = append(identifiers, pParam)
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
