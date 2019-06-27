package sample_test

import (
	"reflect"
	"samplesort/sample"
	"testing"
)

func TestItShouldFlattenData(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": 0.001,
		},
	}
	expected := []float64{0.001}
	SUT := sample.New("")
	SUT.Flatten(input)
	actual := SUT.Values()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldExtractFloats(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": 0.001,
			"baz": "hello",
		},
	}
	expected := []float64{0.001}
	SUT := sample.New("")
	SUT.Flatten(input)
	actual := SUT.Values()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldSortValuesByKey(t *testing.T) {
	input := map[string]interface{}{
		"zzz": 4.,
		"foo": map[string]interface{}{
			"bbb": 2.,
			"zzz": 3.,
			"aaa": 1.,
		},
	}
	expected := []float64{1, 2, 3, 4}
	SUT := sample.New("")
	SUT.Flatten(input)
	actual := SUT.Values()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldSortKeysByKey(t *testing.T) {
	input := map[string]interface{}{
		"zzz": 4.,
		"foo": map[string]interface{}{
			"bbb": 2.,
			"zzz": 3.,
			"aaa": 1.,
		},
	}
	expected := []string{"foo.aaa", "foo.bbb", "foo.zzz", "zzz"}
	SUT := sample.New("")
	SUT.Flatten(input)
	actual := SUT.Keys()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}
