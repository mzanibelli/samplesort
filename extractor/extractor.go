package extractor

type Extractor struct {
	fs     Storage
	exec   RunnerFunc
	decode DecodeFunc
	format string
	stdout chan *payload
	stderr chan error
}

type Storage interface {
	ReadAll(name string) ([]byte, error)
	Exists(name string) bool
}

type RunnerFunc func(src, dst string) error
type DecodeFunc func(content []byte) []map[string]interface{}

func New(fs Storage, exec RunnerFunc, decode DecodeFunc, format string) *Extractor {
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
		e.stderr <- err
		return
	}

	err = e.load(p, dst)
	if err != nil {
		e.stderr <- err
		return
	}

	e.stdout <- p
}

func (e *Extractor) load(p *payload, path string) error {
	content, err := e.fs.ReadAll(path)
	if err != nil {
		return err
	}
	p.data = e.decode(content)
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
