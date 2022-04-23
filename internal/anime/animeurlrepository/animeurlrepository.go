package animeurlrepository

import (
	"gobot/internal/database"
	"gobot/pkg/animeurlservice"
)

type AnimeUrlRepository interface {
	GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error)
	AddAnimeUrls(...animeurlservice.AnimeUrlInfo) error
}

type animeUrlRepository struct {
	db database.Database
}

var _ AnimeUrlRepository = (*animeUrlRepository)(nil)

func NewAnimeUrlRepository(db database.Database) AnimeUrlRepository {
	return &animeUrlRepository{
		db: db,
	}
}

func (ar *animeUrlRepository) GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error) {
	return ar.db.GetAnimeUrlByName(name)
}

func (ar *animeUrlRepository) AddAnimeUrls(urls ...animeurlservice.AnimeUrlInfo) error {
	return ar.db.AddAnimeUrls(urls...)
}
