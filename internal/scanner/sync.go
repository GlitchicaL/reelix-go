package scanner

import (
	"fmt"
	"log"

	"reelix-go/internal/db"
)

func SyncVaults(vaults []db.Vault) ([]db.Vault, error) {
	// We use make() here because at this point we know the size of
	// the vault and we won't need to reallocate memory if we were
	// to just loop and append.
	names := make([]string, len(vaults))

	for i, v := range vaults {
		names[i] = v.Name
	}

	dbVaults, err := db.CreateVaults(names)

	if err != nil {
		return nil, fmt.Errorf("db vault insert error: %v", err)
	}

	return dbVaults, nil
}

func SyncGalleries(galleries []db.Gallery) error {
	titles := make([]string, len(galleries))
	slugs := make([]string, len(galleries))
	imageCounts := make([]int, len(galleries))
	vaultIds := make([]int, len(galleries))

	for i, g := range galleries {
		titles[i] = g.Title
		slugs[i] = g.Slug
		imageCounts[i] = g.ImageCount
		vaultIds[i] = g.VaultID
	}

	_, err := db.CreateGallery(titles, slugs, imageCounts, vaultIds)

	if err != nil {
		return fmt.Errorf("db gallery insert error: %v", err)
	}

	return nil
}

func SyncCollections(collections []db.Collection) ([]db.Collection, error) {
	names := make([]string, len(collections))
	paths := make([]string, len(collections))
	vaultIds := make([]int, len(collections))

	for i, c := range collections {
		names[i] = c.Name
		paths[i] = c.Path
		vaultIds[i] = c.VaultID
	}

	dbCollections, err := db.CreateCollections(names, paths, vaultIds)

	if err != nil {
		return nil, fmt.Errorf("db vault insert error: %v", err)
	}

	return dbCollections, nil
}

func SyncVideos(videos []db.Video) error {
	for _, v := range videos {
		err := db.CreateVideo(v)

		if err != nil {
			return fmt.Errorf("db sync insert error: %v", err)
		}

		log.Printf("synced video: %v (collection: %v)", v.Title, v.CollectionID)
	}

	return nil
}

func SyncActors(actors []db.Actor) error {
	for _, a := range actors {
		_, err := db.CreateActor(a)

		if err != nil {
			return fmt.Errorf("db sync insert error: %v", err)
		}

		log.Printf("synced actor: %v", a.Name)
	}

	return nil
}
