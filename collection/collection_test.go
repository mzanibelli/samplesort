package collection_test

import (
	"reflect"
	"samplesort/collection"
	"testing"
)

func TestCollection(t *testing.T) {
	t.Run("it should compute features", func(t *testing.T) {
		SUT := makeSUT()
		expected := [][]float64{{1}, {2}}
		actual := SUT.Features()
		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
	t.Run("it should sort entities", func(t *testing.T) {
		SUT := makeSUT()
		SUT.Sort([]int{1, 0})
		expected := "barfoo"
		actual := SUT.String()
		if expected != actual {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}

func makeSUT() *collection.Collection {
	SUT := collection.New()
	SUT.Append(testEntity{"foo", 1.0})
	SUT.Append(testEntity{"bar", 2.0})
	return SUT
}

type testEntity struct {
	key string
	val float64
}

func (t testEntity) Values() []float64 { return []float64{t.val} }
func (t testEntity) String() string    { return t.key }
