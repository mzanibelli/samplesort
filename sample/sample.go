package sample

import (
	"fmt"
	"sort"

	"github.com/jeremywohl/flatten"
)

type Sample struct {
	path     string
	features []feature
}

func New(path string) *Sample {
	return &Sample{path, make([]feature, 0)}
}

func (s *Sample) String() string { return s.path + "\n" }

func (s *Sample) Flatten(chunks ...map[string]interface{}) {
	defer s.sort()
	s.features = make([]feature, 0)
	for _, chunk := range chunks {
		flat, _ := flatten.Flatten(chunk, "", flatten.DotStyle)
		for key, val := range flat {
			if n, ok := val.(float64); ok {
				s.features = append(s.features, feature{key, n})
			}
		}
	}
}

func (s *Sample) Values() []float64 {
	res := make([]float64, len(s.features))
	for i := range res {
		res[i] = s.features[i].value
	}
	return res
}

func (s *Sample) sort() {
	sort.Slice(s.features, func(i, j int) bool {
		switch {
		case s.features[i].key < s.features[j].key:
			return true
		case s.features[i].key > s.features[j].key:
			return false
		default:
			return true
		}
	})
}

type feature struct {
	key   string
	value float64
}

func (f feature) Key() string    { return f.key }
func (f feature) Value() float64 { return f.value }

type errNotFound struct{ key string }

func (e errNotFound) String() string { return fmt.Sprintf("feature not found: %s", e.key) }
func (e errNotFound) Error() string  { return e.String() }
