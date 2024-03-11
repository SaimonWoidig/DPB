package numbers

import (
	"fmt"
	"slices"
)

// Factorize returns a list of prime factors of the input number.
// The input must be a positive integer.
func Factorize(n int) ([]int, error) {
	if n < 1 {
		return nil, fmt.Errorf("invalid input: %d", n)
	}
	if n == 1 {
		return []int{1}, nil
	}

	var factors []int

	for i := 2; i*i <= n; i++ {
		for n%i == 0 {
			factors = append(factors, i)
			n /= i
		}
	}
	if n > 1 {
		factors = append(factors, n)
	}

	slices.Sort(factors)
	return factors, nil
}
