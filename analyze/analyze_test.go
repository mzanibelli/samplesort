package analyze_test

import (
	"errors"
	"samplesort/analyze"
	"testing"
)

func TestItShouldNotFailIfEverythingWorks(t *testing.T) {
	analyze.TestOptionFeatures = [][]float64{{1, 2, 3}, {4, 5, 6}}
	analyze.TestOptionCenters = []int{0, 1}
	analyze.TestOptionFetchError = nil
	SUT := analyze.MakeSUT()
	if err := SUT.Analyze(); err != nil {
		t.Error("should not fail")
	}
}

func TestItShouldNotFailIfNoDataIsFound(t *testing.T) {
	analyze.TestOptionFeatures = [][]float64{}
	analyze.TestOptionCenters = []int{0, 1}
	analyze.TestOptionFetchError = nil
	SUT := analyze.MakeSUT()
	if err := SUT.Analyze(); err != nil {
		t.Error("should not fail")
	}
}

func TestItShouldFailIfCacheFetchFails(t *testing.T) {
	analyze.TestOptionFeatures = [][]float64{{1, 2, 3}, {4, 5, 6}}
	analyze.TestOptionCenters = []int{0, 1}
	analyze.TestOptionFetchError = errors.New("foo")
	SUT := analyze.MakeSUT()
	if err := SUT.Analyze(); err == nil {
		t.Error("should fail")
	}
}
