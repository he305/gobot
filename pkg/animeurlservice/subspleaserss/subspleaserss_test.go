package subspleaserss

import (
	"gobot/pkg/animeurlservice"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func newTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:subsplease="https://subsplease.org/rss"><channel><title>SubsPlease RSS</title><description>RSS feed for SubsPlease releases (1080p)</description><link>https://subsplease.org</link><atom:link href="http://subsplease.org/rss" rel="self" type="application/rss+xml"/><item><title>[SubsPlease] Hanyou no Yashahime (25-48) (1080p) [Batch]</title><link>https://nyaa.si/view/1511823/torrent</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri, 08 Apr 2022 08:47:24 +0000</pubDate><category>Hanyou no Yashahime - 1080</category><subsplease:size>32.95 GiB</subsplease:size></item><item><title>[SubsPlease] Shuumatsu no Harem (01-11) (1080p) [Batch]</title><link>https://nyaa.si/view/1511812/torrent</link><guid isPermaLink="false">VFJUCN7IVULFF4YMRZCKHCLABCM4DANA</guid><pubDate>Fri, 08 Apr 2022 07:04:57 +0000</pubDate><category>Shuumatsu no Harem - 1080</category><subsplease:size>14.79 GiB</subsplease:size></item><item><title>[SubsPlease] Girls&apos; Frontline (01-12) (1080p) [Batch]</title><link>https://nyaa.si/view/1511759/torrent</link><guid isPermaLink="false">ZSUMK22DN7YBEFG6UE4M7SG3OOJFJMZL</guid><pubDate>Fri, 08 Apr 2022 01:20:40 +0000</pubDate><category>Girls' Frontline - 1080</category><subsplease:size>13.08 GiB</subsplease:size></item></channel></rss>`))
	})

	mux.HandleFunc("/brokentime", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:subsplease="https://subsplease.org/rss"><channel><title>SubsPlease RSS</title><description>RSS feed for SubsPlease releases (1080p)</description><link>https://subsplease.org</link><atom:link href="http://subsplease.org/rss" rel="self" type="application/rss+xml"/><item><title>[SubsPlease] Hanyou no Yashahime (25-48) (1080p) [Batch]</title><link>https://nyaa.si/view/1511823/torrent</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri, 08 Apr 2022 08:47:24</pubDate><category>Hanyou no Yashahime - 1080</category><subsplease:size>32.95 GiB</subsplease:size></item></channel></rss>`))
	})

	mux.HandleFunc("/latest", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
		<rss version="2.0" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:subsplease="https://subsplease.org/rss"><channel><title>SubsPlease RSS</title><description>RSS feed for SubsPlease releases (1080p)</description><link>https://subsplease.org</link><atom:link href="http://subsplease.org/rss" rel="self" type="application/rss+xml"/>
		<item><title>Hanyou no Yashahime 1</title><link>https://1</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri, 05 Apr 2022 08:47:24 +0000</pubDate><category>Hanyou no Yashahime - 1080</category><subsplease:size>32.95 GiB</subsplease:size></item>
		<item><title>Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri, 07 Apr 2022 08:47:24 +0000</pubDate><category>Hanyou no Yashahime - 1080</category><subsplease:size>32.95 GiB</subsplease:size></item>
		<item><title>Hanyou no Yashahime 2</title><link>https://2</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri, 06 Apr 2022 08:47:24 +0000</pubDate><category>Hanyou no Yashahime - 1080</category><subsplease:size>32.95 GiB</subsplease:size></item></channel></rss>`))
	})

	return httptest.NewServer(mux)
}

func newTestPrepareScrapper() (*subspleaserss, *httptest.Server) {
	server := newTestServer()
	return &subspleaserss{
		parser:      gofeed.NewParser(),
		logger:      zap.L().Sugar(),
		cachedFeed:  nil,
		updateTimer: 0,
		lastUpdated: time.Unix(0, 0),
		feedUrl:     server.URL}, server
}

func TestNormalizeRssTitles(t *testing.T) {
	data := []string{
		"[SubsPlease] Yami Shibai 10 - 13 (1080p) [A0A563BF].mkv",
		"[SubsPlease] Magia Record Final Season (09-12) (1080p) [Batch]",
		"[SubsPlease] Baraou no Souretsu - 12.5 (1080p) [19D386FA].mkv",
		"[SubsPlease] Gaikotsu Kishi-sama, Tadaima Isekai e Odekakechuu - 03 (1080p) [14C1B0BA].mkv",
	}

	expected := []string{
		"yami shibai 10 - 13",
		"magia record final season (09-12)",
		"baraou no souretsu - 12.5",
		"gaikotsu kishi-sama, tadaima isekai e odekakechuu - 03",
	}

	var actual []string
	for _, d := range data {
		actual = append(actual, normalizeRssTitle(d))
	}

	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], actual[i])
		}
	}
}

func TestSubsPleaseTimeParse(t *testing.T) {
	data := "Tue, 05 Apr 2022 04:40:55 +0000"

	expected := time.Date(2022, time.April, 5, 4, 40, 55, 0, time.UTC)
	actual, err := parseSubsPleaseTime(data)
	if err != nil {
		t.Errorf("expected %v, got %v", expected, err.Error())
	}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestIsRssMathingTitlesTrueSimpleMatching(t *testing.T) {
	dataRss := "Shingeki no kyoujin"
	dataTitle := "sHiNgEkI"

	expected := true
	actual := isRssMatchingTitle(dataRss, dataTitle)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestIsRssMathingTitlesTrueLevenshtein(t *testing.T) {
	dataRss := "shingeki+no+kyoujin"
	dataTitle := "shingeki no kyoujin"

	expected := true
	actual := isRssMatchingTitle(dataRss, dataTitle)

	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestUpdateFeedTooSoon(t *testing.T) {
	sp, _ := newTestPrepareScrapper()

	sp.lastUpdated = time.Now().Add(3 * time.Hour)

	err := sp.updateFeed()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestUpdateFeedBrokenUrl(t *testing.T) {
	sp, _ := newTestPrepareScrapper()

	sp.feedUrl = "/"

	err := sp.updateFeed()
	if err == nil {
		t.Errorf("expected error, got none")
	}
}

func TestUpdateFeedOk(t *testing.T) {
	sp, _ := newTestPrepareScrapper()

	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetAllEntriesOk(t *testing.T) {
	sp, _ := newTestPrepareScrapper()

	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	allEntries := sp.getAllRssEntries()
	if len(allEntries) != 3 {
		t.Errorf("expected 3 entries, got %v", len(allEntries))
	}
}

func TestGetAllEntriesBrokenTime(t *testing.T) {
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/brokentime"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	allEntries := sp.getAllRssEntries()
	if len(allEntries) != 0 {
		t.Errorf("expected 0 entries, got %v", len(allEntries))
	}
}

func TestGetNormalizedEntries(t *testing.T) {
	sp, _ := newTestPrepareScrapper()

	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expected := []string{
		"hanyou no yashahime (25-48)",
		"shuumatsu no harem (01-11)",
		"girls' frontline (01-12)",
	}

	allEntries := sp.getNormalizedRssEntries()
	if len(allEntries) != 3 {
		t.Errorf("expected 3 entries, got %v", len(allEntries))
	}

	for i := range allEntries {
		if allEntries[i].Text != expected[i] {
			t.Errorf("expected %v, got %v", allEntries[i], expected)
		}
	}
}

func TestGetFilteredEntries(t *testing.T) {
	sp, _ := newTestPrepareScrapper()

	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//  Girls&apos; Frontline (01-12) (1080p) [Batch]</title><link>https://nyaa.si/view/1511759/torrent</link><guid isPermaLink="false">ZSUMK22DN7YBEFG6UE4M7SG3OOJFJMZL</guid><pubDate>Fri, 08 Apr 2022 01:20:40 +0000
	expected := subsPleaseRssEntry{
		Text:        "girls' frontline (01-12)",
		TimeUpdated: time.Date(2022, time.April, 8, 1, 20, 40, 0, time.UTC), //08 Apr 2022 01:20:40 +0000
		Url:         "https://nyaa.si/view/1511759/torrent",
		RealText:    "[SubsPlease] Girls' Frontline (01-12) (1080p) [Batch]",
	}
	data := "girls' frontline"

	allEntries := sp.getNormalizedRssEntries()

	actual := sp.filterEntriesByTitles(allEntries, data)

	if len(actual) != 1 {
		t.Errorf("expected 1 entry, got %v", len(allEntries))
	}

	actualEn := actual[0]

	if expected.Text != actualEn.Text ||
		expected.RealText != actualEn.RealText ||
		!expected.TimeUpdated.Equal(actualEn.TimeUpdated) ||
		expected.Url != actualEn.Url {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestFindLatest(t *testing.T) {
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/latest"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24
	expected := subsPleaseRssEntry{
		Text:        "hanyou no yashahime 3",
		TimeUpdated: time.Date(2022, time.April, 7, 8, 47, 24, 0, time.UTC), //08 Apr 2022 01:20:40 +0000
		Url:         "https://3",
		RealText:    "Hanyou no Yashahime 3",
	}
	data := "hanyou no yashahime"

	allEntries := sp.getNormalizedRssEntries()

	filtered := sp.filterEntriesByTitles(allEntries, data)

	if len(filtered) != 3 {
		t.Errorf("expected 3 entry, got %v", len(allEntries))
	}

	actualEn := findLatestPageEntry(filtered)

	if expected.Text != actualEn.Text ||
		expected.RealText != actualEn.RealText ||
		!expected.TimeUpdated.Equal(actualEn.TimeUpdated) ||
		expected.Url != actualEn.Url {
		t.Errorf("expected %v, got %v", expected, actualEn)
	}
}

func TestFindLatestEmpty(t *testing.T) {
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/latest"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24
	data := "some irellevant data"

	allEntries := sp.getNormalizedRssEntries()

	filtered := sp.filterEntriesByTitles(allEntries, data)

	if len(filtered) != 0 {
		t.Errorf("expected 0 entry, got %v", len(allEntries))
	}

	actual := findLatestPageEntry(filtered)
	if actual.Text != "" ||
		actual.RealText != "" ||
		actual.Url != "" {
		t.Errorf("expected %v, got %v", nil, actual)
	}
}

func TestGetLatestReleaseValid(t *testing.T) {
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/latest"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24
	expected := animeurlservice.AnimeUrlInfo{
		TimeUpdated: time.Date(2022, time.April, 7, 8, 47, 24, 0, time.UTC), //08 Apr 2022 01:20:40 +0000
		Url:         "https://3",
		Title:       "Hanyou no Yashahime 3",
	}
	data := "hanyou no yashahime"

	actualEn := sp.GetLatestUrlForTitle(data)

	if expected.Title != actualEn.Title ||
		!expected.TimeUpdated.Equal(actualEn.TimeUpdated) ||
		expected.Url != actualEn.Url {
		t.Errorf("expected %v, got %v", expected, actualEn)
	}
}

func TestGetLatestReleaseEmptyInputs(t *testing.T) {
	assert := assert.New(t)
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/latest"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24
	expected := animeurlservice.AnimeUrlInfo{}

	actualEn := sp.GetLatestUrlForTitle(nil...)

	assert.Equal(expected.Title, actualEn.Title)
	assert.Equal(expected.Url, actualEn.Url)
}

func TestGetLatestReleaseEmptyStringInput(t *testing.T) {
	assert := assert.New(t)
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/latest"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24
	expected := animeurlservice.AnimeUrlInfo{}

	actualEn := sp.GetLatestUrlForTitle("")

	assert.Equal(expected.Title, actualEn.Title)
	assert.Equal(expected.Url, actualEn.Url)
}

func TestGetLatestReleaseNoFiltered(t *testing.T) {
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = sp.feedUrl + "/latest"
	err := sp.updateFeed()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24
	expected := animeurlservice.AnimeUrlInfo{}
	data := "some shit"

	actualEn := sp.GetLatestUrlForTitle(data)

	if expected.Title != actualEn.Title ||
		!expected.TimeUpdated.Equal(actualEn.TimeUpdated) ||
		expected.Url != actualEn.Url {
		t.Errorf("expected %v, got %v", expected, actualEn)
	}
}

func TestGetLatestReleaseBrokenUrl(t *testing.T) {
	sp, _ := newTestPrepareScrapper()
	sp.feedUrl = "/"
	//Hanyou no Yashahime 3</title><link>https://3</link><guid isPermaLink="false">YC7Q3L2TRKWLIIWKZ57VI3ZKE73IHMS5</guid><pubDate>Fri,
	//07 Apr 2022 08:47:24

	data := "hanyou no yashahime"

	actualEn := sp.GetLatestUrlForTitle(data)

	if "" != actualEn.Title ||
		"" != actualEn.Url {
		t.Errorf("expected %v, got %v", nil, actualEn)
	}
}

func TestConstructor(t *testing.T) {
	url := "/test"
	duration := time.Duration(3 * time.Minute)
	NewSubsPleaseRss(url, duration, zap.L().Sugar())
}

func TestIsMathingRss(t *testing.T) {
	assert := assert.New(t)

	titles := []string{
		"Re:Zero kara Hajimeru Isekai Seikatsu 2nd Season",
		"Re: Life in a different world from zero 2nd Season",
		"ReZero 2nd Season",
		"Re:Zero - Starting Life in Another World 2",
		"Re:ZERO -Starting Life in Another World- Season 2",
		"Re：ゼロから始める異世界生活",
	}

	rssTitle := "gaikotsu kishi-sama, tadaima isekai e odekakechuu - 03"

	for _, title := range titles {
		assert.False(isRssMatchingTitle(rssTitle, title))
	}
}
