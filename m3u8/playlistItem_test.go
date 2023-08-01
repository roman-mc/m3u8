package m3u8

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistItem_Parse(t *testing.T) {
	line := `#EXT-X-STREAM-INF:CODECS="avc",BANDWIDTH=540,
PROGRAM-ID=1,RESOLUTION=1920x1080,FRAME-RATE=23.976,
AVERAGE-BANDWIDTH=550,AUDIO="test",VIDEO="test2",STABLE-VARIANT-ID="1234"
SUBTITLES="subs",CLOSED-CAPTIONS="caps",URI="test.url",
NAME="1080p",HDCP-LEVEL=TYPE-0`

	pi := NewPlaylistItem(line, false)
	assertNotNilEqual(t, "1", pi.ProgramID)
	assertNotNilEqual(t, "avc", pi.Codecs)
	assert.Equal(t, 540, pi.Bandwidth)
	assertNotNilEqual(t, 550, pi.AverageBandwidth)
	assertNotNilEqual(t, 1920, pi.Width)
	assertNotNilEqual(t, 1080, pi.Height)
	assertNotNilEqual(t, 23.976, pi.FrameRate)
	assertNotNilEqual(t, "test", pi.Audio)
	assertNotNilEqual(t, "test2", pi.Video)
	assertNotNilEqual(t, "subs", pi.Subtitles)
	assertNotNilEqual(t, "caps", pi.ClosedCaptions)
	assert.Equal(t, "test.url", pi.URI)
	assertNotNilEqual(t, "1080p", pi.Name)
	assert.False(t, pi.IFrame)
	assertNotNilEqual(t, "TYPE-0", pi.HDCPLevel)
	assertNotNilEqual(t, "1234", pi.StableVariantID)
}

func TestPlaylistItem_ToString(t *testing.T) {
	// No codecs specified
	p := &PlaylistItem{
		Bandwidth: 540,
		URI:       "test.url",
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Level not recognized
	p = &PlaylistItem{
		Bandwidth: 540,
		URI:       "test.url",
		Level:     pointer.ToString("9001"),
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Audio codec recognized but profile not recognized
	p = &PlaylistItem{
		Bandwidth:  540,
		URI:        "test.url",
		Profile:    pointer.ToString("best"),
		Level:      pointer.ToString("9001"),
		AudioCodec: pointer.ToString("aac-lc"),
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Profile and level not set, Audio codec recognized
	p = &PlaylistItem{
		Bandwidth:  540,
		URI:        "test.url",
		AudioCodec: pointer.ToString("aac-lc"),
	}
	assert.Contains(t, p.String(), "CODECS")

	// Profile and level recognized, audio codec not recognized
	p = &PlaylistItem{
		Bandwidth:  540,
		URI:        "test.url",
		Profile:    pointer.ToString("high"),
		Level:      pointer.ToString("4.1"),
		AudioCodec: pointer.ToString("fuzzy"),
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Audio codec not set
	p = &PlaylistItem{
		Bandwidth: 540,
		URI:       "test.url",
		Profile:   pointer.ToString("high"),
		Level:     pointer.ToString("4.1"),
	}
	assert.Contains(t, p.String(), `CODECS="avc1.640029"`)

	// Audio codec recognized
	p = &PlaylistItem{
		Bandwidth:  540,
		URI:        "test.url",
		Profile:    pointer.ToString("high"),
		Level:      pointer.ToString("4.1"),
		AudioCodec: pointer.ToString("aac-lc"),
	}
	assert.Contains(t, p.String(), `CODECS="avc1.640029,mp4a.40.2"`)
}

func TestPlaylistItem_ToString_2(t *testing.T) {
	// All fields set
	p := &PlaylistItem{
		Codecs:           pointer.ToString("avc"),
		Bandwidth:        540,
		URI:              "test.url",
		Audio:            pointer.ToString("test"),
		Video:            pointer.ToString("test2"),
		AverageBandwidth: pointer.ToInt(500),
		Subtitles:        pointer.ToString("subs"),
		FrameRate:        pointer.ToFloat64(30),
		ClosedCaptions:   pointer.ToString("caps"),
		Name:             pointer.ToString("SD"),
		HDCPLevel:        pointer.ToString("TYPE-0"),
		ProgramID:        pointer.ToString("1"),
		StableVariantID:  pointer.ToString("1234"),
	}

	expected := `#EXT-X-STREAM-INF:PROGRAM-ID=1,CODECS="avc",BANDWIDTH=540,AVERAGE-BANDWIDTH=500,FRAME-RATE=30.000,HDCP-LEVEL=TYPE-0,AUDIO="test",VIDEO="test2",SUBTITLES="subs",CLOSED-CAPTIONS="caps",NAME="SD",STABLE-VARIANT-ID="1234"
test.url`
	assert.Equal(t, expected, p.String())

	// Closed captions is NONE
	p = &PlaylistItem{
		ProgramID:      pointer.ToString("1"),
		Width:          pointer.ToInt(1920),
		Height:         pointer.ToInt(1080),
		Codecs:         pointer.ToString("avc"),
		Bandwidth:      540,
		URI:            "test.url",
		ClosedCaptions: pointer.ToString("NONE"),
	}

	expected = `#EXT-X-STREAM-INF:PROGRAM-ID=1,RESOLUTION=1920x1080,CODECS="avc",BANDWIDTH=540,CLOSED-CAPTIONS=NONE
test.url`
	assert.Equal(t, expected, p.String())

	// IFrame is true
	p = &PlaylistItem{
		Codecs:           pointer.ToString("avc"),
		Bandwidth:        540,
		URI:              "test.url",
		IFrame:           true,
		Video:            pointer.ToString("test2"),
		AverageBandwidth: pointer.ToInt(550),
	}

	expected = `#EXT-X-I-FRAME-STREAM-INF:CODECS="avc",BANDWIDTH=540,AVERAGE-BANDWIDTH=550,VIDEO="test2",URI="test.url"`
	assert.Equal(t, expected, p.String())
}

func TestPlaylistItem_GenerateCodecs(t *testing.T) {
	assertCodecs(t, "", &PlaylistItem{})
	assertCodecs(t, "test", &PlaylistItem{Codecs: pointer.ToString("test")})
	assertCodecs(t, "mp4a.40.2", &PlaylistItem{AudioCodec: pointer.ToString("aac-lc")})
	assertCodecs(t, "mp4a.40.2", &PlaylistItem{AudioCodec: pointer.ToString("AAC-LC")})
	assertCodecs(t, "mp4a.40.5", &PlaylistItem{AudioCodec: pointer.ToString("he-aac")})
	assertCodecs(t, "", &PlaylistItem{AudioCodec: pointer.ToString("he-aac1")})
	assertCodecs(t, "mp4a.40.34", &PlaylistItem{AudioCodec: pointer.ToString("mp3")})
	assertCodecs(t, "avc1.66.30", &PlaylistItem{
		Profile: pointer.ToString("baseline"),
		Level:   pointer.ToString("3.0"),
	})
	assertCodecs(t, "avc1.66.30,mp4a.40.2", &PlaylistItem{
		Profile:    pointer.ToString("baseline"),
		Level:      pointer.ToString("3.0"),
		AudioCodec: pointer.ToString("aac-lc"),
	})
	assertCodecs(t, "avc1.66.30,mp4a.40.34", &PlaylistItem{
		Profile:    pointer.ToString("baseline"),
		Level:      pointer.ToString("3.0"),
		AudioCodec: pointer.ToString("mp3"),
	})
	assertCodecs(t, "avc1.42001f", &PlaylistItem{
		Profile: pointer.ToString("baseline"),
		Level:   pointer.ToString("3.1"),
	})
	assertCodecs(t, "avc1.42001f,mp4a.40.5", &PlaylistItem{
		Profile:    pointer.ToString("baseline"),
		Level:      pointer.ToString("3.1"),
		AudioCodec: pointer.ToString("he-aac"),
	})
	assertCodecs(t, "avc1.77.30", &PlaylistItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("3.0"),
	})
	assertCodecs(t, "avc1.77.30,mp4a.40.2", &PlaylistItem{
		Profile:    pointer.ToString("main"),
		Level:      pointer.ToString("3.0"),
		AudioCodec: pointer.ToString("aac-lc"),
	})
	assertCodecs(t, "avc1.4d001f", &PlaylistItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("3.1"),
	})
	assertCodecs(t, "avc1.4d0028", &PlaylistItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("4.0"),
	})
	assertCodecs(t, "avc1.4d0029", &PlaylistItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("4.1"),
	})
	assertCodecs(t, "avc1.64001f", &PlaylistItem{
		Profile: pointer.ToString("high"),
		Level:   pointer.ToString("3.1"),
	})
	assertCodecs(t, "avc1.640028", &PlaylistItem{
		Profile: pointer.ToString("high"),
		Level:   pointer.ToString("4.0"),
	})
	assertCodecs(t, "avc1.640029", &PlaylistItem{
		Profile: pointer.ToString("high"),
		Level:   pointer.ToString("4.1"),
	})
}

func TestPlaylistItem_Validate(t *testing.T) {
	line := `#EXT-X-STREAM-INF:CODECS="avc",
PROGRAM-ID=1,FRAME-RATE=23.976,
AUDIO="test",VIDEO="test2",STABLE-VARIANT-ID="1234"
SUBTITLES="subs",CLOSED-CAPTIONS="caps",
NAME="1080p",HDCP-LEVEL=TYPE-0`

	pi := NewPlaylistItem(line, false)
	assertNotNilEqual(t, "1", pi.ProgramID)
	assertNotNilEqual(t, "avc", pi.Codecs)
	assertNotNilEqual(t, 23.976, pi.FrameRate)
	assertNotNilEqual(t, "test", pi.Audio)
	assertNotNilEqual(t, "test2", pi.Video)
	assertNotNilEqual(t, "subs", pi.Subtitles)
	assertNotNilEqual(t, "caps", pi.ClosedCaptions)
	assertNotNilEqual(t, "1080p", pi.Name)
	assert.False(t, pi.IFrame)
	assertNotNilEqual(t, "TYPE-0", pi.HDCPLevel)
	assertNotNilEqual(t, "1234", pi.StableVariantID)

	require.Equal(t,
		[]error{
			fmt.Errorf("%s attribute is not valid", ResolutionTag),
			fmt.Errorf("%s attribute is not valid", AverageBandwidthTag),
			fmt.Errorf("%s attribute is not valid", BandwidthTag),
			fmt.Errorf("%s attribute is not valid", URITag),
		},
		pi.Validate())
}

func assertCodecs(t *testing.T, codecs string, p *PlaylistItem) {
	assert.Equal(t, codecs, p.CodecsString())
}
