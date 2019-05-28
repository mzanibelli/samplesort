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
	checksum  string  = "9c91599c118ad0f2eef14e7bbcc050d8c802d3175b8e1766c820c7ab5ce685f5"
	version   string  = "v2.1_beta2-linux-i686"
	env       string  = "ESSENTIA_EXTRACTOR"
	whitelist string  = ".wav"
	format    string  = ".json"
	precision float64 = 0.05
)

func SampleSort() {
	bin, err := which()
	if err != nil {
		usage(err)
	}

	ext := extractor.New(fs, bin, decode, format)
	par := parser.New(fs, ext, whitelist)
	eng := engine.New(precision)
	data := collection.New(eng)

	if err := par.Parse(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	for e := range ext.Sink() {
		s := sample.New(e.String())
		s.Flatten(e.Data()...)
		data.Append(s)
	}

	done := analyze.New(data).Analyze()

	<-done
}

func decode(content []byte, data []map[string]interface{}) error {
	type tmp struct {
		LowLevel map[string]interface{} `json:"lowlevel"`
		Tonal    map[string]interface{} `json:"tonal"`
	}
	t := new(tmp)
	if err := json.Unmarshal(content, t); err != nil {
		return err
	} else {
		data = append(data, t.LowLevel, t.Tonal)
	}
	return nil
}

func which() (extractor.RunnerFactory, error) {
	nop := func(src, dst string) extractor.Runner { return nil }
	path := os.Getenv(env)
	fd, err := fs.Open(path)
	defer fd.Close()
	switch {
	case err != nil:
		return nop, fmt.Errorf("Error opening executable: %v", err)
	case !crypto.Check(fd, checksum):
		return nop, fmt.Errorf("SHA256 mismatch, expected version %q", version)
	case len(os.Args) != 2:
		return nop, fmt.Errorf("Expected exactly one argument, got: %d", len(os.Args))
	}
	return func(src, dst string) extractor.Runner {
		return exec.Command(path, src, dst)
	}, nil
}

func usage(err error) {
	fmt.Fprintf(os.Stderr, "Usage: %s=xxx %s FILENAME\n", env, os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %s - %s\n", version, checksum)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}
