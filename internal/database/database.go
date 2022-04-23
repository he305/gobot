package database

import (
	"gobot/pkg/animesubs"
	"gobot/pkg/animeurlservice"
)

type Database interface {
	GetAnimeUrlByName(name string) (animeurlservice.AnimeUrlInfo, error)
	AddAnimeUrls(...animeurlservice.AnimeUrlInfo) error
	GetAnimeSubByName(name string) (animesubs.SubsInfo, error)
	AddSubs(...animesubs.SubsInfo) error
}
