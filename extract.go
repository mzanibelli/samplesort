package samplesort

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	EXPECTED_SHA256  string = "9c91599c118ad0f2eef14e7bbcc050d8c802d3175b8e1766c820c7ab5ce685f5"
	EXPECTED_VERSION string = "v2.1_beta2-linux-i686"
	ENV_EXTRACTOR    string = "ESSENTIA_EXTRACTOR"
	EXT_IN           string = ".wav"
	EXT_OUT          string = ".json"
)

func SampleSort() {
	bin, err := which()
	if err != nil {
		usage(err)
	}
	err = filepath.Walk(os.Args[1], extract(bin))
	if err != nil {
		log.Fatal(err)
	}
}

func usage(err error) {
	fmt.Fprintf(os.Stderr, "Usage: %s=xxx %s FILENAME\n", ENV_EXTRACTOR, os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %s - %s\n", EXPECTED_VERSION, EXPECTED_SHA256)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}

func extract(bin string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		out := path + EXT_OUT
		switch {
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case filepath.Ext(path) != EXT_IN:
			return nil
		case exists(out):
			return nil
		default:
			return run(bin, path, out)
		}
	}
}

func run(bin, src, dst string) error {
	cmd := exec.Command(bin, src, dst)
	if err := cmd.Run(); err != nil {
		log.Println("ERR:", dst)
	} else {
		log.Println("NEW:", dst)
	}
	return nil
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
	default:
		return extractor, nil
	}
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
