package releasestorage

import "gobot/internal/anime/animefeeder"

type ReleaseStorage interface {
	UpdateStorage([]animefeeder.LatestReleases) (newEntries []animefeeder.LatestReleases)
}
