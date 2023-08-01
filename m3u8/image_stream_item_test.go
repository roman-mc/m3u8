package m3u8

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

func TestNewImageStreamItem(t *testing.T) {
	attributesKV := []string{
		"BANDWIDTH=540",
		"URI=\"test.url\"",
		"PROGRAM-ID=1",
		"RESOLUTION=1920x1080",
		"AVERAGE-BANDWIDTH=550",
		"CODECS=\"mp4a.40.34\"",
		"VIDEO=\"test2\"",
		"AUDIO=\"test\"",
		"SUBTITLES=\"subs\"",
		"CLOSED-CAPTIONS=\"caps\"",
		"NAME=\"1080p\"",
		"FRAME-RATE=23.976",
		"STABLE-VARIANT-ID=\"1234\"",
		"RANDOM-ATTRIBUTE=123",
		"RANDOM-ATTRIBUTE2=\"123\"",
	}

	line := "EXT-X-IMAGE-STREAM-INF:" + strings.Join(attributesKV, ",")

	item := NewImageStreamItem(line)
	require.Len(t, item.attributes, 2)
	require.Equal(t, "123", item.attributes["RANDOM-ATTRIBUTE"])
	require.Equal(t, "\"123\"", item.attributes["RANDOM-ATTRIBUTE2"])

	require.Equal(t, 540, item.Bandwidth)
	require.Equal(t, "test.url", item.URI)
	require.Equal(t, pointer.ToString("mp4a.40.34"), item.Codecs)
	require.Equal(t, pointer.ToString("1"), item.ProgramID)
	require.Equal(t, &parser.Resolution{Width: 1920, Height: 1080}, item.Resolution)
	require.Equal(t, pointer.ToInt(550), item.AverageBandwidth)
	require.Equal(t, pointer.ToString("test2"), item.Video)
	require.Equal(t, pointer.ToString("test"), item.Audio)
	require.Equal(t, pointer.ToString("subs"), item.Subtitles)
	require.Equal(t, pointer.ToString("caps"), item.ClosedCaptions)
	require.Equal(t, pointer.ToString("1080p"), item.Name)
	require.Equal(t, pointer.ToFloat64(23.976), item.FrameRate)
	require.Equal(t, pointer.ToString("1234"), item.StableVariantID)

	itemEncoded := item.String()
	for _, attributeKV := range attributesKV {
		require.True(t, strings.Contains(itemEncoded, attributeKV), attributeKV, itemEncoded)
	}
}

func TestImageStreamItem_Parse(t *testing.T) {
	line := `#EXT-X-IMAGE-STREAM-INF:CODECS="avc",BANDWIDTH=540,
PROGRAM-ID=1,RESOLUTION=1920x1080,FRAME-RATE=23.976,
AVERAGE-BANDWIDTH=550,AUDIO="test",VIDEO="test2",STABLE-VARIANT-ID="1234"
SUBTITLES="subs",CLOSED-CAPTIONS="caps",URI="test.url",
NAME="1080p",HDCP-LEVEL=TYPE-0`

	pi := NewImageStreamItem(line)
	assertNotNilEqual(t, "1", pi.ProgramID)
	assertNotNilEqual(t, "avc", pi.Codecs)
	assert.Equal(t, 540, pi.Bandwidth)
	assertNotNilEqual(t, 550, pi.AverageBandwidth)
	require.Equal(t, 1920, pi.Resolution.Width)
	require.Equal(t, 1080, pi.Resolution.Height)
	assertNotNilEqual(t, 23.976, pi.FrameRate)
	assertNotNilEqual(t, "test", pi.Audio)
	assertNotNilEqual(t, "test2", pi.Video)
	assertNotNilEqual(t, "subs", pi.Subtitles)
	assertNotNilEqual(t, "caps", pi.ClosedCaptions)
	assert.Equal(t, "test.url", pi.URI)
	assertNotNilEqual(t, "1080p", pi.Name)
	assertNotNilEqual(t, "1234", pi.StableVariantID)
}

func TestImageStreamItem_ToString(t *testing.T) {
	// No codecs specified
	p := &ImageStreamItem{
		Bandwidth: 540,
		URI:       "test.url",
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Level not recognized
	p = &ImageStreamItem{
		Bandwidth: 540,
		URI:       "test.url",
		Level:     pointer.ToString("9001"),
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Audio codec recognized but profile not recognized
	p = &ImageStreamItem{
		Bandwidth:  540,
		URI:        "test.url",
		Profile:    pointer.ToString("best"),
		Level:      pointer.ToString("9001"),
		AudioCodec: pointer.ToString("aac-lc"),
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Profile and level not set, Audio codec recognized
	p = &ImageStreamItem{
		Bandwidth:  540,
		URI:        "test.url",
		AudioCodec: pointer.ToString("aac-lc"),
	}
	assert.Contains(t, p.String(), "CODECS")

	// Profile and level recognized, audio codec not recognized
	p = &ImageStreamItem{
		Bandwidth:  540,
		URI:        "test.url",
		Profile:    pointer.ToString("high"),
		Level:      pointer.ToString("4.1"),
		AudioCodec: pointer.ToString("fuzzy"),
	}
	assert.NotContains(t, p.String(), "CODECS")

	// Audio codec not set
	p = &ImageStreamItem{
		Bandwidth: 540,
		URI:       "test.url",
		Profile:   pointer.ToString("high"),
		Level:     pointer.ToString("4.1"),
	}
	assert.Contains(t, p.String(), `CODECS="avc1.640029"`)

	// Audio codec recognized
	p = &ImageStreamItem{
		Bandwidth:  540,
		URI:        "test.url",
		Profile:    pointer.ToString("high"),
		Level:      pointer.ToString("4.1"),
		AudioCodec: pointer.ToString("aac-lc"),
	}
	assert.Contains(t, p.String(), `CODECS="avc1.640029,mp4a.40.2"`)
}

func TestImageStreamItem_ToString_2(t *testing.T) {
	// All fields set
	p := &ImageStreamItem{
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
		ProgramID:        pointer.ToString("1"),
		StableVariantID:  pointer.ToString("1234"),
	}

	expected := `#EXT-X-IMAGE-STREAM-INF:PROGRAM-ID=1,CODECS="avc",BANDWIDTH=540,AVERAGE-BANDWIDTH=500,FRAME-RATE=30.000,AUDIO="test",VIDEO="test2",SUBTITLES="subs",CLOSED-CAPTIONS="caps",NAME="SD",STABLE-VARIANT-ID="1234",URI="test.url"`
	assert.Equal(t, expected, p.String())

	// Closed captions is NONE
	p = &ImageStreamItem{
		ProgramID: pointer.ToString("1"),
		Resolution: &parser.Resolution{
			Width:  1920,
			Height: 1080,
		},
		Codecs:         pointer.ToString("avc"),
		Bandwidth:      540,
		URI:            "test.url",
		ClosedCaptions: pointer.ToString("NONE"),
	}

	expected = `#EXT-X-IMAGE-STREAM-INF:PROGRAM-ID=1,RESOLUTION=1920x1080,CODECS="avc",BANDWIDTH=540,CLOSED-CAPTIONS=NONE,URI="test.url"`
	assert.Equal(t, expected, p.String())

	p = &ImageStreamItem{
		Codecs:           pointer.ToString("avc"),
		Bandwidth:        540,
		URI:              "test.url",
		Video:            pointer.ToString("test2"),
		AverageBandwidth: pointer.ToInt(550),
	}

	expected = `#EXT-X-IMAGE-STREAM-INF:CODECS="avc",BANDWIDTH=540,AVERAGE-BANDWIDTH=550,VIDEO="test2",URI="test.url"`
	assert.Equal(t, expected, p.String())
}

func TestImageStreamItem_GenerateCodecs(t *testing.T) {
	assertCodecsImageStream(t, "", &ImageStreamItem{})
	assertCodecsImageStream(t, "test", &ImageStreamItem{Codecs: pointer.ToString("test")})
	assertCodecsImageStream(t, "mp4a.40.2", &ImageStreamItem{AudioCodec: pointer.ToString("aac-lc")})
	assertCodecsImageStream(t, "mp4a.40.2", &ImageStreamItem{AudioCodec: pointer.ToString("AAC-LC")})
	assertCodecsImageStream(t, "mp4a.40.5", &ImageStreamItem{AudioCodec: pointer.ToString("he-aac")})
	assertCodecsImageStream(t, "", &ImageStreamItem{AudioCodec: pointer.ToString("he-aac1")})
	assertCodecsImageStream(t, "mp4a.40.34", &ImageStreamItem{AudioCodec: pointer.ToString("mp3")})
	assertCodecsImageStream(t, "avc1.66.30", &ImageStreamItem{
		Profile: pointer.ToString("baseline"),
		Level:   pointer.ToString("3.0"),
	})
	assertCodecsImageStream(t, "avc1.66.30,mp4a.40.2", &ImageStreamItem{
		Profile:    pointer.ToString("baseline"),
		Level:      pointer.ToString("3.0"),
		AudioCodec: pointer.ToString("aac-lc"),
	})
	assertCodecsImageStream(t, "avc1.66.30,mp4a.40.34", &ImageStreamItem{
		Profile:    pointer.ToString("baseline"),
		Level:      pointer.ToString("3.0"),
		AudioCodec: pointer.ToString("mp3"),
	})
	assertCodecsImageStream(t, "avc1.42001f", &ImageStreamItem{
		Profile: pointer.ToString("baseline"),
		Level:   pointer.ToString("3.1"),
	})
	assertCodecsImageStream(t, "avc1.42001f,mp4a.40.5", &ImageStreamItem{
		Profile:    pointer.ToString("baseline"),
		Level:      pointer.ToString("3.1"),
		AudioCodec: pointer.ToString("he-aac"),
	})
	assertCodecsImageStream(t, "avc1.77.30", &ImageStreamItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("3.0"),
	})
	assertCodecsImageStream(t, "avc1.77.30,mp4a.40.2", &ImageStreamItem{
		Profile:    pointer.ToString("main"),
		Level:      pointer.ToString("3.0"),
		AudioCodec: pointer.ToString("aac-lc"),
	})
	assertCodecsImageStream(t, "avc1.4d001f", &ImageStreamItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("3.1"),
	})
	assertCodecsImageStream(t, "avc1.4d0028", &ImageStreamItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("4.0"),
	})
	assertCodecsImageStream(t, "avc1.4d0029", &ImageStreamItem{
		Profile: pointer.ToString("main"),
		Level:   pointer.ToString("4.1"),
	})
	assertCodecsImageStream(t, "avc1.64001f", &ImageStreamItem{
		Profile: pointer.ToString("high"),
		Level:   pointer.ToString("3.1"),
	})
	assertCodecsImageStream(t, "avc1.640028", &ImageStreamItem{
		Profile: pointer.ToString("high"),
		Level:   pointer.ToString("4.0"),
	})
	assertCodecsImageStream(t, "avc1.640029", &ImageStreamItem{
		Profile: pointer.ToString("high"),
		Level:   pointer.ToString("4.1"),
	})
}

func TestImageStreamItem_Validate(t *testing.T) {
	line := `#EXT-X-IMAGE-STREAM-INF:CODECS="avc",
PROGRAM-ID=1,FRAME-RATE=23.976,
AUDIO="test",VIDEO="test2",STABLE-VARIANT-ID="1234"
SUBTITLES="subs",CLOSED-CAPTIONS="caps",
NAME="1080p",HDCP-LEVEL=TYPE-0`

	pi := NewImageStreamItem(line)
	assertNotNilEqual(t, "1", pi.ProgramID)
	assertNotNilEqual(t, "avc", pi.Codecs)
	assertNotNilEqual(t, 23.976, pi.FrameRate)
	assertNotNilEqual(t, "test", pi.Audio)
	assertNotNilEqual(t, "test2", pi.Video)
	assertNotNilEqual(t, "subs", pi.Subtitles)
	assertNotNilEqual(t, "caps", pi.ClosedCaptions)
	assertNotNilEqual(t, "1080p", pi.Name)
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

func assertCodecsImageStream(t *testing.T, codecs string, p *ImageStreamItem) {
	assert.Equal(t, codecs, p.CodecsString())
}
