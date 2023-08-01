package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionKeyItem_Parse(t *testing.T) {
	line := `#EXT-X-SESSION-KEY:METHOD=AES-128,URI="http://test.key",IV=D512BBF,KEYFORMAT="identity",KEYFORMATVERSIONS="1/3"`

	ski := NewSessionKeyItem(line)
	assert.NotNil(t, ski.Encryptable)

	assert.Equal(t, "AES-128", ski.Encryptable.Method)
	assertNotNilEqual(t, "http://test.key", ski.Encryptable.URI)
	assertNotNilEqual(t, "D512BBF", ski.Encryptable.IV)
	assertNotNilEqual(t, "identity", ski.Encryptable.KeyFormat)
	assertNotNilEqual(t, "1/3", ski.Encryptable.KeyFormatVersions)

	assertToString(t, line, ski)
}