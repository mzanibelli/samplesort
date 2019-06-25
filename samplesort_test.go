package samplesort_test

import (
	"bytes"
	"fmt"
	"path/filepath"
	"samplesort"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMissingDirectoryShouldProduceError(t *testing.T) {
	root := "./testdata/missing"
	output := bytes.NewBuffer([]byte{})
	n, err := samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot(root),
		samplesort.WithoutCache(),
	).WriteTo(output)
	if err == nil {
		t.Error("expected error")
	}
	expected := int64(0)
	actual := n
	if expected != actual {
		t.Errorf("expected:%d, actual:%d", expected, actual)
	}
}

func TestEmptyDirectoryShouldDoNothing(t *testing.T) {
	root := "./testdata/empty"
	output := bytes.NewBuffer([]byte{})
	n, err := samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot(root),
		samplesort.WithoutCache(),
	).WriteTo(output)
	if err != nil {
		t.Error(err)
	}
	expected := int64(0)
	actual := n
	if expected != actual {
		t.Errorf("expected:%d, actual:%d", expected, actual)
	}
}

func TestOutputWithSingleSample(t *testing.T) {
	if testing.Short() {
		return
	}
	root := "./testdata/single"
	output := bytes.NewBuffer([]byte{})
	samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot(root),
		samplesort.WithoutCache(),
	).WriteTo(output)
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
	output := bytes.NewBuffer([]byte{})
	samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot("./testdata/duplicates"),
		samplesort.WithSize(4),
		samplesort.WithoutCache(),
	).WriteTo(output)
	b := baseline(output)
	// Two 'x's and a 'z' are duplicates: guard against alphabetical
	// sorting induced luck.
	expected := true
	actual := strings.Contains(b, "xxz") ||
		strings.Contains(b, "xzx") ||
		strings.Contains(b, "zxx")
	if expected != actual {
		t.Errorf("duplicates are not side by side: %s", b)
	}
}

// input:
// /tmp/a.wav
// /tmp/b.wav
// /tmp/c.wav
// output:
// abc
func baseline(s fmt.Stringer) string {
	res := new(strings.Builder)
	parts := strings.Split(s.String(), "\n")
	for _, part := range parts {
		if part == "" {
			continue
		}
		base := filepath.Base(part)
		r, _ := utf8.DecodeRuneInString(base)
		res.WriteRune(r)
	}
	return res.String()
}
