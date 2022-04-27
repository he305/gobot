package animesubsrepository

import (
	"encoding/json"
	"gobot/internal/database"
	"gobot/internal/database/entities"
	"gobot/pkg/animesubs"
	"time"

	"github.com/mitchellh/mapstructure"
)

type AnimeSubsRepository interface {
	GetAnimeSubsByName(name string) (animesubs.SubsInfo, error)
	AddAnimeSubs(...animesubs.SubsInfo) error
}

type animeSubsRepository struct {
	db             database.Database
	collectionName string
}

var _ AnimeSubsRepository = (*animeSubsRepository)(nil)

func NewAnimeSubsRepository(db database.Database, collectionName string) *animeSubsRepository {
	return &animeSubsRepository{
		db:             db,
		collectionName: collectionName,
	}
}

func (ar *animeSubsRepository) GetAnimeSubsByName(name string) (animesubs.SubsInfo, error) {
	data, err := ar.db.GetEntryByName(ar.collectionName, "Title", name)
	if err != nil {
		return animesubs.SubsInfo{}, err
	}

	var result entities.AnimeData
	err = mapstructure.Decode(data, &result)
	if err != nil {
		return animesubs.SubsInfo{}, err
	}

	return animesubs.SubsInfo{
		Title:       result.Title,
		Url:         result.Url,
		TimeUpdated: time.Unix(result.Time, 0),
	}, nil
}

func (ar *animeSubsRepository) AddAnimeSubs(urls ...animesubs.SubsInfo) error {
	if len(urls) == 0 {
		return nil
	}

	animeDataEntries := make([]entities.AnimeData, 0, len(urls))
	for _, url := range urls {
		animeDataEntries = append(animeDataEntries, formFileAnimeDataFromAnimeSubs(url))
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

func formFileAnimeDataFromAnimeSubs(entry animesubs.SubsInfo) entities.AnimeData {
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
