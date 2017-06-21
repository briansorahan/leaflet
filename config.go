package main

import (
	"encoding/json"
	"os"

	"github.com/namsral/flag"
	"github.com/pkg/errors"
)

// Config encapsulates the application's config.
type Config struct {
	Layers     []File
	OutputFile string
	Style      string
	Token      string
	Version    string
}

// ParseConfig parses the application's config.
func ParseConfig() (Config, error) {
	config := Config{}
	flag.StringVar(&config.OutputFile, "o", "index.html", "output file")
	flag.StringVar(&config.Style, "mapbox-style", "mapbox://styles/mapbox/streets-v10", "mapbox map style")
	flag.StringVar(&config.Token, "mapbox-token", "", "mapbox access token")
	flag.StringVar(&config.Version, "mapbox-version", "v0.38.0", "mapbox gl js version")
	flag.Parse()

	layers, err := getLayers(os.Args[1:])
	if err != nil {
		return config, errors.Wrap(err, "getting layers")
	}
	config.Layers = layers

	if len(config.Token) == 0 {
		return config, errors.New("mapbox-token is required")
	}
	return config, nil
}

type File struct {
	Contents interface{}
	Name     string
}

func getLayers(filenames []string) ([]File, error) {
	layers := make([]File, len(filenames))
	for i, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			return nil, errors.Wrap(err, "opening "+filename)
		}
		val := map[string]interface{}{}
		if err := json.NewDecoder(f).Decode(&val); err != nil {
			return nil, errors.Wrap(err, "reading "+filename)
		}
		layers[i] = File{Name: filename, Contents: val}
	}
	return layers, nil
}
