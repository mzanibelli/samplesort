package sample_test

import (
	"reflect"
	"samplesort/sample"
	"testing"
)

func TestSample(t *testing.T) {
	t.Run("it should flatten data",
		func(t *testing.T) {
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
		})
	t.Run("it should extract floats",
		func(t *testing.T) {
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
		})
	t.Run("it should sort features by key",
		func(t *testing.T) {
			input := map[string]interface{}{
				"foo": map[string]interface{}{
					"zzz": 0.001,
					"aaa": 0.05641,
				},
			}
			expected := []float64{0.05641, 0.001}
			SUT := sample.New("")
			SUT.Flatten(input)
			actual := SUT.Values()
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return the value of the given key",
		func(t *testing.T) {
			input := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 0.001,
				},
			}
			expected := 0.001
			SUT := sample.New("")
			SUT.Flatten(input)
			actual := SUT.Get("foo.bar")
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return the list of keys",
		func(t *testing.T) {
			input := map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": 0.001,
				},
			}
			expected := []string{"foo.bar"}
			SUT := sample.New("")
			SUT.Flatten(input)
			actual := SUT.Keys()
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
}
