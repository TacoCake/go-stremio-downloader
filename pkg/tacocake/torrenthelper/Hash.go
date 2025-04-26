package torrenthelper

import (
	"context"
	"fmt"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func GetTorrentByHash(client *torrent.Client, torrentHash string) (*torrent.Torrent, error) {
	// Convert the hash to a torrent.InfoHash
	var infoHash metainfo.Hash
	err := infoHash.FromHexString(torrentHash)
	if err != nil {
		return nil, fmt.Errorf("invalid torrent hash: %w", err)
	}

	// Add the torrent using only the info hash (no .torrent file)
	t, err := client.AddMagnet("magnet:?xt=urn:btih:" + torrentHash)
	if err != nil {
		return nil, fmt.Errorf("failed to add magnet: %w", err)
	}

	// Wait for metadata to be downloaded
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Fetching metadata...")

	select {
	case <-t.GotInfo():
		// Metadata fetched
	case <-ctx.Done():
		return nil, fmt.Errorf("timed out waiting for metadata")
	}

	return t, nil
}
