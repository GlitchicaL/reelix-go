package scanner

import (
	"fmt"
	"log"

	"reelix-go/internal/db"
)

func SyncVaults(vaults []db.Vault) error {
	for _, v := range vaults {
		err := db.CreateVault(v)

		if err != nil {
			return fmt.Errorf("db vault insert error: %v", err)
		}

		log.Printf("synced vault: %v", v.Name)
	}

	return nil
}

func SyncCollections(collections []db.Collection) error {
	for _, c := range collections {
		err := db.CreateCollection(c)

		if err != nil {
			return fmt.Errorf("db collection insert error: %v", err)
		}

		log.Printf("synced collection: %v (vault: %v)", c.Name, c.VaultID)
	}

	return nil
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
