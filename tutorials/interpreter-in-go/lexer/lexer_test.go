package lexer

import (
	"interpreter-in-go/token"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBasicNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		require.Equalf(t, tok.Type, tt.expectedType, "test[%v] - tokentype wrong. Expected=%q, Got=%q", i, tt.expectedType, tok.Type)
		require.Equalf(t, tok.Literal, tt.expectedLiteral, "test[%v] - literal wrong. Expected=%q, Got=%q", i, tt.expectedLiteral, tok.Literal)
	}
}
func TestNormalNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;

	let add = fn(x,y) {
		x+y;
	};

	let result = add(five, ten);
	`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()
		log.Printf("cur val: %v", string(input[i]))

		require.Equalf(t, tok.Type, tt.expectedType, "test[%v] - tokentype wrong. Expected=%q, Got=%q", i, tt.expectedType, tok.Type)
		require.Equalf(t, tok.Literal, tt.expectedLiteral, "test[%v] - literal wrong. Expected=%q, Got=%q", i, tt.expectedLiteral, tok.Literal)

	}
}

func TestLexer_readIdentifier(t *testing.T) {
	type fields struct {
		input        string
		position     int
		readPosition int
		ch           byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
		{
			name: "keywords",
			fields: fields{
				input: "let",
				position: 0,
				readPosition: 3,
				ch: 't',

			},
			want: token.LET,
		},
		{
			name: "keywords",
			fields: fields{
				input: "fn",
				position: 0,
				readPosition: 2,
				ch: 'n',
				
			},
			want: token.FUNCTION,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Lexer{
				input:        tt.fields.input,
				position:     tt.fields.position,
				readPosition: tt.fields.readPosition,
				ch:           tt.fields.ch,
			}
			if got := l.readIdentifier(); got != tt.want {
				t.Errorf("Lexer.readIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}
