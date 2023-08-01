package m3u8

import (
	"fmt"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestPlaylist_New(t *testing.T) {
	p := &Playlist{Master: pointer.ToBool(true)}
	assert.True(t, p.IsMaster())

	p, err := ReadFile("fixtures/master.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsMaster())
	assert.Equal(t, len(p.Items), 8)
}

func TestPlaylist_Duration(t *testing.T) {
	p := &Playlist{
		Items: []Item{
			&SegmentItem{Duration: 10.991, Segment: "test_01.ts"},
			&SegmentItem{Duration: 9.891, Segment: "test_02.ts"},
			&SegmentItem{Duration: 10.556, Segment: "test_03.ts"},
			&SegmentItem{Duration: 8.790, Segment: "test_04.ts"},
		},
	}

	assert.Equal(t, "40.228", fmt.Sprintf("%.3f", p.Duration()))
}

func TestPlaylist_Master(t *testing.T) {
	// Normal master playlist
	p := &Playlist{
		Items: []Item{
			&PlaylistItem{
				ProgramID:  pointer.ToString("1"),
				URI:        "playlist_url",
				Bandwidth:  6400,
				AudioCodec: pointer.ToString("mp3"),
			},
		},
	}
	assert.True(t, p.IsMaster())

	// Media playlist
	p = &Playlist{
		Items: []Item{
			&SegmentItem{Duration: 10.991, Segment: "test_01.ts"},
		},
	}
	assert.False(t, p.IsMaster())

	// Forced master tag
	p = &Playlist{
		Master: pointer.ToBool(true),
	}
	assert.True(t, p.IsMaster())
}

func TestPlaylist_Live(t *testing.T) {
	// Normal master playlist
	p := &Playlist{
		Items: []Item{
			&PlaylistItem{
				ProgramID:  pointer.ToString("1"),
				URI:        "playlist_url",
				Bandwidth:  6400,
				AudioCodec: pointer.ToString("mp3"),
			},
		},
	}
	assert.False(t, p.IsLive())

	// Media playlist set as live
	p = &Playlist{
		Items: []Item{
			&SegmentItem{Duration: 10.991, Segment: "test_01.ts"},
		},
		Live: true,
	}
	assert.True(t, p.IsLive())
}

func TestPlaylist_ToString(t *testing.T) {
	p := &Playlist{
		Items: []Item{
			&PlaylistItem{
				ProgramID:  pointer.ToString("1"),
				URI:        "playlist_url",
				Bandwidth:  6400,
				AudioCodec: pointer.ToString("mp3"),
			},
			&PlaylistItem{
				ProgramID:  pointer.ToString("2"),
				URI:        "playlist_url",
				Bandwidth:  50000,
				Width:      pointer.ToInt(1920),
				Height:     pointer.ToInt(1080),
				Profile:    pointer.ToString("high"),
				Level:      pointer.ToString("4.1"),
				AudioCodec: pointer.ToString("aac-lc"),
			},
		},
	}

	expected := `#EXTM3U
#EXT-X-STREAM-INF:PROGRAM-ID=1,CODECS="mp4a.40.34",BANDWIDTH=6400
playlist_url
#EXT-X-STREAM-INF:PROGRAM-ID=2,RESOLUTION=1920x1080,CODECS="avc1.640029,mp4a.40.2",BANDWIDTH=50000
playlist_url
`

	assert.Equal(t, expected, p.String())

	p = NewPlaylistWithItems(
		[]Item{
			&SegmentItem{Duration: 11.344644, Segment: "1080-7mbps00000.ts"},
			&SegmentItem{Duration: 11.261233, Segment: "1080-7mbps00001.ts"},
		},
	)
	expected = `#EXTM3U
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-TARGETDURATION:10
#EXTINF:11.344644,
1080-7mbps00000.ts
#EXTINF:11.261233,
1080-7mbps00001.ts
#EXT-X-ENDLIST
`

	assert.Equal(t, expected, p.String())
}

func TestPlaylist_Valid(t *testing.T) {
	p := NewPlaylist()
	assert.True(t, p.IsValid())

	p.AppendItem(&PlaylistItem{
		ProgramID: pointer.ToString("1"),
		URI:       "playlist_url",
		Bandwidth: 540,
		Width:     pointer.ToInt(1920),
		Height:    pointer.ToInt(1080),
		Codecs:    pointer.ToString("avc"),
	})

	assert.True(t, p.IsValid())
	assert.Equal(t, 1, len(p.Items))

	p.AppendItem(&PlaylistItem{
		ProgramID: pointer.ToString("1"),
		URI:       "playlist_url",
		Bandwidth: 540,
		Width:     pointer.ToInt(1920),
		Height:    pointer.ToInt(1080),
		Codecs:    pointer.ToString("avc"),
	})

	assert.True(t, p.IsValid())
	assert.Equal(t, 2, len(p.Items))

	p.AppendItem(&SegmentItem{
		Duration: 10.991,
		Segment:  "test.ts",
	})

	assert.False(t, p.IsValid())
}

func TestPlaylist_PlaylistSize(t *testing.T) {
	p := NewPlaylist()
	assert.True(t, p.IsValid())

	p.AppendItem(&PlaylistItem{
		ProgramID: pointer.ToString("1"),
		URI:       "playlist0_url",
		Bandwidth: 540,
		Width:     pointer.ToInt(1920),
		Height:    pointer.ToInt(1080),
		Codecs:    pointer.ToString("avc"),
	})

	p.AppendItem(&PlaylistItem{
		ProgramID: pointer.ToString("1"),
		URI:       "playlist1_url",
		Bandwidth: 540,
		Width:     pointer.ToInt(1920),
		Height:    pointer.ToInt(1080),
		Codecs:    pointer.ToString("avc"),
	})

	assert.Equal(t, 2, p.PlaylistSize())
	pi := p.Playlists()
	assert.Equal(t, "playlist0_url", pi[0].URI)
	assert.Equal(t, "playlist1_url", pi[1].URI)

}

func TestPlaylist_Segments(t *testing.T) {
	p := &Playlist{
		Items: []Item{
			&SegmentItem{Duration: 10.991, Segment: "test_01.ts"},
			&SegmentItem{Duration: 9.891, Segment: "test_02.ts"},
			&SegmentItem{Duration: 10.556, Segment: "test_03.ts"},
			&SegmentItem{Duration: 8.790, Segment: "test_04.ts"},
		},
	}

	assert.Equal(t, 4, p.SegmentSize())
	si := p.Segments()
	assert.Equal(t, "test_01.ts", si[0].Segment)
	assert.Equal(t, "test_02.ts", si[1].Segment)
	assert.Equal(t, 10.556, si[2].Duration)
	assert.Equal(t, 8.790, si[3].Duration)

}
