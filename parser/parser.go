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
	walker    Walker
	extractor Extractor
	whitelist string
}

func New(walker Walker, extractor Extractor, whitelist string) *Parser {
	return &Parser{walker, extractor, whitelist}
}

func (p *Parser) Parse(root string) error {
	defer p.extractor.Close()
	return p.walker.Walk(root, p.visit)
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
		p.extractor.Extract(path)
	}
	return nil
}
