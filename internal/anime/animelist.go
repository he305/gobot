package anime

import (
	as "gobot/pkg/animeservice"
)

type AnimeList struct {
	list []as.AnimeStruct
}

func NewAnimeList() *AnimeList {
	return &AnimeList{list: make([]as.AnimeStruct, 0)}
}

func (alist *AnimeList) FilterByListStatus(statuses ...uint8) []as.AnimeStruct {
	var filteredList []as.AnimeStruct

	for _, entry := range alist.list {
		for _, status := range statuses {
			if entry.ListStatus == status {
				filteredList = append(filteredList, entry)
			}
		}
	}

	return filteredList
}

func (alist *AnimeList) SetNewList(list []as.AnimeStruct) {
	alist.list = list
}

func findMissingEntriesInFirstList(firstList []as.AnimeStruct, secondList []as.AnimeStruct) []as.AnimeStruct {
	var missing []as.AnimeStruct

	for _, entrySecond := range secondList {
		found := false
		for _, entryFirst := range firstList {
			if entrySecond.Id == entryFirst.Id {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, entrySecond)
		}
	}

	return missing
}

func (alist *AnimeList) FindMissingInBothLists(list []as.AnimeStruct) ([]as.AnimeStruct, []as.AnimeStruct) {
	missingInThisList := findMissingEntriesInFirstList(alist.list, list)
	missingInArg := findMissingEntriesInFirstList(list, alist.list)

	return missingInThisList, missingInArg
}
