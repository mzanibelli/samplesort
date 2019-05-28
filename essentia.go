package samplesort

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"samplesort/analyze"
	"samplesort/collection"
	"samplesort/crypto"
	"samplesort/engine"
	"samplesort/parser"
	"samplesort/sample"
)

const (
	EXPECTED_SHA256  string  = "9c91599c118ad0f2eef14e7bbcc050d8c802d3175b8e1766c820c7ab5ce685f5"
	EXPECTED_VERSION string  = "v2.1_beta2-linux-i686"
	ENV_EXTRACTOR    string  = "ESSENTIA_EXTRACTOR"
	EXT_IN           string  = ".wav"
	EXT_OUT          string  = ".json"
	PARAM_PRECISION  float64 = 0.05
)

func SampleSort() {
	_, err := which()
	if err != nil {
		usage(err)
	}
	sink := make(chan *sample.Sample)
	done := process(sink)
	err = parser.New(fs, nil).Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	close(sink)
	<-done
}

func usage(err error) {
	fmt.Fprintf(os.Stderr, "Usage: %s=xxx %s FILENAME\n", ENV_EXTRACTOR, os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %s - %s\n", EXPECTED_VERSION, EXPECTED_SHA256)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}

type loadFunc func(s *sample.Sample, dst string) error
type extractFunc func(path string)

func (f extractFunc) Extract(path string) { f(path) }

func run(bin, src, dst string, load loadFunc) error {
	s := new(sample.Sample)
	s.Path = src
	if fs.Exists(dst) {
		return load(s, dst)
	}
	log.Printf("%s: %s: start", path.Base(bin), src)
	cmd := exec.Command(bin, src, dst)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %s: %v", path.Base(bin), src, err)
	}
	return load(s, dst)
}

func loader(sink chan<- *sample.Sample) loadFunc {
	return func(s *sample.Sample, dst string) error {
		fd, err := fs.Open(dst)
		if err != nil {
			return fmt.Errorf("file: %s: %v", dst, err)
		}
		content, err := ioutil.ReadAll(fd)
		if err != nil {
			return fmt.Errorf("read: %s: %v", dst, err)
		}
		type payload struct {
			LowLevel map[string]interface{} `json:"lowlevel"`
			Tonal    map[string]interface{} `json:"tonal"`
		}
		p := new(payload)
		if err := json.Unmarshal(content, p); err != nil {
			return fmt.Errorf("json: %s: %v", dst, err)
		}
		s.Flatten(p.LowLevel, p.Tonal)
		sink <- s
		return nil
	}
}

func which() (string, error) {
	extractor := os.Getenv(ENV_EXTRACTOR)
	fd, err := fs.Open(extractor)
	defer fd.Close()
	switch {
	case err != nil:
		return "", fmt.Errorf("Error opening executable: %v", err)
	case !crypto.Check(fd, EXPECTED_SHA256):
		return "", fmt.Errorf("SHA256 mismatch, expected version %q", EXPECTED_VERSION)
	case len(os.Args) != 2:
		return "", fmt.Errorf("Expected exactly one argument, got: %d", len(os.Args))
	}
	return extractor, nil
}

func process(input <-chan *sample.Sample) <-chan struct{} {
	stats := engine.New(PARAM_PRECISION)
	data := collection.New(stats)
	for e := range input {
		data.Append(e)
	}
	return analyze.New(data).Analyze()
}
