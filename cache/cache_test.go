package cache_test

import (
	"encoding/json"
	"errors"
	"samplesort/cache"
	"testing"
)

func TestItShouldReturnCachedDataIfFound(t *testing.T) {
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
}

func TestItShouldBuildDataIfNotFound(t *testing.T) {
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
}

func TestItShouldWriteDataIfNotFound(t *testing.T) {
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
}

func TestItShouldNotOverwriteCachedData(t *testing.T) {
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
}

func TestItShouldNotWriteIfBuildFails(t *testing.T) {
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
}

func TestItShouldNotWriteIfEncodingFails(t *testing.T) {
	fs := mkFs("", nil, false)
	target := mkData()
	build := func() (interface{}, error) {
		return mkData("malformed"), nil
	}
	SUT := cache.New(fs, defaultConfig())
	err := SUT.Fetch("foo", target, build)
	t.Log(err)
	expected := 0
	actual := fs.written
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldReturnAnErrorIfReadFails(t *testing.T) {
	fs := mkFs("", errors.New("foo"), true)
	target := mkData()
	SUT := cache.New(fs, defaultConfig())
	err := SUT.Fetch("foo", target, nop)
	expected := true
	actual := err != nil
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldReturnAnErrorIfBuildFails(t *testing.T) {
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
}

func TestItShouldReturnAnErrorIfDecodingFails(t *testing.T) {
	fs := mkFs("", nil, true)
	target := mkData()
	SUT := cache.New(fs, defaultConfig())
	err := SUT.Fetch("foo", target, nop)
	expected := true
	actual := err != nil
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldReturnAnErrorIfWriteFails(t *testing.T) {
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
}

func TestItShouldReturnAnErrorIfPathCannotBeMadeRelative(t *testing.T) {
	fs := mkFs("", nil, false)
	SUT := cache.New(fs, withRoot("/"))
	err := SUT.Fetch("bar/baz", nil, nop)
	expected := true
	actual := err != nil
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldBuildDataIfCacheIsDisabled(t *testing.T) {
	fs := mkFs("", nil, true)
	target := mkData()
	build := func() (interface{}, error) {
		return mkData("foo", "bar"), nil
	}
	SUT := cache.New(fs, noCache())
	err := SUT.Fetch("foo", target, build)
	t.Log(err)
	expected := "bar"
	actual := target.Data["foo"]
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestItShouldAllowTheBuildCommandToHandleWrite(t *testing.T) {
	fs := mkFs("", nil, false)
	target := mkData()
	build := func() (interface{}, error) {
		fs.content = []byte(`{"data":{"foo":"bar"}}`)
		fs.exists = true
		return nil, nil
	}
	SUT := cache.New(fs, defaultConfig())
	err := SUT.Fetch("foo", target, build)
	t.Log(err)
	expected := "bar"
	actual := target.Data["foo"]
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
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
		{
			name:     "it should not go back to current directory",
			root:     "./testdata/duplicates",
			key:      "kmeans",
			expected: "testdata/duplicates/kmeans.json",
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
		return &mockData{map[string]string{}, true}
	case len(vs) == 1:
		return &mockData{nil, false}
	case len(vs)%2 != 0:
		panic("nope")
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
	return &mockData{res, true}
}

type mockData struct {
	Data  map[string]string
	valid bool
}

func (m *mockData) MarshalJSON() ([]byte, error) {
	if !m.valid {
		return []byte("{aa"), errors.New("foo")
	}
	res := make(map[string]interface{})
	res["data"] = m.Data
	return json.Marshal(&res)
}

func (m *mockData) UnmarshalJSON(data []byte) error {
	if !m.valid {
		return errors.New("malformed")
	}
	res := make(map[string]map[string]string)
	err := json.Unmarshal(data, &res)
	if value, ok := res["data"]; ok {
		m.Data = value
	}
	return err
}

func defaultConfig() *mockConfig {
	return &mockConfig{".", ".json", true}
}

func noCache() *mockConfig {
	return &mockConfig{".", ".json", false}
}

func withRoot(root string) *mockConfig {
	return &mockConfig{root, ".json", true}
}

type mockConfig struct {
	root      string
	extension string
	enable    bool
}

func (m *mockConfig) FileSystemRoot() string { return m.root }
func (m *mockConfig) DataFormat() string     { return m.extension }
func (m *mockConfig) EnableCache() bool      { return m.enable }

var nop = func() (interface{}, error) { return nil, nil }
