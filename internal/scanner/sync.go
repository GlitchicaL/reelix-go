package scanner

import (
	"fmt"
	"log"

	"reelix-go/internal/db"
)

func Sync(world World) error {
	for _, v := range world.Vaults {
		dbVaults, err := SyncVaults([]db.Vault{v.Vault})
		if err != nil {
			log.Println("vault sync error:", err)
			continue
		}

		vaultID := dbVaults[0].ID

		if err := SyncActors(v.Actors); err != nil {
			log.Println("actor sync error:", err)
		}

		for i := range v.Galleries {
			v.Galleries[i].VaultID = vaultID
		}
		if err := SyncGalleries(v.Galleries); err != nil {
			log.Println("gallery sync error:", err)
		}

		var collectionsToSync []db.Collection
		for _, c := range v.Collections {
			c.Collection.VaultID = vaultID
			collectionsToSync = append(collectionsToSync, c.Collection)
		}

		dbCollections, err := SyncCollections(collectionsToSync)
		if err != nil {
			log.Println("collection sync error:", err)
			continue
		}

		// Since the name of a collection is unique we can
		// map the name to its ID

		collectionMap := map[string]int{}

		for _, c := range dbCollections {
			collectionMap[c.Name] = c.ID
		}

		for _, c := range v.Collections {

			collectionID := collectionMap[c.Collection.Name]

			for i := range c.Videos {
				c.Videos[i].CollectionID = collectionID
			}

			if err := SyncVideos(c.Videos); err != nil {
				log.Println("video sync error:", err)
			}
		}
	}

	return nil
}

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
