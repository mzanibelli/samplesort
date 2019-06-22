package samplesort

import (
	"log"
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
	logger         *log.Logger
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
		logger:         nil,
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

func WithLogger(value *log.Logger) config {
	return func(p *parameters) error {
		p.logger = value
		return nil
	}
}

func (p *parameters) Log(vs ...interface{}) {
	if p.logger == nil {
		return
	}
	p.logger.Println(vs...)
}
