package m3u8

import "github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"

// UnknownItem
// represents any unknown tag to us, preserve its value as is;
// or if a known tag has a parsing error, we preserve it as is in unknown tag
type UnknownItem struct {
	tagValue string
	err      error
}

func NewUnknownItem(text string, err error) *UnknownItem {
	tagValue := text

	return &UnknownItem{
		tagValue: tagValue,
		err:      err,
	}
}

func (i *UnknownItem) String() string {
	return i.tagValue
}

// GetTagName
// In case we got err during tag initialization we fall back to UnknownItem,
// and if we iterate over all the tags, it's convenient to check unknownTag's name
func (i *UnknownItem) GetTagName() string {
	name := parser.ParseTagName(i.tagValue)
	return name
}
