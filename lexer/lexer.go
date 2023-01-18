package lexer

import (
	"Klang/token"
)

type Lexer struct {
	cursor       int
	source       string
	PeekToken    token.Token
	CurrentToken token.Token
}

func New(source string) *Lexer {
	lex := &Lexer{source: source, cursor: 0}

	// Prime current and peek token
	lex.NextToken()
	lex.NextToken()

	return lex
}

func (l *Lexer) NextToken() {
	l.CurrentToken = l.PeekToken
	l.PeekToken = l.advanceToken()
}

func (l *Lexer) advanceToken() token.Token {
	if l.cursor >= len(l.source) {
		return l.makeToken(token.EOF, "EOF")
	}

	l.skipWhitespaceChar()

	switch l.CurrentChar() {
	case '+':
		tok := l.makeToken(token.PLUS, string(l.source[l.cursor]))
		l.cursor++
		return tok

	case '-':
		tok := l.makeToken(token.MINUS, string(l.source[l.cursor]))
		l.cursor++
		return tok

	case '/':
		tok := l.makeToken(token.SLASH, string(l.source[l.cursor]))
		l.cursor++
		return tok

	case '*':
		tok := l.makeToken(token.STAR, string(l.source[l.cursor]))
		l.cursor++
		return tok

	case '!':
		var tok token.Token

		if l.isPeekChar('=') {
			l.cursor++
			tok = l.makeToken(token.EQUAL_NOT, string(l.source[l.cursor]))
		} else {
			tok = l.makeToken(token.BANG, string(l.source[l.cursor]))
		}

		l.cursor++
		return tok

	case '=':
		var tok token.Token

		if l.isPeekChar('=') {
			pos := l.cursor
			l.cursor++
			tok = l.makeToken(token.EQUAL, string(l.source[pos:l.cursor+1]))
		} else {
			tok = l.makeToken(token.ASSIGN, string(l.source[l.cursor]))
		}

		l.cursor++
		return tok

	case '>':
		var tok token.Token

		if l.isPeekChar('=') {
			l.cursor++
			tok = l.makeToken(token.GREATER_EQUAL, string(l.source[l.cursor]))
		} else {
			tok = l.makeToken(token.GREATER, string(l.source[l.cursor]))
		}

		l.cursor++
		return tok

	case '<':
		var tok token.Token
		if l.isPeekChar('=') {
			l.cursor++
			tok = l.makeToken(token.LESSER_EQUAL, string(l.source[l.cursor]))
		} else {
			tok = l.makeToken(token.LESSER, string(l.source[l.cursor]))
		}

		l.cursor++
		return tok

	default:
		return token.Token{}
	}
}

func (l *Lexer) makeToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) isPeekChar(char byte) bool {
	if l.cursor+1 >= len(l.source) {
		return false
	}

	return l.PeekChar() == char
}

func (l *Lexer) PeekChar() byte {
	if l.cursor+1 >= len(l.source) {
		return 0
	}

	return l.source[l.cursor+1]
}

func (l *Lexer) CurrentChar() byte {
	if l.cursor >= len(l.source) {
		return 0
	}

	return l.source[l.cursor]
}

func (l *Lexer) skipWhitespaceChar() {
	for l.CurrentChar() == ' ' || l.CurrentChar() == '\t' ||
		l.CurrentChar() == '\n' || l.CurrentChar() == '\r' ||
		l.CurrentChar() == '\f' || l.CurrentChar() == '\v' {
		l.cursor++
	}

}
