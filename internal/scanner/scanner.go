package scanner

import (
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

			vaultCollectionsPath := filepath.Join(root, "vaults", v.Slug, "collections")
			vaultPicturesPath := filepath.Join(root, "vaults", v.Slug, "pictures")

			actors, err := scanActors(vaultPicturesPath)

			if err != nil {
				log.Println(err)
			}

			vaultState.Actors = actors

			galleries, err := scanGalleries(vaultPicturesPath)

			if err != nil {
				log.Println(err)
			}

			vaultState.Galleries = galleries

			collections, err := scanCollections(vaultCollectionsPath)

			if err != nil {
				log.Println(err)
			}

			tagsSeen := make(map[string]struct{})

			for _, c := range collections {
				cs := CollectionState{Collection: c}

				collectionPath := filepath.Join(vaultCollectionsPath, c.Slug)

				videos, tags, err := scanVideos(collectionPath)

				if err != nil {
					log.Println(err)
					continue
				}

				cs.Videos = videos

				for _, tag := range tags {
					if _, exists := tagsSeen[tag]; !exists {
						mu.Lock()

						tagsSeen[tag] = struct{}{}
						vaultState.Tags = append(vaultState.Tags, db.Tag{Name: tag})

						mu.Unlock()
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

			log.Printf("%v added to world state", v.Slug)
		}(vault)
	}

	wg.Wait()

	return world, nil
}

func scanVaults(rootPath string) ([]db.Vault, error) {
	vaultsPath := filepath.Join(rootPath, "vaults")
	entries, err := os.ReadDir(vaultsPath)

	if err != nil {
		return nil, fmt.Errorf("failed to read vaults: %w", err)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no vaults detected")
	}

	/*
		Since we know the maximum amount of entries,
		we can preallocate the slice.
	*/

	vaults := make([]db.Vault, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			vaultName := entry.Name()
			vaultPath := filepath.Join(vaultsPath, vaultName)

			/*
				Since the vault is the root for future scans,
				if the dir is empty, no need to append it and
				add more loops.

				We don't care about the function returning
				an error. If there is an error, we treat it
				as an empty dir and continue
			*/

			isEmpty, _ := utils.IsDirEmpty(vaultPath)

			if isEmpty {
				continue
			}

			vaults = append(vaults, db.Vault{
				Name: utils.ToTitle(vaultName),
				Slug: vaultName,
			})
		}
	}

	log.Printf("vaults scanned: %v", len(vaults))

	return vaults, nil
}

func scanGalleries(picturePath string) ([]db.Gallery, error) {
	entries, err := os.ReadDir(picturePath)

	if err != nil {
		return nil, fmt.Errorf("scan galleries: %w", err)
	}

	/*
		Since we know the maximum amount of entries,
		we can preallocate the slice.
	*/

	galleries := make([]db.Gallery, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			galleryName := entry.Name()

			/*
				We ignore the actors/ folder as there
				is a separate scanning flow for actors.
			*/

			if galleryName == "actors" {
				continue
			}

			galleryPath := filepath.Join(picturePath, galleryName)
			galleryEntries, err := os.ReadDir(galleryPath)

			if err != nil {
				return nil, fmt.Errorf("scan galleries: %w", err)
			}

			galleryImageCount := 0

			for _, galleryEntry := range galleryEntries {
				if !galleryEntry.IsDir() {
					galleryImageCount++
				}
			}

			galleries = append(galleries, db.Gallery{
				Title:      utils.ToTitle(galleryName),
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
		return nil, fmt.Errorf("scanning actors: %w", err)
	}

	/*
		Since we know the maximum amount of entries,
		we can preallocate the slice.
	*/

	actors := make([]db.Actor, 0, len(entries))

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
		return nil, fmt.Errorf("scan collections: %w", err)
	}

	/*
		Since we know the maximum amount of entries,
		we can preallocate the slice.
	*/

	collections := make([]db.Collection, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()

			collections = append(collections, db.Collection{
				Name: utils.ToTitle(entry.Name()),
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
		return nil, nil, fmt.Errorf("scan videos: %w", err)
	}

	/*
		Since we know the maximum amount of entries,
		we can preallocate the slice.

		Since tags can vary per video, we don't know
		how many tags we may end up with.
	*/

	videos := make([]db.Video, 0, len(entries))
	var tags []string

	for _, entry := range entries {
		if entry.IsDir() {
			folderName := entry.Name()
			nfoPath := filepath.Join(collectionPath, folderName, folderName+".nfo")

			/*
				Since we may have multiple videos, some video folders
				may not have a .nfo file or it may parse incorrectly.
				If this is the case, we can just continue in order
				for scanning to continue and not disrupt other videos.
			*/

			if _, err := os.Stat(nfoPath); err != nil {
				continue
			}

			metadata, err := utils.ParseNfoFile(nfoPath)

			if err != nil {
				continue
			}

			videos = append(videos, db.Video{
				Title:  metadata.Title,
				Slug:   utils.TitleToSnake(folderName),
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
