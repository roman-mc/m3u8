package m3u8

import (
	"fmt"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// KeyItem represents a set of EXT-X-KEY attributes
type KeyItem struct {
	Encryptable *Encryptable
}

// NewKeyItem parses a text line and returns a *KeyItem
func NewKeyItem(text string) *KeyItem {
	attributes := parser.ParseAttributes(text)

	return &KeyItem{
		Encryptable: NewEncryptable(attributes),
	}
}

func (ki *KeyItem) String() string {
	return fmt.Sprintf("%s:%v", KeyItemTag, ki.Encryptable.String())
}

func (ki *KeyItem) Validate() []error {
	return nil
}
