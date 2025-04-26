package stremio

import (
	"encoding/json"
	"os"
)

/*
{
    "metas": [
        {
            "type": "series",
            "id": "pp_onepace",
            "name": "One Pace",
            "poster": "https://i.pinimg.com/originals/eb/85/c4/eb85c4376b474030b80afa80ad1cd13a.jpg",
            "genres": [
                "Adventure",
                "Fantasy"
            ]
        }
    ]
}
*/

type CatalogMeta struct {
	Type   string   `json:"type"`
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Poster string   `json:"poster"`
	Genres []string `json:"genres"`
}

type Catalog struct {
	Metas []CatalogMeta `json:"metas"`
}

func LoadSeriesCatalog(catalogPath string) (*Catalog, error) {
	data, err := os.ReadFile(catalogPath)
	if err != nil {
		return nil, err
	}
	var catalog Catalog
	err = json.Unmarshal(data, &catalog)
	if err != nil {
		return nil, err
	}
	return &catalog, nil
}
