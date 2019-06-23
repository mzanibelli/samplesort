package samplesort_test

import (
	"bytes"
	"path/filepath"
	"samplesort"
	"strings"
	"testing"
)

func TestOutputWithSingleSample(t *testing.T) {
	if testing.Short() {
		return
	}
	root := "./testdata/single"
	output := bytes.NewBuffer([]byte{})
	s := samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot(root),
		samplesort.WithoutCache(),
	)
	s.WriteTo(output)
	expected := strings.Join([]string{
		filepath.Join(root, "sample.wav"),
	}, "\n")
	actual := strings.Trim(output.String(), "\n")
	if expected != actual {
		t.Errorf("\n-> expected:\n%s\n-> actual:\n%s", expected, actual)
	}
}

func TestSameSamplesShouldBeSideBySide(t *testing.T) {
	if testing.Short() {
		return
	}
	root := "./testdata/duplicates"
	output := bytes.NewBuffer([]byte{})
	s := samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot(root),
		samplesort.WithSize(4),
		samplesort.WithoutCache(),
	)
	s.WriteTo(output)
	expected := strings.Join([]string{
		filepath.Join(root, "f.wav"),
		filepath.Join(root, "b.wav"),
		filepath.Join(root, "d.wav"),
		filepath.Join(root, "a.wav"),
		filepath.Join(root, "c.wav"),
		filepath.Join(root, "e.wav"),
	}, "\n")
	actual := strings.Trim(output.String(), "\n")
	if expected != actual {
		t.Errorf("\n-> expected:\n%s\n-> actual:\n%s", expected, actual)
	}
}
