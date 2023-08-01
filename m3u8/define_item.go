package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// DefineItem represents a set of EXT-X-DEFINE attributes
type DefineItem struct {
	Name       string
	Value      *string
	Import     *string
	QueryParam *string
	attributes map[string]string
}

// NewDefineItem parses a text line and returns a *DefineItem
func NewDefineItem(text string) *DefineItem {
	attributes := parser.ParseAttributes(text)

	defer deleteKeys(attributes,
		AttributeName,
		AttributeValue,
		AttributeImport,
		AttributeQueryParam,
	)

	return &DefineItem{
		Name:       parser.SanitizeAttributeValue(attributes[AttributeName]),
		Value:      parser.PointerTo(attributes, AttributeValue),
		Import:     parser.PointerTo(attributes, AttributeImport),
		QueryParam: parser.PointerTo(attributes, AttributeQueryParam),
		attributes: attributes,
	}
}

func (i *DefineItem) String() string {
	var attributes []string

	if i.Name != "" {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, AttributeName, i.Name))
	}

	if i.Value != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, AttributeValue, *i.Value))
	}
	if i.Import != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, AttributeImport, *i.Import))

	}
	if i.QueryParam != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, AttributeQueryParam, *i.QueryParam))
	}

	for attributeKey, attribute := range i.attributes {
		attributes = append(attributes, fmt.Sprintf(parser.FormatString, attributeKey, attribute))
	}

	return fmt.Sprintf("%s:%v", DefineTag, strings.Join(attributes, ","))
}

func (i *DefineItem) Validate() []error {
	return nil
}
