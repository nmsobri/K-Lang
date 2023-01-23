package lexer

import (
	"Klang/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
    +
    -
    /
    *
    {
    }
    !
    !=
    =
    ==
    >
    >=
    <
    <=
    :
    ;
    55
    55.67
    foo
    foo99
    fn
    return
    let
    true
    false
    if
    else

    let five = 5;
	  let ten = 10;
    let salaray = 100.33;
    foobar;

    99foo
    55.xx
    55.67xx
  `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.STAR, "*"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.BANG, "!"},
		{token.EQUAL_NOT, "!="},
		{token.ASSIGN, "="},
		{token.EQUAL, "=="},
		{token.GREATER, ">"},
		{token.GREATER_EQUAL, ">="},
		{token.LESSER, "<"},
		{token.LESSER_EQUAL, "<="},
		{token.COLON, ":"},
		{token.SEMICOLON, ";"},
		{token.INTEGER, "55"},
		{token.FLOATING, "55.67"},
		{token.IDENTIFIER, "foo"},
		{token.IDENTIFIER, "foo99"},
		{token.FUNCTION, "fn"},
		{token.RETURN, "return"},
		{token.LET, "let"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.IF, "if"},
		{token.ELSE, "else"},

		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INTEGER, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INTEGER, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENTIFIER, "salaray"},
		{token.ASSIGN, "="},
		{token.FLOATING, "100.33"},
		{token.SEMICOLON, ";"},

		{token.IDENTIFIER, "foobar"},
		{token.SEMICOLON, ";"},

		{token.ILLEGAL, "99foo"},
		{token.ILLEGAL, "55.xx"},
		{token.ILLEGAL, "55.67xx"},
		{token.EOF, "EOF"},
	}

	l := New(input)

	for _, test := range tests {
		tok := l.NextToken()

		if tok.Type != test.expectedType {
			t.Fatalf("Token type is not matching expected. want=`%q`, got=`%q`", test.expectedType, tok.Type)
		}

		if tok.Literal != test.expectedLiteral {
			t.Fatalf("Token literal is not matching expected. want=`%s`, got=`%s`", test.expectedLiteral, tok.Literal)
		}
	}
}
