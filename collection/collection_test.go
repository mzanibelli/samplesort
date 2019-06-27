package collection_test

import (
	"reflect"
	"samplesort/collection"
	"testing"
)

func TestItShouldComputeFeatures(t *testing.T) {
	SUT := makeSUT()
	expected := [][]float64{{1}, {2}}
	actual := SUT.Features()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldSortEntities(t *testing.T) {
	SUT := makeSUT()
	SUT.Append(testEntity{"alice", []string{"foo"}, []float64{3.0}})
	SUT.Append(testEntity{"bob", []string{"foo"}, []float64{4.0}})
	SUT.Sort([]int{2, 0, 1, 3})
	expected := "doealicejohnbob"
	actual := SUT.String()
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldDisregardUncommonKeys(t *testing.T) {
	SUT := makeSUT()
	SUT.Append(testEntity{"alice", []string{"foo", "baz"}, []float64{3.0, 12.65}})
	SUT.Append(testEntity{"bob", []string{"qux", "foo"}, []float64{66.51, 4.0}})
	expected := [][]float64{{1}, {2}, {3}, {4}}
	actual := SUT.Features()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldReturnNothingIfNoCommonKeysAreFound(t *testing.T) {
	SUT := makeSUT()
	SUT.Append(testEntity{"alice", []string{"foo", "baz"}, []float64{3.0, 12.65}})
	SUT.Append(testEntity{"bobby", []string{"qux"}, []float64{66.51}})
	expected := [][]float64{{}, {}, {}, {}}
	actual := SUT.Features()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldNotFailOnEmptyInput(t *testing.T) {
	SUT := collection.New()
	actual := SUT.Features()
	if actual != nil {
		t.Errorf("expected nil, got: %v", actual)
	}
}

func makeSUT() *collection.Collection {
	SUT := collection.New()
	SUT.Append(testEntity{"john", []string{"foo"}, []float64{1.0}})
	SUT.Append(testEntity{"doe", []string{"foo"}, []float64{2.0}})
	return SUT
}

type testEntity struct {
	name   string
	keys   []string
	values []float64
}

func (t testEntity) Keys() []string    { return t.keys }
func (t testEntity) Values() []float64 { return t.values }
func (t testEntity) String() string    { return t.name }
