package parser

import (
	"A-Plus-Plus/AST"
	"A-Plus-Plus/lexer"
	"A-Plus-Plus/token"
	"fmt"
	"strconv"
)


type (
    // This function gets called when we encountered a prefix operator
    prefixParseFn func() ast.Expression

    // This function gets called when we encountered an infix operator. Since we need to know the Expression before the operator,
    // we take an Expression as argument
    infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
    l *lexer.Lexer

    errors []string

    curToken token.Token
    peekToken token.Token

    // Map the token to the associated prefix parse function
    prefixParseFns map[token.TokenType]prefixParseFn

    // Map the token to the associated infix parse function
    infixParseFns map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{
        l: l,
        errors: []string{},
    }

    // Read 2 tokens so curToken and peekToken are both populated
    p.nextToken()
    p.nextToken()

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

    // associate function p.parseIdentifier with token type IDENT
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)
    p.registerPrefix(token.BANG, p.parsePrefixExpression)
    p.registerPrefix(token.MINUS, p.parsePrefixExpression)
    p.registerPrefix(token.TRUE, p.parseBoolean)
    p.registerPrefix(token.FALSE, p.parseBoolean)
    p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
    p.registerPrefix(token.IF, p.parseIfExpression)
    p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)


    p.infixParseFns = make(map[token.TokenType]infixParseFn)
    p.registerInfix(token.PLUS, p.parseInfixExpression)
    p.registerInfix(token.MINUS, p.parseInfixExpression)
    p.registerInfix(token.SLASH, p.parseInfixExpression)
    p.registerInfix(token.ASTERISK, p.parseInfixExpression)
    p.registerInfix(token.EQ, p.parseInfixExpression)
    p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
    p.registerInfix(token.LT, p.parseInfixExpression)
    p.registerInfix(token.GT, p.parseInfixExpression)
    p.registerInfix(token.LPAREN, p.parseCallExpression)
    return p
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
    exp := &ast.CallExpression{Token: p.curToken, Function: function}
    exp.Arguments = p.parseCallArguments()
    return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
    args := []ast.Expression{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return args
    }
    p.nextToken()
    args = append(args, p.parseExpression(LOWEST))
    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()
        args = append(args, p.parseExpression(LOWEST))
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return args
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
    lit := &ast.FunctionLiteral{Token: p.curToken}
    if !p.expectPeek(token.LPAREN) {
        return nil
    }

    lit.Parameters = p.parseFunctionParameters()

    if !p.expectPeek(token.LBRACE) {
        return nil
    }

    lit.Body = p.parseBlockStatement()

    return lit
}


func (p *Parser) parseFunctionParameters() []*ast.Identifier {
    identifiers := []*ast.Identifier{}

    if p.peekTokenIs(token.RPAREN) {
        p.nextToken()
        return identifiers
    }

    p.nextToken()

    ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
    identifiers = append(identifiers, ident)

    for p.peekTokenIs(token.COMMA) {
        p.nextToken()
        p.nextToken()
        ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
        identifiers = append(identifiers, ident)
    }

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return identifiers
}

func (p *Parser) parseIfExpression() ast.Expression {
    expression := &ast.IfExpression{Token: p.curToken}

    if !p.expectPeek(token.LPAREN) { // after 'if' keyword should be (
        return nil
    }

    p.nextToken()
    expression.Condition = p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAREN) { // after condition should be )
        return nil
    }

    if !p.expectPeek(token.LBRACE) { // should be {
        return nil
    }

    expression.Consequence = p.parseBlockStatement()

    if p.peekTokenIs(token.ELSE) {
        p.nextToken()
        if !p.expectPeek(token.LBRACE) {
            return nil
        }

        expression.Alternative = p.parseBlockStatement()
    }

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

func (p *Parser) parseGroupedExpression() ast.Expression {
    p.nextToken()

    exp := p.parseExpression(LOWEST)

    if !p.expectPeek(token.RPAREN) {
        return nil
    }

    return exp
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// Add entries to the prefix parse function map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

// Add entries to the infix parse function map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
    return p.errors
}

// Print error when peek token is not what we expect
func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
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
    case token.LET:
        return p.parseLetStatement()
    case token.RETURN:
        return p.parseReturnStatement()
    default:
        return p.parseExpressionStatement()
    }
}

const (
    _ int = iota // incrementing numbers as values. blank '_' takes value 0, the consts below take 1 to 7
    LOWEST
    EQUALS          // ==
    LESSGREATER     // > or <
    SUM             // +
    PRODUCT         // *
    PREFIX          // -X or !X
    CALL            // myFunction(X)
)

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
    msg := fmt.Sprintf("no prefix parse function for %s found", t)
    p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement {Token: p.curToken}
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
        infixFn := p.infixParseFns[p.peekToken.Type]
        if infixFn == nil {
            return leftExp
        }

        p.nextToken()

        leftExp = infixFn(leftExp)
    }

    return leftExp
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



func (p *Parser) parseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token: p.curToken}

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
    }  else {
        p.peekError(t)
        return false
    }
}


func (p *Parser) parseIntegerLiteral() ast.Expression {
    lit := &ast.IntegerLiteral{Token: p.curToken}

    value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
        p.errors = append(p.errors, msg)
        return nil
    }

    lit.Value = value

    return lit
}



func (p *Parser) parsePrefixExpression() ast.Expression {
    expression := &ast.PrefixExpression{
        Token:      p.curToken,
        Operator:   p.curToken.Literal,
    }

    p.nextToken()

    expression.Right = p.parseExpression(PREFIX)

    return expression
}


var precedences = map[token.TokenType]int {
    token.EQ:       EQUALS,
    token.NOT_EQ:   EQUALS,
    token.LT:       LESSGREATER,
    token.GT:       LESSGREATER,
    token.PLUS:     SUM,
    token.MINUS:    SUM,
    token.SLASH:    PRODUCT,
    token.ASTERISK: PRODUCT,
    token.LPAREN:   CALL,
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

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    expression := &ast.InfixExpression{
        Token:      p.curToken,
        Operator:   p.curToken.Literal,
        Left:       left,
    }

    precedence := p.curPrecedence()
    p.nextToken()
    expression.Right = p.parseExpression(precedence)

    return expression
}

func (p *Parser) parseBoolean() ast.Expression {
    return &ast.Boolean{
        Token: p.curToken,
        Value: p.curTokenIs(token.TRUE),
    }
}
