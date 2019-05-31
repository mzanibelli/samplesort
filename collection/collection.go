package collection

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type entity interface {
	fmt.Stringer
	Values() []float64
}

type Collection struct {
	entities []entity
	mu       *sync.Mutex
}

func New() *Collection {
	return &Collection{make([]entity, 0), new(sync.Mutex)}
}

func (c *Collection) Append(e entity) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entities = append(c.entities, e)
}

func (c *Collection) Features() [][]float64 {
	result := make([][]float64, len(c.entities), len(c.entities))
	for i := range result {
		result[i] = c.entities[i].Values()
	}
	return result
}

func (c *Collection) Sort(centers []int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entities) != len(centers) {
		panic("dataset and analysis size mismatch")
	}
	sort.Slice(c.entities, func(i, j int) bool {
		return centers[i] < centers[j]
	})
}

func (c *Collection) String() string {
	b := new(strings.Builder)
	for _, e := range c.entities {
		b.WriteString(e.String())
	}
	return b.String()
}

func (c *Collection) Size() int { return len(c.entities) }
