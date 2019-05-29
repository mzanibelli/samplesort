package samplesort

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

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

	input     string  = ".wav"
	output    string  = ".json"
	precision float64 = 0.05
	size      int     = 100
	threshold int     = 0
)

func SampleSort(root, executable string, loggers ...*log.Logger) (fmt.Stringer, error) {
	bin, err := which(executable)
	if err != nil {
		return nil, err
	}

	ext := extractor.New(fs, bin, decode, output)
	par := parser.New(fs, ext, input)
	eng := engine.New(precision)
	col := collection.New(eng)

	go par.Parse(root)

	go func() {
		for err := range ext.Err() {
			for _, l := range loggers {
				l.Println(err)
			}
		}
	}()

	for e := range ext.Out() {
		s := sample.New(e.String())
		s.Flatten(e.Data()...)
		col.Append(s)
	}

	analyze.New(col, size, threshold).Analyze()

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

func decode(content []byte) []map[string]interface{} {
	type tmp struct {
		LowLevel map[string]interface{} `json:"lowlevel"`
		Tonal    map[string]interface{} `json:"tonal"`
	}
	t := new(tmp)
	res := make([]map[string]interface{}, 0, 2)
	if err := json.Unmarshal(content, t); err == nil {
		res = append(res, t.LowLevel, t.Tonal)
	}
	return res
}

func nop(src, dst string) error { return nil }
