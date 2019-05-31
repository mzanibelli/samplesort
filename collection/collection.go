package collection

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type entity interface {
	fmt.Stringer
	Keys() []string
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
	type tmp struct {
		count   int
		indices []int
	}
	cnts := make(map[string]*tmp, 0)
	for i, e := range c.entities {
		for j, key := range e.Keys() {
			_, ok := cnts[key]
			if !ok {
				cnts[key] = &tmp{0, make(
					[]int,
					len(c.entities), len(c.entities),
				)}
			}
			cnts[key].count++
			cnts[key].indices[i] = j
		}
	}
	res := make([][]float64, len(c.entities), len(c.entities))
	for i := range res {
		res[i] = make([]float64, 0)
		values := c.entities[i].Values()
		for _, cnt := range cnts {
			if cnt.count == len(c.entities) {
				res[i] = append(res[i], values[cnt.indices[i]])
			}
		}
	}
	return res
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
