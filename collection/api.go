package collection

import "fmt"

type entity interface {
	fmt.Stringer
	Keys() []string
	Values() []float64
}
