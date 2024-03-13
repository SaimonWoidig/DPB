package numbers

import (
	"fmt"
	"slices"
)

// Factorize returns a list of prime factors of the input number.
// The input must be a positive integer.
func Factorize(n int) ([]int, error) {
	if n <= 1 {
		return nil, fmt.Errorf("invalid input")
	}

	// calculate prime factors
	var factors []int
	// iterate from 2 to sqrt(n)
	for i := 2; i*i <= n; i++ {
		// if n is divisible by i, add it to the list of prime factors
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
