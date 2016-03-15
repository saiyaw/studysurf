package browser

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saiyawang/studysurf/jar"
	"github.com/stretchr/testify/assert"
)

func TestBrowserForm(t *testing.T) {
	ts := setupTestServer(htmlForm, t)
	defer ts.Close()

	bow := &Browser{}
	bow.headers = make(http.Header, 10)
	bow.history = jar.NewMemoryHistory()

	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	f, err := bow.Form("[name='default']")
	assert.Nil(t, err)

	f.Input("age", "55")
	f.Input("gender", "male")
	err = f.Click("submit2")
	assert.Nil(t, err)
	assert.Contains(t, bow.Body(), "age=55")
	assert.Contains(t, bow.Body(), "gender=male")
	assert.Contains(t, bow.Body(), "submit2=submitted2")

}

func TestBrowserFormClickByValue(t *testing.T) {
	ts := setupTestServer(htmlFormClick, t)
	defer ts.Close()

	bow := &Browser{}
	bow.headers = make(http.Header, 10)
	bow.history = jar.NewMemoryHistory()

	err := bow.Open(ts.URL)
	assert.Nil(t, err)

	f, err := bow.Form("[name='default']")
	assert.Nil(t, err)

	f.Input("age", "55")
	err = f.ClickByValue("submit", "submitted2")
	assert.Nil(t, err)
	assert.Contains(t, bow.Body(), "age=55")
	assert.Contains(t, bow.Body(), "submit=submitted2")
}

func setupTestServer(html string, t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(w, html)
		} else {
			r.ParseForm()
			fmt.Fprint(w, r.Form.Encode())
		}
	}))

	return ts
}

var htmlForm = `<!doctype html>
<html>
	<head>
		<title>Echo Form</title>
	</head>
	<body>
		<form method="post" action="/" name="default">
			<input type="text" name="age" value="" />
			<input type="radio" name="gender" value="male" />
			<input type="radio" name="gender" value="female" />
			<input type="submit" name="submit1" value="submitted1" />
			<input type="submit" name="submit2" value="submitted2" />
		</form>
	</body>
</html>
`

var htmlFormClick = `<!doctype html>
<html>
	<head>
		<title>Echo Form</title>
	</head>
	<body>
		<form method="post" action="/" name="default">
			<input type="text" name="age" value="" />
			<input type="submit" name="submit" value="submitted1" />
			<input type="submit" name="submit" value="submitted2" />
		</form>
	</body>
</html>
`
