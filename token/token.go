package token

type TokenType string

type Token struct {
    Type TokenType 
    Literal string
}

// TokenTypes
const (
    ILLEGAL = "ILLEGAL"
    EOF = "EOF"

    // Identifiers + literals
    IDENT = "IDENT"
    INT = "INT"
    
    // Operators
    ASSIGN = "="
    PLUS = "+"
    
    // Delimiters
    COMMA = ","
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"

    // Keywords
    FUNCTION = "FUNCTION"
    LET = "LET"
)

var keywords = map[string]TokenType {
    "fn": FUNCTION,
    "let": LET,
}

func Lookup_ident(ident string) TokenType {
    if tok, ok := keywords[ident]; ok { // ok is a bool that returns true if the key exists in the map
        return tok
    }
    return IDENT
}
