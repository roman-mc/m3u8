package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// DateRangeItem represents a #EXT-X-DATERANGE tag
type DateRangeItem struct {
	ID               string
	Class            *string
	Cue              *string
	StartDate        string
	EndDate          *string
	Duration         *float64
	PlannedDuration  *float64
	Scte35Cmd        *string
	Scte35Out        *string
	Scte35In         *string
	EndOnNext        bool
	ClientAttributes map[string]string
}

// NewDateRangeItem parses a text line in playlist and returns a *DateRangeItem
func NewDateRangeItem(text string) *DateRangeItem {
	attributes := parser.ParseAttributes(text)

	defer deleteKeys(attributes,
		DurationTag,
		PlannedDurationTag,
		IDTag,
		ClassTag,
		StartDateTag,
		CueTag,
		EndDateTag,
		Scte35CmdTag,
		Scte35OutTag,
		Scte35InTag,
		EndOnNextTag,
	)

	return &DateRangeItem{
		ID:               parser.SanitizeAttributeValue(attributes[IDTag]),
		Class:            parser.PointerTo(attributes, ClassTag),
		StartDate:        parser.SanitizeAttributeValue(attributes[StartDateTag]),
		Cue:              parser.PointerTo(attributes, CueTag),
		EndDate:          parser.PointerTo(attributes, EndDateTag),
		Duration:         parser.PointerToFloat(attributes, DurationTag),
		PlannedDuration:  parser.PointerToFloat(attributes, PlannedDurationTag),
		Scte35Cmd:        parser.PointerTo(attributes, Scte35CmdTag),
		Scte35Out:        parser.PointerTo(attributes, Scte35OutTag),
		Scte35In:         parser.PointerTo(attributes, Scte35InTag),
		EndOnNext:        parser.AttributeExists(EndOnNextTag, attributes),
		ClientAttributes: attributes,
	}
}

func (dri *DateRangeItem) String() string {
	var slice []string

	slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, IDTag, dri.ID))
	if dri.Class != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, ClassTag, *dri.Class))
	}
	if dri.Cue != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, CueTag, *dri.Cue))
	}
	slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, StartDateTag, dri.StartDate))
	if dri.EndDate != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, EndDateTag, *dri.EndDate))
	}
	if dri.Duration != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, DurationTag, *dri.Duration))
	}
	if dri.PlannedDuration != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, PlannedDurationTag, *dri.PlannedDuration))
	}
	clientAttributes := formatClientAttributes(dri.ClientAttributes)
	slice = append(slice, clientAttributes...)

	if dri.Scte35Cmd != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, Scte35CmdTag, *dri.Scte35Cmd))
	}
	if dri.Scte35Out != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, Scte35OutTag, *dri.Scte35Out))
	}
	if dri.Scte35In != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, Scte35InTag, *dri.Scte35In))
	}
	if dri.EndOnNext {
		slice = append(slice, fmt.Sprintf(`%s=YES`, EndOnNextTag))
	}

	return fmt.Sprintf("%s:%s", DateRangeItemTag, strings.Join(slice, ","))
}

func (dri *DateRangeItem) Validate() []error {
	return nil
}

func formatClientAttributes(ca map[string]string) []string {
	var slice []string

	for key, value := range ca {
		slice = append(slice, fmt.Sprintf(parser.FormatString, key, value))
	}

	return slice
}
