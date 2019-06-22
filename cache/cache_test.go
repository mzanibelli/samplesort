package cache_test

import (
	"encoding/json"
	"errors"
	"samplesort/cache"
	"testing"
)

func TestCache(t *testing.T) {
	t.Run("it should return cached data if found",
		func(t *testing.T) {
			fs := mkFs(`{"data":{"foo":"bar"}}`, nil, true)
			target := mkData()
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, nop)
			t.Log(err)
			expected := "bar"
			actual := target.Data["foo"]
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should build data if not found",
		func(t *testing.T) {
			fs := mkFs("", nil, false)
			target := mkData()
			build := func() (interface{}, error) {
				return mkData("foo", "bar"), nil
			}
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, build)
			t.Log(err)
			expected := "bar"
			actual := target.Data["foo"]
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should write data if not found",
		func(t *testing.T) {
			fs := mkFs("", nil, false)
			target := mkData()
			build := func() (interface{}, error) {
				return mkData("foo", "bar"), nil
			}
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, build)
			t.Log(err)
			b, _ := json.Marshal(target)
			expected := len(b)
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not (over)write data if cached data is found",
		func(t *testing.T) {
			fs := mkFs(`{"data":{"foo":"bar"}}`, nil, true)
			target := mkData()
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, nop)
			t.Log(err)
			expected := 0
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not write data if there was an error during build",
		func(t *testing.T) {
			fs := mkFs("", nil, false)
			target := mkData()
			build := func() (interface{}, error) {
				return mkData("foo", "bar"), errors.New("foo")
			}
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, build)
			t.Log(err)
			expected := 0
			actual := fs.written
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should not write data if there was an error during encoding",
		func(t *testing.T) {
			t.Skip("need to make mockData produce invalid json")
		})
	t.Run("it should return an error if we couldn't read the cached data",
		func(t *testing.T) {
			fs := mkFs("", errors.New("foo"), true)
			target := mkData()
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, nop)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if we couldn't build data",
		func(t *testing.T) {
			fs := mkFs("", nil, false)
			target := mkData()
			build := func() (interface{}, error) {
				return mkData(), errors.New("foo")
			}
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, build)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if decoding failed",
		func(t *testing.T) {
			fs := mkFs("", nil, true)
			target := mkData()
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, nop)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error if writing failed",
		func(t *testing.T) {
			fs := mkFs("", errors.New("foo"), false)
			target := mkData()
			build := func() (interface{}, error) {
				return mkData("foo", "bar"), nil
			}
			SUT := cache.New(fs, defaultConfig())
			err := SUT.Fetch("foo", target, build)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
	t.Run("it should return an error key and root cannot be made relative",
		func(t *testing.T) {
			fs := mkFs("", nil, false)
			SUT := cache.New(fs, withRoot("/"))
			err := SUT.Fetch("bar/baz", nil, nop)
			expected := true
			actual := err != nil
			if expected != actual {
				t.Errorf("expected: %v, actual: %v", expected, actual)
			}
		})
}

func TestPath(t *testing.T) {
	cases := []struct{ name, root, key, expected string }{
		{
			name:     "it should use the absolute key path if possible",
			root:     "/tmp",
			key:      "/foo.wav",
			expected: "/foo.wav.json",
		},
		{
			name:     "it should make the final path relative to the root",
			root:     "/tmp",
			key:      "foo.wav",
			expected: "/tmp/foo.wav.json",
		},
		{
			name:     "it should dedupe common path components",
			root:     "./data",
			key:      "./data/foo.wav",
			expected: "data/foo.wav.json",
		},
	}
	for _, c := range cases {
		t.Run(c.name,
			func(t *testing.T) {
				fs := pathChecker(func(actual string) {
					expected := c.expected
					if expected != actual {
						t.Errorf("expected: %s, actual: %s", expected, actual)
					}
				})
				SUT := cache.New(fs, withRoot(c.root))
				SUT.Fetch(c.key, nil, nop)
			})
	}
}

func mkFs(data string, err error, exists bool) *mockFS {
	return &mockFS{[]byte(data), err, exists, 0, func(name string) {}}
}

func pathChecker(f func(name string)) *mockFS {
	return &mockFS{nil, nil, false, 0, f}
}

type mockFS struct {
	content   []byte
	err       error
	exists    bool
	written   int
	checkPath func(name string)
}

func (m *mockFS) ReadAll(name string) ([]byte, error) { return m.content, m.err }

func (m *mockFS) Exists(name string) bool {
	m.checkPath(name)
	return m.exists
}

func (m *mockFS) WriteAll(name string, data []byte) error {
	if m.err == nil {
		m.written = len(data)
		m.content = data
		m.exists = true
	}
	return m.err
}

func mkData(vs ...string) *mockData {
	switch {
	case len(vs) == 0:
		return &mockData{map[string]string{}}
	case len(vs)%2 != 0:
		panic("use key-value pairs to build data mock")
	}
	res := make(map[string]string, len(vs)/2)
	var key string
	for i, v := range vs {
		if i%2 == 0 {
			key = v
		} else {
			res[key] = v
		}
	}
	return &mockData{res}
}

type mockData struct {
	Data map[string]string `json:"data"`
}

func defaultConfig() *mockConfig {
	return &mockConfig{".", ".json"}
}

func withRoot(root string) *mockConfig {
	return &mockConfig{root, ".json"}
}

type mockConfig struct {
	root      string
	extension string
}

func (m *mockConfig) FileSystemRoot() string { return m.root }
func (m *mockConfig) DataFormat() string     { return m.extension }

var nop = func() (interface{}, error) { return nil, nil }
