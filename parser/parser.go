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
	PRODUCT     // * or /
	PREFIX      // -x or !x
	CALL        // myFunction(x), post-increment ++, post-decrement --, array indexing x[y]
)

var precedences = map[token.TokenType]int{
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
	token.AND:          LOGICAL_AND,
	token.OR:           LOGICAL_OR,
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: ASSIGN,
	token.LPAREN:       CALL,
	token.LBRACKET:     CALL,
	token.PLUS_PLUS:    CALL,
	token.MINUS_MINUS:  CALL,
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
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignment)
	p.registerInfix(token.COLON_ASSIGN, p.parseAssignment)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerInfix(token.MINUS_MINUS, p.parsePostfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []error {
	return p.errors
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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be '%s', got '%s' instead", t, p.peekToken.Type)
	p.errors = append(p.errors, errors.NewSyntaxError(msg, p.peekToken.Line, p.peekToken.Col))
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

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

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
		// Stop if we encounter '{' to facilitate unparenthesized control flows
		if p.peekTokenIs(token.LBRACE) {
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("No prefix parse function for '%s' found", t)
	p.errors = append(p.errors, errors.NewSyntaxError(msg, p.curToken.Line, p.curToken.Col))
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, errors.NewSyntaxError(msg, p.curToken.Line, p.curToken.Col))
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as float", p.curToken.Literal)
		p.errors = append(p.errors, errors.NewSyntaxError(msg, p.curToken.Line, p.curToken.Col))
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
	tok := p.curToken
	elements := p.parseExpressionList(token.RBRACKET)
	return &ast.ArrayLiteral{
		Token:    tok,
		Elements: elements,
	}
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	var list []ast.Expression
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if p.peekTokenIs(end) {
			break
		}
		p.nextToken()
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
	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, errors.NewSyntaxError(
			fmt.Sprintf("postfix operator '%s' must be applied to an identifier, got %T", p.curToken.Literal, left),
			p.curToken.Line, p.curToken.Col,
		))
		return nil
	}
	return &ast.PostfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Operand:  ident,
	}
}

func (p *Parser) parseAssignment(left ast.Expression) ast.Expression {
	tok := p.curToken
	p.nextToken()
	val := p.parseExpression(LOWEST)
	return &ast.AssignmentExpression{
		Token:    tok,
		Left:     left,
		Operator: tok.Literal,
		Value:    val,
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
			elseIf := p.parseIfExpression()
			block := &ast.BlockStatement{
				Token:      elseIf.GetToken(),
				Statements: []ast.Statement{&ast.ExpressionStatement{Token: elseIf.GetToken(), Expression: elseIf}},
			}
			expression.Alternative = block
		} else if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			expression.Alternative = p.parseBlockStatement()
		} else {
			p.errors = append(p.errors, errors.NewSyntaxError("Expected '{' or 'if' after else", p.curToken.Line, p.curToken.Col))
			return nil
		}
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	tok := p.curToken

	var name *ast.Identifier
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	params := p.parseFunctionParameters()

	// Parse optional return type annotation before the opening brace
	for !p.curTokenIs(token.LBRACE) && !p.curTokenIs(token.EOF) {
		if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			break
		}
		p.nextToken()
	}

	if !p.curTokenIs(token.LBRACE) {
		p.errors = append(p.errors, errors.NewSyntaxError("Expected '{' for function body", p.curToken.Line, p.curToken.Col))
		return nil
	}

	body := p.parseBlockStatement()

	return &ast.FunctionLiteral{
		Token:      tok,
		Name:       name,
		Parameters: params,
		Body:       body,
	}
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	if p.peekTokenIs(token.IDENT) {
		p.nextToken() // consume type name
	}

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		if p.peekTokenIs(token.RPAREN) {
			break
		}
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)

		if p.peekTokenIs(token.IDENT) {
			p.nextToken() // consume type name
		}
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
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

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	tok := p.curToken
	args := p.parseExpressionList(token.RPAREN)
	return &ast.CallExpression{
		Token:     tok,
		Function:  left,
		Arguments: args,
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	tok := p.curToken
	p.nextToken()
	index := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return &ast.IndexExpression{
		Token: tok,
		Left:  left,
		Index: index,
	}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	stmt := &ast.WhileStatement{Token: p.curToken}

	if !p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
	} else {
		stmt.Condition = &ast.BooleanLiteral{Token: token.Token{Type: token.TRUE, Literal: "true"}, Value: true}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() ast.Statement {
	stmt := &ast.ForStatement{Token: p.curToken}

	if p.isThreePartForLoop() {
		p.nextToken()

		if p.curTokenIs(token.VAR) {
			stmt.Init = p.parseVarStatement()
		} else if !p.curTokenIs(token.SEMICOLON) {
			stmt.Init = p.parseExpressionStatement()
		}

		if !p.curTokenIs(token.SEMICOLON) {
			if !p.expectPeek(token.SEMICOLON) {
				return nil
			}
		}

		p.nextToken()

		if !p.curTokenIs(token.SEMICOLON) {
			stmt.Condition = p.parseExpression(LOWEST)
		}

		if !p.expectPeek(token.SEMICOLON) {
			return nil
		}

		p.nextToken()

		if !p.curTokenIs(token.LBRACE) {
			stmt.Post = p.parseExpressionStatement()
		}
	} else {
		if !p.peekTokenIs(token.LBRACE) {
			p.nextToken()
			stmt.Condition = p.parseExpression(LOWEST)
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) isThreePartForLoop() bool {
	input := p.l.Input()
	pos := p.l.Position()

	inString := false
	for i := pos; i < len(input); i++ {
		ch := input[i]
		if ch == '"' {
			inString = !inString
		}
		if !inString {
			if ch == ';' {
				return true
			}
			if ch == '{' {
				return false
			}
		}
	}
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
