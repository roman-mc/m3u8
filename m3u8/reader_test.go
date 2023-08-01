package m3u8

import (
	"github.com/stretchr/testify/require"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	p, err := ReadFile("fixtures/master.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.True(t, p.IsMaster())
	assert.Nil(t, p.DiscontinuitySequence)
	assert.True(t, p.IndependentSegments)

	item := p.Items[0]
	assert.IsType(t, &SessionKeyItem{}, item)
	keyItem := item.(*SessionKeyItem)
	assert.Equal(t, "AES-128", keyItem.Encryptable.Method)
	assertNotNilEqual(t, "https://priv.example.com/key.php?r=52", keyItem.Encryptable.URI)

	item = p.Items[1]
	assert.IsType(t, &PlaybackStart{}, item)
	psi := item.(*PlaybackStart)
	assert.Equal(t, 20.2, psi.TimeOffset)

	item = p.Items[2]
	assert.IsType(t, &PlaylistItem{}, item)
	pi := item.(*PlaylistItem)
	assert.Equal(t, "hls/1080-7mbps/1080-7mbps.m3u8", pi.URI)
	assertNotNilEqual(t, "1", pi.ProgramID)
	assertNotNilEqual(t, 1920, pi.Width)
	assertNotNilEqual(t, 1080, pi.Height)
	assert.Equal(t, "1920x1080", pi.Resolution.String())
	assert.Equal(t, "avc1.640028,mp4a.40.2", pi.CodecsString())
	assert.Equal(t, 5042000, pi.Bandwidth)
	assert.False(t, pi.IFrame)
	assert.Nil(t, pi.AverageBandwidth)

	item = p.Items[7]
	assert.IsType(t, &PlaylistItem{}, item)
	pi = item.(*PlaylistItem)
	assert.Equal(t, "hls/64k/64k.m3u8", pi.URI)
	assertNotNilEqual(t, "1", pi.ProgramID)
	assert.Nil(t, pi.Height)
	assert.Nil(t, pi.Width)
	assert.Empty(t, pi.Resolution.String())
	assert.Equal(t, 6400, pi.Bandwidth)
	assert.False(t, pi.IFrame)
	assert.Nil(t, pi.AverageBandwidth)

	assert.Equal(t, 8, p.ItemSize())

	item = p.Items[len(p.Items)-1]
	assert.IsType(t, &PlaylistItem{}, item)
	pi = item.(*PlaylistItem)
	assert.Empty(t, pi.Resolution.String())
}

func TestReader_IFrame(t *testing.T) {
	p, err := ReadFile("fixtures/masterIframes.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.True(t, p.IsMaster())

	assert.Equal(t, 7, p.ItemSize())

	item := p.Items[1]
	assert.IsType(t, &PlaylistItem{}, item)
	pi := item.(*PlaylistItem)
	assert.Equal(t, "low/iframe.m3u8", pi.URI)
	assert.Equal(t, 86000, pi.Bandwidth)
	assert.True(t, pi.IFrame)
}

func TestReader_MediaPlaylist(t *testing.T) {
	p, err := ReadFile("fixtures/playlist.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.False(t, p.IsMaster())

	assertNotNilEqual(t, 4, p.Version)
	assert.Equal(t, 1, p.Sequence)
	assertNotNilEqual(t, 8, p.DiscontinuitySequence)
	assertNotNilEqual(t, false, p.Cache)
	assert.Equal(t, 12, p.Target)
	assertNotNilEqual(t, "VOD", p.Type)

	item := p.Items[0]
	assert.IsType(t, &SegmentItem{}, item)
	si := item.(*SegmentItem)
	assert.Equal(t, 11.344644, si.Duration)
	assert.Nil(t, si.Comment)

	item = p.Items[4]
	assert.IsType(t, &TimeItem{}, item)
	ti := item.(*TimeItem)
	assert.Equal(t, "2010-02-19T14:54:23Z", FormatTime(ti.Time))

	assert.Equal(t, 140, p.ItemSize())
}

func TestReader_PlaylistLiveCheck(t *testing.T) {
	p, err := ReadFile("fixtures/playlist.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.False(t, p.IsLive())

	p, err = ReadFile("fixtures/playlist-live.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.True(t, p.IsLive())
}
func TestReader_IFramePlaylist(t *testing.T) {
	p, err := ReadFile("fixtures/iframes.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())

	assert.True(t, p.IFramesOnly)
	assert.Equal(t, 3, p.ItemSize())

	item := p.Items[0]
	assert.IsType(t, &SegmentItem{}, item)
	si := item.(*SegmentItem)
	assert.Equal(t, 4.12, si.Duration)
	assert.NotNil(t, si.ByteRange)
	assertNotNilEqual(t, 9400, si.ByteRange.Length)
	assertNotNilEqual(t, 376, si.ByteRange.Start)
	assert.Equal(t, "segment1.ts", si.Segment)

	item = p.Items[1]
	assert.IsType(t, &SegmentItem{}, item)
	si = item.(*SegmentItem)
	assert.NotNil(t, si.ByteRange)
	assertNotNilEqual(t, 7144, si.ByteRange.Length)
	assert.Nil(t, si.ByteRange.Start)
}

func TestReader_PlaylistWithComments(t *testing.T) {
	p, err := ReadFile("fixtures/playlistWithComments.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())

	assert.False(t, p.IsMaster())
	assertNotNilEqual(t, 4, p.Version)
	assert.Equal(t, 1, p.Sequence)
	assertNotNilEqual(t, false, p.Cache)
	assert.Equal(t, 12, p.Target)
	assertNotNilEqual(t, "VOD", p.Type)

	item := p.Items[0]
	assert.IsType(t, &SegmentItem{}, item)
	si := item.(*SegmentItem)

	assert.Equal(t, 11.344644, si.Duration)
	assertNotNilEqual(t, "anything", si.Comment)

	item = p.Items[1]
	assert.IsType(t, &DiscontinuityItem{}, item)

	assert.Equal(t, 139, p.ItemSize())
}

func TestReader_VariantAudio(t *testing.T) {
	p, err := ReadFile("fixtures/variantAudio.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.True(t, p.IsMaster())
	assert.Equal(t, 10, p.ItemSize())

	item := p.Items[0]
	assert.IsType(t, &MediaItem{}, item)
	mi := item.(*MediaItem)

	assert.Equal(t, "AUDIO", mi.Type)
	assert.Equal(t, "audio-lo", mi.GroupID)
	assert.Equal(t, "English", mi.Name)
	assertNotNilEqual(t, "eng", mi.Language)
	assertNotNilEqual(t, "spoken", mi.AssocLanguage)
	assertNotNilEqual(t, true, mi.AutoSelect)
	assertNotNilEqual(t, true, mi.Default)
	assertNotNilEqual(t, "englo/prog_index.m3u8", mi.URI)
	assertNotNilEqual(t, true, mi.Forced)
}

func TestReader_VariantAngles(t *testing.T) {
	p, err := ReadFile("fixtures/variantAngles.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.True(t, p.IsMaster())
	assert.Equal(t, 11, p.ItemSize())

	item := p.Items[1]
	assert.IsType(t, &MediaItem{}, item)
	mi := item.(*MediaItem)

	assert.Equal(t, "VIDEO", mi.Type)
	assert.Equal(t, "200kbs", mi.GroupID)
	assert.Equal(t, "Angle2", mi.Name)
	assert.Nil(t, mi.Language)
	assertNotNilEqual(t, true, mi.AutoSelect)
	assertNotNilEqual(t, false, mi.Default)
	assertNotNilEqual(t, "Angle2/200kbs/prog_index.m3u8", mi.URI)

	item = p.Items[9]
	assert.IsType(t, &PlaylistItem{}, item)
	pi := item.(*PlaylistItem)
	assertNotNilEqual(t, 300001, pi.AverageBandwidth)
	assertNotNilEqual(t, "aac", pi.Audio)
	assertNotNilEqual(t, "200kbs", pi.Video)
	assertNotNilEqual(t, "captions", pi.ClosedCaptions)
	assertNotNilEqual(t, "subs", pi.Subtitles)
}

func TestReader_SessionData(t *testing.T) {
	p, err := ReadFile("fixtures/sessionData.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.Equal(t, 3, p.ItemSize())

	item := p.Items[0]
	assert.IsType(t, &SessionDataItem{}, item)
	sdi := item.(*SessionDataItem)

	assert.Equal(t, "com.example.lyrics", sdi.DataID)
	assertNotNilEqual(t, "lyrics.json", sdi.URI)
}

func TestReader_Encrypted(t *testing.T) {
	p, err := ReadFile("fixtures/encrypted.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.Equal(t, 6, p.ItemSize())

	item := p.Items[0]
	assert.IsType(t, &KeyItem{}, item)
	ki := item.(*KeyItem)

	assert.Equal(t, "AES-128", ki.Encryptable.Method)
	assertNotNilEqual(t, "https://priv.example.com/key.php?r=52", ki.Encryptable.URI)
}

func TestReader_Map(t *testing.T) {
	p, err := ReadFile("fixtures/mapPlaylist.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.Equal(t, 1, p.ItemSize())

	item := p.Items[0]
	assert.IsType(t, &MapItem{}, item)
	mi := item.(*MapItem)

	assert.Equal(t, "frelo/prog_index.m3u8", mi.URI)
	assert.NotNil(t, mi.ByteRange)
	assertNotNilEqual(t, 4500, mi.ByteRange.Length)
	assertNotNilEqual(t, 600, mi.ByteRange.Start)
}

func TestReader_Timestamp(t *testing.T) {
	p, err := ReadFile("fixtures/timestampPlaylist.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.Equal(t, 6, p.ItemSize())

	item := p.Items[0]
	assert.IsType(t, &SegmentItem{}, item)
	si := item.(*SegmentItem)

	assert.NotNil(t, si.ProgramDateTime)
	assert.Equal(t, "2016-04-11T15:24:31Z", FormatTime(si.ProgramDateTime.Time))
}

func TestReader_DateRange(t *testing.T) {
	p, err := ReadFile("fixtures/dateRangeScte35.m3u8")
	assert.Nil(t, err)
	assert.True(t, p.IsValid())
	assert.Equal(t, 5, p.ItemSize())

	dateRange := &DateRangeItem{}
	segment := &SegmentItem{}
	assert.IsType(t, dateRange, p.Items[0])
	assert.IsType(t, dateRange, p.Items[4])
	assert.IsType(t, segment, p.Items[1])
	assert.IsType(t, segment, p.Items[2])
	assert.IsType(t, segment, p.Items[3])
}

func TestReader_Invalid(t *testing.T) {
	_, err := ReadFile("path/to/file")
	assert.NotNil(t, err)

	validPlaylist := `#EXT-X-I-FRAME-STREAM-INF:CODECS="avc",BANDWIDTH=540,PROGRAM-ID=1,RESOLUTION=1920x1080,FRAME-RATE=23.976,AVERAGE-BANDWIDTH=550,AUDIO="test",VIDEO="test2",STABLE-VARIANT-ID="1234"SUBTITLES="subs",CLOSED-CAPTIONS="caps",URI="test.url",NAME="1080p",HDCP-LEVEL=TYPE-0`
	testCases := []struct {
		manifest string
		err      error
	}{
		{
			manifest: strings.Join([]string{HeaderTag, HeaderTag}, "\n"),
			err:      ErrPlaylistInvalid,
		},
		{
			manifest: strings.Join([]string{VersionTag + ":" + "7", HeaderTag}, "\n"),
			err:      ErrPlaylistInvalid,
		},
		{
			manifest: strings.Join([]string{
				HeaderTag,
				VersionTag + ":" + "7",
				validPlaylist,
				"#EXTINF:10.991,",
				"test.ts",
			}, "\n"),
			err: ErrPlaylistInvalidType,
		},
	}

	for _, tc := range testCases {
		pl, err := ReadString(tc.manifest)
		require.Nil(t, pl)
		require.Equal(t, tc.err, err)
	}
}

func TestReader_InvalidItems(t *testing.T) {
	s := strings.Join([]string{
		HeaderTag,
		VersionTag + ":" + "x",
		PlaybackStartTag,
	}, "\n")

	pl, err := ReadString(s)
	require.NoError(t, err)
	log.Println(pl)
	require.Len(t, pl.Items, 2)
	require.IsType(t, &UnknownItem{}, pl.Items[0])
	require.IsType(t, &UnknownItem{}, pl.Items[1])
	require.Equal(t, VersionTag+":"+"x", pl.Items[0].String())
	require.Equal(t, PlaybackStartTag, pl.Items[1].String())
}
