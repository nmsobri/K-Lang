package token

const (
	// Single character token
	PLUS      = "PLUS"      // `+`
	MINUS     = "MINUS"     // `-`
	STAR      = "STAR"      // `*`
	SLASH     = "SLASH"     // `/`
	ASSIGN    = "ASSIGN"    // `=`
	GREATER   = "GREATER"   // `>`
	LESSER    = "LESSER"    // `<`
	BANG      = "BANG"      // `!`
	SEMICOLON = "SEMICOLON" // `;`
	COLON     = "COLON"     // `:`
	LBRACE    = "LBRACE"    // `{`
	RBRACE    = "RBRACE"    // `}`
	LPAREN    = "LPAREN"    // `(`
	RPAREN    = "RPAREN"    // `)`
	COMMA     = "COMMA"     // `,`
	LBRACKET  = "LBRACKET"  // `[`
	RBRACKET  = "RBRACKET"  // `]`

	// Double character token
	EQUAL         = "EQUAL"         // `==`
	EQUAL_NOT     = "EQUAL_NOT"     // `!=`
	GREATER_EQUAL = "GREATER_EQUAL" // `>=`
	LESSER_EQUAL  = "LESSER_EQUAL"  // `<=`

	// Multiple character token
	INTEGER    = "INTEGER"    // `[0-9]+`
	FLOATING   = "FLOATING"   // `[0-9]+\.[0-9]+`
	IDENTIFIER = "IDENTIFIER" // `[a-zA-Z_][a-zA-Z_0-9]+`
	STRING     = "STRING"     // `"[^"]+"`

	// Special token
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	// Keywords
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	WHILE    = "WHILE"
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"while":  WHILE,
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func LookupIdent(literal string) TokenType {
	if tok, ok := keywords[literal]; ok {
		return tok
	}

	return IDENTIFIER
}
