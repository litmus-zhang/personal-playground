package calculator

import (
	"fmt"
)

type Engine struct{}

func (e *Engine) Add(a, b float64) float64 {
	return a + b
}

func (e *Engine) Subtract(a, b float64) float64 {
	return a - b
}

func (e *Engine) Multiply(a, b float64) float64 {
	return a * b
}

func (e *Engine) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return a / b, nil
}
