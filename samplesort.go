package samplesort

import (
	"fmt"
	"io"
	"os/exec"

	"samplesort/analyze"
	"samplesort/cache"
	"samplesort/collection"
	"samplesort/engine"
	"samplesort/extractor"
	"samplesort/parser"
	"samplesort/sample"
)

const (
	Checksum string = "9c91599c118ad0f2eef14e7bbcc050d8c802d3175b8e1766c820c7ab5ce685f5"
	Version  string = "v2.1_beta2-linux-i686"
)

type sampleSort struct {
	cache      *cache.Cache
	extractor  *extractor.Extractor
	parser     *parser.Parser
	collection *collection.Collection
	engine     *engine.Engine
	analyze    *analyze.Analyze
	config     *parameters
}

func New(executable string, configs ...config) *sampleSort {
	s := &sampleSort{
		config:     newConfig(configs...),
		collection: collection.New(),
	}
	bin := which(executable, s.config.DataFormat())
	s.cache = cache.New(fs, s.config)
	s.extractor = extractor.New(s.cache, bin)
	s.parser = parser.New(fs, s.extractor, s.config)
	s.engine = engine.New(s.config)
	s.analyze = analyze.New(s.collection, s.engine, s.cache, s.config)
	return s
}

func (s *sampleSort) WriteTo(output io.Writer) (int64, error) {
	go s.parser.Parse(s.config.FileSystemRoot())
	go func() {
		for err := range s.extractor.Err() {
			s.config.Log(err)
		}
	}()
	for e := range s.extractor.Out() {
		smp := sample.New(e.String())
		smp.Flatten(e.Data())
		s.collection.Append(smp)
	}
	if err := s.analyze.Analyze(); err != nil {
		return 0, err
	}
	if s.collection.Len() > 0 {
		written, err := fmt.Fprintln(output, s.collection)
		return int64(written), err
	}
	return 0, nil
}

func which(bin, extension string) func(src string) (interface{}, error) {
	return func(src string) (interface{}, error) {
		return nil, exec.Command(bin, src, src+extension).Run()
	}
}
