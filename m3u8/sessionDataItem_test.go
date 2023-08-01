package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionDataItem_Parse(t *testing.T) {
	line := `#EXT-X-SESSION-DATA:DATA-ID="com.test.movie.title",VALUE="Test",LANGUAGE="en"`

	sdi := NewSessionDataItem(line)

	assert.Equal(t, "com.test.movie.title", sdi.DataID)
	assertNotNilEqual(t, "Test", sdi.Value)
	assert.Nil(t, sdi.URI)
	assertNotNilEqual(t, "en", sdi.Language)
	assertToString(t, line, sdi)

	line = `#EXT-X-SESSION-DATA:DATA-ID="com.test.movie.title",URI="http://test",LANGUAGE="en"`
	sdi = NewSessionDataItem(line)

	assert.Equal(t, "com.test.movie.title", sdi.DataID)
	assert.Nil(t, sdi.Value)
	assertNotNilEqual(t, "http://test", sdi.URI)
	assertNotNilEqual(t, "en", sdi.Language)
	assertToString(t, line, sdi)
}
