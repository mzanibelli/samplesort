lib = $(wildcard *.go) \
	  $(wildcard analyze/*.go) \
	  $(wildcard cache/*.go) \
	  $(wildcard collection/*.go) \
	  $(wildcard crypto/*.go) \
	  $(wildcard engine/*.go) \
	  $(wildcard extractor/*.go) \
	  $(wildcard parser/*.go) \
	  $(wildcard sample/*.go)

bin/samplesort: test cmd/samplesort/main.go $(lib)
	go build -o bin/samplesort samplesort/cmd/samplesort

.PHONY: test
test:
	o test -race ./...
