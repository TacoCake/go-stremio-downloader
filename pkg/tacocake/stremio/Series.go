package stremio

import (
	"encoding/json"
	"os"
)

type StremioVideo struct {
	Season  int    `json:"season"`
	Episode int    `json:"episode"`
	Id      string `json:"id"`
	Title   string `json:"title"`
}

type StremioSeriesMeta struct {
	Id          string         `json:"id"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Poster      string         `json:"poster"`
	Genres      []string       `json:"genres"`
	Description string         `json:"description"`
	Director    []string       `json:"director"`
	Logo        string         `json:"logo"`
	Background  string         `json:"background"`
	Videos      []StremioVideo `json:"videos"`
}

type StremioSeriesManifest struct {
	Meta StremioSeriesMeta `json:"meta"`
}

func LoadSeriesMetaData(manifestFile string) (*StremioSeriesManifest, error) {
	data, err := os.ReadFile(manifestFile)
	if err != nil {
		return nil, err
	}
	var stremioSeriesManifest StremioSeriesManifest
	err = json.Unmarshal(data, &stremioSeriesManifest)
	if err != nil {
		return nil, err
	}
	return &stremioSeriesManifest, nil
}
