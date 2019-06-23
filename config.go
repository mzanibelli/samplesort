package samplesort

import (
	"fmt"
	"log"
	"strings"
)

// TODO: default size makes the test fail.
// Is it because of RNG and the fact that KMeans is not a deterministic
// algorithm?
const (
	defaultFileSystemRoot string  = ""
	defaultAudioFormat    string  = ".wav"
	defaultDataFormat     string  = ".json"
	defaultSize           int     = 5
	defaultMaxIterations  int     = 100
	defaultMaxZScore      float64 = 0.5
	defaultEnableCache    bool    = true
)

type parameters struct {
	fileSystemRoot string
	audioFormat    string
	dataFormat     string
	size           int
	maxIterations  int
	maxZScore      float64
	enableCache    bool
	logger         *log.Logger
}

func (p *parameters) FileSystemRoot() string { return p.fileSystemRoot }
func (p *parameters) AudioFormat() string    { return p.audioFormat }
func (p *parameters) DataFormat() string     { return p.dataFormat }
func (p *parameters) Size() int              { return p.size }
func (p *parameters) MaxIterations() int     { return p.maxIterations }
func (p *parameters) MaxZScore() float64     { return p.maxZScore }
func (p *parameters) EnableCache() bool      { return p.enableCache }

type config func(p *parameters) error

func newConfig(configs ...config) *parameters {
	params := &parameters{
		fileSystemRoot: defaultFileSystemRoot,
		audioFormat:    defaultAudioFormat,
		dataFormat:     defaultDataFormat,
		size:           defaultSize,
		maxIterations:  defaultMaxIterations,
		maxZScore:      defaultMaxZScore,
		enableCache:    defaultEnableCache,
		logger:         nil,
	}
	for _, setConfigTo := range configs {
		setConfigTo(params)
	}
	return params
}

func (p *parameters) String() string {
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

func WithoutCache() config {
	return func(p *parameters) error {
		p.enableCache = false
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
