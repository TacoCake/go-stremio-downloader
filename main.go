package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"tacocake/go-stremio-downloader/pkg/tacocake/jellyfin"
	"tacocake/go-stremio-downloader/pkg/tacocake/stremio"

	"github.com/anacrolix/torrent"
)

func main() {

	var stremioPackagePath = flag.String("p", "", "The path to the Stremio package")
	var outputPath = flag.String("o", "", "The output folder path")
	flag.Parse()

	if stremioPackagePath == nil || *stremioPackagePath == "" {
		log.Fatal("Missing required flag -p <path to Stremio package>")
		return
	}

	// Load Series Catalog
	catalog, err := stremio.LoadSeriesCatalog(path.Join(*stremioPackagePath, "catalog", "series", "seriesCatalog.json"))
	if err != nil {
		fmt.Println("Error loading series catalog:", err)
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	// Create a torrent client
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = "/tmp/go-torrent"
	client, err := torrent.NewClient(clientConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create torrent client: %w", err))
	}
	defer client.Close()

	for _, m := range catalog.Metas {
		// Get Series
		fmt.Printf("Series Name: %s\n", m.Name)

		manifest, err := stremio.LoadSeriesMetaData(path.Join(*stremioPackagePath, "meta", "series", m.Id+".json"))
		if err != nil {
			fmt.Printf("Error loading series manifest for %s: %v\n", m.Name, err)
			continue
		}

		fmt.Printf("Number of episodes: %d\n", len(manifest.Meta.Videos))

		seriesBaseFolder := path.Join(*outputPath, m.Name)

		series := jellyfin.Series{
			Title: manifest.Meta.Name,
			Plot:  manifest.Meta.Description,
		}

		series.MakeNfo(path.Join(seriesBaseFolder, "tvshow.nfo"))

		// Sort videos by season and episode number
		sort.Slice(manifest.Meta.Videos, func(i, j int) bool {
			if manifest.Meta.Videos[i].Season != manifest.Meta.Videos[j].Season {
				return manifest.Meta.Videos[i].Season < manifest.Meta.Videos[j].Season
			}
			return manifest.Meta.Videos[i].Episode < manifest.Meta.Videos[j].Episode
		})

		// Prep seasons
		markFiles := make([]string, 0)
		seasonFolders := make(map[int]string)
		for _, v := range manifest.Meta.Videos {
			seasonFolderName := fmt.Sprintf("Season %d", v.Season)
			seasonFolder := path.Join(seriesBaseFolder, seasonFolderName)
			seasonFolders[v.Season] = seasonFolder
			if v.Episode == 1 {
				// mkdir season dir
				os.MkdirAll(seasonFolder, 0775)

				season := jellyfin.Season{
					SeasonNumber: v.Season,
				}
				nfoPath := path.Join(seasonFolder, "season.nfo")
				if _, err := os.Stat(nfoPath); err != nil {
					season.MakeNfo(nfoPath)
				}
			}

			markFiles = append(markFiles, processEpisode(client, clientConfig, stremioPackagePath, seriesBaseFolder, seasonFolder, v, *manifest))
		}

		for _, seasonFolder := range seasonFolders {
			// loop all files with .downloaded in seasonFolder
			files, err := os.ReadDir(seasonFolder)
			if err != nil {
				log.Fatal(err)
			}

			for _, file := range files {
				if !strings.HasSuffix(file.Name(), ".downloaded") {
					continue
				}

				// check if in list of markFiles
				found := false
				for _, markFile := range markFiles {
					if strings.HasSuffix(markFile, file.Name()) {
						found = true
						break
					}
				}

				// if not found, delete
				if found {
					continue
				}

				p := path.Join(seasonFolder, file.Name())
				os.Remove(p)
			}
		}
	}
}

func processEpisode(client *torrent.Client, clientConfig *torrent.ClientConfig, stremioPackagePath *string, seriesBaseFolder string, seasonFolder string, v stremio.StremioVideo, manifest stremio.StremioSeriesManifest) string {
	// Get Stream Manifest
	streamManifest, err := stremio.LoadSeriesStream(path.Join(*stremioPackagePath, "stream", "series", fmt.Sprintf("%s.json", v.Id)))
	if err != nil {
		log.Printf("Failed to load stream manifest for %s: %v", v.Id, err)
		return ""
	}
	if len(streamManifest.Streams) != 1 {
		log.Fatal(fmt.Errorf("expected exactly one stream in manifest for %s", v.Id))
	}

	// Get Stream
	episodeStream := streamManifest.Streams[0]
	markFile := path.Join(seasonFolder, fmt.Sprintf(".%s_%d.downloaded", episodeStream.InfoHash, episodeStream.FileIdx))
	if _, err := os.Stat(markFile); err == nil {
		log.Printf("Skipping downloaded episode S%02dE%02d\n", v.Season, v.Episode)
		return markFile
	}

	torrentFile, err := episodeStream.DownloadEpisode(client, seasonFolder)
	defer func() {
		p := path.Join(clientConfig.DataDir, torrentFile.Path())
		log.Println("Deleting file:", p)
		// delete file
		os.Remove(p)
	}()
	if err != nil {
		log.Printf("Error downloading episode %d for season %d: %v\n", v.Episode, v.Season, err)
		return ""
	}

	// Rename file to match season and episode number
	torrentFilePath := path.Base(torrentFile.Path())
	oldName := path.Join(seasonFolder, torrentFilePath)
	ext := filepath.Ext(torrentFilePath)

	// Keep everything that is inside square brackets
	extraInformation := ""
	re := regexp.MustCompile(`\[.*?\]|\(.*?\)`)
	matches := re.FindAllString(torrentFilePath, -1)
	if len(matches) > 0 {
		extraInformation = strings.Join(matches, "")
	}

	newName := path.Join(seasonFolder, fmt.Sprintf("%s - S%02dE%02d %s%s", manifest.Meta.Name, v.Season, v.Episode, extraInformation, ext))
	os.Rename(oldName, newName)

	episodeFileName := strings.TrimSuffix(newName, filepath.Ext(newName))

	// Add episode with filename
	episode := jellyfin.Episode{
		Title:       v.Title,
		Season:      v.Season,
		Episode:     v.Episode,
		TorrentHash: episodeStream.InfoHash,
	}
	nfoFilename := fmt.Sprintf("%s.nfo", episodeFileName)
	episode.MakeNfo(nfoFilename)

	// Add file to mark downloaded hash
	os.WriteFile(markFile, []byte{}, 0755)
	return markFile
}
