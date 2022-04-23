package animemessageprovider

import (
	"fmt"
	"gobot/internal/anime/animefeeder"
)

type AnimeMessageProvider interface {
	GetMessage() (string, error)
}

type animeMessageProvider struct {
	feeder animefeeder.AnimeFeeder
}

var _ AnimeMessageProvider = (*animeMessageProvider)(nil)

func NewAnimeMessageProvider(feeder animefeeder.AnimeFeeder) AnimeMessageProvider {
	return &animeMessageProvider{
		feeder: feeder,
	}
}

func (amp *animeMessageProvider) formUpdatedListMessage() (string, error) {
	var returnString string
	missingInCached, missingInNew, err := amp.feeder.UpdateList()
	if err != nil {
		return returnString, err
	}

	if len(missingInNew) != 0 {
		returnString = returnString + fmt.Sprintf("%d entries were deleted:\n", len(missingInNew))
		for _, entry := range missingInNew {
			returnString += entry.VerboseOutput() + "\n"
		}
	}

	if len(missingInCached) != 0 {
		returnString = returnString + fmt.Sprintf("%d entries were added:\n", len(missingInCached))
		for _, entry := range missingInCached {
			returnString += entry.VerboseOutput() + "\n"
		}
	}

	return returnString, nil
}

func (amp *animeMessageProvider) formReleasesMessage() string {
	var returnedString string

	releases := amp.feeder.FindLatestReleases()

	for _, release := range releases {
		st := fmt.Sprintf("New release for %v:\n", release.Title)
		if release.AnimeUrl.Title != "" {
			st += fmt.Sprintf("New series torrent:\nTitle: %v, Url: %v\n", release.AnimeUrl.Title, release.AnimeUrl.Url)
		}
		if release.SubsUrl.Title != "" {
			st += fmt.Sprintf("New series subs:\nTitle: %v, Url: %v\n", release.SubsUrl.Title, release.SubsUrl.Url)
		}

		returnedString += st
	}

	return returnedString
}

func (amp *animeMessageProvider) GetMessage() (string, error) {
	var returnedString string

	mess, err := amp.formUpdatedListMessage()
	if err != nil {
		return "", err
	}
	returnedString = returnedString + mess

	mess = amp.formReleasesMessage()
	returnedString = returnedString + mess

	return returnedString, nil
}
