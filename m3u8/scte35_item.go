package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// SCTE35Item #EXT-X-SCTE35
type SCTE35Item struct {
	Cue        string
	Duration   *float64
	Elapsed    *float64
	ID         *string
	Time       *float64
	Type       *int
	UPID       *string
	Blackout   *string
	CueOut     *string
	CueIn      *string
	Segne      *string
	attributes map[string]string
}

func NewSCTE35Item(text string) (*SCTE35Item, error) {
	attributes := parser.ParseAttributes(text)

	defer deleteKeys(attributes,
		CueTag,
		DurationTag,
		ElapsedAttribute,
		IDTag,
		TimeAttribute,
		TypeTag,
		DurationTag,
		UPIDAttribute,
		BlackoutAttribute,
		CueOutAttribute,
		CueInAttribute,
		SegneAttribute,
	)

	return &SCTE35Item{
		Cue:        parser.SanitizeAttributeValue(attributes[CueTag]),
		Duration:   parser.PointerToFloat(attributes, DurationTag),
		Elapsed:    parser.PointerToFloat(attributes, ElapsedAttribute),
		ID:         parser.PointerTo(attributes, IDTag),
		Time:       parser.PointerToFloat(attributes, TimeAttribute),
		Type:       parser.PointerToInt(attributes, TypeTag),
		UPID:       parser.PointerTo(attributes, UPIDAttribute),
		Blackout:   parser.PointerTo(attributes, BlackoutAttribute),
		CueOut:     parser.PointerTo(attributes, CueOutAttribute),
		CueIn:      parser.PointerTo(attributes, CueInAttribute),
		Segne:      parser.PointerTo(attributes, SegneAttribute),
		attributes: attributes,
	}, nil

}

func (i *SCTE35Item) String() string {
	attributes := []string{fmt.Sprintf(parser.QuotedFormatString, CueTag, i.Cue)}

	if i.Duration != nil {
		attributes = append(attributes, fmt.Sprintf("%s=%.12f", DurationTag, *i.Duration))
	}

	if i.Elapsed != nil {
		attributes = append(attributes, fmt.Sprintf("%s=%.3f", ElapsedAttribute, *i.Elapsed))
	}

	if i.ID != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, IDTag, *i.ID))
	}

	if i.Time != nil {
		attributes = append(attributes, fmt.Sprintf(parser.FrameRateFormatString, TimeAttribute, *i.Time))
	}

	if i.Type != nil {
		attributes = append(attributes, fmt.Sprintf("%s=0x%x", TypeTag, *i.Type))
	}

	if i.UPID != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, UPIDAttribute, *i.UPID))
	}

	if i.Blackout != nil {
		attributes = append(attributes, fmt.Sprintf(parser.FormatString, BlackoutAttribute, *i.Blackout))
	}

	if i.CueOut != nil {
		attributes = append(attributes, fmt.Sprintf(parser.FormatString, CueOutAttribute, *i.CueOut))
	}

	if i.CueIn != nil {
		attributes = append(attributes, fmt.Sprintf(parser.FormatString, CueInAttribute, *i.CueIn))
	}

	if i.Segne != nil {
		attributes = append(attributes, fmt.Sprintf(parser.QuotedFormatString, SegneAttribute, *i.Segne))
	}

	attributes = attributesJoinMap(attributes, i.attributes)

	return fmt.Sprintf("%s:%s", SCTE35Tag, strings.Join(attributes, ","))
}
