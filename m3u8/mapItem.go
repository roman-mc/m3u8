package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// MapItem represents a EXT-X-MAP tag which specifies how to obtain the Media
// Initialization Section
type MapItem struct {
	URI        string
	ByteRange  *ByteRange
	attributes map[string]string
}

// NewMapItem parses a text line and returns a *MapItem
func NewMapItem(text string) *MapItem {
	attributes := parser.ParseAttributes(text)
	br, _ := NewByteRange(parser.SanitizeAttributeValue(attributes[ByteRangeTag]))

	defer deleteKeys(
		attributes,
		ByteRangeTag,
		URITag,
	)

	return &MapItem{
		URI:        parser.SanitizeAttributeValue(attributes[URITag]),
		ByteRange:  br,
		attributes: attributes,
	}
}

func (mi *MapItem) String() string {
	var attributes []string
	attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, URITag, mi.URI))

	if mi.ByteRange != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, ByteRangeTag, mi.ByteRange))
	}

	for attributeKey, attribute := range mi.attributes {
		attributes = append(attributes, fmt.Sprintf(parser.FormatString, attributeKey, attribute))
	}

	return fmt.Sprintf(`%s:%s`, MapItemTag, strings.Join(attributes, ","))
}

func (mi *MapItem) Validate() []error {
	return nil
}
