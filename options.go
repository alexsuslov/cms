package cms

import "os"

import (
	"gopkg.in/yaml.v3"
)

type Options map[string]interface{}

func Load(filename string) (opts *Options, err error) {
	opts = &Options{}

	f, err := os.Open(filename)
	if err != nil {
		return
	}
	err = yaml.NewDecoder(f).Decode(opts)

	return opts.
		Set("Generator", "github.com/alexsuslov/cms").
		Set("Authtor", "AlexSuslov"), nil
}

func Check(data []byte) error {
	opts := &Options{}
	return yaml.Unmarshal(data, opts)
}

func (Options *Options) Set(key string, value interface{}) *Options {
	o := *Options
	o[key] = value
	return &o
}

func (Options *Options) Refresh(data []byte) error {
	return yaml.Unmarshal(data, Options)
}

func (Options Options) Extend(m Options) Options {
	for key, val := range Options {
		m[key] = val
	}
	return m
}
