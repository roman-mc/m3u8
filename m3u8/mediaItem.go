package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// MediaItem represents a set of EXT-X-MEDIA attributes
type MediaItem struct {
	Type              string
	GroupID           string
	Name              string
	Language          *string
	AssocLanguage     *string
	AutoSelect        *bool
	Default           *bool
	Forced            *bool
	URI               *string
	InStreamID        *string
	Characteristics   *string
	Channels          *string
	StableRenditionId *string
	attributes        map[string]string
}

// NewMediaItem parses a text line and returns a *MediaItem
func NewMediaItem(text string) *MediaItem {
	attributes := parser.ParseAttributes(text)
	defer deleteKeys(attributes,
		TypeTag,
		GroupIDTag,
		NameTag,
		LanguageTag,
		AssocLanguageTag,
		AutoSelectTag,
		DefaultTag,
		ForcedTag,
		URITag,
		InStreamIDTag,
		CharacteristicsTag,
		ChannelsTag,
		StableRenditionIDTag,
	)

	return &MediaItem{
		Type:              parser.SanitizeAttributeValue(attributes[TypeTag]),
		GroupID:           parser.SanitizeAttributeValue(attributes[GroupIDTag]),
		Name:              parser.SanitizeAttributeValue(attributes[NameTag]),
		Language:          parser.PointerTo(attributes, LanguageTag),
		AssocLanguage:     parser.PointerTo(attributes, AssocLanguageTag),
		AutoSelect:        parser.ParseYesNo(attributes, AutoSelectTag),
		Default:           parser.ParseYesNo(attributes, DefaultTag),
		Forced:            parser.ParseYesNo(attributes, ForcedTag),
		URI:               parser.PointerTo(attributes, URITag),
		InStreamID:        parser.PointerTo(attributes, InStreamIDTag),
		Characteristics:   parser.PointerTo(attributes, CharacteristicsTag),
		Channels:          parser.PointerTo(attributes, ChannelsTag),
		StableRenditionId: parser.PointerTo(attributes, StableRenditionIDTag),
		attributes:        attributes,
	}
}

func (mi *MediaItem) String() string {
	slice := []string{
		fmt.Sprintf(parser.FormatString, TypeTag, mi.Type),
		fmt.Sprintf(parser.QuotedFormatString, GroupIDTag, mi.GroupID),
	}

	if mi.Language != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, LanguageTag, *mi.Language))
	}
	if mi.AssocLanguage != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, AssocLanguageTag, *mi.AssocLanguage))
	}
	slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, NameTag, mi.Name))
	if mi.AutoSelect != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, AutoSelectTag, parser.FormatYesNo(*mi.AutoSelect)))
	}
	if mi.Default != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, DefaultTag, parser.FormatYesNo(*mi.Default)))
	}
	if mi.URI != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, URITag, *mi.URI))
	}
	if mi.Forced != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, ForcedTag, parser.FormatYesNo(*mi.Forced)))
	}
	if mi.InStreamID != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, InStreamIDTag, *mi.InStreamID))
	}
	if mi.Characteristics != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, CharacteristicsTag, *mi.Characteristics))
	}
	if mi.Channels != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, ChannelsTag, *mi.Channels))
	}
	if mi.StableRenditionId != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, StableRenditionIDTag, *mi.StableRenditionId))
	}
	for attributeKey, attribute := range mi.attributes {
		slice = append(slice, fmt.Sprintf(parser.FormatString, attributeKey, attribute))
	}

	return fmt.Sprintf("%s:%s", MediaItemTag, strings.Join(slice, ","))
}

func (mi *MediaItem) Validate() []error {
	return nil
}
