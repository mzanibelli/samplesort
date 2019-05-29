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
	res, err := samplesort.SampleSort(
		os.Args[1],
		os.Getenv(env),
		logger,
	)
	usage(err)
	fmt.Fprintln(os.Stdout, res)
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
