package webpagekeeper

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

var data = `
<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<!-- saved from url=(0023)https://kitsunekko.net/ -->
<html><head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8">		
<meta http-equiv="Content-Security-Policy" content="upgrade-insecure-requests">
<title>kitsuneko.net - Japanese and English anime subtitles</title>
<meta name="keywords" content="">
<meta name="description" content="">
<style>
body {font-family: tahoma;}
table {font-size: 10pt;}
hr {color: gray;}
</style>
<!-- Global site tag (gtag.js) - Google Analytics -->
<script type="text/javascript" async="" src="./kitsuneko.net - Japanese and English anime subtitles_files/watch.js.Без названия"></script><script async="" src="./kitsuneko.net - Japanese and English anime subtitles_files/js"></script>
<script>
window.dataLayer = window.dataLayer || [];
function gtag(){dataLayer.push(arguments);}
gtag('js', new Date());	
gtag('config', 'UA-3419527-1');
</script>	
</head>	
<body topmargin="50%" leftmargin="50%" rightmargin="50%" bottommargin="50%" marginheight="0" marginwidth="0" text="black" link="#B10006" alink="red" vlink="maroon">
<table width="80%" cellpadding="0" cellspacing="15" border="0" align="center">
<tbody><tr>
<td align="right"></td>
</tr>
<tr><td>
<table width="100%" cellpadding="0" cellspacing="15" border="0" align="center">
<tbody><tr>	
<td valign="top">
<strong>Download files:</strong>
<p>
<a href="https://kitsunekko.net/dirlist.php?dir=subtitles%2F">English subtitles</a>
</p><p>
<a href="https://kitsunekko.net/dirlist.php?dir=subtitles%2Fjapanese%2F">Japanese subtitles</a>
</p><p>
<a href="https://kitsunekko.net/dirlist.php?dir=subtitles%2Fchinese%2F">Chinese subtitles</a>
</p><p>
<a href="https://kitsunekko.net/dirlist.php?dir=subtitles%2Fkorean%2F">Korean subtitles</a>
</p></td>
<td style="font-size: x-small;" width="60%" valign="top">
<p style="color:red;">
Download all subtitles: 
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
<a target="_blank" href="https://drive.google.com/file/d/1CWm3C12b_xmnNX1NYDDDgaeUAteXvYkX/view?usp=sharing">English</a>
&nbsp;&nbsp;&nbsp;&nbsp;
<a target="_blank" href="https://drive.google.com/file/d/1Z3myG9t3OcxsOPgPKjNXCvCf4axFRuNI/view?usp=sharing">Japanese</a>
</p>
<br><br><br>
<p>
Please visit forum if you have any questions <a href="http://forum.kitsunekko.net/">forum.kitsunekko.net</a>
</p><p>J-Drama subtitles <a href="http://jpsubbers.xyz/Japanese-Subtitles/%40Mains/">jpsubbers.xyz</a>
</p><p>Sorted anime subtitles <a href="https://itazuraneko.neocities.org/library/sub.html">itazuraneko</a>
</p><p>
</p><p>Watch anime with Japanese subtitles&nbsp;&nbsp; <a href="https://animelon.com/">animelon.com</a> &nbsp;&nbsp; <a href="https://anjsub.com/">anjsub.com</a>    
</p><p>&nbsp;
</p></td>
</tr>
</tbody></table>
</td></tr>
<tr>
<td>
</td>
</tr>
<tr>
<td align="right" style="font-size: xx-small; font-weight: bold;">
</td>
</tr>
</tbody></table>
<!-- Yandex.Metrika counter -->
<script type="text/javascript">
(function (d, w, c) {
(w[c] = w[c] || []).push(function() {
try {
w.yaCounter23071486 = new Ya.Metrika({id:23071486,
trackLinks:true});
} catch(e) { }
});
var n = d.getElementsByTagName("script")[0],
s = d.createElement("script"),
f = function () { n.parentNode.insertBefore(s, n); };
s.type = "text/javascript";
s.async = true;
s.src = (d.location.protocol == "https:" ? "https:" : "http:") + "//mc.yandex.ru/metrika/watch.js";
if (w.opera == "[object Opera]") {
d.addEventListener("DOMContentLoaded", f, false);
} else { f(); }
})(document, window, "yandex_metrika_callbacks");
</script>
<noscript><div><img src="//mc.yandex.ru/watch/23071486" style="position:absolute; left:-9999px;" alt="" /></div></noscript>
<!-- /Yandex.Metrika counter -->
</body></html>
`

const duration = 1 * time.Second

func newTestSever() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(data))
	})
	return httptest.NewServer(mux)
}

func prepareWebPageKeeper() (*webpagekeeper, *httptest.Server) {
	server := newTestSever()
	wpk := &webpagekeeper{
		timeToUpdate:   duration,
		cachedWebPages: []webPage{},
		logger:         zap.L().Sugar(),
	}

	return wpk, server
}

func TestGetUrlBody(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	expected := data
	actual, err := wpk.GetUrlBody(serv.URL+"/data", false)
	assert.NoError(err)

	assert.Equal(len(actual), len(expected))
}

func TestGetWebPageValid(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	expected := webPage{
		url:         serv.URL + "/data",
		timeUpdated: time.Now(),
		body:        []byte(data),
	}

	actual, err := wpk.getWebPage(serv.URL + "/data")
	assert.NoError(err)

	assert.Equal(expected.url, actual.url)
	assert.Equal(expected.body, actual.body)
	assert.LessOrEqual(expected.timeUpdated, actual.timeUpdated)
}

func TestGetWebPageWrongUrl(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	expected := webPage{}

	actual, err := wpk.getWebPage(serv.URL + "/crap")
	assert.Error(err)

	assert.Equal(expected.url, actual.url)
	assert.Equal(expected.body, actual.body)
	assert.LessOrEqual(expected.timeUpdated, actual.timeUpdated)
}

func TestGetWebPageNonExistingUrl(t *testing.T) {
	assert := assert.New(t)

	wpk, _ := prepareWebPageKeeper()

	expected := webPage{}

	actual, err := wpk.getWebPage("crap")
	assert.Error(err)

	assert.Equal(expected.url, actual.url)
	assert.Equal(expected.body, actual.body)
	assert.LessOrEqual(expected.timeUpdated, actual.timeUpdated)
}

func TestGetUrlBodyError(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	var expected []byte
	actual, err := wpk.GetUrlBody(serv.URL+"/crap", false)
	assert.Error(err)

	assert.Equal(len(actual), len(expected))
}

func TestGetUrlTestSave(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	expected := data
	_, err := wpk.GetUrlBody(serv.URL+"/data", true)
	assert.NoError(err)

	actual, err := wpk.GetUrlBody(serv.URL+"/data", true)
	assert.NoError(err)

	assert.Equal(len(actual), len(expected))
}

func TestGetUrlTestSaveUpdate(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	expected := data
	_, err := wpk.GetUrlBody(serv.URL+"/data", true)
	assert.NoError(err)

	time.Sleep(duration + 10*time.Millisecond)
	actual, err := wpk.GetUrlBody(serv.URL+"/data", true)
	assert.NoError(err)

	assert.Equal(len(actual), len(expected))
}

func TestGetUrlTestSaveUpdateError(t *testing.T) {
	assert := assert.New(t)

	wpk, serv := prepareWebPageKeeper()

	var expected []byte
	_, err := wpk.GetUrlBody(serv.URL+"/data", true)
	assert.NoError(err)

	serv.Close()
	time.Sleep(duration + 10*time.Millisecond)
	actual, err := wpk.GetUrlBody(serv.URL+"/data", true)
	assert.Error(err)

	assert.Equal(len(actual), len(expected))
}
