package calc

import "errors"

func Add(a, b int) int      { return a + b }
func Subtract(a, b int) int { return a - b }

var ErrDivideByZero = errors.New("division by zero")

func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
