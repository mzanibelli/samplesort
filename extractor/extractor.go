package extractor

import "fmt"

type storage interface {
	ReadAll(name string) ([]byte, error)
	Exists(name string) bool
}

type runnerFunc func(src, dst string) error
type decodeFunc func(content []byte) ([]map[string]interface{}, error)

type Extractor struct {
	fs     storage
	exec   runnerFunc
	decode decodeFunc
	format string
	stdout chan *payload
	stderr chan error
}

func New(fs storage, exec runnerFunc, decode decodeFunc, format string) *Extractor {
	return &Extractor{fs, exec, decode, format, make(chan *payload), make(chan error)}
}

func (e *Extractor) Extract(src string) {
	var err error
	p := &payload{
		path: src,
		data: make([]map[string]interface{}, 0),
	}

	dst := src + e.format
	if !e.fs.Exists(dst) {
		err = e.exec(src, dst)
	}
	if err != nil {
		e.stderr <- fmt.Errorf("%s: %v", src, err)
		return
	}

	err = e.load(p, dst)
	if err != nil {
		e.stderr <- fmt.Errorf("%s: %v", dst, err)
		return
	}

	e.stdout <- p
}

func (e *Extractor) load(p *payload, path string) error {
	content, err := e.fs.ReadAll(path)
	if err != nil {
		return err
	}
	p.data, err = e.decode(content)
	if err != nil {
		return err
	}
	return nil
}

func (e *Extractor) Out() <-chan *payload { return e.stdout }
func (e *Extractor) Err() <-chan error    { return e.stderr }

func (e *Extractor) Close() {
	close(e.stdout)
	close(e.stderr)
}

type payload struct {
	path string
	data []map[string]interface{}
}

func (p *payload) String() string                 { return p.path }
func (p *payload) Data() []map[string]interface{} { return p.data }
