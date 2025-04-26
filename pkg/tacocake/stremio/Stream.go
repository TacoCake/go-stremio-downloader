package stremio

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"tacocake/one_pace/pkg/tacocake/torrenthelper"

	"github.com/anacrolix/torrent"
)

type StremioSeriesStream struct {
	InfoHash string `json:"infoHash"`
	FileIdx  int    `json:"fileIdx"`
}

type StremioSeriesStreamManifest struct {
	Streams []StremioSeriesStream `json:"streams"`
}

func LoadSeriesStream(manifestPath string) (*StremioSeriesStreamManifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	var manifest StremioSeriesStreamManifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}
	return &manifest, nil
}

func (sss *StremioSeriesStream) DownloadEpisode(client *torrent.Client, folderPath string) (*torrent.File, error) {
	hash := sss.InfoHash

	torrentItem, err := torrenthelper.GetTorrentByHash(client, hash)
	if err != nil {
		log.Fatal(err)
	}

	// Download episode torrent file
	torrentItem.DownloadAll()
	var downloadedTorrent *torrent.File
	for index, torrentFile := range torrentItem.Files() {
		if index != sss.FileIdx {
			torrentFile.SetPriority(torrent.PiecePriorityNone)
			continue
		}

		if downloadedTorrent != nil {
			log.Fatal(fmt.Errorf("Multiple files would be downloaded%s", hash))
		}
		downloadedTorrent = torrentFile
		torrentFile.Download()

		torrentfilePath := path.Base(torrentFile.Path())

		fmt.Println("Downloading:", torrentfilePath)
		r := torrentFile.NewReader()
		defer r.Close()

		// Copy file to disk
		f, err := os.OpenFile(path.Join(folderPath, torrentfilePath), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0775)
		if err != nil {
			log.Fatal(fmt.Errorf("Failed to open file for writing: %w", err))
		}
		defer f.Close()

		_, err = io.Copy(f, r)
		if err != nil {
			log.Fatal(fmt.Errorf("Failed to write file: %w", err))
		}

		fmt.Println("Downloaded:", torrentfilePath)
	}
	return downloadedTorrent, nil
}
