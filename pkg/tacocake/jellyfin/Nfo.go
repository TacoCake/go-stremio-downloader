package jellyfin

import (
	"encoding/xml"
	"os"
)

type Episode struct {
	Season      int    `xml:"season"`
	Episode     int    `xml:"episode"`
	Title       string `xml:"title,omitempty"`
	Plot        string `xml:"plot,omitempty"`
	TorrentHash string `xml:"torrentHash,omitempty"` // Non-jellyfin field
}

type Season struct {
	SeasonNumber int    `xml:"seasonNumber"`
	Title        string `xml:"title,omitempty"`
}

type Series struct {
	Title string `xml:"title"`
	Plot  string `xml:"plot"`
}

func (s *Series) MakeNfo(nfoPath string) error {
	tmp := struct {
		Series
		XMLName struct{} `xml:"tvshow"`
	}{Series: *s}

	xmlData, err := xml.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(nfoPath, xmlData, 0775)
	if err != nil {
		return err
	}
	return nil
}

func (s *Season) MakeNfo(nfoPath string) error {
	tmp := struct {
		Season
		XMLName struct{} `xml:"season"`
	}{Season: *s}

	xmlData, err := xml.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(nfoPath, xmlData, 0775)
	if err != nil {
		return err
	}
	return nil
}
func (e *Episode) MakeNfo(nfoPath string) error {
	tmp := struct {
		Episode
		XMLName struct{} `xml:"episode"`
	}{Episode: *e}
	xmlData, err := xml.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(nfoPath, xmlData, 0775)
	if err != nil {
		return err
	}
	return nil
}
