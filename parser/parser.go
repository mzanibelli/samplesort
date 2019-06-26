package parser

import (
	"os"
	"path/filepath"
)

type Parser struct {
	walker    walker
	extractor extractor
	cfg       config
}

func New(walker walker, extractor extractor, cfg config) *Parser {
	return &Parser{walker, extractor, cfg}
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
	case filepath.Ext(path) != p.cfg.AudioFormat():
		return nil
	default:
		p.extractor.Extract(path)
	}
	return nil
}
