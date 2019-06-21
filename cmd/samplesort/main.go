package main

import (
	"fmt"
	"log"
	"os"
	"samplesort"
)

const (
	env string = "ESSENTIA_EXTRACTOR"
)

func main() {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	_, err := samplesort.SampleSort(
		os.Getenv(env),
		samplesort.WithFileSystemRoot(os.Args[1]),
		samplesort.WithLoggers(logger),
	)
	usage(err)
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
