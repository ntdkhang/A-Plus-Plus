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
    STRING = "STRING"

    // Operators
    ASSIGN = "="
    PLUS = "+"
    MINUS = "-"
    BANG = "!"
    ASTERISK = "*"
    SLASH = "/"

    LT = "<"
    GT = ">"

    EQ = "=="
    NOT_EQ = "!="

    // Delimiters
    COMMA = ","
    SEMICOLON = ";"

    LPAREN = "("
    RPAREN = ")"
    LBRACE = "{"
    RBRACE = "}"

    // Keywords
    FUNCTION = "FUNCTION"
    LET      = "LET"
    TRUE     = "TRUE"
    FALSE    = "FALSE"
    IF       = "IF"
    ELSE     = "ELSE"
    RETURN   = "RETURN"
)

var keywords = map[string]TokenType {
    "fn": FUNCTION,
    "let": LET,
    "true": TRUE,
    "false":  FALSE,
    "if":     IF,
    "else":   ELSE,
    "return": RETURN,
}

var two_char_operators = map[string]TokenType {
    "==": EQ,
    "!=": NOT_EQ,
}

func LookupOperator(op string) TokenType {
    if tok, ok := two_char_operators[op]; ok {
        return tok
    }
    return ILLEGAL
}

func LookupIdent(ident string) TokenType {
    if tok, ok := keywords[ident]; ok { // ok is a bool that returns true if the key exists in the map
        return tok
    }
    return IDENT
}
