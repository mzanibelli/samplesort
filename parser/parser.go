package parser

import (
	"os"
	"path/filepath"
)

type Walker interface {
	Walk(name string, f filepath.WalkFunc) error
}

type Extractor interface {
	Extract(name string)
	Close()
}

type Parser struct {
	w         Walker
	e         Extractor
	whitelist string
}

func New(w Walker, e Extractor, whitelist string) *Parser {
	return &Parser{w, e, whitelist}
}

func (p *Parser) Parse(root string) error {
	defer p.e.Close()
	return p.w.Walk(root, p.visit)
}

func (p *Parser) visit(path string, info os.FileInfo, err error) error {
	switch {
	case err != nil:
		return err
	case info.IsDir():
		return nil
	case filepath.Ext(path) != p.whitelist:
		return nil
	default:
		p.e.Extract(path)
	}
	return nil
}
