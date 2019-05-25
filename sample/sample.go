package sample

import (
	"log"

	"github.com/jeremywohl/flatten"
)

type Unit interface {
	Extract(chunks ...map[string]interface{})
	Data() map[string]float64
}

type Sample struct {
	Path string
	data map[string]float64
}

func (s *Sample) Extract(chunks ...map[string]interface{}) { s.Flatten(chunks...) }
func (s *Sample) Data() map[string]float64                 { return s.Flatten() }
func (s *Sample) String() string                           { return s.Path + "\n" }

func (s *Sample) Flatten(chunks ...map[string]interface{}) map[string]float64 {
	switch {
	case s.data != nil && len(chunks) == 0:
		return s.data
	case s.data == nil:
		s.data = make(map[string]float64, 0)
		break
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
			flat, err := flatten.Flatten(chunk, "", flatten.DotStyle)
			if err != nil {
				log.Println("flatten:", err)
				return
			}
			send(flat)
		}
	}()
	return c
}
