package lexer

// TODO: Support Unicode for indentifier

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
    
    l.skip_whitespace()

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
    default:
        if is_letter(l.ch) {
            tok.Literal = l.read_identifier()
            tok.Type = token.Lookup_ident(tok.Literal)
            return tok
        } else if is_digit(l.ch) {
            tok.Type = token.INT 
            tok.Literal = l.read_number()
            return tok
        } else {
            tok = new_token(token.ILLEGAL, l.ch)
        }
    }

    l.read_char()
    return tok
}

func (l *Lexer) read_identifier() string {
    position := l.position 
    for is_letter(l.ch) {
        l.read_char()
    }
    return l.input[position:l.position]
}


func new_token(token_type token.TokenType, ch byte) token.Token {
    return token.Token{Type: token_type, Literal: string(ch)}
}

func (l *Lexer) skip_whitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.read_char()
    }
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

func (l *Lexer) read_number() string {
    position := l.position 
    for is_digit(l.ch) {
        l.read_char()
    }
    return l.input[position:l.position]
}


func is_digit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func is_letter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
