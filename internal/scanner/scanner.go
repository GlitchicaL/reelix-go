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

func Scan() {
	root := "/videos"

	vaults, err := scanVaults(root)

	if err != nil {
		log.Println("scan vaults error: %w", err)
	}

	dbVaults, err := SyncVaults(vaults)

	if err != nil {
		log.Println("vault sync error:", err)
	}

	for _, v := range dbVaults {
		vaultPath := filepath.Join(root, "/Vaults", v.Name, "/Videos")

		actors, err := scanActors(root, v.Name)

		if err != nil {
			log.Println("actor scan error: %w", err)
		}

		if err := SyncActors(actors); err != nil {
			log.Println("actors sync error:", err)
		}

		galleryPath := filepath.Join(root, "/Vaults", v.Name, "/Pictures")
		galleries, err := scanGalleries(galleryPath, v.ID)

		if err != nil {
			log.Println("scan gallery error: %w", err)
		}

		if err := SyncGalleries(galleries); err != nil {
			log.Println("gallery sync error:", err)
		}

		collections, err := scanCollections(vaultPath, v.ID)

		if err != nil {
			log.Println("scan collections error: %w", err)
		}

		dbCollections, err := SyncCollections(collections)

		if err != nil {
			log.Println("collection sync error:", err)
		}

		for _, c := range dbCollections {
			collectionPath := filepath.Join(root, "/Vaults", v.Name, "/Videos", c.Name)

			videos, err := scanVideos(collectionPath, v.ID, c.ID)

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
	vaultsPath := filepath.Join(rootPath, "/Vaults")
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

func scanGalleries(picturePath string, vaultId int) ([]db.Gallery, error) {
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
				VaultID:    vaultId,
			})
		}
	}

	return galleries, nil
}

func scanActors(rootPath string, vaultName string) ([]db.Actor, error) {
	actorsPath := filepath.Join(rootPath, "/Vaults/", vaultName, "/Pictures/actors/")
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
