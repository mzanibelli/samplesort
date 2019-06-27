package analyze_test

import (
	"errors"
	"samplesort/analyze"
	"testing"
)

func TestItShouldNotFailIfEverythingWorks(t *testing.T) {
	SUT := analyze.MakeSUT(
		[][]float64{{1, 2, 3}, {4, 5, 6}},
		[]int{0, 1},
		nil,
	)
	if err := SUT.Analyze(); err != nil {
		t.Error("should not fail")
	}
}

func TestItShouldNotFailIfNoDataIsFound(t *testing.T) {
	SUT := analyze.MakeSUT(
		[][]float64{},
		[]int{0, 1},
		nil,
	)
	if err := SUT.Analyze(); err != nil {
		t.Error("should not fail")
	}
}

func TestItShouldFailIfCacheFetchFails(t *testing.T) {
	SUT := analyze.MakeSUT(
		[][]float64{{1, 2, 3}, {4, 5, 6}},
		[]int{0, 1},
		errors.New("foo"),
	)
	if err := SUT.Analyze(); err == nil {
		t.Error("should fail")
	}
}
