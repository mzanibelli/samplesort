package engine

import "gonum.org/v1/gonum/stat"

type featStat struct {
	values []float64
	mean   float64
	std    float64
}

func newFeatStat(size int) *featStat {
	return &featStat{
		values: make([]float64, size, size),
		mean:   0,
		std:    0,
	}
}

func (s *featStat) setMeanStd() {
	s.mean, s.std = stat.MeanStdDev(s.values, s.weights())
}

// TODO: how to smartly weight features?
func (s *featStat) weights() []float64 { return nil }
