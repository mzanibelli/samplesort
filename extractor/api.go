package extractor

type cache interface {
	Fetch(key string, target interface{}, build func() (interface{}, error)) error
}

type buildFunc func(src string) (interface{}, error)
