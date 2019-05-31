package extractor_test

import (
	"errors"
	"reflect"
	"samplesort/extractor"
	"testing"
)

func TestExtractor(t *testing.T) {
	t.Run("it should execute the runner if no data file is found",
		func(t *testing.T) {
			hasRun := false
			SUT := extractor.New(
				&mockFS{[]byte{}, nil, false},
				func(src, dst string) error {
					hasRun = true
					return nil
				},
				func(content []byte) ([]map[string]interface{}, error) {
					return []map[string]interface{}{}, nil
				},
				"",
			)
			go SUT.Extract("")
			select {
			case <-SUT.Out():
				if !hasRun {
					t.Error("runner was not invoked")
				}
				break
			case err := <-SUT.Err():
				t.Error(err)
				break
			}
		})
	t.Run("it should not execute the loader if there was an error with the runner",
		func(t *testing.T) {
			SUT := extractor.New(
				&mockFS{[]byte{}, nil, false},
				func(src, dst string) error {
					return errors.New("foo")
				},
				func(content []byte) ([]map[string]interface{}, error) {
					t.Error("loader was wrongfully invoked")
					return []map[string]interface{}{}, nil
				},
				"",
			)
			go SUT.Extract("")
			select {
			case <-SUT.Out():
				t.Error("received output instead of error")
				break
			case <-SUT.Err():
				break
			}
		})
	t.Run("it should not execute the loader if there was an error reading the data file",
		func(t *testing.T) {
			SUT := extractor.New(
				&mockFS{[]byte{}, errors.New("foo"), true},
				func(src, dst string) error {
					return nil
				},
				func(content []byte) ([]map[string]interface{}, error) {
					t.Error("loader was wrongfully invoked")
					return []map[string]interface{}{}, nil
				},
				"",
			)
			go SUT.Extract("")
			select {
			case <-SUT.Out():
				t.Error("received output instead of error")
				break
			case <-SUT.Err():
				break
			}
		})
	t.Run("it should decode the content of the data file",
		func(t *testing.T) {
			expected := []byte("hello world")
			SUT := extractor.New(
				&mockFS{expected, nil, true},
				func(src, dst string) error {
					return nil
				},
				func(actual []byte) ([]map[string]interface{}, error) {
					if !reflect.DeepEqual(expected, actual) {
						t.Errorf("expected: %v, actual: %v", expected, actual)
					}
					return []map[string]interface{}{}, nil
				},
				"",
			)
			go SUT.Extract("")
			select {
			case <-SUT.Out():
				break
			case err := <-SUT.Err():
				t.Error(err)
				break
			}
		})
	t.Run("it should not send data if there was a decoding error",
		func(t *testing.T) {
			SUT := extractor.New(
				&mockFS{[]byte{}, nil, true},
				func(src, dst string) error {
					return nil
				},
				func(content []byte) ([]map[string]interface{}, error) {
					return []map[string]interface{}{}, errors.New("foo")
				},
				"",
			)
			go SUT.Extract("")
			select {
			case <-SUT.Out():
				t.Error("received output instead of error")
				break
			case <-SUT.Err():
				break
			}
		})
}

type mockFS struct {
	content []byte
	err     error
	exists  bool
}

func (m *mockFS) ReadAll(name string) ([]byte, error) { return m.content, m.err }
func (m *mockFS) Exists(name string) bool             { return m.exists }
