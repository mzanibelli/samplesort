package extractor

type Extractor struct {
	fs     Storage
	exec   RunnerFunc
	decode DecodeFunc
	format string
	sink   chan *payload
	err    error
}

type Storage interface {
	ReadAll(name string) ([]byte, error)
	Exists(name string) bool
}

type RunnerFunc func(src, dst string) error
type DecodeFunc func(content []byte, data []map[string]interface{}) error

func New(fs Storage, exec RunnerFunc, decode DecodeFunc, format string) *Extractor {
	return &Extractor{fs, exec, decode, format, make(chan *payload), nil}
}

func (e *Extractor) Extract(src string) {
	s := &payload{
		path: src,
		data: make([]map[string]interface{}, 0),
	}
	dst := src + e.format
	if !e.fs.Exists(dst) {
		e.err = e.exec(src, dst)
	}
	e.load(s, dst)
	if e.err == nil {
		e.sink <- s
	}
	e.err = nil
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
	e.err = e.decode(content, p.data)
}

func (e *Extractor) Sink() <-chan *payload { return e.sink }
func (e *Extractor) Close()                { close(e.sink) }

type payload struct {
	path string
	data []map[string]interface{}
}

func (p *payload) String() string                 { return p.path }
func (p *payload) Data() []map[string]interface{} { return p.data }
