package collection

import (
	"fmt"
	"sort"
	"strings"
)

type entity interface {
	fmt.Stringer
	Keys() []string
	Values() []float64
}

type Collection struct {
	entities []entity
	centers  []int
	scores   map[string]*featScore
}

func New() *Collection {
	return &Collection{
		entities: make([]entity, 0),
		scores:   nil,
		centers:  nil,
	}
}

func (c *Collection) Append(e entity) {
	c.entities = append(c.entities, e)
}

func (c *Collection) Features() [][]float64 {
	if len(c.entities) == 0 {
		return nil
	}
	c.computeScores()
	data := make([][]float64, len(c.entities), len(c.entities))
	for i := range data {
		data[i] = c.orderedValues(i)
	}
	return data
}

func (c *Collection) Sort(centers []int) {
	if len(c.entities) != len(centers) {
		panic(fmt.Sprintf("dataset and analysis size mismatch: %d / %d",
			len(c.entities), len(centers)))
	}
	c.centers = centers
	sort.Sort(c)
}

func (c *Collection) String() string {
	b := new(strings.Builder)
	for _, e := range c.entities {
		b.WriteString(e.String())
	}
	return b.String()
}

func (c *Collection) Len() int {
	return len(c.entities)
}

func (c *Collection) Swap(i, j int) {
	c.entities[i], c.entities[j] = c.entities[j], c.entities[i]
	c.centers[i], c.centers[j] = c.centers[j], c.centers[i]
}

func (c *Collection) Less(i, j int) bool {
	return c.centers[i] < c.centers[j]
}

func newScore(size int) *featScore {
	return &featScore{0, make([]int, size, size), false}
}

type featScore struct {
	count   int
	indices []int
	valid   bool
}

func (c *Collection) computeScores() {
	c.scores = make(map[string]*featScore)
	for i, e := range c.entities {
		for j, key := range e.Keys() {
			c.updateScore(key, i, j)
		}
	}
}

func (c *Collection) updateScore(key string, i, j int) {
	if _, ok := c.scores[key]; !ok {
		c.scores[key] = newScore(len(c.entities))
	}
	c.scores[key].count++
	c.scores[key].indices[i] = j
	if c.scores[key].count == len(c.entities) {
		c.scores[key].valid = true
	}
}

func (c *Collection) orderedValues(i int) []float64 {
	res := make([]float64, 0, len(c.scores))
	values := c.entities[i].Values()
	for j, key := range c.entities[i].Keys() {
		if c.isCommon(key) {
			res = append(res, values[j])
		}
	}
	return res
}

func (c *Collection) isCommon(key string) bool {
	score, ok := c.scores[key]
	return ok && score.valid
}
