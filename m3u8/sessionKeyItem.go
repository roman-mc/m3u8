package m3u8

import (
	"fmt"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// SessionKeyItem represents a set of EXT-X-SESSION-KEY attributes
type SessionKeyItem struct {
	Encryptable *Encryptable
}

// NewSessionKeyItem parses a text line and returns a *SessionKeyItem
func NewSessionKeyItem(text string) *SessionKeyItem {
	attributes := parser.ParseAttributes(text)
	return &SessionKeyItem{
		Encryptable: NewEncryptable(attributes),
	}
}

func (ski *SessionKeyItem) String() string {
	return fmt.Sprintf("%s:%v", SessionKeyItemTag, ski.Encryptable.String())
}

func (ski *SessionKeyItem) Validate() []error {
	return nil
}
