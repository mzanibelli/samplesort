package samplesort_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"samplesort"
	"sort"
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
		samplesort.WithoutCache(),
	).DumpConfig(output)
	expected, actual := formatForAssertion(output,
		"data: foo",
		"input: bar",
		"output: baz",
		"size: 1337",
		"threshold: 42",
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

func TestFeaturesShouldBeTheRawSortedData(t *testing.T) {
	root := "./testdata/consistency"
	output := bytes.NewBuffer([]byte{})
	samplesort.New(
		"./bin/streaming_extractor_music",
		samplesort.WithFileSystemRoot(root),
		samplesort.WithoutCache(),
	).WriteTo(output)
	flat := make(map[string]float64)
	fromJSON("./testdata/consistency/flat.json", &flat)
	expected := sortByKey(flat)
	actual := make([][]float64, 0)
	fromJSON(filepath.Join(root, "features.json"), &actual)
	if !reflect.DeepEqual(expected, actual) {
		for i := range expected {
			for j := range expected[i] {
				if expected[i][j] == actual[i][j] {
					continue
				}
				t.Log("at:", i, j)
				t.Log("expected:", expected[i][j])
				t.Log("actual:", actual[i][j])
				break
			}
		}
		t.Error("features and raw data are different")
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

func sortByKey(input map[string]float64) [][]float64 {
	keys := make([]string, 0, len(input))
	for key := range input {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	features := make([]float64, len(keys), len(keys))
	for i, key := range keys {
		features[i] = input[key]
	}
	res := make([][]float64, 1, 1)
	res[0] = features
	return res
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

func formatForAssertion(console fmt.Stringer, lines ...string) (expected, actual string) {
	actual = strings.Trim(console.String(), "\n")
	expected = strings.Join(lines, "\n")
	return expected, actual
}

func fromJSON(path string, target interface{}) {
	fd, err := os.Open(path)
	must(err)
	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	must(err)
	must(json.Unmarshal(content, target))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
