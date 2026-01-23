package scanner

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
}

type CollectionState struct {
	Collection db.Collection
	Videos     []db.Video
}

func Scan(root string) (World, error) {
	world := World{}

	vaults, err := scanVaults(root)

	if err != nil {
		return world, err
	}

	for _, vault := range vaults {
		vaultState := VaultState{Vault: vault}

		vaultVideosPath := filepath.Join(root, "vaults", vault.Name, "videos")
		vaultPicturesPath := filepath.Join(root, "vaults", vault.Name, "pictures")

		actors, _ := scanActors(vaultPicturesPath)
		vaultState.Actors = actors

		galleries, _ := scanGalleries(vaultPicturesPath)
		vaultState.Galleries = galleries

		collections, _ := scanCollections(vaultVideosPath)

		for _, c := range collections {
			cs := CollectionState{Collection: c}

			collectionPath := filepath.Join(vaultVideosPath, c.Slug)

			videos, err := scanVideos(collectionPath)
			if err != nil {
				log.Println("video scan error:", err)
				continue
			}

			cs.Videos = videos
			vaultState.Collections = append(vaultState.Collections, cs)
		}

		world.Vaults = append(world.Vaults, vaultState)
	}

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

			log.Printf("vaults: %v", vaults)
		}
	}

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

	return galleries, nil
}

func scanActors(path string) ([]db.Actor, error) {
	actorsPath := filepath.Join(path, "actors")
	entries, err := os.ReadDir(actorsPath)

	log.Printf("path %v", actorsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read actors: %w", err)
	}

	var actors []db.Actor

	for _, entry := range entries {
		actor := strings.TrimSuffix(entry.Name(), ".jpg")

		log.Printf("scanned actor: %v", actor)

		actors = append(actors, db.Actor{
			Name: utils.SnakeToTitle(actor),
			Slug: actor,
		})
	}

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

			log.Printf("collections: %v", collections)
		}
	}
	return collections, nil
}

func scanVideos(collectionPath string) ([]db.Video, error) {
	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection: %w", err)
	}

	var videos []db.Video

	for _, entry := range entries {
		if entry.IsDir() {
			folderName := entry.Name()
			nfoPath := filepath.Join(collectionPath, folderName, folderName+".nfo")

			if _, err := os.Stat(nfoPath); err != nil {
				return nil, fmt.Errorf("missing .nfo file for folder %v", folderName)
			}

			metadata, err := parseNfoFile(nfoPath)
			if err != nil {
				return nil, fmt.Errorf("failed to parse .nfo for %v: %w", folderName, err)
			}

			videos = append(videos, db.Video{
				Title:  metadata.Title,
				Slug:   folderName,
				Studio: metadata.Studio,
				Tags:   metadata.Tags,
				Actors: metadata.Actors,
			})
		}
	}

	return videos, nil
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
