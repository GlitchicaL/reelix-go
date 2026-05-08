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

		_, err = SyncActors(v.Actors)

		if err != nil {
			log.Println("actor sync error:", err)
		}

		_, err = SyncTags(v.Tags)

		if err != nil {
			log.Println("tags sync error:", err)
		}

		for i := range v.Galleries {
			v.Galleries[i].VaultID = vaultID
		}

		if _, err := SyncGalleries(v.Galleries); err != nil {
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

			_, err := SyncVideos(c.Videos)

			if err != nil {
				log.Println("video sync error:", err)
			}

			if err := SyncVideoTags(c.Videos, v.Tags); err != nil {
				log.Println("video tag sync error:", err)
			}

			if err := SyncVideoActors(c.Videos, v.Actors); err != nil {
				log.Println("video tag sync error:", err)
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

func SyncGalleries(galleries []db.Gallery) ([]db.Gallery, error) {
	dbGalleries, err := db.CreateGallery(galleries)

	if err != nil {
		return nil, fmt.Errorf("db galleries sync error: %v", err)
	}

	return dbGalleries, nil
}

func SyncCollections(collections []db.Collection) ([]db.Collection, error) {
	dbCollections, err := db.CreateCollections(collections)

	if err != nil {
		return nil, fmt.Errorf("db collections sync error: %v", err)
	}

	return dbCollections, nil
}

func SyncVideos(videos []db.Video) ([]db.Video, error) {
	dbVideos, err := db.CreateVideos(videos)

	if err != nil {
		return nil, fmt.Errorf("db videos sync error: %v", err)
	}

	/*
		Now that the videos has been inserted, we can update
		the world state with the associated video ID.

		Since video slugs are unique, we can map it with the
		ID so we know what video to update.
	*/

	videosMap := map[string]int{}

	for _, v := range dbVideos {
		videosMap[v.Slug] = v.ID
	}

	for i := range videos {
		videos[i].ID = videosMap[videos[i].Slug]
	}

	log.Printf("videos sync: %v", len(dbVideos))

	return videos, nil
}

func SyncTags(tags []db.Tag) ([]db.Tag, error) {
	dbTags, err := db.CreateTags(tags)

	if err != nil {
		return nil, fmt.Errorf("db tags sync error: %v", err)
	}

	/*
		Now that the tags has been inserted, we can update
		the world state with the associated tag ID.

		Since tag names are unique, we can map it with the
		ID so we know what tag to update.
	*/

	tagsMap := map[string]int{}

	for _, t := range dbTags {
		tagsMap[t.Name] = t.ID
	}

	for i := range tags {
		tags[i].ID = tagsMap[tags[i].Name]
	}

	log.Printf("tags sync: %v", len(dbTags))

	return tags, nil
}

func SyncVideoTags(videos []db.Video, tags []db.Tag) error {
	err := db.LinkVideoTags(videos, tags)

	if err != nil {
		return fmt.Errorf("db actors sync error: %v", err)
	}

	return nil
}

func SyncActors(actors []db.Actor) ([]db.Actor, error) {
	dbActors, err := db.CreateActors(actors)

	if err != nil {
		return nil, fmt.Errorf("db actors sync error: %v", err)
	}

	/*
		Now that the actors has been inserted, we can update
		the world state with the associated actor ID.

		Since actor names are unique, we can map it with the
		ID so we know what actor to update.
	*/

	actorsMap := map[string]int{}

	for _, a := range dbActors {
		actorsMap[a.Name] = a.ID
	}

	for i := range actors {
		actors[i].ID = actorsMap[actors[i].Name]
	}

	log.Printf("actors sync: %v", len(dbActors))

	return actors, nil
}

func SyncVideoActors(videos []db.Video, actors []db.Actor) error {
	err := db.LinkVideoActors(videos, actors)

	if err != nil {
		return fmt.Errorf("db actors sync error: %v", err)
	}

	return nil
}
