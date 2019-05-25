package collection_test

import (
	"reflect"
	"samplesort/collection"
	"testing"
)

func TestCollection(t *testing.T) {
	t.Run("it should compute features", func(t *testing.T) {
		SUT, _ := makeSUT()
		expected := [][]float64{{2.0}, {4.0}}
		actual := SUT.Features()
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
	t.Run("it should trigger updates", func(t *testing.T) {
		_, e := makeSUT()
		expected := map[string]float64{"foo": 1.0, "bar": 2.0}
		actual := e.updates
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
	t.Run("it should sort entities", func(t *testing.T) {
		SUT, _ := makeSUT()
		SUT.Sort([]int{1, 0})
		expected := "barfoo"
		actual := SUT.String()
		if expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func makeSUT() (*collection.Collection, *testEngine) {
	e := new(testEngine)
	e.updates = make(map[string]float64)
	SUT := collection.New(e)
	SUT.Append(testEntity{"foo", 1.0})
	SUT.Append(testEntity{"bar", 2.0})
	return SUT, e
}

type testEntity struct {
	key string
	val float64
}

func (t testEntity) Data() map[string]float64 {
	return map[string]float64{t.key: t.val}
}

func (t testEntity) String() string {
	return t.key
}

type testEngine struct {
	updates map[string]float64
}

func (t *testEngine) Update(key string, val float64) {
	t.updates[key] = val
}

func (t *testEngine) Normalize(key string, n float64) float64 {
	return n * 2
}
