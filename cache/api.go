package cache

type storage interface {
	Exists(name string) bool
	ReadAll(name string) ([]byte, error)
	WriteAll(name string, data []byte) error
}

type config interface {
	FileSystemRoot() string
	DataFormat() string
	EnableCache() bool
}
