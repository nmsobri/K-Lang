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

	// Double character token
	EQUAL         = "EQUAL"         // `==`
	EQUAL_NOT     = "EQUAL_NOT"     // `!=`
	GREATER_EQUAL = "GREATER_EQUAL" // `>=`
	LESSER_EQUAL  = "LESSER_EQUAL"  // `<=`

	// Multiple character token
	INTEGER    = "INTEGER"    // `[0-9]`
	FLOATING   = "FLOATING"   // `[0-9]\.[0-9]`
	IDENTIFIER = "IDENTIFIER" // `[a-z][A-Z][0-9]`

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
)

var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
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
