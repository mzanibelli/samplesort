package extractor

import "fmt"

type Extractor struct {
	cache  cache
	build  buildFunc
	stdout chan *payload
	stderr chan error
}

func New(cache cache, build buildFunc) *Extractor {
	return &Extractor{cache, build, make(chan *payload), make(chan error)}
}

func (e *Extractor) Extract(src string) {
	p := &payload{
		path: src,
		data: make(map[string]interface{}, 0),
	}
	err := e.cache.Fetch(src, &(p.data), func() (interface{}, error) {
		return e.build(src)
	})
	if err != nil {
		e.stderr <- fmt.Errorf("%s: %v", src, err)
		return
	}
	e.stdout <- p
}

func (e *Extractor) Out() <-chan *payload { return e.stdout }
func (e *Extractor) Err() <-chan error    { return e.stderr }

func (e *Extractor) Close() {
	close(e.stdout)
	close(e.stderr)
}

type payload struct {
	path string
	data map[string]interface{}
}

func (p *payload) String() string               { return p.path }
func (p *payload) Data() map[string]interface{} { return p.data }
