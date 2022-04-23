package kitsunekkov2

import (
	"gobot/pkg/animesubs"
	"gobot/pkg/webpagekeeper"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var data = `
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
`

var brokenDataWithoutUrl = `
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
<tr><td colspan="2"><strong>Shingeki No Kyojin</strong></td> <td class="tdright" title="Apr 03 2022 06:38:03 PM" > 4&nbsp;days </td></tr>
</tbody></table>
</div>
</body>
</html>
`

var brokenDataWithoutDate = `
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
<tr><td colspan="2"><a href="/titosi" class=""><strong>Shingeki No Kyojin</strong></a></td> <td class="tdright"> 4&nbsp;days </td></tr>
</tbody></table>
</div>
</body>
</html>
`

var brokenDataWrongTime = `
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
<tr><td colspan="2"><a href="/titosi" class=""><strong>Shingeki No Kyojin</strong></a></td> <td class="tdright" title="Apr 03 2022 06:38:03" > 4&nbsp;days </td></tr>
</tbody></table>
</div>
</body>
</html>
`

var titosi = `
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
</html>`

const duration = 1 * time.Second

func newTestSever() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(data))
	})
	mux.HandleFunc("/brokenDataWithoutUrl", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(brokenDataWithoutUrl))
	})
	mux.HandleFunc("/brokenDataWithoutDate", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(brokenDataWithoutDate))
	})
	mux.HandleFunc("/brokenDataWrongTime", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(brokenDataWrongTime))
	})
	mux.HandleFunc("/titosi", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(titosi))
	})
	return httptest.NewServer(mux)
}

func prepareForTests() (*kitsunekkoScrapperV2, *httptest.Server) {
	serv := newTestSever()
	kserv := &kitsunekkoScrapperV2{
		logger:     zap.L().Sugar(),
		timeUpdate: duration,
		webkeeper:  webpagekeeper.NewWebPageKeeper(duration, zap.L().Sugar()),
	}

	return kserv, serv
}

func TestGetAllEntriesValid(t *testing.T) {
	assert := assert.New(t)

	kits, serv := prepareForTests()

	body, err := kits.webkeeper.GetUrlBody(serv.URL+"/data", false)
	assert.NoError(err)

	entries, err := kits.getAllEntries(body)
	assert.NoError(err)

	assert.Equal(len(entries), 6)
	parsedTime, err := parseKitsunekkoTime("Apr 03 2022 06:38:03 PM")
	assert.NoError(err)
	expected := pageEntry{
		url:         "/titosi",
		text:        "Shingeki No Kyojin",
		timeUpdated: parsedTime,
	}

	assert.Equal(entries[2], expected)
}

func TestGetLatestEntry(t *testing.T) {
	assert := assert.New(t)
	kits, server := prepareForTests()
	// 2022-04-03 18:38:03 +0000 UTC
	kitsunekkoJapBaseUrl = server.URL + "/data"
	actual, err := kits.getRequiredAnimeUrl([]string{"shingeki no kyoujin"})
	assert.NoError(err)

	expected := "/titosi"

	assert.Equal(actual, expected)
}

func TestGetLatestEntryNotFound(t *testing.T) {
	assert := assert.New(t)
	kits, server := prepareForTests()
	// 2022-04-03 18:38:03 +0000 UTC
	kitsunekkoJapBaseUrl = server.URL + "/data"
	actual, err := kits.getRequiredAnimeUrl([]string{"something else"})
	assert.NoError(err)

	expected := ""

	assert.Equal(actual, expected)
}

func TestGetAllEntriesNoUrl(t *testing.T) {
	assert := assert.New(t)

	kits, serv := prepareForTests()
	body, err := kits.webkeeper.GetUrlBody(serv.URL+"/brokenDataWithoutUrl", false)
	assert.NoError(err)
	entries, err := kits.getAllEntries(body)
	assert.NoError(err)

	assert.Equal(len(entries), 0)
}
func TestGetAllEntriesNoDate(t *testing.T) {
	assert := assert.New(t)

	kits, serv := prepareForTests()
	body, err := kits.webkeeper.GetUrlBody(serv.URL+"/brokenDataWithoutDate", false)
	assert.NoError(err)
	entries, err := kits.getAllEntries(body)
	assert.NoError(err)

	assert.Equal(len(entries), 0)
}
func TestGetAllEntriesNoTime(t *testing.T) {
	assert := assert.New(t)

	kits, serv := prepareForTests()
	body, err := kits.webkeeper.GetUrlBody(serv.URL+"/brokenDataWrongTime", false)
	assert.NoError(err)
	entries, err := kits.getAllEntries(body)
	assert.NoError(err)

	assert.Equal(len(entries), 0)
}

func TestGetLatestAnimeEntry(t *testing.T) {
	assert := assert.New(t)

	kits, serv := prepareForTests()
	actual, err := kits.getLatestAnimeEntry(serv.URL + "/titosi")
	assert.NoError(err)

	rawTime := "Apr 03 2022 06:38:04 PM"
	parsedTime, err := parseKitsunekkoTime(rawTime)
	assert.NoError(err)
	expected := pageEntry{
		text:        "Shingeki No Kyojin 097.srt",
		url:         "subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_097.srt",
		timeUpdated: parsedTime,
	}

	assert.Equal(actual, expected)
}

func TestGetUrlLatestSubForAnimeValid(t *testing.T) {
	assert := assert.New(t)

	kits, serv := prepareForTests()

	kitsunekkoJapBaseUrl = serv.URL + "/data"
	kitsunekkoBaseUrl = serv.URL

	rawTime := "Apr 03 2022 06:38:04 PM"
	parsedTime, err := parseKitsunekkoTime(rawTime)
	assert.NoError(err)

	expected := animesubs.SubsInfo{
		Title:       "Shingeki No Kyojin 097.srt",
		Url:         kitsunekkoBaseUrl + "/subtitles/japanese/Shingeki_No_Kyojin/Shingeki_No_Kyojin_097.srt",
		TimeUpdated: parsedTime,
	}

	actual := kits.GetUrlLatestSubForAnime([]string{"shingeki no kyoujin"})

	assert.Equal(expected, actual)
}
