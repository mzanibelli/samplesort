package samplesort

import (
	"fmt"
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

func SampleSort(executable string, configs ...config) ([][]float64, error) {
	cfg := newConfig(configs...)

	bin, err := which(executable, cfg.DataFormat())
	if err != nil {
		return nil, err
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

	normalizedFeatures, err := ana.Analyze()
	if err != nil {
		return nil, err
	}

	fmt.Println(col)

	return normalizedFeatures, nil
}

func which(path, output string) (func(src string) ([]byte, error), error) {
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
	return func(src string) ([]byte, error) {
		dst := path + output
		err := exec.Command(path, src, dst).Run()
		if err != nil {
			return []byte{}, err
		}
		return fs.ReadAll(dst)
	}, nil
}

func nop(string) ([]byte, error) {
	return []byte{}, nil
}
