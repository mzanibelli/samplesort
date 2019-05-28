package parser_test

import (
	"errors"
	"os"
	"path/filepath"
	"samplesort/parser"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	t.Run("it should stop in case of previous error",
		func(t *testing.T) {
			fs := &mockFS{t, []fixture{
				{"foo", mkfile("foo"), errors.New("foo")},
				{"bar", mkfile("bar"), nil},
			}, 0}
			SUT := parser.New(fs, nil)
			SUT.Parse("")
			expected := 1
			actual := fs.count
			if expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}
		})
	t.Run("it should do nothing if the entry is a directory",
		func(t *testing.T) {
			fs := &mockFS{t, []fixture{
				{"foo", mkdir("foo"), nil},
				{"bar", mkfile("bar"), nil},
			}, 0}
			SUT := parser.New(fs, nil)
			SUT.Parse("")
			expected := 2
			actual := fs.count
			if expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}
		})
	t.Run("it should do nothing if the entry is not a supported input file",
		func(t *testing.T) {
			fs := &mockFS{t, []fixture{
				{"foo.txt", mkfile("foo.txt"), nil},
				{"bar.pdf", mkfile("bar.pdf"), nil},
			}, 0}
			SUT := parser.New(fs, nil)
			SUT.Parse("")
			expected := 2
			actual := fs.count
			if expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}
		})
	t.Run("it should run the extractor if the entry is a supported input file",
		func(t *testing.T) {
			fs := &mockFS{t, []fixture{
				{"foo.txt", mkfile("foo.txt"), nil},
				{"bar.wav", mkfile("bar.wav"), nil},
			}, 0}
			ext := new(mockExtractor)
			SUT := parser.New(fs, ext)
			SUT.Parse("")
			expected := 1
			actual := ext.count
			if expected != actual {
				t.Errorf("expected: %d, actual: %d", expected, actual)
			}
		})
}

func mkfile(id string) *mockFI { return &mockFI{id, false} }
func mkdir(id string) *mockFI  { return &mockFI{id, true} }

type mockFS struct {
	t        *testing.T
	fixtures []fixture
	count    int
}

type fixture struct {
	path string
	info os.FileInfo
	err  error
}

func (m *mockFS) Walk(name string, f filepath.WalkFunc) error {
	m.count = 0
	for _, file := range m.fixtures {
		m.count++
		err := f(file.path, file.info, file.err)
		if err != nil {
			return err
		}
	}
	return nil
}

type mockFI struct {
	id    string
	isDir bool
}

func (f *mockFI) Name() string       { return f.id }
func (f *mockFI) Size() int64        { return 0 }
func (f *mockFI) Mode() os.FileMode  { return 0 }
func (f *mockFI) ModTime() time.Time { return time.Now() }
func (f *mockFI) IsDir() bool        { return f.isDir }
func (f *mockFI) Sys() interface{}   { return nil }

type mockExtractor struct {
	count int
}

func (e *mockExtractor) Extract(path string) { e.count++ }
