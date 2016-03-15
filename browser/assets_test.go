package browser

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {

	out := &bytes.Buffer{}
	u, _ := url.Parse("http://i.imgur.com/HW4bJtY.jpg")
	asset := NewImageAsset(u, "", "", "")
	l, err := DownloadAsset(asset, out)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, int(l))
	assert.Equal(t, int(l), out.Len())
}

func TestDownloadAsync(t *testing.T) {
	ch := make(AsyncDownloadChannel, 1)
	u1, _ := url.Parse("http://i.imgur.com/HW4bJtY.jpg")
	u2, _ := url.Parse("http://i.imgur.com/HkPOzEH.jpg")
	asset1 := NewImageAsset(u1, "", "", "")
	asset2 := NewImageAsset(u2, "", "", "")
	out1 := &bytes.Buffer{}
	out2 := &bytes.Buffer{}

	queue := 2
	DownloadAssetAsync(asset1, out1, ch)
	DownloadAssetAsync(asset2, out2, ch)

	for {
		select {
		case result := <-ch:
			assert.NotEqual(t, 0, int(result.Size))
			if result.Asset == asset1 {
				assert.Equal(t, int(result.Size), out1.Len())
			} else if result.Asset == asset2 {
				assert.Equal(t, int(result.Size), out2.Len())
			} else {
				t.Failed()
			}
			queue--
			if queue == 0 {
				goto COMPLETE
			}
		}
	}

COMPLETE:
	close(ch)
	assert.Equal(t, 0, queue)
}
