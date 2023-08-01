package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyItem_Parse(t *testing.T) {
	line := `#EXT-X-KEY:METHOD=AES-128,URI="http://test.key",IV=D512BBF,KEYFORMAT="identity",KEYFORMATVERSIONS="1/3"`

	ki := NewKeyItem(line)
	assert.NotNil(t, ki.Encryptable)
	assert.Equal(t, "AES-128", ki.Encryptable.Method)
	assertNotNilEqual(t, "http://test.key", ki.Encryptable.URI)
	assertNotNilEqual(t, "D512BBF", ki.Encryptable.IV)
	assertNotNilEqual(t, "identity", ki.Encryptable.KeyFormat)
	assertNotNilEqual(t, "1/3", ki.Encryptable.KeyFormatVersions)

	assertToString(t, line, ki)
}
