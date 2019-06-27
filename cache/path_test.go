package cache_test

import (
	"errors"
	"testing"
)

func TestPath(t *testing.T) {
	cases := []struct {
		name     string
		root     string
		key      string
		format   string
		expected string
		err      error
	}{
		{
			name:     "it should use the absolute key path if possible",
			root:     "/tmp",
			key:      "/foo.wav",
			format:   ".json",
			expected: "/foo.wav.json",
			err:      nil,
		},
		{
			name:     "it should make the final path relative to the root",
			root:     "/tmp",
			key:      "foo.wav",
			format:   ".json",
			expected: "/tmp/foo.wav.json",
			err:      nil,
		},
		{
			name:     "it should dedupe common path components",
			root:     "./data",
			key:      "./data/foo.wav",
			format:   ".json",
			expected: "data/foo.wav.json",
			err:      nil,
		},
		{
			name:     "it should not go back to current directory",
			root:     "./testdata/duplicates",
			key:      "kmeans",
			format:   ".json",
			expected: "testdata/duplicates/kmeans.json",
			err:      nil,
		},
		{
			name:     "it should return an error if path cannot be made relative",
			root:     "/",
			key:      "bar/baz",
			format:   ".json",
			expected: "testdata/duplicates/kmeans.json",
			err:      errors.New("foo"),
		},
	}
	for _, c := range cases {
		t.Run(c.name,
			func(t *testing.T) {
				actual, err := cache.Path(c.root, c.key, c.format)
				if c.err != err {
					t.Log("got:", err)
					t.Fatal("error assertion failed")
				}
				if c.expected != actual {
					t.Errorf("expected: %s, actual: %s", expected, actual)
				}
			})
	}
}
