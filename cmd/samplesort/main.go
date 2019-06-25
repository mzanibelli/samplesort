package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"samplesort"
)

const (
	env string = "ESSENTIA_EXTRACTOR"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	executable := os.Getenv(env)
	usage(check(executable))
	s := samplesort.New(
		executable,
		samplesort.WithSize(4),
		samplesort.WithFileSystemRoot(os.Args[1]),
		samplesort.WithLogger(logger),
	)
	if _, err := s.WriteTo(os.Stdout); err != nil {
		logger.Fatal(err)
	}
}

func check(path string) error {
	fd, err := os.Open(path)
	defer fd.Close()
	switch {
	case err != nil:
		return fmt.Errorf("Error opening executable: %v", err)
	case !sha256Verify(fd, samplesort.Checksum):
		return fmt.Errorf("SHA256 mismatch, expected version %q", samplesort.Version)
	case len(os.Args) != 2:
		return fmt.Errorf("Expected exactly one argument, got: %d", len(os.Args)-1)
	}
	return nil
}

func usage(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "Usage: %s=xxx %s FILENAME\n", env, os.Args[0])
	fmt.Fprintf(os.Stderr, "Version: %s - %s\n", samplesort.Version, samplesort.Checksum)
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}

func sha256Verify(r io.Reader, sum string) bool {
	h := sha256.New()
	n, err := io.Copy(h, r)
	switch {
	case err != nil:
		return false
	case n == 0:
		return false
	default:
		return sum == fmt.Sprintf("%x", h.Sum(nil))
	}
}
