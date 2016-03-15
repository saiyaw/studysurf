package surf

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saiyawang/studysurf/jar"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/page1" {
			fmt.Fprint(w, htmlPage1)
		} else if req.URL.Path == "/page2" {
			fmt.Fprint(w, htmlPage2)
		}
	}))
	defer ts.Close()

	bow := NewBrowser()

	err := bow.Open(ts.URL + "/page1")
	assert.Nil(t, err)
	assert.Equal(t, "Surf Page 1", bow.Title())
	assert.Contains(t, bow.Body(), "<p>Hello, Surf!</p>")

	err = bow.Open(ts.URL + "/page2")
	assert.Nil(t, err)
	assert.Equal(t, "Surf Page 2", bow.Title())

	ok := bow.Back()
	assert.True(t, ok)
	assert.Equal(t, "Surf Page 1", bow.Title())

	ok = bow.Back()
	assert.False(t, ok)
	assert.Equal(t, "Surf Page 1", bow.Title())

}

func TestPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(500)
			w.Write([]byte("ERROR"))
		}
	}))
	defer ts.Close()

	bow := NewBrowser()
	bow.Post(ts.URL, "application/x-www-form-urlencoded", nil)
	assert.Equal(t, 200, bow.StatusCode())
}

func TestHead(t *testing.T) {
	var r *http.Request

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		r = req
	}))
	defer ts.Close()

	bow := NewBrowser()

	err := bow.Head(ts.URL + "/page1")
	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func TestDownload(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, req.UserAgent())
	}))
	defer ts.Close()

	bow := NewBrowser()
	bow.Open(ts.URL)

	buff := &bytes.Buffer{}
	l, err := bow.Download(buff)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, int(l))
	assert.Equal(t, int(l), buff.Len())
}

func TestDownloadContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		b := bytes.NewBufferString("Hello")
		fmt.Fprint(w, b)
	}))
	defer ts.Close()

	bow := NewBrowser()
	bow.Open(ts.URL)

	buff := &bytes.Buffer{}
	bow.Download(buff)
	assert.Equal(t, "Hello", buff.String())
}

func TestUserAgent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, req.UserAgent())
	}))
	defer ts.Close()

	bow := NewBrowser()
	bow.SetUserAgent("Testing/1.0")
	err := bow.Open(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, "Testing/1.0", bow.Body())
}

func TestHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, req.Header.Get("X-Testing-1"))
		fmt.Fprint(w, req.Header.Get("X-Testing-2"))
	}))
	defer ts.Close()

	bow := NewBrowser()
	bow.AddRequestHeader("X-Testing-1", "Testing-1")
	bow.AddRequestHeader("X-Testing-2", "Testing-2")
	err := bow.Open(ts.URL)
	assert.Nil(t, err)
	assert.Contains(t, bow.Body(), "Testing-1")
	assert.Contains(t, bow.Body(), "Testing-2")
}

// TestHeadersSet
// See: https://github.com/headzoo/surf/pull/19
func TestHeadersBug19(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, req.Header.Get("X-Testing"))
	}))
	defer ts.Close()

	bow := NewBrowser()
	bow.AddRequestHeader("X-Testing", "Testing-1")
	bow.AddRequestHeader("X-Testing", "Testing-2")
	err := bow.Open(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, "Testing-2", bow.Body())
}

func TestBookmarks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlPage1)
	}))
	defer ts.Close()

	bookmarks := jar.NewMemoryBookmarks()
	bow := NewBrowser()
	bow.SetBookmarksJar(bookmarks)

	bookmarks.Save("test1", ts.URL)
	bow.OpenBookmark("test1")
	assert.Equal(t, "Surf Page 1", bow.Title())
	assert.Contains(t, bow.Body(), "<p>Hello, Surf!</p>")

	err := bow.Bookmark("test2")
	assert.Nil(t, err)
	bow.OpenBookmark("test2")
	assert.Equal(t, "Surf Page 1", bow.Title())
}

func TestClick(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprint(w, htmlPage1)
		} else if r.URL.Path == "/page2" {
			fmt.Fprint(w, htmlPage1)
		}
	}))
	defer ts.Close()

	bow := NewBrowser()
	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	err = bow.Click("a:contains('click')")
	assert.Nil(t, err)
	assert.Contains(t, bow.Body(), "<p>Hello, Surf!</p>")
}

func TestLinks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlPage1)
	}))
	defer ts.Close()

	bow := NewBrowser()
	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	links := bow.Links()

	assert.Equal(t, 2, len(links))
	assert.Equal(t, "", links[0].ID)
	assert.Equal(t, ts.URL+"/page2", links[0].URL.String())
	assert.Equal(t, "click", links[0].Text)
	assert.Equal(t, "page3", links[1].ID)
	assert.Equal(t, ts.URL+"/page3", links[1].URL.String())
	assert.Equal(t, "no clicking", links[1].Text)
}

func TestImages(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlPage1)
	}))
	defer ts.Close()

	bow := NewBrowser()
	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	images := bow.Images()

	assert.Equal(t, 2, len(images))
	assert.Equal(t, "imgur-image", images[0].ID)
	assert.Equal(t, "http://i.imgur.com/HW4bJtY.jpg", images[0].URL.String())
	assert.Equal(t, "", images[0].Alt)
	assert.Equal(t, "It's a...", images[0].Title)

	assert.Equal(t, "", images[1].ID)
	assert.Equal(t, ts.URL+"/Cxagv.jpg", images[1].URL.String())
	assert.Equal(t, "A picture", images[1].Alt)
	assert.Equal(t, "", images[1].Title)

	buff := &bytes.Buffer{}
	l, err := images[0].Download(buff)

	assert.Nil(t, err)
	assert.NotEqual(t, 0, buff.Len())
	assert.Equal(t, int(l), buff.Len())

}

func TestStylesheets(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlPage1)
	}))
	defer ts.Close()

	bow := NewBrowser()
	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	stylesheets := bow.Stylesheets()

	assert.Equal(t, 2, len(stylesheets))
	assert.Equal(t, "http://godoc.org/-/site.css", stylesheets[0].URL.String())
	assert.Equal(t, "all", stylesheets[0].Media)
	assert.Equal(t, "text/css", stylesheets[0].Type)

	assert.Equal(t, ts.URL+"/print.css", stylesheets[1].URL.String())
	assert.Equal(t, "print", stylesheets[1].Media)
	assert.Equal(t, "text/css", stylesheets[1].Type)

	buff := &bytes.Buffer{}
	l, err := stylesheets[0].Download(buff)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, buff.Len())
	assert.Equal(t, int(l), buff.Len())
}

func TestScripts(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlPage1)
	}))
	defer ts.Close()

	bow := NewBrowser()
	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	scripts := bow.Scripts()

	assert.Equal(t, 2, len(scripts))
	assert.Equal(t, "http://godoc.org/-/site.js", scripts[0].URL.String())
	assert.Equal(t, "text/javascript", scripts[0].Type)

	assert.Equal(t, ts.URL+"/jquery.min.js", scripts[1].URL.String())
	assert.Equal(t, "text/javascript", scripts[1].Type)

	buff := &bytes.Buffer{}
	l, err := scripts[0].Download(buff)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, buff.Len())
	assert.Equal(t, int(l), buff.Len())
}

var htmlPage1 = `<!doctype html>
<html>
	<head>
		<title>Surf Page 1</title>
		<link href="/favicon.ico" rel="icon" type="image/x-icon">
		<link href="http://godoc.org/-/site.css" media="all" rel="stylesheet" type="text/css" />
		<link href="/print.css" rel="stylesheet" media="print" />
	</head>
	<body>
		<p>Hello, Surf!</p>
		<img src="http://i.imgur.com/HW4bJtY.jpg" id="imgur-image" title="It's a..." />
		<img src="/Cxagv.jpg" alt="A picture" />

		<p>Click the link below.</p>
		<a href="/page2">click</a>
		<a href="/page3" id="page3">no clicking</a>

		<script src="http://godoc.org/-/site.js" type="text/javascript"></script>
		<script src="/jquery.min.js" type="text/javascript"></script>
		<script type="text/javascript">
			var _gaq = _gaq || [];
		</script>
	</body>
</html>
`

var htmlPage2 = `<!doctype html>
<html>
	<head>
		<title>Surf Page 2</title>
	</head>
	<body>
		<p>Hello, Surf!</p>
	</body>
</html>
`
