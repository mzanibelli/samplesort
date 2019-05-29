package extractor

type Extractor struct {
	fs     Storage
	exec   RunnerFunc
	decode DecodeFunc
	format string
	stdout chan *payload
	stderr chan error
	err    error
}

type Storage interface {
	ReadAll(name string) ([]byte, error)
	Exists(name string) bool
}

type RunnerFunc func(src, dst string) error
type DecodeFunc func(content []byte) []map[string]interface{}

func New(fs Storage, exec RunnerFunc, decode DecodeFunc, format string) *Extractor {
	return &Extractor{fs, exec, decode, format, make(chan *payload), make(chan error), nil}
}

func (e *Extractor) Extract(src string) {
	p := &payload{
		path: src,
		data: make([]map[string]interface{}, 0),
	}
	dst := src + e.format
	if !e.fs.Exists(dst) {
		e.err = e.exec(src, dst)
	}
	e.load(p, dst)
	if e.err == nil {
		e.stdout <- p
	} else {
		e.stderr <- e.err
		e.err = nil
	}
}

func (e *Extractor) load(p *payload, path string) {
	if e.err != nil {
		return
	}
	var content []byte
	content, e.err = e.fs.ReadAll(path)
	if e.err != nil {
		return
	}
	p.data = e.decode(content)
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
