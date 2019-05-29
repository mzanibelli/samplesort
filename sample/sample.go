package sample

import "github.com/jeremywohl/flatten"

type Sample struct {
	path string
	data map[string]float64
}

func New(path string) *Sample {
	return &Sample{path, make(map[string]float64, 0)}
}

func (s *Sample) String() string           { return s.path + "\n" }
func (s *Sample) Data() map[string]float64 { return s.Flatten() }

func (s *Sample) Flatten(chunks ...map[string]interface{}) map[string]float64 {
	if s.data != nil && len(chunks) == 0 {
		return s.data
	}
	for _, chunk := range chunks {
		flat, _ := flatten.Flatten(chunk, "", flatten.DotStyle)
		for key, val := range flat {
			if n, ok := val.(float64); ok {
				s.data[key] = n
			}
		}
	}
	return s.data
}
