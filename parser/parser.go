package parser

import(
    "APlusPlus/ast"
    "APlusPlus/lexer"
    "APlusPlus/token"
)

type Parser struct {
    l *lexer.Lexer

    curToken token.Token
    peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}

    // Read 2 tokens so curToken and peekToken are both populated
    p.nextToken()
    p.nextToken()

    return p
}

func (p *Parser) nextToken() {
    p.curToken = p.peekToken
    p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    return nil
}
