package animeurlrepository

import (
	"encoding/json"
	"fmt"
	"gobot/internal/database"
	"gobot/internal/database/entities"
	"gobot/pkg/animeurlservice"
	"time"

	"github.com/mitchellh/mapstructure"
)

type AnimeUrlRepository interface {
	GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error)
	AddAnimeUrls(...animeurlservice.AnimeUrlInfo) error
}

type animeUrlRepository struct {
	db             database.Database
	collectionName string
}

var _ AnimeUrlRepository = (*animeUrlRepository)(nil)

func NewAnimeUrlRepository(db database.Database, collectionName string) *animeUrlRepository {
	return &animeUrlRepository{
		db:             db,
		collectionName: collectionName,
	}
}

func (ar *animeUrlRepository) GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error) {
	data, err := ar.db.GetEntryByName(ar.collectionName, "Title", name)
	if err != nil {
		return animeurlservice.AnimeUrlInfo{}, err
	}

	var result entities.AnimeData
	err = mapstructure.Decode(data, &result)
	if err != nil {
		return animeurlservice.AnimeUrlInfo{}, err
	}

	return animeurlservice.AnimeUrlInfo{
		Title:       result.Title,
		Url:         result.Url,
		TimeUpdated: time.Unix(result.Time, 0),
	}, nil
}

func (ar *animeUrlRepository) AddAnimeUrls(urls ...animeurlservice.AnimeUrlInfo) error {
	if len(urls) == 0 {
		return nil
	}

	animeDataEntries := make([]entities.AnimeData, 0, len(urls))
	fmt.Println(len(animeDataEntries))
	for _, url := range urls {
		animeDataEntries = append(animeDataEntries, formFileAnimeDataFromAnimeUrl(url))
	}

	entries := make([]map[string]interface{}, 0, len(animeDataEntries))
	for _, animeDataEntry := range animeDataEntries {
		out, err := formJsonFromAnimeData(animeDataEntry)
		if err != nil {
			return err
		}
		entries = append(entries, out)
	}

	return ar.db.AddEntries(ar.collectionName, entries)
}

func formFileAnimeDataFromAnimeUrl(entry animeurlservice.AnimeUrlInfo) entities.AnimeData {
	return entities.AnimeData{
		Title: entry.Title,
		Url:   entry.Url,
		Time:  entry.TimeUpdated.Unix(),
	}
}

func formJsonFromAnimeData(entry entities.AnimeData) (map[string]interface{}, error) {
	var out map[string]interface{}
	bytesMarshalled, err := json.Marshal(entry)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytesMarshalled, &out)
	return out, err
}
