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

func (l *Lexer) NextToken() token.Token {
    var tok token.Token 
    
    l.skipWhitespace()

    switch l.ch {
    case '=':
        if l.peekChar() == '=' {
            tok = l.makeTwoCharTok()
        } else {
            tok = newToken(token.ASSIGN, l.ch)
        }
    case ';':
        tok = newToken(token.SEMICOLON, l.ch)
    case '(':
        tok = newToken(token.LPAREN, l.ch)
    case ')':
        tok = newToken(token.RPAREN, l.ch)
    case '{':
        tok = newToken(token.LBRACE, l.ch)
    case '}':
        tok = newToken(token.RBRACE, l.ch)
    case ',':
        tok = newToken(token.COMMA, l.ch)
    case '+':
        tok = newToken(token.PLUS, l.ch)
    case '-':
        tok = newToken(token.MINUS, l.ch)
    case '<':
        tok = newToken(token.LT, l.ch)
    case '>':
        tok = newToken(token.GT, l.ch)
    case '*':
        tok = newToken(token.ASTERISK, l.ch)
    case '/':
        tok = newToken(token.SLASH, l.ch)
    case '!':
        if l.peekChar() == '=' {
            tok = l.makeTwoCharTok()
        } else {
            tok = newToken(token.BANG, l.ch)
        }
    case 0:
        tok.Literal = ""
        tok.Type = token.EOF
    default:
        if isLetter(l.ch) {
            tok.Literal = l.readIdentifier()
            tok.Type = token.Lookup_ident(tok.Literal)
            return tok
        } else if isDigit(l.ch) {
            tok.Type = token.INT 
            tok.Literal = l.readNumber()
            return tok
        } else {
            tok = newToken(token.ILLEGAL, l.ch)
        }
    }

    l.readChar()
    return tok
}

func (l *Lexer) makeTwoCharTok() token.Token {
    ch := l.ch 
    l.readChar()
    literal := string(ch) + string(l.ch)
    tok := token.Token{Type: token.Lookup_operator(literal), Literal: literal}
    return tok
}

func (l *Lexer) peekChar() byte {
    if l.read_position >= len(l.input) {
        return 0
    } else {
        return l.input[l.read_position]
    }
}

func (l *Lexer) readIdentifier() string {
    position := l.position 
    for isLetter(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}


func newToken(token_type token.TokenType, ch byte) token.Token {
    return token.Token{Type: token_type, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
    for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
        l.readChar()
    }
}

func (l *Lexer) readChar() {
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
    l.readChar()
    return l
}

func (l *Lexer) readNumber() string {
    // TODO: read floats and hex, or even octal
    position := l.position 
    for isDigit(l.ch) {
        l.readChar()
    }
    return l.input[position:l.position]
}


func isDigit(ch byte) bool {
    return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
    return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
