package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// SessionDataItem represents a set of EXT-X-SESSION-DATA attributes
type SessionDataItem struct {
	DataID     string
	Value      *string
	URI        *string
	Language   *string
	attributes map[string]string
}

// NewSessionDataItem parses a text line and returns a *SessionDataItem
func NewSessionDataItem(text string) *SessionDataItem {
	attributes := parser.ParseAttributes(text)

	defer deleteKeys(
		attributes,
		DataIDTag,
		ValueTag,
		URITag,
		LanguageTag,
	)

	return &SessionDataItem{
		DataID:     parser.SanitizeAttributeValue(attributes[DataIDTag]),
		Value:      parser.PointerTo(attributes, ValueTag),
		URI:        parser.PointerTo(attributes, URITag),
		Language:   parser.PointerTo(attributes, LanguageTag),
		attributes: attributes,
	}
}

func (sdi *SessionDataItem) String() string {
	slice := []string{fmt.Sprintf(parser.QuotedFormatString, DataIDTag, sdi.DataID)}

	if sdi.Value != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, ValueTag, *sdi.Value))
	}
	if sdi.URI != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, URITag, *sdi.URI))
	}
	if sdi.Language != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, LanguageTag, *sdi.Language))
	}

	slice = attributesJoinMap(slice, sdi.attributes)

	return fmt.Sprintf(`%s:%s`, SessionDataItemTag, strings.Join(slice, ","))
}

func (sdi *SessionDataItem) Validate() []error {
	return nil
}
