package parser

import (
	"fmt"
	"rustx/ast"
	rxErrors "rustx/errors"
	"rustx/lexer"
	"rustx/token"
	"strconv"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN_PREC // = or :=
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	EQUALS      // ==, !=
	LESSGREATER // >, <, >=, <=
	SUM         // +, -
	PRODUCT     // *, /, %
	PREFIX      // -X, !X
	CALL        // myFunction(X), ++, --
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:       ASSIGN_PREC,
	token.COLON_ASSIGN: ASSIGN_PREC,
	token.OR:           LOGICAL_OR,
	token.AND:          LOGICAL_AND,
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LE:           LESSGREATER,
	token.GE:           LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.MODULO:       PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     INDEX,
	token.INCREMENT:    CALL,
	token.DECREMENT:    CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser parses tokens into an AST.
type Parser struct {
	l *lexer.Lexer

	errors []*rxErrors.SyntaxError

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// New creates a new Parser instance.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []*rxErrors.SyntaxError{},
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
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
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
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.INCREMENT, p.parsePostfixExpression)
	p.registerInfix(token.DECREMENT, p.parsePostfixExpression)

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

func (p *Parser) Errors() []*rxErrors.SyntaxError {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s ('%s') instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	p.errors = append(p.errors, &rxErrors.SyntaxError{
		Message: msg,
		Line:    p.peekToken.Line,
		Column:  p.peekToken.Column,
	})
}

func (p *Parser) addError(msg string, tok token.Token) {
	p.errors = append(p.errors, &rxErrors.SyntaxError{
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

// ParseProgram parses all statements until EOF.
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
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		// Check for := shorthand: ident := expr
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
			return p.parseColonAssignStatement()
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

	// Check if next token is a type annotation or '='
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = p.curToken.Literal
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // eat '='
		p.nextToken() // move to expression start
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseColonAssignStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{Token: token.Token{
		Type:    token.COLON_ASSIGN,
		Literal: ":=",
		Line:    p.peekToken.Line,
		Column:  p.peekToken.Column,
	}}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // advance to ':='
	p.nextToken() // advance to value start

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
	tok := p.curToken
	expr := p.parseExpression(LOWEST)

	// Check if this was an assignment: left = expr
	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // eat '='
		assignTok := p.curToken
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

	stmt := &ast.ExpressionStatement{Token: tok, Expression: expr}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.addError(fmt.Sprintf("No prefix parse function for %s ('%s')", p.curToken.Type, p.curToken.Literal), p.curToken)
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
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	hasParen := p.peekTokenIs(token.LPAREN)
	if hasParen {
		p.nextToken() // eat '('
	}

	p.nextToken() // start of condition
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
		p.nextToken() // eat 'else'

		if p.peekTokenIs(token.IF) {
			p.nextToken() // eat 'if'
			nestedIf := p.parseIfExpression()
			expression.Alternative = &ast.BlockStatement{
				Token: expression.Token,
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Token: expression.Token, Expression: nestedIf},
				},
			}
		} else if p.peekTokenIs(token.LBRACE) {
			p.nextToken() // eat '{'
			expression.Alternative = p.parseBlockStatement()
		} else {
			p.peekError(token.LBRACE)
			return nil
		}
	}

	return expression
}

func (p *Parser) parseWhileExpression() ast.Expression {
	expression := &ast.WhileExpression{Token: p.curToken}

	hasParen := p.peekTokenIs(token.LPAREN)
	if hasParen {
		p.nextToken()
	}

	p.nextToken()
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

	hasParen := p.peekTokenIs(token.LPAREN)
	if hasParen {
		p.nextToken()
	}

	// Check if this is a 3-part for loop (init; cond; post) or condition-only (for cond {}) or infinite (for {})
	if p.peekTokenIs(token.LBRACE) {
		// Infinite loop: for { body }
		p.nextToken()
		expression.Body = p.parseBlockStatement()
		return expression
	}

	// Scan ahead to check if there are semicolons in the header
	isThreePart := p.peekHasSemicolon()

	if isThreePart {
		p.nextToken() // move to init
		if !p.curTokenIs(token.SEMICOLON) {
			if p.curTokenIs(token.VAR) {
				expression.Init = p.parseVarStatement()
			} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
				expression.Init = p.parseColonAssignStatement()
			} else {
				expression.Init = p.parseExpressionStatement()
			}
		}

		// Ensure we are past the first semicolon
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// Condition
		if !p.curTokenIs(token.SEMICOLON) {
			expression.Condition = p.parseExpression(LOWEST)
			p.nextToken() // eat to semicolon
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken() // move past second semicolon
		}

		// Post statement
		if !p.curTokenIs(token.RPAREN) && !p.curTokenIs(token.LBRACE) {
			if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
				expression.Post = p.parseColonAssignStatement()
			} else {
				expression.Post = p.parseExpressionStatement()
			}
		}

		if hasParen {
			if p.peekTokenIs(token.RPAREN) {
				p.nextToken()
			}
		}
	} else {
		// Single condition for loop
		p.nextToken()
		expression.Condition = p.parseExpression(LOWEST)

		if hasParen {
			if !p.expectPeek(token.RPAREN) {
				return nil
			}
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()
	return expression
}

// peekHasSemicolon checks if upcoming tokens before '{' contain a semicolon
func (p *Parser) peekHasSemicolon() bool {
	tempLexer := lexer.New(p.l.Input()[p.l.Position():])
	braceCount := 0
	parenCount := 0
	for {
		tok := tempLexer.NextToken()
		if tok.Type == token.EOF {
			return false
		}
		if tok.Type == token.LPAREN {
			parenCount++
		} else if tok.Type == token.RPAREN {
			parenCount--
		} else if tok.Type == token.LBRACE && parenCount <= 0 {
			return false
		}
		if tok.Type == token.SEMICOLON && braceCount == 0 {
			return true
		}
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = p.curToken.Literal
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type
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

func (p *Parser) parseFunctionParameters() []*ast.Parameter {
	identifiers := []*ast.Parameter{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	param := &ast.Parameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	identifiers = append(identifiers, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // eat ','
		p.nextToken() // move to param name
		pParam := &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			pParam.Type = p.curToken.Literal
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
