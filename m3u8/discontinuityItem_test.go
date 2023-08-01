package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscontinuityItem_Parse(t *testing.T) {
	di, err := NewDiscontinuityItem()
	assert.Nil(t, err)
	assert.Equal(t, DiscontinuityItemTag, di.String())
}
