package parser

import "path/filepath"

type walker interface {
	Walk(name string, f filepath.WalkFunc) error
}

type extractor interface {
	Extract(name string)
	Close()
}

type config interface {
	AudioFormat() string
}
