package lexer

import (
	"Klang/token"
)

type Lexer struct {
	source          string
	currentChar     byte
	readPosition    int
	currentPosition int
}

func New(source string) *Lexer {
	lex := &Lexer{
		source:          source,
		currentPosition: 0,
		readPosition:    0,
		currentChar:     0,
	}

	return lex
}

func (l *Lexer) ReadChar() byte {
	if l.readPosition >= len(l.source) {
		return 0
	}

	l.currentChar = l.source[l.readPosition]

	l.currentPosition = l.readPosition
	l.readPosition++
	return l.currentChar
}

func (l *Lexer) NextToken() token.Token {
	if l.readPosition >= len(l.source) {
		return l.makeToken(token.EOF, "EOF")
	}

	l.skipWhitespaceChar()

	switch l.ReadChar() {
	case '+':
		return l.makeToken(token.PLUS, string(l.CurrentChar()))

	case '-':
		return l.makeToken(token.MINUS, string(l.CurrentChar()))

	case '/':
		return l.makeToken(token.SLASH, string(l.CurrentChar()))

	case '*':
		return l.makeToken(token.STAR, string(l.CurrentChar()))

	case '!':
		var tok token.Token

		if l.isPeekChar('=') {
			l.readPosition++
			tok = l.makeToken(token.EQUAL_NOT, string(l.source[l.readPosition]))
		} else {
			tok = l.makeToken(token.BANG, string(l.source[l.readPosition]))
		}

		l.readPosition++
		return tok

	case '=':
		var tok token.Token

		if l.isPeekChar('=') {
			pos := l.readPosition
			l.readPosition++
			tok = l.makeToken(token.EQUAL, string(l.source[pos:l.readPosition+1]))
		} else {
			tok = l.makeToken(token.ASSIGN, string(l.source[l.readPosition]))
		}

		l.readPosition++
		return tok

	case '>':
		var tok token.Token

		if l.isPeekChar('=') {
			l.readPosition++
			tok = l.makeToken(token.GREATER_EQUAL, string(l.source[l.readPosition]))
		} else {
			tok = l.makeToken(token.GREATER, string(l.source[l.readPosition]))
		}

		l.readPosition++
		return tok

	case '<':
		var tok token.Token
		if l.isPeekChar('=') {
			l.readPosition++
			tok = l.makeToken(token.LESSER_EQUAL, string(l.source[l.readPosition]))
		} else {
			tok = l.makeToken(token.LESSER, string(l.source[l.readPosition]))
		}

		l.readPosition++
		return tok

	default:
		if l.isNumber(l.CurrentChar()) {
			pos := l.readPosition

			for l.isNumber(l.CurrentChar()) {
				l.readPosition++
			}

			num := l.source[pos:l.readPosition]
			return l.makeToken(token.INTEGER, num)

		} else if l.isAlphabet(l.CurrentChar()) {
			pos := l.readPosition

			for l.isAlphaNum(l.CurrentChar()) {
				l.readPosition++
			}

			identifier := l.source[pos:l.readPosition]
			return l.makeToken(token.IDENTIFIER, identifier)
		}

		return l.makeToken(token.ILLEGAL, "ILLEGAL")
	}
}

func (l *Lexer) makeToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) isPeekChar(char byte) bool {
	if l.readPosition+1 >= len(l.source) {
		return false
	}

	return l.PeekChar() == char
}

func (l *Lexer) PeekChar() byte {
	if l.readPosition+1 >= len(l.source) {
		return 0
	}

	return l.source[l.readPosition+1]
}

func (l *Lexer) CurrentChar() byte {
	return l.currentChar
}

func (l *Lexer) skipWhitespaceChar() {
	for l.CurrentChar() == ' ' || l.CurrentChar() == '\t' ||
		l.CurrentChar() == '\n' || l.CurrentChar() == '\r' ||
		l.CurrentChar() == '\f' || l.CurrentChar() == '\v' {
		l.readPosition++
	}

}

func (l *Lexer) isNumber(char byte) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexer) isAlphabet(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z')
}

func (l *Lexer) isAlphaNum(char byte) bool {
	return l.isNumber(char) || l.isAlphabet(char)
}
