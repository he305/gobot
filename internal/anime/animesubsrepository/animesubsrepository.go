package animesubsrepository

import (
	"gobot/internal/database"
	"gobot/pkg/animesubs"
)

type AnimeSubsRepository interface {
	GetAnimeSubsByName(name string) (animesubs.SubsInfo, error)
	AddAnimeSubs(...animesubs.SubsInfo) error
}

type animeSubsRepository struct {
	db database.Database
}

var _ AnimeSubsRepository = (*animeSubsRepository)(nil)

func NewAnimeSubsRepository(db database.Database) AnimeSubsRepository {
	return &animeSubsRepository{
		db: db,
	}
}

func (ar *animeSubsRepository) GetAnimeSubsByName(name string) (animesubs.SubsInfo, error) {
	return ar.db.GetAnimeSubByName(name)
}

func (ar *animeSubsRepository) AddAnimeSubs(urls ...animesubs.SubsInfo) error {
	return ar.db.AddSubs(urls...)
}
