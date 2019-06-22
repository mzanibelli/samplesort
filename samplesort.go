package samplesort

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"samplesort/analyze"
	"samplesort/cache"
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
)

func SampleSort(executable string, output io.Writer, configs ...config) error {
	cfg := newConfig(configs...)

	bin, err := which(executable, cfg.DataFormat())
	if err != nil {
		return err
	}

	// DI happens here.
	cac := cache.New(fs, cfg)
	ext := extractor.New(cac, bin)
	par := parser.New(fs, ext, cfg)
	col := collection.New()
	eng := engine.New(cfg)
	ana := analyze.New(col, eng, cac, cfg)

	go par.Parse(cfg.FileSystemRoot())

	go func() {
		for err := range ext.Err() {
			cfg.Log(err)
		}
	}()

	wg := new(sync.WaitGroup)
	for e := range ext.Out() {
		wg.Add(1)
		copy := e
		go func() {
			defer wg.Done()
			s := sample.New(copy.String())
			s.Flatten(copy.Data())
			col.Append(s)
		}()
	}
	wg.Wait()

	err = ana.Analyze()
	if err != nil {
		return err
	}

	fmt.Fprintln(output, col)

	return nil
}

// TODO: is this the way to do things right?
// Ensure we use a compatible version of the binary and contains the
// coupling to the external command execution method.
// We do JSON decoding at this level because this is the one format we
// don't control and which might change over time according to the
// decisions of Essentia developers.
func which(path, output string) (func(src string) (interface{}, error), error) {
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

	return func(src string) (interface{}, error) {
		dst := path + output
		err := exec.Command(path, src, dst).Run()
		if err != nil {
			return nil, err
		}
		res := make(map[string]interface{})
		content, err := fs.ReadAll(dst)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(content, &res)
		return res, err
	}, nil
}

func nop(string) (interface{}, error) { return nil, nil }
