package cache_test

import (
	"encoding/json"
	"errors"
	"samplesort/cache"
	"testing"
)

var (
	defaultFileData  = []byte(`{"foo":"bar"}`)
	defaultBuildData = map[string]string{"foo": "bar"}
)

var cases = []struct {
	name    string
	success bool
	cfg     *cache.TestConfig
}{{
	"it should succeed if enabled, file exists and is successfuly decoded",
	true,
	&cache.TestConfig{
		IsEnabled:   true,
		FileExists:  true,
		FileContent: defaultFileData,
	},
}, {
	"it should fail if enabled, file exists and can't be read",
	false,
	&cache.TestConfig{
		IsEnabled:  true,
		FileExists: true,
		ReadError:  errors.New("foo"),
	},
}, {
	"it should fail if enabled, file exists and decode fails",
	false,
	&cache.TestConfig{
		IsEnabled:      true,
		FileExists:     true,
		UnmarshalError: errors.New("foo"),
	},
}, {
	"it should fail if enabled, file does not exists and build fails",
	false,
	&cache.TestConfig{
		IsEnabled:   true,
		FileContent: defaultFileData,
		BuildError:  errors.New("foo"),
	},
}, {
	"it should succeed if enabled and a well-formed file was built",
	true,
	&cache.TestConfig{
		IsEnabled:    true,
		FileContent:  defaultFileData,
		BuildCreates: true,
	},
}, {
	"it should fail if enabled and an unreadable file was built",
	false,
	&cache.TestConfig{
		IsEnabled:    true,
		BuildCreates: true,
		ReadError:    errors.New("foo"),
	},
}, {
	"it should fail if enabled and a malformed file was built",
	false,
	&cache.TestConfig{
		IsEnabled:      true,
		BuildCreates:   true,
		UnmarshalError: errors.New("foo"),
	},
}, {
	"it should fail if enabled, file does not exist and encoding fails",
	false,
	&cache.TestConfig{
		IsEnabled:   true,
		FileContent: defaultFileData,
		BuildData:   json.RawMessage("{"),
	},
}, {
	"it should fail if enabled, file does not exist and write fails",
	false,
	&cache.TestConfig{
		IsEnabled:   true,
		FileContent: defaultFileData,
		WriteError:  errors.New("foo"),
	},
}, {
	"it should succeed if enabled and the written file is well-formed",
	true,
	&cache.TestConfig{
		IsEnabled:   true,
		FileContent: defaultFileData,
	},
}, {
	"it should fail if enabled and Path() fails",
	false,
	&cache.TestConfig{
		IsEnabled:   true,
		FileContent: defaultFileData,
		FsRoot:      "/",
	},
}, {
	"it should fail if enabled and the written file is unreadable",
	false,
	&cache.TestConfig{
		IsEnabled: true,
		ReadError: errors.New("foo"),
	},
}, {
	"it should fail if enabled and the written file is malformed",
	false,
	&cache.TestConfig{
		IsEnabled:      true,
		UnmarshalError: errors.New("foo"),
	},
}, {
	"it should fail if enabled and the written file is empty",
	false,
	&cache.TestConfig{
		IsEnabled: true,
	},
}, {
	"it should succeed if disabled and a well-formed file was built",
	true,
	&cache.TestConfig{
		FileContent:  defaultFileData,
		BuildCreates: true,
	},
}, {
	"it should fail if disabled and Path() fails",
	false,
	&cache.TestConfig{
		FileContent: defaultFileData,
		FsRoot:      "/",
	},
}, {
	"it should fail if disabled and build fails",
	false,
	&cache.TestConfig{
		BuildError: errors.New("foo"),
	},
}, {
	"it should fail if disabled and an unreadable file was built",
	false,
	&cache.TestConfig{
		BuildCreates: true,
		ReadError:    errors.New("foo"),
	},
}, {
	"it should fail if disabled and a malformed file was built",
	false,
	&cache.TestConfig{
		BuildCreates:   true,
		UnmarshalError: errors.New("foo"),
	},
}, {
	"it should fail if disabled, file does not exist and encoding fails",
	false,
	&cache.TestConfig{
		FileContent: defaultFileData,
		BuildData:   json.RawMessage("{"),
	},
}, {
	"it should fail if disabled, file does not exist and write fails",
	false,
	&cache.TestConfig{
		FileContent: defaultFileData,
		WriteError:  errors.New("foo"),
	},
}, {
	"it should succeed if disabled and the written file is well-formed",
	true,
	&cache.TestConfig{
		FileContent: defaultFileData,
	},
}, {
	"it should fail if disabled and the written file is unreadable",
	false,
	&cache.TestConfig{
		ReadError: errors.New("foo"),
	},
}, {
	"it should fail if disabled and the written file is malformed",
	false,
	&cache.TestConfig{
		UnmarshalError: errors.New("foo"),
	},
},
}

func TestFetch(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := cache.FetchTest(c.cfg)
			t.Log(err)
			if (err == nil) != c.success {
				t.Fatal("error assertion failed")
			}
		})
	}
}
