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
	ASSIGN       // = or :=
	LOGICAL_OR   // ||
	LOGICAL_AND  // &&
	EQUALS       // == !=
	LESSGREATER  // > < >= <=
	SUM          // + -
	PRODUCT      // * / %
	PREFIX       // !X or -X
	CALL         // myFunction(X) or arr[x] or i++
)

var precedences = map[token.TokenType]int{
	token.EQ:           EQUALS,
	token.NOT_EQ:        EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LT_EQ:        LESSGREATER,
	token.GT_EQ:        LESSGREATER,
	token.AND:          LOGICAL_AND,
	token.OR:           LOGICAL_OR,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.PERCENT:      PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     CALL,
	token.PLUS_PLUS:    CALL,
	token.MINUS_MINUS:  CALL,
	token.ASSIGN:       ASSIGN,
	token.COLON_ASSIGN: ASSIGN,
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
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

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
	p.registerInfix(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerInfix(token.MINUS_MINUS, p.parsePostfixExpression)

	// Read two tokens so curToken and peekToken are both set
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
	msg := fmt.Sprintf("expected next token to be %s, got %s ('%s') instead",
		t, p.peekToken.Type, p.peekToken.Literal)
	p.addError(msg, p.peekToken)
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
		// Check for shorthand assignment: ident := expr
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

	// Check if there is an optional type annotation before '='
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		stmt.Type = &ast.TypeAnnotation{Token: p.curToken, Name: p.curToken.Literal}
	}

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken()
		p.nextToken()
		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseShortVarStatement() *ast.ShortVarStatement {
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	p.nextToken() // move to :=
	colonAssignTok := p.curToken

	p.nextToken() // move to expression start
	val := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return &ast.ShortVarStatement{
		Token: colonAssignTok,
		Name:  name,
		Value: val,
	}
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

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // move to '='
		assignTok := p.curToken
		p.nextToken() // move to rhs
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
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

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found ('%s')", t, p.curToken.Literal)
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
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if p.peekTokenIs(token.IF) {
			p.nextToken()
			// parse 'else if' as an IfExpression wrapped in a block statement
			elseIfExp := p.parseIfExpression()
			if elseIfExp != nil {
				expression.Alternative = &ast.BlockStatement{
					Token: expression.Token,
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Token:      elseIfExp.GetToken(),
							Expression: elseIfExp,
						},
					},
				}
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		} else {
			return nil
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

func (p *Parser) isThreePartFor() bool {
	// Create a temp lexer on the rest of input from current position to look for semicolons before LBRACE
	subLex := lexer.New(p.l.Input()[p.l.Position():])
	braceCount := 0
	parenCount := 0
	for {
		tok := subLex.NextToken()
		if tok.Type == token.EOF {
			return false
		}
		if tok.Type == token.LPAREN {
			parenCount++
		} else if tok.Type == token.RPAREN {
			parenCount--
		} else if tok.Type == token.LBRACE && parenCount <= 0 {
			return false
		} else if tok.Type == token.SEMICOLON && parenCount <= 0 && braceCount <= 0 {
			return true
		}
	}
}

func (p *Parser) parseForExpression() ast.Expression {
	expression := &ast.ForExpression{Token: p.curToken}

	// Check if this is a 3-part for loop (contains semicolons before '{') or a simple conditional for loop
	if p.isThreePartFor() {
		p.nextToken() // advance past 'for'

		// 1. Init
		if !p.curTokenIs(token.SEMICOLON) {
			if p.curTokenIs(token.VAR) {
				expression.Init = p.parseVarStatement()
			} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
				expression.Init = p.parseShortVarStatement()
			} else {
				expression.Init = p.parseExpressionStatement()
			}
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// 2. Condition
		if !p.curTokenIs(token.SEMICOLON) {
			expression.Condition = p.parseExpression(LOWEST)
			p.nextToken()
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		}

		// 3. Post
		if !p.curTokenIs(token.LBRACE) {
			expression.Post = p.parseExpressionStatement()
		}

		if !p.curTokenIs(token.LBRACE) {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
		}

		expression.Body = p.parseBlockStatement()
		return expression
	}

	// Conditional for: `for condition { ... }` or infinite loop `for { ... }`
	if p.peekTokenIs(token.LBRACE) {
		p.nextToken()
		expression.Body = p.parseBlockStatement()
		return expression
	}

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Body = p.parseBlockStatement()
	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Check if there is an optional function name: func name(...) ...
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Check for optional return type: func (...) ReturnType {
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.ReturnType = &ast.TypeAnnotation{Token: p.curToken, Name: p.curToken.Literal}
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

	param := &ast.Parameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = &ast.TypeAnnotation{Token: p.curToken, Name: p.curToken.Literal}
	}
	parameters = append(parameters, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip comma
		p.nextToken() // move to next param name
		pParam := &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			pParam.Type = &ast.TypeAnnotation{Token: p.curToken, Name: p.curToken.Literal}
		}
		parameters = append(parameters, pParam)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
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
		p.nextToken() // skip ','
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
