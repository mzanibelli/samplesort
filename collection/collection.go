package collection

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type entity interface {
	fmt.Stringer
	Data() map[string]float64
}

type engine interface {
	Update(key string, n float64)
	Normalize(key string, n float64) float64
}

type Collection struct {
	entities []entity
	engine   engine
	mu       *sync.Mutex
}

func New(e engine) *Collection {
	return &Collection{make([]entity, 0), e, new(sync.Mutex)}
}

func (c *Collection) Append(e entity) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entities = append(c.entities, e)
	for key, val := range e.Data() {
		c.engine.Update(key, val)
	}
}

func (c *Collection) Features() [][]float64 {
	result := make([][]float64, len(c.entities), len(c.entities))
	for i, e := range c.entities {
		for key, val := range e.Data() {
			result[i] = append(result[i], c.engine.Normalize(key, val))
		}
	}
	return result
}

func (c *Collection) Sort(centers []int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entities) != len(centers) {
		return
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
