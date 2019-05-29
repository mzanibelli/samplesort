package extractor_test

import (
	"errors"
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
				func(content []byte, data []map[string]interface{}) error {
					return nil
				},
				"",
			)
			go SUT.Extract("")
			<-SUT.Out()
			if !hasRun {
				t.Errorf("runner was not invoked")
			}
		})
	t.Run("it should not execute the loader if there was an error with the runner",
		func(t *testing.T) {
			hasRun := false
			SUT := extractor.New(
				&mockFS{[]byte{}, nil, false},
				func(src, dst string) error {
					return errors.New("foo")
				},
				func(content []byte, data []map[string]interface{}) error {
					hasRun = true
					return nil
				},
				"",
			)
			go SUT.Extract("")
			<-SUT.Err()
			if hasRun {
				t.Errorf("loader was wrongfully invoked")
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
