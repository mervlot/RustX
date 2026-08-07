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
	OR          // ||
	AND         // &&
	EQUALS      // == or !=
	LESSGREATER // > or < or >= or <=
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // myFunction(X) or postfix ++/--
	INDEX       // array[index]
)

var precedences = map[token.TokenType]int{
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
	token.LPAREN:       CALL,
	token.INC:          CALL,
	token.DEC:          CALL,
	token.LBRACKET:     INDEX,
	token.ASSIGN:       LOWEST,
	token.COLON_ASSIGN: LOWEST,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l        *lexer.Lexer
	filename string

	curToken  token.Token
	peekToken token.Token
	errors    []error

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer, filename string) *Parser {
	p := &Parser{
		l:        l,
		filename: filename,
		errors:   []error{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)

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
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

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
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, &errors.SyntaxError{
		Message: msg,
		File:    p.filename,
		Line:    p.peekToken.Line,
		Column:  p.peekToken.Column,
	})
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
		// skip optional semicolons
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}
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
	default:
		// Could be a shorthand assignment: name := expression
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
			return p.parseShorthandAssignStatement()
		}
		// Could be standard assignment: name = expression or array[idx] = expression
		// To distinguish, standard assignment parsed as ExpressionStatement could also be converted,
		// or we can parse expression and then check if the next is =
		return p.parseExpressionOrAssignStatement()
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

func (p *Parser) parseShorthandAssignStatement() *ast.VarStatement {
	// curToken is IDENT
	identTok := p.curToken
	p.nextToken() // now curToken is :=

	stmt := &ast.VarStatement{Token: p.curToken}
	stmt.Name = &ast.Identifier{Token: identTok, Value: identTok.Literal}

	p.nextToken() // now curToken is start of expression

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionOrAssignStatement() ast.Statement {
	// We parse an expression first
	expr := p.parseExpression(LOWEST)

	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // consume the '='
		assignTok := p.curToken
		p.nextToken() // move to right-hand side

		stmt := &ast.AssignStatement{
			Token: assignTok,
			Name:  expr,
			Value: p.parseExpression(LOWEST),
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	if p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.RBRACE) || p.curTokenIs(token.EOF) {
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

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
		// Stop parsing expression if we see a '{' (for unparenthesized conditions)
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
	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.errors = append(p.errors, &errors.SyntaxError{
		Message: msg,
		File:    p.filename,
		Line:    p.curToken.Line,
		Column:  p.curToken.Column,
	})
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

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, &errors.SyntaxError{
			Message: msg,
			File:    p.filename,
			Line:    p.curToken.Line,
			Column:  p.curToken.Column,
		})
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
		p.errors = append(p.errors, &errors.SyntaxError{
			Message: msg,
			File:    p.filename,
			Line:    p.curToken.Line,
			Column:  p.curToken.Column,
		})
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
		Operator: p.curToken.Literal,
		Left:     left,
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
		p.nextToken() // consume current expression end
		p.nextToken() // consume ',' and point to next
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
			p.nextToken() // move to 'if'
			// Parse the 'else if' block as a single Statement containing an IfExpression,
			// but we can also represent it by embedding the alternate block statement.
			// Let's create a block statement with the next IfExpression.
			nestedIf := p.parseIfExpression()
			expression.Alternative = &ast.BlockStatement{
				Token: nestedIf.GetToken(),
				Statements: []ast.Statement{
					&ast.ExpressionStatement{Token: nestedIf.GetToken(), Expression: nestedIf},
				},
			}
		} else if p.expectPeek(token.LBRACE) {
			expression.Alternative = p.parseBlockStatement()
		} else {
			return nil
		}
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// Optional Name (for recursion or named function syntax)
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type
	if p.peekTokenIs(token.TYPE_INT) || p.peekTokenIs(token.TYPE_FLOAT) || p.peekTokenIs(token.TYPE_BOOL) || p.peekTokenIs(token.TYPE_STRING) {
		p.nextToken()
		lit.ReturnType = p.curToken.Literal
	} else if p.peekTokenIs(token.IDENT) { // Fallback if type annotation is some other ident
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

	// first parameter
	param := &ast.Parameter{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	// check if followed by type
	if p.peekTokenIs(token.TYPE_INT) || p.peekTokenIs(token.TYPE_FLOAT) || p.peekTokenIs(token.TYPE_BOOL) || p.peekTokenIs(token.TYPE_STRING) {
		p.nextToken()
		param.Type = p.curToken.Literal
	} else if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = p.curToken.Literal
	}
	identifiers = append(identifiers, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume current
		p.nextToken() // consume ',' and point to next name
		nextParam := &ast.Parameter{
			Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
		}
		if p.peekTokenIs(token.TYPE_INT) || p.peekTokenIs(token.TYPE_FLOAT) || p.peekTokenIs(token.TYPE_BOOL) || p.peekTokenIs(token.TYPE_STRING) {
			p.nextToken()
			nextParam.Type = p.curToken.Literal
		} else if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			nextParam.Type = p.curToken.Literal
		}
		identifiers = append(identifiers, nextParam)
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

func (p *Parser) parseWhileExpression() ast.Expression {
	expr := &ast.WhileExpression{Token: p.curToken}

	p.nextToken()
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

func (p *Parser) parseForExpression() ast.Expression {
	expr := &ast.ForExpression{Token: p.curToken}

	p.nextToken() // Move past 'for'

	// Semicolons in loop header distinguish it from standard while loop
	// We can check if we have a full Go-style 3-part loop `for init; cond; post {`
	// To do this, let's look ahead to see if there are semicolons.
	// We scan from the lexer's current lookahead without consuming.
	// Since we are in Pratt parser, we can parse the first statement/expression first,
	// and check if a semicolon follows it.

	// If the current token is '{', this is a for { ... } infinite loop (not requested but nice to support)
	if p.curTokenIs(token.LBRACE) {
		expr.Body = p.parseBlockStatement()
		return expr
	}

	// We can tentatively parse an expression/statement
	// Since init can be a declaration like `i := 0` or standard assign `i = 0` or expression `i`
	var firstStmt ast.Statement
	if p.curTokenIs(token.VAR) {
		firstStmt = p.parseVarStatement()
	} else if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
		firstStmt = p.parseShorthandAssignStatement()
	} else {
		// Could be standard assign or expression statement
		firstStmt = p.parseExpressionOrAssignStatement()
	}

	// Now check what follows firstStmt.
	// If the next token is a semicolon, then we have a 3-part loop: `for init; cond; post {`
	// Wait, if firstStmt was already parsed, it might have consumed the semicolon (some of our helpers do that)
	// Let's check `p.curToken` or `p.peekToken`
	// If the previous parser step consumed the semicolon, or if we are currently at semicolon:
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken() // consume semicolon
	}

	// If we are indeed in a 3-part loop, we need condition and post
	// We can check if there's a semicolon in the lexer input before the next brace, or we can simply
	// parse condition, then semicolon, then post.
	// Let's assume it is 3-part if we had a semicolon.
	// To be robust: let's verify if the next part is condition.
	expr.Init = firstStmt

	if !p.curTokenIs(token.LBRACE) {
		// Condition
		if !p.curTokenIs(token.SEMICOLON) {
			expr.Condition = p.parseExpression(LOWEST)
			if p.peekTokenIs(token.SEMICOLON) {
				p.nextToken()
			}
		}
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken() // consume semicolon
		}

		// Post
		if !p.curTokenIs(token.LBRACE) {
			if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
				expr.Post = p.parseShorthandAssignStatement()
			} else {
				expr.Post = p.parseExpressionOrAssignStatement()
			}
		}
	} else {
		// If no semicolon followed, this is a `for condition {` loop (behaving as while loop)
		// Our firstStmt was parsed as an ExpressionStatement.
		// We convert it to condition expression.
		if exprStmt, ok := firstStmt.(*ast.ExpressionStatement); ok {
			expr.Condition = exprStmt.Expression
			expr.Init = nil
			expr.Post = nil
		} else {
			// This is a syntax error
			msg := "Invalid loop condition"
			p.errors = append(p.errors, &errors.SyntaxError{
				Message: msg,
				File:    p.filename,
				Line:    p.curToken.Line,
				Column:  p.curToken.Column,
			})
			return nil
		}
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}
