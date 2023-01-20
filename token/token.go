package token

const (
	// Single character token
	PLUS    = "PLUS"    // `+`
	MINUS   = "MINUS"   // `-`
	STAR    = "STAR"    // `*`
	SLASH   = "SLASH"   // `/`
	ASSIGN  = "ASSIGN"  // `=`
	GREATER = "GREATER" // `>`
	LESSER  = "LESSER"  // `<`
	BANG    = "BANG"    // `!`

	// Double character token
	EQUAL         = "EQUAL"         // `==`
	EQUAL_NOT     = "EQUAL_NOT"     // `!=`
	GREATER_EQUAL = "GREATER_EQUAL" // `>=`
	LESSER_EQUAL  = "LESSER_EQUAL"  // `<=`

	// Multiple character token
	INTEGER    = "INTEGER"    // `[0-9]`
	FLOATING   = "FLOATING"   // `[0-9]\.[0-9]`
	IDENTIFIER = "IDENTIFIER" // `[0-9]\.[0-9]`

	// Special token
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
