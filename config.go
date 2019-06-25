package samplesort

import (
	"fmt"
	"log"
	"strings"
)

const (
	defaultFileSystemRoot string  = ""
	defaultAudioFormat    string  = ".wav"
	defaultDataFormat     string  = ".json"
	defaultSize           int     = 1
	defaultMaxIterations  int     = 1
	defaultMaxZScore      float64 = 2
	defaultEnableCache    bool    = true
)

type config struct {
	fileSystemRoot string
	audioFormat    string
	dataFormat     string
	size           int
	maxIterations  int
	maxZScore      float64
	enableCache    bool
	logger         *log.Logger
}

func (p *config) FileSystemRoot() string { return p.fileSystemRoot }
func (p *config) AudioFormat() string    { return p.audioFormat }
func (p *config) DataFormat() string     { return p.dataFormat }
func (p *config) Size() int              { return p.size }
func (p *config) MaxIterations() int     { return p.maxIterations }
func (p *config) MaxZScore() float64     { return p.maxZScore }
func (p *config) EnableCache() bool      { return p.enableCache }

type option func(p *config) error

func newConfig(options ...option) *config {
	params := &config{
		fileSystemRoot: defaultFileSystemRoot,
		audioFormat:    defaultAudioFormat,
		dataFormat:     defaultDataFormat,
		size:           defaultSize,
		maxIterations:  defaultMaxIterations,
		maxZScore:      defaultMaxZScore,
		enableCache:    defaultEnableCache,
		logger:         nil,
	}
	for _, applyOption := range options {
		applyOption(params)
	}
	return params
}

func (p *config) String() string {
	b := new(strings.Builder)
	b.WriteString(fmt.Sprintf(
		"data: %s\n", p.fileSystemRoot,
	))
	b.WriteString(fmt.Sprintf(
		"input: %s\n", p.audioFormat,
	))
	b.WriteString(fmt.Sprintf(
		"output: %s\n", p.dataFormat,
	))
	b.WriteString(fmt.Sprintf(
		"size: %d\n", p.size,
	))
	b.WriteString(fmt.Sprintf(
		"threshold: %d\n", p.maxIterations,
	))
	b.WriteString(fmt.Sprintf(
		"zscore: %.2f\n", p.maxZScore,
	))
	b.WriteString(fmt.Sprintf(
		"cache: %t\n", p.enableCache,
	))
	return b.String()
}

func WithFileSystemRoot(value string) option {
	return func(p *config) error {
		p.fileSystemRoot = value
		return nil
	}
}

func WithAudioFormat(value string) option {
	return func(p *config) error {
		p.audioFormat = value
		return nil
	}
}

func WithDataFormat(value string) option {
	return func(p *config) error {
		p.dataFormat = value
		return nil
	}
}

func WithSize(value int) option {
	return func(p *config) error {
		p.size = value
		return nil
	}
}

func WithMaxIterations(value int) option {
	return func(p *config) error {
		p.maxIterations = value
		return nil
	}
}

func WithMaxZScore(value float64) option {
	return func(p *config) error {
		p.maxZScore = value
		return nil
	}
}

func WithoutCache() option {
	return func(p *config) error {
		p.enableCache = false
		return nil
	}
}

func WithLogger(value *log.Logger) option {
	return func(p *config) error {
		p.logger = value
		return nil
	}
}

func (p *config) Log(vs ...interface{}) {
	if p.logger == nil {
		return
	}
	p.logger.Println(vs...)
}
