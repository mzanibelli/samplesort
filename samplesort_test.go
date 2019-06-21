package samplesort_test

import (
	"samplesort"
	"testing"
)

func TestIntegration(t *testing.T) {
	samplesort.SampleSort(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot("./testdata"),
		samplesort.WithSize(2),
		samplesort.WithMaxIterations(1),
	)
}
