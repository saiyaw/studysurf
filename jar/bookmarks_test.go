package jar

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryBookmarks(t *testing.T) {
	b := NewMemoryBookmarks()

	assertBookmarks(t, b)
}

func TestFileBookmarks(t *testing.T) {

	b, err := NewFileBookmarks("./bookmarks.json")
	assert.Nil(t, err)

	defer func() {
		err = os.Remove("./bookmarks.json")
	}()

	assertBookmarks(t, b)
}

// assertBookmarks tests the given bookmark jar.
func assertBookmarks(t *testing.T, b BookmarksJar) {
	err := b.Save("test1", "http://localhost")
	assert.Nil(t, err)

	err = b.Save("test2", "http://127.0.0.1")
	assert.Nil(t, err)

	err = b.Save("test1", "http://localhost")
	assert.NotNil(t, err)

	url, err := b.Read("test1")
	assert.Nil(t, err)
	assert.Equal(t, "http://localhost", url)

	url, err = b.Read("test2")
	assert.Equal(t, "http://127.0.0.1", url)

	url, err = b.Read("test3")
	assert.NotNil(t, err)

	r := b.Remove("test2")
	assert.True(t, r)

	r = b.Remove("test3")
	assert.False(t, r)

	r = b.Has("test1")
	assert.True(t, r)

	r = b.Has("test4")
	assert.False(t, r)

}
