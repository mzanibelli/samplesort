package samplesort

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"samplesort/collection"
	"samplesort/engine"
	"samplesort/sample"

	"github.com/bugra/kmeans"
)

const (
	EXPECTED_SHA256  string  = "9c91599c118ad0f2eef14e7bbcc050d8c802d3175b8e1766c820c7ab5ce685f5"
	EXPECTED_VERSION string  = "v2.1_beta2-linux-i686"
	ENV_EXTRACTOR    string  = "ESSENTIA_EXTRACTOR"
	EXT_IN           string  = ".wav"
	EXT_OUT          string  = ".json"
	PARAM_SIZE       int     = 100
	PARAM_THRESHOLD  int     = 0
	PARAM_PRECISION  float64 = 0.05
)

func SampleSort() {
	bin, err := which()
	if err != nil {
		usage(err)
	}
	sink := make(chan *sample.Sample)
	done := make(chan struct{})
	go analyze(sink, done)
	err = filepath.Walk(os.Args[1], extract(bin, loader(sink)))
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

func extract(bin string, load loadFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		out := path + EXT_OUT
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case filepath.Ext(path) != EXT_IN:
			return nil
		}
		if err := run(bin, path, out, load); err != nil {
			log.Println(err)
		}
		return nil
	}
}

func run(bin, src, dst string, load loadFunc) error {
	s := new(sample.Sample)
	s.Path = src
	if exists(dst) {
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
		fd, err := os.Open(dst)
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
	switch {
	case extractor == "":
		return "", fmt.Errorf("Please set %q environment variable", ENV_EXTRACTOR)
	case !exists(extractor):
		return "", fmt.Errorf("File %q not found", extractor)
	case !checksum(extractor):
		return "", fmt.Errorf("SHA256 mismatch, expected version %q", EXPECTED_VERSION)
	case len(os.Args) != 2:
		return "", fmt.Errorf("Expected exactly one argument, got: %d", len(os.Args))
	}
	return extractor, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checksum(path string) bool {
	fd, err := os.Open(path)
	if err != nil {
		return false
	}
	defer fd.Close()
	h := sha256.New()
	_, err = io.Copy(h, fd)
	if err != nil {
		return false
	}
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum == EXPECTED_SHA256
}

func analyze(input <-chan *sample.Sample, done chan<- struct{}) {
	defer close(done)
	coll := collection.New(engine.New(PARAM_PRECISION))
	for e := range input {
		coll.Append(e)
	}
	means, err := kmeans.Kmeans(coll.Features(), PARAM_SIZE,
		kmeans.HammingDistance, PARAM_THRESHOLD)
	if err != nil {
		log.Println("could not compute kmeans:", err)
	}
	coll.Sort(means)
	fmt.Fprintln(os.Stdout, coll)
}
