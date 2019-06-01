package cache_test

import (
	"errors"
	"samplesort/cache"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("it should return cached data if found",
		func(t *testing.T) {
			fs := &mockFS{[]byte(`{"data":{"foo":"bar"}}`), nil, true, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte{}, nil
			}
			SUT := cache.New(fs, ".json")
			SUT.Fetch("foo", target, build)
			expected := "bar"
			actual := target.Data["foo"]
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should build data if not found",
		func(t *testing.T) {
			fs := &mockFS{[]byte(``), nil, false, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte(`{"data":{"foo":"bar"}}`), nil
			}
			SUT := cache.New(fs, ".json")
			SUT.Fetch("foo", target, build)
			expected := "bar"
			actual := target.Data["foo"]
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should write data if not found",
		func(t *testing.T) {
			fs := &mockFS{[]byte(``), nil, false, 0}
			target := &mockData{map[string]string{}}
			foo := []byte(`{"data":{"foo":"bar"}}`)
			build := func() ([]byte, error) {
				return foo, nil
			}
			SUT := cache.New(fs, ".json")
			SUT.Fetch("foo", target, build)
			expected := len(foo)
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not (over)write data if cached data is found",
		func(t *testing.T) {
			fs := &mockFS{[]byte(`{"data":{"foo":"bar"}}`), nil, true, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte{}, errors.New("foo")
			}
			SUT := cache.New(fs, ".json")
			SUT.Fetch("foo", target, build)
			expected := 0
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not write data if there was an error during build",
		func(t *testing.T) {
			fs := &mockFS{[]byte{}, nil, false, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte(`{"data":{"foo":"bar"}}`), errors.New("foo")
			}
			SUT := cache.New(fs, ".json")
			SUT.Fetch("foo", target, build)
			expected := 0
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not write data if there was an error during decoding",
		func(t *testing.T) {
			fs := &mockFS{[]byte{}, nil, false, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte{}, nil
			}
			SUT := cache.New(fs, ".json")
			SUT.Fetch("foo", target, build)
			expected := 0
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if we couldn't read the cached data",
		func(t *testing.T) {
			fs := &mockFS{[]byte{}, errors.New("foo"), true, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte{}, nil
			}
			SUT := cache.New(fs, ".json")
			err := SUT.Fetch("foo", target, build)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if we couldn't build data",
		func(t *testing.T) {
			fs := &mockFS{[]byte{}, nil, false, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte{}, errors.New("foo")
			}
			SUT := cache.New(fs, ".json")
			err := SUT.Fetch("foo", target, build)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if decoding failed",
		func(t *testing.T) {
			fs := &mockFS{[]byte{}, nil, true, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte{}, nil
			}
			SUT := cache.New(fs, ".json")
			err := SUT.Fetch("foo", target, build)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if writing failed",
		func(t *testing.T) {
			fs := &mockFS{[]byte{}, errors.New("foo"), false, 0}
			target := &mockData{map[string]string{}}
			build := func() ([]byte, error) {
				return []byte(`{"data":{"foo":"bar"}}`), nil
			}
			SUT := cache.New(fs, ".json")
			err := SUT.Fetch("foo", target, build)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
}

type mockFS struct {
	content []byte
	err     error
	exists  bool
	written int
}

func (m *mockFS) Exists(name string) bool             { return m.exists }
func (m *mockFS) ReadAll(name string) ([]byte, error) { return m.content, m.err }
func (m *mockFS) WriteAll(name string, data []byte) error {
	m.written = len(data)
	return m.err
}

type mockData struct {
	Data map[string]string `json:"data"`
}
