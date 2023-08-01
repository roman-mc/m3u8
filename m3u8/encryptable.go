package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// Encryptable is common representation for KeyItem and SessionKeyItem
type Encryptable struct {
	Method            string
	URI               *string
	IV                *string
	KeyFormat         *string
	KeyFormatVersions *string
	KeyID             *string
	attributes        map[string]string
}

// NewEncryptable takes an attributes map and returns an *Encryptable
func NewEncryptable(attributes map[string]string) *Encryptable {
	defer deleteKeys(
		attributes,
		MethodTag,
		URITag,
		IVTag,
		KeyFormatTag,
		KeyFormatVersionsTag,
		KeyID,
	)

	return &Encryptable{
		Method:            parser.SanitizeAttributeValue(attributes[MethodTag]),
		URI:               parser.PointerTo(attributes, URITag),
		IV:                parser.PointerTo(attributes, IVTag),
		KeyFormat:         parser.PointerTo(attributes, KeyFormatTag),
		KeyFormatVersions: parser.PointerTo(attributes, KeyFormatVersionsTag),
		KeyID:             parser.PointerTo(attributes, KeyID),
		attributes:        attributes,
	}
}

func (e *Encryptable) String() string {
	var slice []string

	slice = append(slice, fmt.Sprintf(parser.FormatString, MethodTag, e.Method))
	if e.URI != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, URITag, *e.URI))
	}
	if e.IV != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, IVTag, *e.IV))
	}
	if e.KeyFormat != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, KeyFormatTag, *e.KeyFormat))
	}
	if e.KeyFormatVersions != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, KeyFormatVersionsTag, *e.KeyFormatVersions))
	}
	if e.KeyID != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, KeyID, *e.KeyID))
	}

	for attributeKey, attribute := range e.attributes {
		slice = append(slice, fmt.Sprintf(parser.FormatString, attributeKey, attribute))
	}

	return strings.Join(slice, ",")
}
