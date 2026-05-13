package scanner

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"reelix-go/internal/db"
	"reelix-go/internal/utils"
)

type World struct {
	Vaults []VaultState
}

type VaultState struct {
	Vault       db.Vault
	Collections []CollectionState
	Galleries   []db.Gallery
	Actors      []db.Actor
	Tags        []db.Tag
}

type CollectionState struct {
	Collection db.Collection
	Videos     []db.Video

	// Needs research:
	// This should give us O(1) lookups while preserving order in the above slice.
	VideoSlugToID map[string]int
}

func Scan(root string) (World, error) {
	world := World{}

	vaults, err := scanVaults(root)

	if err != nil {
		return world, err
	}

	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	sem := make(chan struct{}, 2)

	for _, vault := range vaults {
		sem <- struct{}{}
		wg.Add(1)

		go func(v db.Vault) {
			defer wg.Done()
			defer func() { <-sem }()

			vaultState := VaultState{Vault: v}

			vaultCollectionsPath := filepath.Join(root, "vaults", v.Name, "collections")
			vaultPicturesPath := filepath.Join(root, "vaults", v.Name, "pictures")

			actors, _ := scanActors(vaultPicturesPath)
			vaultState.Actors = actors

			galleries, _ := scanGalleries(vaultPicturesPath)
			vaultState.Galleries = galleries

			collections, _ := scanCollections(vaultCollectionsPath)

			tagsSeen := make(map[string]struct{})

			for _, c := range collections {
				cs := CollectionState{Collection: c}

				collectionPath := filepath.Join(vaultCollectionsPath, c.Slug)

				videos, tags, err := scanVideos(collectionPath)
				if err != nil {
					log.Println("video scan error:", err)
					continue
				}

				cs.Videos = videos

				for _, tag := range tags {
					if _, exists := tagsSeen[tag]; !exists {
						tagsSeen[tag] = struct{}{}
						vaultState.Tags = append(vaultState.Tags, db.Tag{Name: tag})
					}
				}

				vaultState.Collections = append(vaultState.Collections, cs)
			}

			mu.Lock()

			/*
				Since multiple goroutines will append to the
				world state, we lock to avoid race conditions
			*/

			world.Vaults = append(world.Vaults, vaultState)

			mu.Unlock()

			log.Printf("%v added to world state", v.Name)
		}(vault)
	}

	wg.Wait()

	return world, nil
}

func scanVaults(rootPath string) ([]db.Vault, error) {
	vaultsPath := filepath.Join(rootPath, "vaults")
	entries, err := os.ReadDir(vaultsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	var vaults []db.Vault

	for _, entry := range entries {
		if entry.IsDir() {
			vaults = append(vaults, db.Vault{
				Name: entry.Name(),
			})
		}
	}

	log.Printf("vaults scanned: %v", len(vaults))

	return vaults, nil
}

func scanGalleries(picturePath string) ([]db.Gallery, error) {
	entries, err := os.ReadDir(picturePath)

	if err != nil {
		return nil, fmt.Errorf("failed to read collection: %w", err)
	}

	var galleries []db.Gallery

	for _, entry := range entries {
		if entry.IsDir() {
			galleryName := entry.Name()

			// We ignore the actors/ folder as there
			// is a separate scanning/syncing flow for actors.
			if galleryName == "actors" {
				continue
			}

			galleryPath := filepath.Join(picturePath, galleryName)
			galleryEntries, err := os.ReadDir(galleryPath)

			if err != nil {
				return nil, err
			}

			galleryImageCount := 0

			for _, galleryEntry := range galleryEntries {
				if !galleryEntry.IsDir() {
					galleryImageCount++
				}
			}

			galleries = append(galleries, db.Gallery{
				Title:      utils.SnakeToTitle(galleryName),
				Slug:       galleryName,
				ImageCount: galleryImageCount,
			})
		}
	}

	log.Printf("galleries scanned: %v", len(galleries))

	return galleries, nil
}

func scanActors(path string) ([]db.Actor, error) {
	actorsPath := filepath.Join(path, "actors")
	entries, err := os.ReadDir(actorsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read actors: %w", err)
	}

	var actors []db.Actor

	for _, entry := range entries {
		actor := strings.TrimSuffix(entry.Name(), ".jpg")

		actors = append(actors, db.Actor{
			Name: utils.SnakeToTitle(actor),
			Slug: actor,
		})
	}

	log.Printf("actors scanned: %v", len(actors))

	return actors, nil
}

func scanCollections(vaultPath string) ([]db.Collection, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	var collections []db.Collection

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()

			collections = append(collections, db.Collection{
				Name: utils.SnakeToTitle(entry.Name()),
				Slug: name,
				Path: filepath.Join(vaultPath, name),
			})
		}
	}

	log.Printf("collections scanned: %v", len(collections))

	return collections, nil
}

func scanVideos(collectionPath string) ([]db.Video, []string, error) {
	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read collection: %w", err)
	}

	var videos []db.Video
	var tags []string

	for _, entry := range entries {
		if entry.IsDir() {
			folderName := entry.Name()
			nfoPath := filepath.Join(collectionPath, folderName, folderName+".nfo")

			if _, err := os.Stat(nfoPath); err != nil {
				return nil, nil, fmt.Errorf("missing .nfo file for folder %v", folderName)
			}

			metadata, err := parseNfoFile(nfoPath)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse .nfo for %v: %w", folderName, err)
			}

			videos = append(videos, db.Video{
				Title:  metadata.Title,
				Slug:   folderName,
				Studio: metadata.Studio,
				Tags:   metadata.Tags,
				Actors: metadata.Actors,
			})

			for _, tag := range metadata.Tags {
				tags = append(tags, tag)
			}
		}
	}

	log.Printf("videos scanned: %v", len(videos))
	log.Printf("tags scanned: %v", len(tags))

	return videos, tags, nil
}

type VideoMetadata struct {
	Title  string     `xml:"title"`
	Studio string     `xml:"studio"`
	Tags   []string   `xml:"tag"`
	Actors []db.Actor `xml:"actor"`
}

func parseNfoFile(nfoPath string) (VideoMetadata, error) {
	data, err := os.ReadFile(nfoPath)
	if err != nil {
		return VideoMetadata{}, err
	}

	var metadata VideoMetadata
	err = xml.Unmarshal(data, &metadata)
	if err != nil {
		return VideoMetadata{}, err
	}

	return metadata, nil
}
