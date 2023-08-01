package m3u8

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// PlaybackStart represents a #EXT-X-START tag and attributes
type PlaybackStart struct {
	TimeOffset float64
	Precise    *bool
	attributes map[string]string
}

// NewPlaybackStart parses a text line and returns a *PlaybackStart
func NewPlaybackStart(text string) (*PlaybackStart, error) {
	attributes := parser.ParseAttributes(text)

	timeOffset, err := strconv.ParseFloat(parser.SanitizeAttributeValue(attributes[TimeOffsetTag]), 64)
	if err != nil {
		return nil, err
	}

	defer deleteKeys(attributes,
		TimeOffsetTag,
		PreciseTag,
	)

	return &PlaybackStart{
		TimeOffset: timeOffset,
		Precise:    parser.ParseYesNo(attributes, PreciseTag),
		attributes: attributes,
	}, nil
}

func (ps *PlaybackStart) String() string {
	slice := []string{fmt.Sprintf(parser.FormatString, TimeOffsetTag, ps.TimeOffset)}
	if ps.Precise != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, PreciseTag, parser.FormatYesNo(*ps.Precise)))
	}

	for attributeKey, attribute := range ps.attributes {
		slice = append(slice, fmt.Sprintf(parser.FormatString, attributeKey, attribute))
	}

	return fmt.Sprintf(`%s:%s`, PlaybackStartTag, strings.Join(slice, ","))
}

func (ps *PlaybackStart) Validate() []error {
	return nil
}
