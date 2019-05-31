package samplesort

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"samplesort/analyze"
	"samplesort/collection"
	"samplesort/crypto"
	"samplesort/engine"
	"samplesort/extractor"
	"samplesort/parser"
	"samplesort/sample"
)

const (
	Checksum string = "9c91599c118ad0f2eef14e7bbcc050d8c802d3175b8e1766c820c7ab5ce685f5"
	Version  string = "v2.1_beta2-linux-i686"

	input     string = ".wav"
	output    string = ".json"
	size      int    = 100
	threshold int    = 0
)

type result interface {
	fmt.Stringer
	Size() int
	Features() [][]float64
}

func SampleSort(root, executable string, loggers ...*log.Logger) (result, error) {
	bin, err := which(executable)
	if err != nil {
		return nil, err
	}

	ext := extractor.New(fs, bin, decode, output)
	par := parser.New(fs, ext, input)
	col := collection.New()
	eng := engine.New()
	ana := analyze.New(col, eng, size, threshold)

	go par.Parse(root)

	go func() {
		for err := range ext.Err() {
			for _, l := range loggers {
				l.Println(err)
			}
		}
	}()

	wg := new(sync.WaitGroup)
	for e := range ext.Out() {
		wg.Add(1)
		copy := e
		go func() {
			defer wg.Done()
			s := sample.New(copy.String())
			s.Flatten(copy.Data()...)
			col.Append(s)
		}()
	}
	wg.Wait()

	ana.Analyze()

	return col, nil
}

func which(path string) (func(src, dst string) error, error) {
	fd, err := fs.Open(path)
	defer fd.Close()
	switch {
	case err != nil:
		return nop,
			fmt.Errorf("Error opening executable: %v", err)
	case !crypto.Check(fd, Checksum):
		return nop,
			fmt.Errorf("SHA256 mismatch, expected version %q", Version)
	case len(os.Args) != 2:
		return nop,
			fmt.Errorf("Expected exactly one argument, got: %d", len(os.Args)-1)
	}
	return func(src, dst string) error {
		return exec.Command(path, src, dst).Run()
	}, nil
}

func decode(content []byte) ([]map[string]interface{}, error) {
	type tmp struct {
		LowLevel map[string]interface{} `json:"lowlevel"`
		Tonal    map[string]interface{} `json:"tonal"`
	}
	t := new(tmp)
	res := make([]map[string]interface{}, 0, 2)
	if err := json.Unmarshal(content, t); err != nil {
		return res, err
	}
	res = append(res, t.LowLevel, t.Tonal)
	return res, nil
}

func nop(src, dst string) error { return nil }
