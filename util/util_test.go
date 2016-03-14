package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	assert.True(t, FileExists("./util_test.go"), "util_test.go should be existed.")

	assert.False(t, FileExists("./util.txt"), "util.txt should not be existed.")
}
