package collection

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
