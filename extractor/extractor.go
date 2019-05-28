package extractor

import (
	"log"
)

type Extractor struct {
	fs     Storage
	cmd    RunnerFactory
	decode DecodeFunc
	sink   chan *payload
	format string
	err    error
}

type Storage interface {
	ReadAll(name string) ([]byte, error)
	Exists(name string) bool
}

type RunnerFactory func(src, dst string) Runner

type DecodeFunc func(content []byte, data []map[string]interface{}) error

type Runner interface {
	Run() error
}

func New(fs Storage, cmd RunnerFactory, decode DecodeFunc, format string) *Extractor {
	return &Extractor{fs, cmd, decode, make(chan *payload), format, nil}
}

func (e *Extractor) Extract(src string) {
	s := new(payload)
	s.path = src
	s.data = make([]map[string]interface{}, 0)
	dst := src + e.format
	if !e.fs.Exists(dst) {
		e.err = e.cmd(src, dst).Run()
	}
	e.load(s, dst)
	if e.err == nil {
		e.sink <- s
	}
	if e.err != nil {
		log.Println("extract:", e.err)
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
