package samplesort_test

import (
	"bytes"
	"path/filepath"
	"samplesort"
	"strings"
	"testing"
)

func TestOutputWithSingleSample(t *testing.T) {
	root := "./testdata/single"
	output := bytes.NewBuffer([]byte{})
	samplesort.SampleSort(
		"./bin/streaming_extractor_music",
		output,
		samplesort.WithFileSystemRoot(root),
	)
	expected := strings.Join([]string{
		filepath.Join(root, "sample.wav"),
	}, "\n")
	actual := strings.Trim(output.String(), "\n")
	if expected != actual {
		t.Errorf("\n-> expected:\n%s\n-> actual:\n%s", expected, actual)
	}
}

func TestSameSamplesShouldBeSideBySide(t *testing.T) {
	t.Skip("not ready yet")
	root := "./testdata/duplicates"
	output := bytes.NewBuffer([]byte{})
	samplesort.SampleSort(
		"./bin/streaming_extractor_music",
		output,
		samplesort.WithFileSystemRoot(root),
		samplesort.WithSize(2),
	)
	t.Error()
	expected := strings.Join([]string{
		filepath.Join(root, "a.wav"),
		filepath.Join(root, "b.wav"),
		filepath.Join(root, "c.wav"),
		filepath.Join(root, "d.wav"),
		filepath.Join(root, "e.wav"),
		filepath.Join(root, "f.wav"),
	}, "\n")
	actual := strings.Trim(output.String(), "\n")
	if expected != actual {
		t.Errorf("\n-> expected:\n%s\n-> actual:\n%s", expected, actual)
	}
}
