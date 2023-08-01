package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// ImageStreamItem represents a set of EXT-X-IMAGE-STREAM-INF ImageStreamItemTag
// All attributes defined for the EXT-X-I-FRAME-STREAM-INF tag are also defined for the EXT-X-IMAGE-STREAM-INF tag, except for HDCP-LEVEL and VIDEO-RANGE, which are not applicable.
type ImageStreamItem struct {
	Bandwidth int
	URI       string

	Name             *string
	AverageBandwidth *int
	ProgramID        *string
	Codecs           *string
	AudioCodec       *string // TODO: review this attribute, it isn't set and perhaps breaks some logic
	Profile          *string // TODO: review this attribute, it isn't set and perhaps breaks some logic
	Level            *string // TODO: review this attribute, it isn't set and perhaps breaks some logic
	Video            *string
	Audio            *string
	Subtitles        *string
	ClosedCaptions   *string
	FrameRate        *float64
	Resolution       *parser.Resolution
	StableVariantID  *string
	attributes       map[string]string
}

// NewImageStreamItem parses a text line and returns a *ImageStreamItem
func NewImageStreamItem(text string) *ImageStreamItem {
	attributes := parser.ParseAttributes(text)
	resolution, _ := parser.ParseResolution(attributes, ResolutionTag)
	averageBandwidth, _ := parser.ParseInt(attributes, AverageBandwidthTag)
	frameRate, _ := parser.ParseFloat(attributes, FrameRateTag)
	bandwidth, _ := parser.ParseBandwidth(attributes, BandwidthTag)

	if frameRate != nil && *frameRate <= 0 {
		frameRate = nil
	}

	defer deleteKeys(attributes,
		ResolutionTag,
		AverageBandwidthTag,
		FrameRateTag,
		BandwidthTag,
		ProgramIDTag,
		CodecsTag,
		AudioTag,
		VideoTag,
		URITag,
		SubtitlesTag,
		ClosedCaptionsTag,
		NameTag,
		StableVariantIDTag,
	)

	return &ImageStreamItem{
		ProgramID:        parser.PointerTo(attributes, ProgramIDTag),
		Codecs:           parser.PointerTo(attributes, CodecsTag),
		Bandwidth:        bandwidth,
		AverageBandwidth: averageBandwidth,
		FrameRate:        frameRate,
		Video:            parser.PointerTo(attributes, VideoTag),
		Audio:            parser.PointerTo(attributes, AudioTag),
		URI:              parser.SanitizeAttributeValue(attributes[URITag]),
		Subtitles:        parser.PointerTo(attributes, SubtitlesTag),
		ClosedCaptions:   parser.PointerTo(attributes, ClosedCaptionsTag),
		Name:             parser.PointerTo(attributes, NameTag),
		Resolution:       resolution,
		StableVariantID:  parser.PointerTo(attributes, StableVariantIDTag),
		attributes:       attributes,
	}
}

func (i *ImageStreamItem) String() string {
	var slice []string

	if i.ProgramID != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, ProgramIDTag, *i.ProgramID))
	}
	if i.Resolution != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, ResolutionTag, i.Resolution.String()))
	}
	codecs := i.formatCodecs()
	if codecs != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, CodecsTag, *codecs))
	}
	slice = append(slice, fmt.Sprintf(parser.FormatString, BandwidthTag, i.Bandwidth))
	if i.AverageBandwidth != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, AverageBandwidthTag, *i.AverageBandwidth))
	}
	if i.FrameRate != nil {
		slice = append(slice, fmt.Sprintf(parser.FrameRateFormatString, FrameRateTag, *i.FrameRate))
	}
	if i.Audio != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, AudioTag, *i.Audio))
	}
	if i.Video != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, VideoTag, *i.Video))
	}
	if i.Subtitles != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, SubtitlesTag, *i.Subtitles))
	}
	if i.ClosedCaptions != nil {
		cc := *i.ClosedCaptions
		fs := parser.QuotedFormatString
		if cc == parser.NoneValue {
			fs = parser.FormatString
		}
		slice = append(slice, fmt.Sprintf(fs, ClosedCaptionsTag, cc))
	}
	if i.Name != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, NameTag, *i.Name))
	}
	if i.StableVariantID != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, StableVariantIDTag, *i.StableVariantID))
	}

	slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, URITag, i.URI))

	slice = attributesJoinMap(slice, i.attributes)
	attributesString := strings.Join(slice, ",")

	return fmt.Sprintf("%s:%s", ImageStreamItemTag, attributesString)
}

func (i *ImageStreamItem) formatCodecs() *string {
	if i.Codecs != nil {
		return i.Codecs
	}

	videoCodecPtr := videoCodec(i.Profile, i.Level)
	// profile or level were specified but not recognized any codecs
	if !(i.Profile == nil && i.Level == nil) && videoCodecPtr == nil {
		return nil
	}

	audioCodecPtr := audioCodec(i.AudioCodec)
	// audio codec was specified but not recognized
	if !(i.AudioCodec == nil) && audioCodecPtr == nil {
		return nil
	}

	var slice []string
	if videoCodecPtr != nil {
		slice = append(slice, *videoCodecPtr)
	}
	if audioCodecPtr != nil {
		slice = append(slice, *audioCodecPtr)
	}

	if len(slice) <= 0 {
		return nil
	}

	value := strings.Join(slice, ",")
	return &value
}

// CodecsString returns the string representation of codecs for a ImageStreamItem
func (i *ImageStreamItem) CodecsString() string {
	codecsPtr := i.formatCodecs()
	if codecsPtr == nil {
		return ""
	}

	return *codecsPtr
}

func (i *ImageStreamItem) Validate() []error {
	var errs []error

	if i.Resolution == nil {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", ResolutionTag))
	}
	if i.AverageBandwidth == nil {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", AverageBandwidthTag))
	}
	if i.Bandwidth == 0 {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", BandwidthTag))
	}
	if len(i.URI) == 0 {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", URITag))
	}

	return errs
}
