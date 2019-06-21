package samplesort

import (
	"io"
	"log"
	"os"
)

const (
	defaultFileSystemRoot string  = ""
	defaultAudioFormat    string  = ".wav"
	defaultDataFormat     string  = ".json"
	defaultSize           int     = 5
	defaultMaxIterations  int     = 10000
	defaultMaxZScore      float64 = 0.5
)

type parameters struct {
	fileSystemRoot string
	audioFormat    string
	dataFormat     string
	size           int
	maxIterations  int
	maxZScore      float64
	stdout         io.Writer
	stderr         io.Writer
	loggers        []*log.Logger
}

func (p *parameters) FileSystemRoot() string { return p.fileSystemRoot }
func (p *parameters) AudioFormat() string    { return p.audioFormat }
func (p *parameters) DataFormat() string     { return p.dataFormat }
func (p *parameters) Size() int              { return p.size }
func (p *parameters) MaxIterations() int     { return p.maxIterations }
func (p *parameters) MaxZScore() float64     { return p.maxZScore }

type config func(p *parameters) error

func newConfig(configs ...config) *parameters {
	params := &parameters{
		fileSystemRoot: defaultFileSystemRoot,
		audioFormat:    defaultAudioFormat,
		dataFormat:     defaultDataFormat,
		size:           defaultSize,
		maxIterations:  defaultMaxIterations,
		maxZScore:      defaultMaxZScore,
		stdout:         os.Stdout,
		stderr:         os.Stderr,
		loggers:        make([]*log.Logger, 0),
	}
	for _, setConfigTo := range configs {
		setConfigTo(params)
	}
	return params
}

func WithFileSystemRoot(value string) config {
	return func(p *parameters) error {
		p.fileSystemRoot = value
		return nil
	}
}

func WithAudioFormat(value string) config {
	return func(p *parameters) error {
		p.audioFormat = value
		return nil
	}
}

func WithDataFormat(value string) config {
	return func(p *parameters) error {
		p.dataFormat = value
		return nil
	}
}

func WithSize(value int) config {
	return func(p *parameters) error {
		p.size = value
		return nil
	}
}

func WithMaxIterations(value int) config {
	return func(p *parameters) error {
		p.maxIterations = value
		return nil
	}
}

func WithMaxZScore(value float64) config {
	return func(p *parameters) error {
		p.maxZScore = value
		return nil
	}
}

func WithLoggers(loggers ...*log.Logger) config {
	return func(p *parameters) error {
		p.loggers = loggers
		return nil
	}
}

// TODO: why is the config taking care of the logging?
// This is convenient but weird.
func (p *parameters) Out(vs ...interface{}) {
	for _, l := range p.loggers {
		l.SetOutput(p.stdout)
		l.Println(vs...)
	}
}

func (p *parameters) Err(vs ...interface{}) {
	for _, l := range p.loggers {
		l.SetOutput(p.stderr)
		l.Println(vs...)
	}
}
