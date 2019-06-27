package cache

func FetchTest(cfg *TestConfig) error {
	return New(cfg, cfg).Fetch("foo/bar", cfg, cfg.build)
}

type TestConfig struct {
	FileExists     bool
	FileContent    []byte
	WriteError     error
	ReadError      error
	BuildCreates   bool
	BuildData      interface{}
	BuildError     error
	FsRoot         string
	Format         string
	IsEnabled      bool
	MarshalError   error
	UnmarshalError error
}

func (t *TestConfig) ReadAll(name string) ([]byte, error) { return t.FileContent, t.ReadError }
func (t *TestConfig) Exists(name string) bool             { return t.FileExists }
func (t *TestConfig) FileSystemRoot() string              { return t.FsRoot }
func (t *TestConfig) DataFormat() string                  { return ".json" }
func (t *TestConfig) EnableCache() bool                   { return t.IsEnabled }
func (t *TestConfig) MarshalJSON() ([]byte, error)        { return nil, t.MarshalError }
func (t *TestConfig) UnmarshalJSON(data []byte) error     { return t.UnmarshalError }

func (t *TestConfig) build() (interface{}, error) {
	if t.BuildCreates && t.BuildError == nil {
		t.FileExists = true
	}
	return t.BuildData, t.BuildError
}

func (t *TestConfig) WriteAll(name string, data []byte) error {
	if t.WriteError == nil {
		t.FileExists = true
	}
	return t.WriteError
}
