package scanner

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"reelix-go/internal/db"
)

func Scan() {
	root := "/videos/"

	vaults, err := scanVaults(root)

	if err != nil {
		log.Println("scan vaults error: %w", err)
	}

	if err := SyncVaults(vaults); err != nil {
		log.Println("vault sync error:", err)
	}

	for i, v := range vaults {
		vaultPath := filepath.Join(root, v.Name)

		collections, err := scanCollections(vaultPath, i+1)

		if err != nil {
			log.Println("scan collections error: %w", err)
		}

		if err := SyncCollections(collections); err != nil {
			log.Println("collection sync error:", err)
		}

		for j, c := range collections {
			collectionPath := filepath.Join(root, v.Name, c.Name)

			videos, err := scanVideos(collectionPath, i+1, j+1)

			if err != nil {
				log.Println("scan videos error: %w", err)
			}

			if err := SyncVideos(videos); err != nil {
				log.Println("video sync error:", err)
			}
		}
	}

}

func scanVaults(rootPath string) ([]db.Vault, error) {
	entries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	var vaults []db.Vault

	for _, entry := range entries {
		if entry.IsDir() {
			// vaultPath := filepath.Join(rootPath, entry.Name())

			vaults = append(vaults, db.Vault{
				Name: entry.Name(),
			})

			log.Printf("vaults: %v", vaults)
		}
	}

	return vaults, nil
}

func scanCollections(vaultPath string, vaultID int) ([]db.Collection, error) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault: %w", err)
	}

	var collections []db.Collection

	for _, entry := range entries {
		if entry.IsDir() {
			collections = append(collections, db.Collection{
				Name:    entry.Name(),
				Path:    filepath.Join(vaultPath, entry.Name()),
				VaultID: vaultID,
			})

			log.Printf("collections: %v", collections)
		}
	}
	return collections, nil
}

func scanVideos(collectionPath string, vaultID int, collectionID int) ([]db.Video, error) {
	entries, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection: %w", err)
	}

	var videos []db.Video

	for _, entry := range entries {
		if entry.IsDir() {
			folderName := entry.Name()
			nfoPath := filepath.Join(collectionPath, folderName, folderName+".nfo")

			// Check if .nfo file exists
			if _, err := os.Stat(nfoPath); err != nil {
				return nil, fmt.Errorf("missing .nfo file for folder %v", folderName)
			}

			// Parse .nfo file
			metadata, err := parseNfoFile(nfoPath)
			if err != nil {
				return nil, fmt.Errorf("failed to parse .nfo for %v: %w", folderName, err)
			}

			videos = append(videos, db.Video{
				Title:        metadata.Title,
				Slug:         folderName,
				Studio:       metadata.Studio,
				Tags:         metadata.Tags,
				Actors:       metadata.Actors,
				VaultID:      vaultID,
				CollectionID: collectionID,
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
