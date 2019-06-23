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
	expected := "xxx" // 3 similar files with name starting with 'x'
	actual := baseline(output)
	if !strings.Contains(actual, expected) {
		t.Errorf("%q does not contain %q", actual, expected)
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
