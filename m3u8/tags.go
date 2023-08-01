package m3u8

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// Item tags

	SessionKeyItemTag    = `#EXT-X-SESSION-KEY`
	KeyItemTag           = `#EXT-X-KEY`
	DiscontinuityItemTag = `#EXT-X-DISCONTINUITY`
	TimeItemTag          = `#EXT-X-PROGRAM-DATE-TIME`
	DateRangeItemTag     = `#EXT-X-DATERANGE`
	MapItemTag           = `#EXT-X-MAP`
	SessionDataItemTag   = `#EXT-X-SESSION-DATA`
	SegmentItemTag       = `#EXTINF`
	ByteRangeItemTag     = `#EXT-X-BYTERANGE`
	PlaybackStartTag     = `#EXT-X-START`
	MediaItemTag         = `#EXT-X-MEDIA`
	PlaylistItemTag      = `#EXT-X-STREAM-INF`
	PlaylistIframeTag    = `#EXT-X-I-FRAME-STREAM-INF`
	DefineTag            = "#EXT-X-DEFINE"
	SCTE35Tag            = "#EXT-X-SCTE35"
	ImageStreamItemTag   = "#EXT-X-IMAGE-STREAM-INF"

	// Playlist tags

	HeaderTag                = `#EXTM3U`
	FooterTag                = `#EXT-X-ENDLIST`
	TargetDurationTag        = `#EXT-X-TARGETDURATION`
	CacheTag                 = `#EXT-X-ALLOW-CACHE`
	DiscontinuitySequenceTag = `#EXT-X-DISCONTINUITY-SEQUENCE`
	IndependentSegmentsTag   = `#EXT-X-INDEPENDENT-SEGMENTS`
	PlaylistTypeTag          = `#EXT-X-PLAYLIST-TYPE`
	IFramesOnlyTag           = `#EXT-X-I-FRAMES-ONLY`
	MediaSequenceTag         = `#EXT-X-MEDIA-SEQUENCE`
	VersionTag               = `#EXT-X-VERSION`

	// ByteRange tags

	ByteRangeTag = "BYTERANGE"

	// Encryptable tags

	MethodTag            = "METHOD"
	URITag               = "URI"
	IVTag                = "IV"
	KeyFormatTag         = "KEYFORMAT"
	KeyFormatVersionsTag = "KEYFORMATVERSIONS"
	KeyID                = "KEYID"

	// DateRangeItem tags

	IDTag              = "ID"
	ClassTag           = "CLASS"
	CueTag             = "CUE"
	StartDateTag       = "START-DATE"
	EndDateTag         = "END-DATE"
	DurationTag        = "DURATION"
	PlannedDurationTag = "PLANNED-DURATION"
	Scte35CmdTag       = "SCTE35-CMD"
	Scte35OutTag       = "SCTE35-OUT"
	Scte35InTag        = "SCTE35-IN"
	EndOnNextTag       = "END-ON-NEXT"

	// SCTE35Item tags (SCTE35Item has DateRangeItem tags + its own)

	ElapsedAttribute  = "ELAPSED"
	TimeAttribute     = "TIME"
	UPIDAttribute     = "UPID"
	BlackoutAttribute = "BLACKOUT"
	CueOutAttribute   = "CUE-OUT"
	CueInAttribute    = "CUE-IN"
	SegneAttribute    = "SEGNE"

	// PlaybackStart tags

	TimeOffsetTag = "TIME-OFFSET"
	PreciseTag    = "PRECISE"

	// SessionDataItem tags

	DataIDTag   = "DATA-ID"
	ValueTag    = "VALUE"
	LanguageTag = "LANGUAGE"

	// Define tags

	AttributeName       = "NAME"
	AttributeValue      = "VALUE"
	AttributeImport     = "IMPORT"
	AttributeQueryParam = "QUERYPARAM"

	// MediaItem tags

	TypeTag              = "TYPE"
	GroupIDTag           = "GROUP-ID"
	AssocLanguageTag     = "ASSOC-LANGUAGE"
	NameTag              = "NAME"
	AutoSelectTag        = "AUTOSELECT"
	DefaultTag           = "DEFAULT"
	ForcedTag            = "FORCED"
	InStreamIDTag        = "INSTREAM-ID"
	CharacteristicsTag   = "CHARACTERISTICS"
	ChannelsTag          = "CHANNELS"
	StableRenditionIDTag = "STABLE-RENDITION-ID"

	// PlaylistItem tags

	ResolutionTag       = "RESOLUTION"
	ProgramIDTag        = "PROGRAM-ID"
	CodecsTag           = "CODECS"
	BandwidthTag        = "BANDWIDTH"
	AverageBandwidthTag = "AVERAGE-BANDWIDTH"
	FrameRateTag        = "FRAME-RATE"
	VideoTag            = "VIDEO"
	AudioTag            = "AUDIO"
	SubtitlesTag        = "SUBTITLES"
	ClosedCaptionsTag   = "CLOSED-CAPTIONS"
	HDCPLevelTag        = "HDCP-LEVEL"
	StableVariantIDTag  = "STABLE-VARIANT-ID"
)

var (
	ErrTagIsNotImplemented = errors.New("tag is not implemented")
)

var (
	notImplementedReadLine = func(line string, pl *Playlist, st *state) error {
		return ErrTagIsNotImplemented
	}
	noopReadLine = func(line string, pl *Playlist, st *state) error {
		return nil
	}
)

var tagsMap = map[string]struct {
	Attributes []string
	ReadLine   func(line string, pl *Playlist, st *state) error
}{
	"#EXT-X-BITRATE": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	ByteRangeItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			value := strings.Replace(line, ByteRangeItemTag+":", "", -1)
			value = strings.Replace(value, "\n", "", -1)
			br, err := NewByteRange(value)
			if err != nil {
				return parseError(line, err)
			}
			mit, ok := st.currentItem.(*MapItem)
			if ok {
				mit.ByteRange = br
				st.currentItem = mit
			} else {
				sit, ok := st.currentItem.(*SegmentItem)
				if ok {
					sit.ByteRange = br
					st.currentItem = sit
				}
			}

			return nil
		},
	},
	"#EXT-X-CONTENT-STEERING": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	DateRangeItemTag: {
		Attributes: []string{
			IDTag,
			ClassTag,
			StartDateTag,
			CueTag,
			EndDateTag,
			DurationTag,
			PlannedDurationTag,
		},
		ReadLine: func(line string, pl *Playlist, st *state) error {
			dri := NewDateRangeItem(line)
			pl.Items = append(pl.Items, dri)
			return nil
		},
	},
	DefineTag: {
		Attributes: []string{
			AttributeName,
			AttributeValue,
			AttributeImport,
			AttributeQueryParam,
		},
		ReadLine: func(line string, pl *Playlist, st *state) error {
			item := NewDefineItem(line)
			pl.Items = append(pl.Items, item)
			return nil
		},
	},
	DiscontinuityItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			st.master = false
			st.open = false
			item, err := NewDiscontinuityItem()
			if err != nil {
				return parseError(line, err)
			}
			pl.Items = append(pl.Items, item)
			return nil
		},
	},
	DiscontinuitySequenceTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			var err error
			pl.DiscontinuitySequence, err = parseIntPtr(line, DiscontinuitySequenceTag)
			return err
		},
	},
	FooterTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			pl.Live = false
			return nil
		},
	},
	"#EXT-X-GAP": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	PlaylistIframeTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			st.master = true
			st.open = false
			pi := NewPlaylistItem(line, true)
			pl.Items = append(pl.Items, pi)
			st.currentItem = pi
			return nil
		},
	},
	IFramesOnlyTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			pl.IFramesOnly = true
			return nil
		},
	},
	IndependentSegmentsTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			pl.IndependentSegments = true
			return nil
		},
	},
	KeyItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			item := NewKeyItem(line)
			pl.Items = append(pl.Items, item)
			return nil
		},
	},
	MapItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			item := NewMapItem(line)
			pl.Items = append(pl.Items, item)

			return nil
		},
	},
	MediaItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			st.open = false
			mi := NewMediaItem(line)
			pl.Items = append(pl.Items, mi)
			return nil
		},
	},
	MediaSequenceTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			var err error
			pl.Sequence, err = parseIntValue(line, MediaSequenceTag)
			return err
		},
	},
	"#EXT-X-PART": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	"#EXT-X-PART-INF": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	PlaylistTypeTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			pl.Type = parseStringPtr(line, PlaylistTypeTag)
			return nil
		},
	},
	"#EXT-X-PRELOAD-HINT": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	TimeItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			pdt, err := NewTimeItem(line)
			if err != nil {
				return parseError(line, err)
			}
			if st.open {
				item, ok := st.currentItem.(*SegmentItem)
				if !ok {
					return parseError(line, ErrSegmentItemInvalid)
				}
				item.ProgramDateTime = pdt
			} else {
				pl.Items = append(pl.Items, pdt)
			}
			return nil
		},
	},
	"#EXT-X-RENDITION-REPORT": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	"#EXT-X-SERVER-CONTROL": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	SessionDataItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			sdi := NewSessionDataItem(line)
			pl.Items = append(pl.Items, sdi)
			return nil
		},
	},
	SessionKeyItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			ski := NewSessionKeyItem(line)
			pl.Items = append(pl.Items, ski)
			return nil
		},
	},
	"#EXT-X-SKIP": {
		ReadLine: notImplementedReadLine,
	}, // TODO
	PlaybackStartTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			ps, err := NewPlaybackStart(line)
			if err != nil {
				return parseError(line, err)
			}
			pl.Items = append(pl.Items, ps)
			return nil
		},
	},
	PlaylistItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			st.master = true
			st.open = true
			pi := NewPlaylistItem(line, false)
			st.currentItem = pi

			return nil
		},
	},
	TargetDurationTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			var err error
			pl.Target, err = parseIntValue(line, TargetDurationTag)

			return err
		},
	},
	VersionTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			var err error
			var version *int
			version, err = parseIntPtr(line, VersionTag)
			if err != nil {
				return err
			}

			pl.Version = version
			return nil
		},
	},
	SegmentItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			var err error
			st.currentItem, err = NewSegmentItem(line)
			st.master = false
			st.open = true

			return err

		},
	},
	HeaderTag: {ReadLine: noopReadLine},
	// CacheTag is a non-standard tag
	CacheTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			ptr := ParseYesNoPtr(line, CacheTag)
			pl.Cache = ptr
			return nil
		},
	},
	// SCTE35Tag is a non-standard tag
	SCTE35Tag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			item, err := NewSCTE35Item(line)
			if err != nil {
				return parseError(line, err)
			}

			pl.Items = append(pl.Items, item)

			return nil
		},
	},
	// ImageStreamItemTag is a non-standard tag
	ImageStreamItemTag: {
		ReadLine: func(line string, pl *Playlist, st *state) error {
			item := NewImageStreamItem(line)
			pl.Items = append(pl.Items, item)
			return nil
		},
	},
}

type Tag interface {
	fmt.Stringer
	Validate() []error
}
