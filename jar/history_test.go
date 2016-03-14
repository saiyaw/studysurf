package jar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryHistory(t *testing.T) {
	stack := NewMemoryHistory()

	page1 := &State{}
	stack.Push(page1)
	assert.Equal(t, 1, stack.Len())
	assert.Equal(t, page1, stack.Top())

	page2 := &State{}
	stack.Push(page2)
	assert.Equal(t, 2, stack.Len())
	assert.Equal(t, page2, stack.Top())

	page := stack.Pop()
	assert.Equal(t, page, page2)
	assert.Equal(t, 1, stack.Len())
	assert.Equal(t, page1, stack.Top())

	page = stack.Pop()
	assert.Equal(t, page, page1)
	assert.Equal(t, 0, stack.Len())
}
