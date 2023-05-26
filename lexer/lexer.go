package lexer

// TODO: Support Unicode

import(
    "APlusPlus/token"
)

type Lexer struct {
    input string 
    position int // position of the current char
    read_position int // position where we are currently reading after the current char (since we need to peek further into the input)
    ch byte // current char 
}

func (l *Lexer) next_token() token.Token {
    var tok token.Token 
    
    switch l.ch {
    case '=':
        tok = new_token(token.ASSIGN, l.ch)
    case ';':
        tok = new_token(token.SEMICOLON, l.ch)
    case '(':
        tok = new_token(token.LPAREN, l.ch)
    case ')':
        tok = new_token(token.RPAREN, l.ch)
    case '{':
        tok = new_token(token.LBRACE, l.ch)
    case '}':
        tok = new_token(token.RBRACE, l.ch)
    case ',':
        tok = new_token(token.COMMA, l.ch)
    case '+':
        tok = new_token(token.PLUS, l.ch)
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    }

    l.read_char()
    return tok
}

func new_token(token_type token.TokenType, ch byte) token.Token {
    return token.Token{Type: token_type, Literal: string(ch)}
}

func (l *Lexer) read_char() {
    // Change l.ch to next char and update position and read_position
    if l.read_position >= len(l.input) {
        l.ch = 0 // NUL
    } else {
        l.ch = l.input[l.read_position]
    }
    l.position = l.read_position // where we last read
    l.read_position += 1 // where we are going to read
}

func New(input string) *Lexer {
    l := &Lexer{input: input}
    l.read_char()
    return l
}



