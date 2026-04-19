package calc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdd_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, Add(tt.a, tt.b))
		})
	}
}

func TestSubtract_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"both positive numbers", 10, 3, 7},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -10, 3, -13},
		{"both negative", -2, -3, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, Subtract(tt.a, tt.b))
		})
	}
}

func TestDivide_Success(t *testing.T) {
	got, err := Divide(10, 2)
	require.NoError(t, err)
	require.Equal(t, 5, got)
}

func TestDivide_DivideByZero(t *testing.T) {
	_, err := Divide(10, 0)
	require.ErrorIs(t, err, ErrDivideByZero)
}
