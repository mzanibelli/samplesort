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

func TestConfigIsCorrectlyApplied(t *testing.T) {
	output := bytes.NewBuffer([]byte{})
	samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot("foo"),
		samplesort.WithAudioFormat("bar"),
		samplesort.WithDataFormat("baz"),
		samplesort.WithSize(1337),
		samplesort.WithMaxIterations(42),
		samplesort.WithMaxZScore(777),
		samplesort.WithoutCache(),
	).DumpConfig(output)
	expected, actual := formatForAssertion(output,
		"data: foo",
		"input: bar",
		"output: baz",
		"size: 1337",
		"threshold: 42",
		"zscore: 777.00",
		"cache: false",
	)
	if expected != actual {
		t.Errorf("\n-> expected:\n%s\n-> actual:\n%s", expected, actual)
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
	expected, actual := formatForAssertion(output,
		filepath.Join(root, "sample.wav"),
	)
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

func TestTheResultsAreConsistentWithHumanEar(t *testing.T) {
	t.Skip("TODO: establish well known comparisons between two similar and one very dissimilar sample")
}

// The first letter for each file name of a given directory allows quick
// comparison of the order.
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

func formatForAssertion(s fmt.Stringer, lines ...string) (expected, actual string) {
	actual = strings.Trim(s.String(), "\n")
	expected = strings.Join(lines, "\n")
	return expected, actual
}
