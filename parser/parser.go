package parser

import (
	"os"
	"path/filepath"
)

const (
	whitelist string = ".wav"
)

type Walker interface {
	Walk(name string, f filepath.WalkFunc) error
}

type Extractor interface {
	Extract(name string)
}

type Parser struct {
	w Walker
	e Extractor
}

func New(w Walker, e Extractor) *Parser {
	return &Parser{w, e}
}

func (p *Parser) Parse(root string) error {
	return p.w.Walk(root, p.visit)
}

func (p *Parser) visit(path string, info os.FileInfo, err error) error {
	switch {
	case err != nil:
		return err
	case info.IsDir():
		return nil
	case filepath.Ext(path) != whitelist:
		return nil
	default:
		p.e.Extract(path)
	}
	return nil
}
