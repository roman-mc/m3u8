package m3u8

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeItem_New(t *testing.T) {
	timeVar, err := ParseTime("2010-02-19T14:54:23.031Z")
	assert.Nil(t, err)
	ti := &TimeItem{
		Time: timeVar,
	}

	assert.Equal(t, "#EXT-X-PROGRAM-DATE-TIME:2010-02-19T14:54:23.031Z", ti.String())
}

func TestTimeItem_Parse(t *testing.T) {
	ti, err := NewTimeItem("#EXT-X-PROGRAM-DATE-TIME:2010-02-19T14:54:23.031Z")
	assert.Nil(t, err)

	expected, err := time.Parse(time.RFC3339Nano, "2010-02-19T14:54:23.031Z")
	assert.Nil(t, err)

	assert.Equal(t, expected, ti.Time)
}

func TestTimeItem_Error(t *testing.T) {
	ti, err := NewTimeItem("#EXT-X-PROGRAM-DATE-TIME:23120312312")
	assert.Error(t, err)
	assert.Nil(t, ti)

}
