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

	lex.ReadChar()

	return lex
}

func (l *Lexer) ReadChar() {
	if l.readPosition >= len(l.source) {
		l.currentChar = 0
	} else {
		l.currentChar = l.source[l.readPosition]
	}

	l.currentPosition = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() token.Token {
	if l.readPosition > len(l.source) {
		return l.makeToken(token.EOF, "EOF")
	}

	var tok token.Token
	l.skipWhitespaceChar()

	switch l.CurrentChar() {
	case '+':
		tok = l.makeToken(token.PLUS, string(l.CurrentChar()))

	case '-':
		tok = l.makeToken(token.MINUS, string(l.CurrentChar()))

	case '/':
		tok = l.makeToken(token.SLASH, string(l.CurrentChar()))

	case '*':
		tok = l.makeToken(token.STAR, string(l.CurrentChar()))

	case '!':
		if l.isPeekChar('=') {
			l.ReadChar()
			tok = l.makeToken(token.EQUAL_NOT, string(l.source[l.currentPosition-1:l.readPosition]))
		} else {
			tok = l.makeToken(token.BANG, string(l.CurrentChar()))
		}

		l.readPosition++
		return tok

	case '=':
		if l.isPeekChar('=') {
			l.ReadChar()
			tok = l.makeToken(token.EQUAL, string(l.source[l.currentPosition-1:l.readPosition]))
		} else {
			tok = l.makeToken(token.ASSIGN, string(l.CurrentChar()))
		}

	case '>':
		if l.isPeekChar('=') {
			l.ReadChar()
			tok = l.makeToken(token.GREATER_EQUAL, string(l.source[l.currentPosition-1:l.readPosition]))
		} else {
			tok = l.makeToken(token.GREATER, string(l.CurrentChar()))
		}

	case '<':
		var tok token.Token
		if l.isPeekChar('=') {
			l.ReadChar()
			tok = l.makeToken(token.LESSER_EQUAL, string(l.source[l.currentPosition-1:l.readPosition]))
		} else {
			tok = l.makeToken(token.LESSER, string(l.CurrentChar()))
		}

		l.readPosition++
		return tok

	case ':':
		tok = l.makeToken(token.COLON, string(l.CurrentChar()))

	case ';':
		tok = l.makeToken(token.SEMICOLON, string(l.CurrentChar()))

	case 0:
		tok = l.makeToken(token.EOF, "EOF")

	default:
		if l.isNumber(l.CurrentChar()) {
			pos := l.currentPosition

			for l.isNumber(l.CurrentChar()) {
				l.ReadChar()
			}

			// integer
			if l.isWhitespaceChar(l.CurrentChar()) || l.isEof(l.CurrentChar()) {
				num := l.source[pos:l.currentPosition]

				// early exit. when we arrive here, we already sit at non number character
				// and at the bottom of this function, we read next character.
				// so we want to prevent double read of next character
				return l.makeToken(token.INTEGER, num)
			}

			// floating point
			if l.CurrentChar() == '.' && l.isNumber(l.PeekChar()) {
				l.ReadChar() // consume the `.`

				for l.isNumber(l.CurrentChar()) {
					l.ReadChar()
				}

				if l.isWhitespaceChar(l.CurrentChar()) || l.isEof(l.CurrentChar()) {
					num := l.source[pos:l.currentPosition]

					// early exit. when we arrive here, we already sit at non number character
					// and at the bottom of this function, we read next character.
					// so we want to prevent double read of next character
					return l.makeToken(token.FLOATING, num)
				}

				// illegal floating point ( contain non numeric character )
				for !l.isWhitespaceChar(l.CurrentChar()) && !l.isEof(l.CurrentChar()) {
					l.ReadChar()
				}

				// early exit. when we arrive here, we already sit at non number character
				// and at the bottom of this function, we read next character.
				// so we want to prevent double read of next character
				return l.makeToken(token.ILLEGAL, string(l.source[pos:l.currentPosition]))

			}

			// illegal number ( contain non numeric character )
			for !l.isWhitespaceChar(l.CurrentChar()) && !l.isEof(l.CurrentChar()) {
				l.ReadChar()
			}

			// early exit. when we arrive here, we already sit at non number character
			// and at the bottom of this function, we read next character.
			// so we want to prevent double read of next character
			return l.makeToken(token.ILLEGAL, string(l.source[pos:l.currentPosition]))

		} else if l.isAlphabet(l.CurrentChar()) {
			pos := l.currentPosition

			for l.isAlphaNum(l.CurrentChar()) {
				l.ReadChar()
			}

			// identifier/keyword
			if l.isWhitespaceChar(l.CurrentChar()) || l.isEof(l.CurrentChar()) {
				literal := l.source[pos:l.currentPosition]
				tokType := token.LookupIdent(literal)

				// early exit. when we arrive here, we already sit at non alphabet character
				// and at the bottom of this function, we read next character.
				// so we want to prevent double read of next character
				return l.makeToken(tokType, literal)
			}

			// illegal identifier/keyword ( contain non alpha numeric character like `!&^` )
			for !l.isWhitespaceChar(l.CurrentChar()) && !l.isEof(l.CurrentChar()) {
				l.ReadChar()
			}

			// early exit. when we arrive here, we already sit at non number character
			// and at the bottom of this function, we read next character.
			// so we want to prevent double read of next character
			return l.makeToken(token.ILLEGAL, string(l.source[pos:l.currentPosition]))

		} else {
			literal := l.CurrentChar()
			tok = l.makeToken(token.ILLEGAL, string(literal))
		}
	}

	l.ReadChar()
	return tok
}

func (l *Lexer) makeToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) isPeekChar(char byte) bool {
	if l.readPosition >= len(l.source) {
		return false
	}

	return l.PeekChar() == char
}

func (l *Lexer) PeekChar() byte {
	if l.readPosition >= len(l.source) {
		return 0
	}

	return l.source[l.readPosition]
}

func (l *Lexer) CurrentChar() byte {
	return l.currentChar
}

func (l *Lexer) skipWhitespaceChar() {
	for l.isWhitespaceChar(l.CurrentChar()) {
		l.ReadChar()
	}
}

func (l *Lexer) isWhitespaceChar(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' ||
		ch == '\r' || ch == '\f' || ch == '\v'
}

func (l *Lexer) isNumber(char byte) bool {
	return char >= '0' && char <= '9'
}

func (l *Lexer) isAlphabet(char byte) bool {
	return (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') || (char == '_')
}

func (l *Lexer) isAlphaNum(char byte) bool {
	return l.isNumber(char) || l.isAlphabet(char)
}

func (l *Lexer) isEof(ch byte) bool {
	return ch == 0
}
