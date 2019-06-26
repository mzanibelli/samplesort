package analyze

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type payload struct {
	dense *mat.Dense
}

func newPayload(data [][]float64) *payload {
	if len(data) == 0 {
		return nil
	}
	r, c := len(data), len(data[0])
	p := &payload{mat.NewDense(r, c, make([]float64, r*c, r*c))}
	for i := range data {
		p.dense.SetRow(i, data[i])
	}
	return p
}

func (p *payload) data() [][]float64 {
	r, c := p.dense.Dims()
	res := make([][]float64, r, r)
	for i := range res {
		row := make([]float64, c, c)
		copy(row, p.dense.RawRowView(i))
		res[i] = row
	}
	return res
}

func (p *payload) apply(f func(i, j int, v float64) float64) *payload {
	p.dense.Apply(f, p.dense)
	return p
}

func (p *payload) String() string {
	f := mat.Formatted(p.dense, mat.Prefix(""), mat.Squeeze())
	return fmt.Sprintf("%v\n", f)
}
