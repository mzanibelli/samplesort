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
	size      int     = 100
	threshold int     = 0
)

func SampleSort() {
	bin := which()
	ext := extractor.New(fs, bin, decode, format)
	par := parser.New(fs, ext, whitelist)
	eng := engine.New(precision)
	col := collection.New(eng)

	go par.Parse(os.Args[1])

	go func() {
		for err := range ext.Err() {
			log.Println(err)
		}
	}()

	for e := range ext.Out() {
		s := sample.New(e.String())
		s.Flatten(e.Data()...)
		col.Append(s)
	}

	analyze.New(col, size, threshold).Analyze()

	fmt.Fprintln(os.Stdout, col)
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

func which() func(src, dst string) error {
	path := os.Getenv(env)
	fd, err := fs.Open(path)
	defer fd.Close()
	switch {
	case err != nil:
		usage(fmt.Errorf("Error opening executable: %v", err))
	case !crypto.Check(fd, checksum):
		usage(fmt.Errorf("SHA256 mismatch, expected version %q", version))
	case len(os.Args) != 2:
		usage(fmt.Errorf("Expected exactly one argument, got: %d", len(os.Args)))
	}
	return func(src, dst string) error {
		return exec.Command(path, src, dst).Run()
	}
}

func usage(err error) {
	fmt.Fprintf(os.Stderr, "Usage: %s=xxx %s FILENAME\n", env, os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %s - %s\n", version, checksum)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}
