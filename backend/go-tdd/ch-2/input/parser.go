package input

import "github.com/litmus-zhang/go-tdd/ch-2/calculator"

type Parser struct {
	engine    *calculator.Engine
	validator *Validator
}

func (p *Parser) ProcessExpression(expr string) (string, error) {
	return " ", nil
}
