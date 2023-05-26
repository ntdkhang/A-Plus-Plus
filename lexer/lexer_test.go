package lexer

import( 
    "testing"
    "APlusPlus/token"
)

func TestNextToken(t *testing.T) {
    input := `=+(){},;`

    tests := []struct {
        expected_type token.TokenType
        expected_literal string
    }{
        {token.ASSIGN, "="}, 
        {token.PLUS, "+"}, 
        {token.LPAREN, "("}, 
        {token.RPAREN, ")"}, 
        {token.LBRACE, "{"}, 
        {token.RBRACE, "}"}, 
        {token.COMMA, ","}, 
        {token.SEMICOLON, ";"}, 
        {token.EOF, ""}, 
    }

    l := New(input)

    for i, tt := range tests {
        tok := l.next_token()
        if tok.Type != tt.expected_type {
            t.Fatalf("tests[%d]: wrong TokenType. Expected: %q, got: %q", i, tt.expected_type, tok.Type)
        }
        
        if tok.Literal != tt.expected_literal {
            t.Fatalf("tests[%d]: wrong Literal. Expected: %q, got: %q", i, tt.expected_literal, tok.Literal)
        }
    }

}



