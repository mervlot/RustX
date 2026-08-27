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
	LOGICAL_OR   // ||
	LOGICAL_AND  // &&
	EQUALS       // == !=
	LESSGREATER  // > < >= <=
	SUM          // + -
	PRODUCT      // * / %
	PREFIX       // !X -X
	CALL         // fn(X), arr[X], X++, X--
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
	token.LTE:          LESSGREATER,
	token.GTE:          LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.MODULO:       PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     CALL,
	token.INC:          CALL,
	token.DEC:          CALL,
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
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
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
	p.registerInfix(token.LTE, p.parseInfixExpression)
	p.registerInfix(token.GTE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.COLON_ASSIGN, p.parseColonAssignExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []*errors.SyntaxError {
	return p.errors
}

func (p *Parser) addError(msg string, tok token.Token) {
	err := errors.NewSyntaxError(msg, tok.Line, tok.Column)
	p.errors = append(p.errors, err)
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
	p.addError(fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type), p.peekToken)
	return false
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Type]; ok {
		return prec
	}
	return LOWEST
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
	case token.FUNCTION:
		// Top level or statement level `func name(...)`
		if p.peekTokenIs(token.IDENT) {
			expr := p.parseFunctionLiteral()
			if expr != nil {
				return &ast.ExpressionStatement{Tok: expr.GetToken(), Expression: expr}
			}
			return nil
		}
		return p.parseExpressionStatement()
	case token.SEMICOLON:
		return nil
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Tok: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}

	// Optional type annotation: e.g. `var age int = 16`
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Tok: p.curToken}

	p.nextToken()

	if !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt.ReturnValue = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Tok: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Tok: p.curToken}
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

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{Tok: p.curToken}

	p.nextToken() // move past 'while'

	// condition
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	stmt := &ast.ForStatement{Tok: p.curToken}

	// Determine if this is a 3-part for loop (`for init; cond; post`) or single cond (`for cond`) or infinite (`for`)
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		stmt.Body = p.parseBlockStatement()
		return stmt
	}

	// Check if there is a semicolon before the opening brace '{'
	isThreePart := p.hasSemicolonBeforeBrace()

	if isThreePart {
		p.nextToken() // move to init statement/expression

		if p.curTokenIs(token.VAR) {
			stmt.Init = p.parseVarStatement()
		} else if !p.curTokenIs(token.SEMICOLON) {
			stmt.Init = p.parseExpressionStatement()
		}

		if p.curTokenIs(token.SEMICOLON) {
			// Init statement consumed trailing semicolon
		} else if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		p.nextToken() // move to condition

		if !p.curTokenIs(token.SEMICOLON) {
			stmt.Condition = p.parseExpression(LOWEST)
		}

		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		p.nextToken() // move to post statement

		if !p.curTokenIs(token.LBRACE) {
			stmt.Post = p.parseExpressionStatement()
		}

		if p.peekTokenIs(token.LBRACE) {
			p.nextToken()
		}
		stmt.Body = p.parseBlockStatement()
	} else {
		// Single condition: `for condition { body }`
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		stmt.Body = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) hasSemicolonBeforeBrace() bool {
	// Lookahead in input from lexer position to see if ';' appears before '{'
	input := p.l.Input()
	pos := p.l.Position()
	braceCount := 0
	parenCount := 0

	for i := pos; i < len(input); i++ {
		ch := input[i]
		if ch == '(' {
			parenCount++
		} else if ch == ')' {
			if parenCount > 0 {
				parenCount--
			}
		} else if ch == '{' && parenCount == 0 {
			if braceCount == 0 {
				return false
			}
		} else if ch == ';' && parenCount == 0 {
			return true
		}
	}
	return false
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type), p.curToken)
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
	return &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Tok: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = val
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Tok: p.curToken}

	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as float", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = val
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Tok: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Tok: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Tok:      p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Tok:      p.curToken,
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
		Tok:      p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}
}

func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.addError("left side of assignment must be an identifier", left.GetToken())
		return nil
	}

	expression := &ast.AssignExpression{
		Tok:  p.curToken,
		Name: ident,
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

	return expression
}

func (p *Parser) parseColonAssignExpression(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.addError("left side of shorthand assignment must be an identifier", left.GetToken())
		return nil
	}

	expression := &ast.ColonAssignExpression{
		Tok:  p.curToken,
		Name: ident,
	}

	p.nextToken()
	expression.Value = p.parseExpression(LOWEST)

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
	expression := &ast.IfExpression{Tok: p.curToken}

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
			// else if block wrapper
			nestedIf := p.parseIfExpression()
			if nestedIf != nil {
				expression.Alternative = &ast.BlockStatement{
					Tok:        nestedIf.GetToken(),
					Statements: []ast.Statement{&ast.ExpressionStatement{Tok: nestedIf.GetToken(), Expression: nestedIf}},
				}
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Tok: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type annotation: e.g. `func add(a int, b int) int`
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.ReturnType = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	parameters := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()

	param := p.parseParameter()
	parameters = append(parameters, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := p.parseParameter()
		parameters = append(parameters, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseParameter() *ast.Parameter {
	param := &ast.Parameter{Tok: p.curToken}
	param.Name = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = &ast.Identifier{Tok: p.curToken, Value: p.curToken.Literal}
	}

	return param
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Tok: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Tok: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Tok: p.curToken}
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
