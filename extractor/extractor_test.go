package extractor_test

import (
	"errors"
	"samplesort/extractor"
	"testing"
)

func TestExtractor(t *testing.T) {
	t.Run("it should pass src to build func",
		func(t *testing.T) {
			ok := false
			SUT := extractor.New(
				&mockCache{nil},
				func(src string) (map[string]interface{}, error) {
					ok = true
					if src != "hello" {
						t.Error("path not provided to build func")
					}
					return nil, nil
				},
			)
			go SUT.Extract("hello")
			select {
			case <-SUT.Out():
				break
			case <-SUT.Err():
				t.Error("received error instead of output")
				break
			}
			if !ok {
				t.Error("build func was not called")
			}
		})
	t.Run("it should send an error if fetch fails",
		func(t *testing.T) {
			SUT := extractor.New(
				&mockCache{errors.New("foo")},
				func(string) (map[string]interface{}, error) { return nil, nil },
			)
			go SUT.Extract("hello")
			select {
			case <-SUT.Out():
				t.Error("received output instead of error")
				break
			case <-SUT.Err():
				break
			}
		})
	t.Run("it should send output if fetch works",
		func(t *testing.T) {
			SUT := extractor.New(
				&mockCache{nil},
				func(string) (map[string]interface{}, error) { return nil, nil },
			)
			go SUT.Extract("hello")
			select {
			case <-SUT.Out():
				break
			case <-SUT.Err():
				t.Error("received error instead of output")
				break
			}
		})
}

type mockCache struct{ err error }

func (m *mockCache) Fetch(
	key string,
	target interface{},
	build func() (interface{}, error),
) error {
	if m.err == nil {
		build()
	}
	return m.err
}
