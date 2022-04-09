package kitsunekko

import (
	"fmt"
	"gobot/pkg/animesubs"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"
)

type mockfileio struct{}

func (f *mockfileio) SaveResponseToFile(data *http.Response, filePath string) error {
	return nil
}

type mockfileioerror struct{}

func (f *mockfileioerror) SaveResponseToFile(data *http.Response, filePath string) error {
	return fmt.Errorf("")
}

func newTestSever() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/mainpage", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div id="maindiv">

<br/>
<div id="breadcrumbs"> <a href="/">kitsunekko.net</a> 
&gt; <a href="/dirlist.php?dir=subtitles%2Fjapanese%2F">Japanese subtitles</a></div>
<br/>
<div id="listingheader">
Sort by: <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=name&amp;order=desc'>File↑</a> <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=date&amp;order=desc'>Date</a> </div>
<br/>
<table id="flisttable" cellpadding="0" cellspacing="0"><tbody>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki%21+Kyojin+Chuugakkou%2F" class=""><strong>Shingeki! Kyojin Chuugakkou</strong> </a></td> <td class="tdright" title="Aug 18 2017 12:51:28 AM" > 4&nbsp;years </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Bahamut+Genesis%2F" class=""><strong>Shingeki no Bahamut Genesis</strong> </a></td> <td class="tdright" title="Apr 26 2019 11:54:03 PM" > 2&nbsp;years </td></tr>
<tr><td colspan="2"><a href="/titosi" class=""><strong>Shingeki No Kyojin</strong> </a></td> <td class="tdright" title="Apr 03 2022 06:38:03 PM" > 4&nbsp;days </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Kyojin%2F" class=""><strong>Shingeki no Kyojin</strong> </a></td> <td class="tdright" title="Feb 05 2022 02:52:36 AM" > 2&nbsp;months </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Kyojin+OVA%2F" class=""><strong>Shingeki no Kyojin OVA</strong> </a></td> <td class="tdright" title="Jun 21 2019 10:08:59 PM" > 2&nbsp;years </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki_No_Bahamut_Virgin_Soul%2F" class=""><strong>Shingeki No Bahamut Virgin Soul</strong> </a></td> <td class="tdright" title="Apr 27 2019 12:40:37 AM" > 2&nbsp;years </td></tr>
</tbody></table>
</div>
</body>
</html>
`))
	})

	mux.HandleFunc("/mainpagewithouttdright", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div id="maindiv">

<br/>
<div id="breadcrumbs"> <a href="/">kitsunekko.net</a> 
&gt; <a href="/dirlist.php?dir=subtitles%2Fjapanese%2F">Japanese subtitles</a></div>
<br/>
<div id="listingheader">
Sort by: <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=name&amp;order=desc'>File↑</a> <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=date&amp;order=desc'>Date</a> </div>
<br/>
<table id="flisttable" cellpadding="0" cellspacing="0"><tbody>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki%21+Kyojin+Chuugakkou%2F" class=""><strong>Shingeki! Kyojin Chuugakkou</strong> </a></td> </tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Bahamut+Genesis%2F" class=""><strong>Shingeki no Bahamut Genesis</strong> </a></td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki_No_Kyojin%2F" class=""><strong>Shingeki No Kyojin</strong> </a></td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Kyojin%2F" class=""><strong>Shingeki no Kyojin</strong> </a></td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Kyojin+OVA%2F" class=""><strong>Shingeki no Kyojin OVA</strong> </a></td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki_No_Bahamut_Virgin_Soul%2F" class=""><strong>Shingeki No Bahamut Virgin Soul</strong> </a></td> </tr>
</tbody></table>
</div>
</body>
</html>
`))
	})

	mux.HandleFunc("/mainpagebrokentime", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div id="maindiv">

<br/>
<div id="breadcrumbs"> <a href="/">kitsunekko.net</a> 
&gt; <a href="/dirlist.php?dir=subtitles%2Fjapanese%2F">Japanese subtitles</a></div>
<br/>
<div id="listingheader">
Sort by: <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=name&amp;order=desc'>File↑</a> <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=date&amp;order=desc'>Date</a> </div>
<br/>
<table id="flisttable" cellpadding="0" cellspacing="0"><tbody>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki%21+Kyojin+Chuugakkou%2F" class=""><strong>Shingeki! Kyojin Chuugakkou</strong> </a></td> <td class="tdright" title="Aug 18 2017 12:51:28" > 4&nbsp;years </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Bahamut+Genesis%2F" class=""><strong>Shingeki no Bahamut Genesis</strong> </a></td> <td class="tdright" title="Apr 26 2019 11:54:03" > 2&nbsp;years </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki_No_Kyojin%2F" class=""><strong>Shingeki No Kyojin</strong> </a></td> <td class="tdright" title="Apr 03 2022 06:38:03" > 4&nbsp;days </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Kyojin%2F" class=""><strong>Shingeki no Kyojin</strong> </a></td> <td class="tdright" title="Feb 05 2022 02:52:36" > 2&nbsp;months </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Kyojin+OVA%2F" class=""><strong>Shingeki no Kyojin OVA</strong> </a></td> <td class="tdright" title="Jun 21 2019 10:08:59" > 2&nbsp;years </td></tr>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki_No_Bahamut_Virgin_Soul%2F" class=""><strong>Shingeki No Bahamut Virgin Soul</strong> </a></td> <td class="tdright" title="Apr 27 2019 12:40:37" > 2&nbsp;years </td></tr>
</tbody></table>
</div>
</body>
</html>
`))
	})

	mux.HandleFunc("/titosi", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div id="maindiv">
	<br/>
<div id="breadcrumbs"> <a href="/">kitsunekko.net</a> 
	&gt; <a href="/dirlist.php?dir=subtitles%2Fjapanese%2F">Japanese subtitles</a> &gt; <a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki_No_Kyojin%2F">Shingeki No Kyojin</a></div>
<br/>
<div id="listingheader">
	Sort by: <a href='/dirlist.php?dir=subtitles/japanese/Shingeki_No_Kyojin/&amp;sort=name&amp;order=desc'>File↑</a> <a href='/dirlist.php?dir=subtitles/japanese/Shingeki_No_Kyojin/&amp;sort=size&amp;order=asc'>Size</a> <a href='/dirlist.php?dir=subtitles/japanese/Shingeki_No_Kyojin/&amp;sort=date&amp;order=desc'>Date</a> </div>
<br/>
<table id="flisttable" cellpadding="0" cellspacing="0"><tbody>
	<tr><td><a href="subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_094.srt" class=""><strong>Shingeki No Kyojin 094.srt</strong> </a></td> <td class="tdleft"  title="30212"  > 30&nbsp;KB </td> <td class="tdright" title="Mar 06 2022 06:53:04 PM" > 1&nbsp;month </td></tr>	
	<tr><td><a href="subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_097.srt" class=""><strong>Shingeki No Kyojin 097.srt</strong> </a></td> <td class="tdleft"  title="25410"  > 25&nbsp;KB </td> <td class="tdright" title="Apr 03 2022 06:38:04 PM" > 4&nbsp;days </td></tr>		
	<tr><td><a href="subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_095.srt" class=""><strong>Shingeki No Kyojin 095.srt</strong> </a></td> <td class="tdleft"  title="27140"  > 27&nbsp;KB </td> <td class="tdright" title="Mar 13 2022 06:38:04 PM" > 3&nbsp;weeks </td></tr>	
	<tr><td><a href="subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_096.srt" class=""><strong>Shingeki No Kyojin 096.srt</strong> </a></td> <td class="tdleft"  title="24101"  > 24&nbsp;KB </td> <td class="tdright" title="Mar 20 2022 06:38:05 PM" > 2&nbsp;weeks </td></tr>	
</tbody></table>
</div>
</body>
</html>
		`))
	})

	mux.HandleFunc("/mainpageempty", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
</head>
<body>
<div id="maindiv">

<br/>
<div id="breadcrumbs"> <a href="/">kitsunekko.net</a> 
&gt; <a href="/dirlist.php?dir=subtitles%2Fjapanese%2F">Japanese subtitles</a></div>
<br/>
<div id="listingheader">
Sort by: <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=name&amp;order=desc'>File↑</a> <a href='/dirlist.php?dir=subtitles/japanese/&amp;sort=date&amp;order=desc'>Date</a> </div>
<br/>
<table id="flisttable" cellpadding="0" cellspacing="0"><tbody>
<tr><td colspan="2"><a href="/dirlist.php?dir=subtitles%2Fjapanese%2FShingeki+no+Bahamut+Genesis%2F" class=""><strong>Shingeki no Bahamut Genesis</strong> </a></td> <td class="tdright" title="Apr 26 2019 11:54:03 PM" > 2&nbsp;years </td></tr>
</tbody></table>
</div>
</body>
</html>
`))
	})

	return httptest.NewServer(mux)
}

func newTestPrepareScrapper() (*kitsunekkoScrapper, *httptest.Server) {
	serv := newTestSever()
	return &kitsunekkoScrapper{logger: zap.L().Sugar(),
		collector:       configureKitsunekkoCollyCollector(),
		updateTimer:     0,
		cachedFilePath:  serv.URL + "/mainpage",
		lastTimeUpdated: time.Unix(0, 0),
		cachedVisitUrl:  serv.URL + "/mainpage",
		fileIo:          &mockfileio{}}, serv
}

func TestGetLatestEntryFilterValid(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	// 2022-04-03 18:38:03 +0000 UTC
	actual := ws.getLatestEntry(server.URL+"/mainpage", "shingeki no kyoujin")
	expected := pageEntry{
		Url:         server.URL + "/titosi",
		Text:        "Shingeki No Kyojin",
		TimeUpdated: time.Date(2022, time.April, 3, 18, 38, 3, 0, time.UTC),
	}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestEntryFilterWithoutTDRight(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	// 2022-04-03 18:38:03 +0000 UTC
	actual := ws.getLatestEntry(server.URL+"/mainpagewithouttdright", "shingeki no kyoujin")
	expected := pageEntry{}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}

}

func TestGetLatestEntryFilterBrokenTime(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	// 2022-04-03 18:38:03 +0000 UTC
	actual := ws.getLatestEntry(server.URL+"/mainpagebrokentime", "shingeki no kyoujin")
	expected := pageEntry{}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestEntryFilterBrokenUrl(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	// 2022-04-03 18:38:03 +0000 UTC
	actual := ws.getLatestEntry(server.URL+"/crap", "shingeki no kyoujin")
	expected := pageEntry{}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestEntry(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	// 2022-04-03 18:38:03 +0000 UTC
	actual := ws.getLatestEntry(server.URL + "/mainpage")
	expected := pageEntry{
		Url:         server.URL + "/titosi",
		Text:        "Shingeki No Kyojin",
		TimeUpdated: time.Date(2022, time.April, 3, 18, 38, 3, 0, time.UTC),
	}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestEntryAnimePage(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	// 2022-04-03 18:38:03 +0000 UTC
	actual := ws.getLatestEntry(server.URL + "/titosi")
	expected := pageEntry{
		Url:         server.URL + "/subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_097.srt",
		Text:        "Shingeki No Kyojin 097.srt",
		TimeUpdated: time.Date(2022, time.April, 3, 18, 38, 4, 0, time.UTC),
	}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestSubForAnimeValid(t *testing.T) {
	ws, server := newTestPrepareScrapper()
	ws.lastTimeUpdated = time.Now().Add(1 * time.Hour)
	kitsunekkoBaseUrl = ""
	kitsunekkoJapBaseUrl = ws.cachedVisitUrl
	actual := ws.GetUrlLatestSubForAnime([]string{"shingeki no kyoujin"})

	expected := animesubs.SubsInfo{
		Url:         server.URL + "/subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_097.srt",
		Title:       "Shingeki No Kyojin 097.srt",
		TimeUpdated: time.Date(2022, time.April, 3, 18, 38, 4, 0, time.UTC),
	}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestSubForAnimeNotFoundEntries(t *testing.T) {
	ws, serv := newTestPrepareScrapper()
	ws.cachedVisitUrl = serv.URL + "/mainpageempty"
	ws.lastTimeUpdated = time.Now().Add(1 * time.Hour)
	kitsunekkoBaseUrl = ""
	kitsunekkoJapBaseUrl = ws.cachedVisitUrl
	actual := ws.GetUrlLatestSubForAnime([]string{"shingeki no kyoujin"})

	expected := animesubs.SubsInfo{}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestSubForAnimeFailedCacheButContinue(t *testing.T) {
	ws, serv := newTestPrepareScrapper()
	ws.cachedVisitUrl = "blabla"
	ws.fileIo = &mockfileioerror{}
	kitsunekkoBaseUrl = ""
	kitsunekkoJapBaseUrl = serv.URL + "/mainpageempty"
	actual := ws.GetUrlLatestSubForAnime([]string{"shingeki no kyoujin"})

	expected := animesubs.SubsInfo{}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestGetLatestSubForAnimeWithMockSave(t *testing.T) {
	ws, serv := newTestPrepareScrapper()
	ws.cachedVisitUrl = serv.URL + "/mainpageempty"
	kitsunekkoBaseUrl = ""
	kitsunekkoJapBaseUrl = ws.cachedVisitUrl
	actual := ws.GetUrlLatestSubForAnime([]string{"shingeki no kyoujin"})

	expected := animesubs.SubsInfo{}

	if !actual.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

func TestUpdateCacheBrokenUrl(t *testing.T) {
	ws, serv := newTestPrepareScrapper()
	ws.cachedVisitUrl = serv.URL + "/mainpageempty"
	kitsunekkoBaseUrl = ""
	kitsunekkoJapBaseUrl = "blablabla"

	err := ws.updateCache()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestUpdateCacheErrorWhileSaving(t *testing.T) {
	ws, serv := newTestPrepareScrapper()
	ws.cachedVisitUrl = serv.URL + "/mainpageempty"
	kitsunekkoBaseUrl = ""
	ws.fileIo = &mockfileioerror{}
	kitsunekkoJapBaseUrl = ws.cachedVisitUrl

	err := ws.updateCache()
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestParseKitsunekkoTimeWithoutError(t *testing.T) {
	data := "Apr 05 2022 02:37:35 PM"
	expected := time.Date(2022, time.April, 5, 14, 37, 35, 0, time.UTC)

	got, err := parseKitsunekkoTime(data)

	if err != nil {
		t.Error("Error parsing kitsunekko time using default time layout")
	}

	if !got.Equal(expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestParseKitsunekkoTimeWithError(t *testing.T) {
	data := "Apr 05 2022 02:37:35"
	expected := time.Time{}

	got, err := parseKitsunekkoTime(data)

	if err == nil {
		t.Error("Expected error while parsing kitsunekko time using default time layout")
	}

	if !got.Equal(expected) {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestFindLatestMainPageEntryEmptyArray(t *testing.T) {
	var data []pageEntry

	expected := pageEntry{}
	got := findLatestPageEntry(data)

	if got != expected {
		t.Errorf("expected: %v, got: %v", expected, got)
	}
}

func TestFindLatestMainPageEntryValid(t *testing.T) {
	data := []pageEntry{
		{
			Text:        "test1",
			TimeUpdated: time.Now(),
			Url:         "test",
		},
		{
			Text:        "test2",
			TimeUpdated: time.Now().Add(1 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test3",
			TimeUpdated: time.Now().Add(2 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test4",
			TimeUpdated: time.Now().Add(3 * time.Hour),
			Url:         "test",
		},
	}

	actual := data[len(data)-1]
	got := findLatestPageEntry(data)

	if !(actual.Text == got.Text &&
		actual.TimeUpdated.Equal(got.TimeUpdated) &&
		actual.Url == got.Url) {
		t.Errorf("expected: %v, got: %v", actual, got)
	}
}

func TestFilterMainPageEntriesByTitlesValid(t *testing.T) {
	entries := []pageEntry{
		{
			Text:        "THIS",
			TimeUpdated: time.Now(),
			Url:         "test",
		},
		{
			Text:        "tHiS",
			TimeUpdated: time.Now().Add(1 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test3",
			TimeUpdated: time.Now().Add(2 * time.Hour),
			Url:         "test",
		},
		{
			Text:        "test4",
			TimeUpdated: time.Now().Add(3 * time.Hour),
			Url:         "test",
		},
	}

	titles := []string{
		"irellevant",
		"THISA",
		"another text",
	}

	expected := []pageEntry{entries[0], entries[1]}
	got := filterPageEntriesByTitles(entries, titles)

	for _, expectedEn := range expected {
		found := false
		for _, gotEn := range got {
			if gotEn.Text == expectedEn.Text {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Entry %v not found in got list", expectedEn)
		}
	}
}

func TestConstructor(t *testing.T) {
	serv := NewKitsunekkoScrapper(&mockfileio{}, "", 1*time.Minute)
	if serv == nil {
		t.Errorf("Returned nil object")
	}
}

func TestCollyCollectorSettings(t *testing.T) {
	collector := configureKitsunekkoCollyCollector()
	if !collector.AllowURLRevisit ||
		len(collector.URLFilters) != 0 {
		t.Errorf("Wrong colly collector configuration")
	}
}
