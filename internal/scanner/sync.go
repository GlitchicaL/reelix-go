package scanner

import (
	"fmt"
	"log"

	"reelix-go/internal/db"
)

func SyncVaults(vaults []db.Vault) ([]db.Vault, error) {
	dbVaults, err := db.CreateVaults(vaults)

	if err != nil {
		return nil, fmt.Errorf("db vaults sync error: %v", err)
	}

	return dbVaults, nil
}

func SyncGalleries(galleries []db.Gallery) error {
	_, err := db.CreateGallery(galleries)

	if err != nil {
		return fmt.Errorf("db galleries sync error: %v", err)
	}

	return nil
}

func SyncCollections(collections []db.Collection) ([]db.Collection, error) {
	dbCollections, err := db.CreateCollections(collections)

	if err != nil {
		return nil, fmt.Errorf("db collections sync error: %v", err)
	}

	return dbCollections, nil
}

func SyncVideos(videos []db.Video) error {
	for _, v := range videos {
		err := db.CreateVideo(v)

		if err != nil {
			return fmt.Errorf("db videos sync error: %v", err)
		}

		log.Printf("synced video: %v (collection: %v)", v.Title, v.CollectionID)
	}

	return nil
}

func SyncActors(actors []db.Actor) error {
	for _, a := range actors {
		_, err := db.CreateActor(a)

		if err != nil {
			return fmt.Errorf("db actors sync error: %v", err)
		}

		log.Printf("synced actor: %v", a.Name)
	}

	return nil
}
