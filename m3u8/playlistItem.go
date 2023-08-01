package m3u8

import (
	"fmt"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// PlaylistItem represents a set of EXT-X-STREAM-INF or
// EXT-X-I-FRAME-STREAM-INF attributes
type PlaylistItem struct {
	Bandwidth int
	URI       string
	IFrame    bool

	Name             *string
	Width            *int // TODO: review this attribute, useless
	Height           *int // TODO: review this attribute, useless
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
	HDCPLevel        *string
	Resolution       *parser.Resolution
	StableVariantID  *string
	attributes       map[string]string
}

// NewPlaylistItem parses a text line and returns a *PlaylistItem
func NewPlaylistItem(text string, isIframe bool) *PlaylistItem {
	attributes := parser.ParseAttributes(text)
	resolution, _ := parser.ParseResolution(attributes, ResolutionTag)

	var width, height *int
	if resolution != nil {
		width = &resolution.Width
		height = &resolution.Height
	}

	averageBandwidth, _ := parser.ParseInt(attributes, AverageBandwidthTag)
	frameRate, _ := parser.ParseFloat(attributes, FrameRateTag)
	if frameRate != nil && *frameRate <= 0 {
		frameRate = nil
	}

	bandwidth, _ := parser.ParseBandwidth(attributes, BandwidthTag)

	defer deleteKeys(attributes,
		ResolutionTag,
		AverageBandwidthTag,
		FrameRateTag,
		BandwidthTag,
		ProgramIDTag,
		CodecsTag,
		VideoTag,
		AudioTag,
		URITag,
		SubtitlesTag,
		ClosedCaptionsTag,
		NameTag,
		HDCPLevelTag,
		StableVariantIDTag,
	)

	return &PlaylistItem{
		ProgramID:        parser.PointerTo(attributes, ProgramIDTag),
		Codecs:           parser.PointerTo(attributes, CodecsTag),
		Width:            width,
		Height:           height,
		Bandwidth:        bandwidth,
		AverageBandwidth: averageBandwidth,
		FrameRate:        frameRate,
		Video:            parser.PointerTo(attributes, VideoTag),
		Audio:            parser.PointerTo(attributes, AudioTag),
		URI:              parser.SanitizeAttributeValue(attributes[URITag]),
		Subtitles:        parser.PointerTo(attributes, SubtitlesTag),
		ClosedCaptions:   parser.PointerTo(attributes, ClosedCaptionsTag),
		Name:             parser.PointerTo(attributes, NameTag),
		HDCPLevel:        parser.PointerTo(attributes, HDCPLevelTag),
		Resolution:       resolution,
		StableVariantID:  parser.PointerTo(attributes, StableVariantIDTag),
		IFrame:           isIframe,
		attributes:       attributes,
	}
}

func (pi *PlaylistItem) String() string {
	var slice []string
	// Check resolution
	if pi.Resolution == nil && pi.Width != nil && pi.Height != nil {
		r := &parser.Resolution{
			Width:  *pi.Width,
			Height: *pi.Height,
		}
		pi.Resolution = r
	}
	if pi.ProgramID != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, ProgramIDTag, *pi.ProgramID))
	}
	if pi.Resolution != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, ResolutionTag, pi.Resolution.String()))
	}
	codecs := formatCodecs(pi)
	if codecs != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, CodecsTag, *codecs))
	}
	slice = append(slice, fmt.Sprintf(parser.FormatString, BandwidthTag, pi.Bandwidth))
	if pi.AverageBandwidth != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, AverageBandwidthTag, *pi.AverageBandwidth))
	}
	if pi.FrameRate != nil {
		slice = append(slice, fmt.Sprintf(parser.FrameRateFormatString, FrameRateTag, *pi.FrameRate))
	}
	if pi.HDCPLevel != nil {
		slice = append(slice, fmt.Sprintf(parser.FormatString, HDCPLevelTag, *pi.HDCPLevel))
	}
	if pi.Audio != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, AudioTag, *pi.Audio))
	}
	if pi.Video != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, VideoTag, *pi.Video))
	}
	if pi.Subtitles != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, SubtitlesTag, *pi.Subtitles))
	}
	if pi.ClosedCaptions != nil {
		cc := *pi.ClosedCaptions
		fs := parser.QuotedFormatString
		if cc == parser.NoneValue {
			fs = parser.FormatString
		}
		slice = append(slice, fmt.Sprintf(fs, ClosedCaptionsTag, cc))
	}
	if pi.Name != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, NameTag, *pi.Name))
	}
	if pi.StableVariantID != nil {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, StableVariantIDTag, *pi.StableVariantID))
	}

	var uriLine string
	itemTag := PlaylistItemTag

	if pi.IFrame {
		slice = append(slice, fmt.Sprintf(parser.QuotedFormatString, URITag, pi.URI))
		itemTag = PlaylistIframeTag
	} else {
		uriLine = "\n" + pi.URI
	}
	slice = attributesJoinMap(slice, pi.attributes)
	attributesString := strings.Join(slice, ",")

	return fmt.Sprintf("%s:%s%s", itemTag, attributesString, uriLine)
}

// CodecsString returns the string representation of codecs for a playlist item
func (pi *PlaylistItem) CodecsString() string {
	codecsPtr := formatCodecs(pi)
	if codecsPtr == nil {
		return ""
	}

	return *codecsPtr
}

func formatCodecs(pi *PlaylistItem) *string {
	if pi.Codecs != nil {
		return pi.Codecs
	}

	videoCodecPtr := videoCodec(pi.Profile, pi.Level)
	// profile or level were specified but not recognized any codecs
	if !(pi.Profile == nil && pi.Level == nil) && videoCodecPtr == nil {
		return nil
	}

	audioCodecPtr := audioCodec(pi.AudioCodec)
	// audio codec was specified but not recognized
	if !(pi.AudioCodec == nil) && audioCodecPtr == nil {
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

func (pi *PlaylistItem) Validate() []error {
	var errs []error

	if pi.Resolution == nil {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", ResolutionTag))
	}
	if pi.AverageBandwidth == nil {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", AverageBandwidthTag))
	}
	if pi.Bandwidth == 0 {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", BandwidthTag))
	}
	if len(pi.URI) == 0 {
		errs = append(errs, fmt.Errorf("%s attribute is not valid", URITag))
	}

	return errs
}
