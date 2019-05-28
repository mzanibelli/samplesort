package sample

import "github.com/jeremywohl/flatten"

type Sample struct {
	Path string
	data map[string]float64
}

func New(path string) *Sample {
	return &Sample{path, make(map[string]float64, 0)}
}

func (s *Sample) String() string           { return s.Path + "\n" }
func (s *Sample) Data() map[string]float64 { return s.Flatten() }

func (s *Sample) Flatten(chunks ...map[string]interface{}) map[string]float64 {
	if s.data != nil && len(chunks) == 0 {
		return s.data
	}
	for p := range floats(chunks...) {
		s.data[p.key] = p.value
	}
	return s.data
}

type pair struct {
	key   string
	value float64
}

func floats(chunks ...map[string]interface{}) <-chan pair {
	c := make(chan pair)
	send := func(flat map[string]interface{}) {
		for key, val := range flat {
			if n, ok := val.(float64); ok {
				c <- pair{key, n}
			}
		}
	}
	go func() {
		defer close(c)
		for _, chunk := range chunks {
			// ignore errors here
			flat, _ := flatten.Flatten(chunk, "", flatten.DotStyle)
			send(flat)
		}
	}()
	return c
}
