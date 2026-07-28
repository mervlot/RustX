package parser

import (
	"fmt"
	"strconv"

	"github.com/username/rustx/ast"
	"github.com/username/rustx/errors"
	"github.com/username/rustx/lexer"
	"github.com/username/rustx/token"
)

// Precedence levels
const (
	_ int = iota
	LOWEST
	ASSIGN      // = or :=
	OR          // ||
	AND         // &&
	EQUALS      // == or !=
	LESSGREATER // > or < or >= or <=
	SUM         // + or -
	PRODUCT     // * or /
	PREFIX      // -X or !X
	CALL        // myFunction(X), array[index], post++ / post--
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
	token.LT_EQ:        LESSGREATER,
	token.GT_EQ:        LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
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
	p.registerPrefix(token.WHILE, p.parseWhileExpression)
	p.registerPrefix(token.FOR, p.parseForExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
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
	p.registerInfix(token.INC, p.parsePostfixExpression)
	p.registerInfix(token.DEC, p.parsePostfixExpression)

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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
	err := &errors.SyntaxError{
		Message: msg,
		Line:    p.peekToken.Line,
		Column:  p.peekToken.Column,
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) syntaxError(msg string, tok token.Token) {
	err := &errors.SyntaxError{
		Message: msg,
		Line:    tok.Line,
		Column:  tok.Column,
	}
	p.errors = append(p.errors, err)
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
	case token.FUNC:
		if p.peekTokenIs(token.IDENT) {
			return p.parseFunctionStatement()
		}
		return p.parseExpressionStatement()
	default:
		// Check for assignment statements, e.g. x = y
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignmentStatement()
		}
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.COLON_ASSIGN) {
			return p.parseShortAssignStatement()
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

	// Optional type annotations are ignored during parsing, just parse them if present
	if p.peekTokenIs(token.IDENT) {
		p.nextToken() // skip type annotation, e.g. 'int'
	}

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

func (p *Parser) parseShortAssignStatement() *ast.ShortAssignStatement {
	stmt := &ast.ShortAssignStatement{Token: p.peekToken} // pointing to ':='

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.nextToken() // now on ':='
	p.nextToken() // now on first token of expression

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignmentStatement() *ast.AssignmentStatement {
	stmt := &ast.AssignmentStatement{Token: p.peekToken} // pointing to '='

	// In the future this can be IndexExpression or Identifier
	left := p.parseExpression(LOWEST)
	stmt.Name = left

	p.nextToken() // now on '='
	p.nextToken() // now on first token of expression

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	stmt := &ast.FunctionStatement{Token: p.curToken}

	p.nextToken() // func -> Name
	name := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Name = name

	// Now we parse function literal but inject the Name
	fnLitExpr := p.parseFunctionLiteral()
	if fnLitExpr == nil {
		return nil
	}
	fnLit, ok := fnLitExpr.(*ast.FunctionLiteral)
	if !ok {
		return nil
	}
	fnLit.Name = name.Value
	stmt.Value = fnLit

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
		p.syntaxError(fmt.Sprintf("No prefix parse function for %s found", p.curToken.Type), p.curToken)
		return nil
	}
	leftExp := prefix()

	// To facilitate unparenthesized control flow, Pratt expression parsing should stop if it encounters '{'
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() && !p.peekTokenIs(token.LBRACE) {
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
		p.syntaxError(fmt.Sprintf("Could not parse %q as integer", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.syntaxError(fmt.Sprintf("Could not parse %q as float", p.curToken.Literal), p.curToken)
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
			// Handle else if via a single statement inside a block or directly as an alternative block
			// But for cleanliness we can just parse the IfExpression directly wrapped in a block statement
			// Or just parsing directly as an Alternative BlockStatement containing the nested IfExpression
			elseIfExp := p.parseIfExpression()
			block := &ast.BlockStatement{
				Token:      elseIfExp.GetToken(),
				Statements: []ast.Statement{&ast.ExpressionStatement{Token: elseIfExp.GetToken(), Expression: elseIfExp}},
			}
			expression.Alternative = block
		} else {
			if !p.expectPeek(token.LBRACE) {
				return nil
			}
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
	p.nextToken() // consume 'for'

	// Semicolons in Go-style loop headers: `for i := 0; i < 10; i++ { ... }`
	// Check if this is a standard `for condition {` loop, or a Go-style three-part loop
	// We can parse the first part.
	var init ast.Statement
	var cond ast.Expression
	var post ast.Statement

	// If the current token is not `{`, and if we can look ahead, we try to parse Go-style `for i := 0; i < 10; i++`
	// Actually we can parse the first part as a statement if we have semicolons coming up.
	// We can use a trick: parse an expression/statement, and if the next is semicolon, it is a three-part loop.
	// Otherwise, it is just a condition.
	if p.curTokenIs(token.LBRACE) {
		p.syntaxError("For loop missing condition/clauses", p.curToken)
		return nil
	}

	// Let's parse the first part as either: VarStatement, ShortAssign, assignment, or expression.
	// To do this cleanly, we can parse standard statements or expression.
	var isThreePart = false

	// Let's see if there is a semicolon in the future of this block.
	// Since we are in Pratt parser, we can do look-ahead, or we can just try to parse.
	// A simpler way: parse the first expression or statement.
	// But `var` or `:=` are statements.
	// Let's look for a semicolon in the token stream (without consuming too far or using a local helper).
	// Let's inspect the tokens. Since l is a Lexer, we don't have infinite backtracking, but we can look at peekToken.
	// If peekToken or some upcoming token is SEMICOLON, then it's a three-part loop.
	// Let's check peekToken first, or do a scan. Since we can't easily scan without cloning lexer,
	// can we know? Yes, if curToken is VAR, or if peekToken is COLON_ASSIGN, or if peekToken is ASSIGN,
	// or if we have an explicit semicolon after the expression.
	// Let's check if the loop is `for condition {` or `for init; cond; post {` by checking if there's a semicolon.
	// Let's check if there is a semicolon in the next few tokens. Since we only have peekToken, we can check if
	// the token parsed is followed by a semicolon.
	// Let's parse the first item as a statement.
	if p.curTokenIs(token.VAR) {
		isThreePart = true
		init = p.parseVarStatement() // this consumes up to the semicolon or end of expression
	} else if p.peekTokenIs(token.COLON_ASSIGN) {
		isThreePart = true
		init = &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: p.parseShortAssignStatement(),
		}
	} else if p.peekTokenIs(token.ASSIGN) {
		isThreePart = true
		init = &ast.ExpressionStatement{
			Token:      p.curToken,
			Expression: p.parseAssignmentStatement(),
		}
	} else {
		// Could be condition OR init expression (like `i = 0` or just `i < 10`).
		// Parse it as expression.
		expr := p.parseExpression(LOWEST)
		if p.peekTokenIs(token.SEMICOLON) {
			isThreePart = true
			init = &ast.ExpressionStatement{
				Token:      expr.GetToken(),
				Expression: expr,
			}
			p.nextToken() // consume the parsed expression's trailing token
		} else {
			cond = expr
		}
	}

	if isThreePart {
		// We expect a semicolon after init if it wasn't already consumed (like in var statement parsing)
		if p.curTokenIs(token.SEMICOLON) {
			p.nextToken()
		} else if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
			p.nextToken()
		}

		// Now parse condition
		if !p.curTokenIs(token.SEMICOLON) {
			cond = p.parseExpression(LOWEST)
			p.nextToken()
		}
		if !p.curTokenIs(token.SEMICOLON) {
			p.syntaxError("Expected ';' in for loop header", p.curToken)
			return nil
		}
		p.nextToken() // consume ';'

		// Now parse post statement/expression
		if !p.curTokenIs(token.LBRACE) {
			// post statement is typically x++, x = x + 1, etc.
			// Let's parse it as an expression and wrap in an ExpressionStatement.
			postExpr := p.parseExpression(LOWEST)
			post = &ast.ExpressionStatement{
				Token:      postExpr.GetToken(),
				Expression: postExpr,
			}
			p.nextToken()
		}
	}

	if !p.curTokenIs(token.LBRACE) {
		p.syntaxError(fmt.Sprintf("Expected '{' after for loop clauses, got %s", p.curToken.Type), p.curToken)
		return nil
	}

	expression.Init = init
	expression.Condition = cond
	expression.Post = post
	expression.Body = p.parseBlockStatement()

	return expression
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

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Optional return type annotation: check if there's a type name or ident before the '{'
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.ReturnType = &ast.TypeAnnotation{Token: p.curToken, Value: p.curToken.Literal}
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

	param := &ast.Parameter{Token: p.curToken}
	param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		param.Type = &ast.TypeAnnotation{Token: p.curToken, Value: p.curToken.Literal}
	}
	identifiers = append(identifiers, param)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		param := &ast.Parameter{Token: p.curToken}
		param.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		if p.peekTokenIs(token.IDENT) {
			p.nextToken()
			param.Type = &ast.TypeAnnotation{Token: p.curToken, Value: p.curToken.Literal}
		}
		identifiers = append(identifiers, param)
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

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}
