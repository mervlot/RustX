package parser

import (
	"fmt"
	"strconv"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/errors"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/token"
)

// Precedences
const (
	_ int = iota
	LOWEST
	ASSIGN_PREC  // = or :=
	LOGICAL_OR   // ||
	LOGICAL_AND  // &&
	EQUALS       // == !=
	LESSGREATER  // > < >= <=
	SUM          // + -
	PRODUCT      // * /
	PREFIX       // !X or -X
	CALL         // fn(X), array[index], x++
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:      ASSIGN_PREC,
	token.DECL_ASSIGN: ASSIGN_PREC,
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
	token.INC:         CALL,
	token.DEC:         CALL,
	token.LPAREN:      CALL,
	token.LBRACKET:    CALL,
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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

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
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)

	// Read two tokens so curToken and peekToken are both set
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
		// Top-level or named function statement: func name() ...
		if p.peekToken.Type == token.IDENT {
			return p.parseFunctionStatement()
		}
		return p.parseExpressionStatement()
	case token.IDENT:
		// Check for x := val or x = val statement
		if p.peekToken.Type == token.DECL_ASSIGN {
			return p.parseDeclAssignStatement()
		} else if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
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

	// Optional type annotation before = or ;
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		stmt.Type = p.curToken.Literal
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

func (p *Parser) parseDeclAssignStatement() *ast.DeclAssignStatement {
	stmt := &ast.DeclAssignStatement{}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // curToken is now :=
	stmt.Token = p.curToken

	p.nextToken() // curToken is now start of expression
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignStatement() *ast.AssignStatement {
	stmt := &ast.AssignStatement{}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // curToken is now =
	stmt.Token = p.curToken

	p.nextToken() // curToken is start of value
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

	p.nextToken() // move to header

	// Distinguish 3-part for loop (init; cond; post) from single-condition for loop (for cond)
	// Look ahead to check if there is a semicolon in the loop header.
	if p.hasSemicolonAheadInHeader() {
		// 3-part for loop
		if p.curToken.Type != token.SEMICOLON {
			stmt.Init = p.parseForInitStatement()
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}

		if p.curToken.Type != token.SEMICOLON {
			stmt.Condition = p.parseExpression(LOWEST)
			p.nextToken()
		}
		if p.curToken.Type == token.SEMICOLON {
			p.nextToken()
		}

		if p.curToken.Type != token.LBRACE {
			stmt.Post = p.parseForPostStatement()
			p.nextToken()
		}
	} else {
		// Single-condition loop: for condition { ... }
		stmt.Condition = p.parseExpression(LOWEST)
		p.nextToken()
	}

	if p.curToken.Type != token.LBRACE {
		p.addError(errors.NewSyntaxError(fmt.Sprintf("Expected '{' for loop body, got %s", p.curToken.Literal), p.curToken.Line, p.curToken.Column))
		return nil
	}

	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) hasSemicolonAheadInHeader() bool {
	// Helper to inspect upcoming tokens up to '{' or EOF
	// We create a temporary lexer scanning remaining input
	remInput := p.l.Input()[p.l.Position():]
	tempLexer := lexer.New(remInput)

	// Include current token and peek token in scanning
	for {
		tok := tempLexer.NextToken()
		if tok.Type == token.SEMICOLON {
			return true
		}
		if tok.Type == token.LBRACE || tok.Type == token.EOF {
			return false
		}
	}
}

func (p *Parser) parseForInitStatement() ast.Statement {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.IDENT:
		if p.peekToken.Type == token.DECL_ASSIGN {
			return p.parseDeclAssignStatement()
		} else if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseForPostStatement() ast.Statement {
	if p.curToken.Type == token.IDENT {
		if p.peekToken.Type == token.ASSIGN {
			return p.parseAssignStatement()
		}
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseFunctionStatement() *ast.FunctionLiteral {
	fl := &ast.FunctionLiteral{Token: p.curToken}

	p.nextToken() // curToken is function name IDENT
	fl.Name = p.curToken.Literal

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fl.Parameters = p.parseFunctionParameters()

	// Optional return type
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		fl.ReturnType = p.curToken.Literal
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fl.Body = p.parseBlockStatement()
	return fl
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type
	if p.peekToken.Type == token.IDENT {
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

	if p.peekToken.Type == token.RPAREN {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	param := &ast.FunctionParameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	if p.peekToken.Type == token.IDENT {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	identifiers = append(identifiers, param)

	for p.peekToken.Type == token.COMMA {
		p.nextToken() // consume ','
		p.nextToken() // move to parameter name
		param := &ast.FunctionParameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekToken.Type == token.IDENT {
			p.nextToken()
			param.Type = p.curToken.Literal
		}
		identifiers = append(identifiers, param)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// Stop parsing expression if peekToken is '{' (e.g., control flow condition header)
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	val, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addError(errors.NewSyntaxError(fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal), p.curToken.Line, p.curToken.Column))
		return nil
	}

	lit.Value = val
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	val, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addError(errors.NewSyntaxError(fmt.Sprintf("Could not parse %q as float", p.curToken.Literal), p.curToken.Line, p.curToken.Column))
		return nil
	}

	lit.Value = val
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	p.nextToken() // move to condition
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			expression.Alternative = &ast.BlockStatement{
				Token: p.curToken,
				Statements: []ast.Statement{
					&ast.ExpressionStatement{
						Token:      p.curToken,
						Expression: p.parseIfExpression(),
					},
				},
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		}
	}

	return expression
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

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.addError(errors.NewSyntaxError(msg, p.peekToken.Line, p.peekToken.Column))
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.addError(errors.NewSyntaxError(msg, p.curToken.Line, p.curToken.Column))
}

func (p *Parser) addError(err error) {
	p.errors = append(p.errors, err)
}
