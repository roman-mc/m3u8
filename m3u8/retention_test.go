// here are tests checking retention (that we don't lose) of all tags and attributes after decoding and encoding back
package m3u8

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

func TestRetentionProperDecoding(t *testing.T) {
	testCases := []struct {
		filePath string
	}{
		{filePath: "fixtures/peacock_master1.m3u8"},
		{filePath: "fixtures/peacock_playlist1.m3u8"},
		{filePath: "fixtures/fer_drm.m3u8"},
		{filePath: "fixtures/fer_with_ads.m3u8"},
		{filePath: "fixtures/fer_without_ads.m3u8"},
		{filePath: "fixtures/live_channels_master.m3u8"},
		{filePath: "fixtures/live_channels_playlist.m3u8"},
		{filePath: "fixtures/sle_drm.m3u8"},
		{filePath: "fixtures/sle_non_drm.m3u8"},
		{filePath: "fixtures/vod_drm.m3u8"},
		{filePath: "fixtures/vod_long_form.m3u8"},
		{filePath: "fixtures/vod_short_form.m3u8"},
		{filePath: "fixtures/vod_non_drm.m3u8"},
		{filePath: "fixtures/vod_non_drm_compact.m3u8"},
		{filePath: "fixtures/vod_pre_roll.m3u8"},
		{filePath: "fixtures/vod_short_form.m3u8"},
	}

	for _, tc := range testCases {
		f, err := os.ReadFile(tc.filePath)
		require.NoError(t, err)

		p, err := Read(bytes.NewReader(f))
		require.NoError(t, err)

		require.True(t, p.IsValid())
		encoded := p.String()
		decodedPlaylist, err := ReadString(encoded)
		require.NoError(t, err)

		require.Equal(t, p.DiscontinuitySequence, p.DiscontinuitySequence)
		require.Equal(t, p.Cache, decodedPlaylist.Cache)
		require.Equal(t, p.Type, decodedPlaylist.Type)
		require.Equal(t, p.Sequence, decodedPlaylist.Sequence)
		require.Equal(t, p.IFramesOnly, decodedPlaylist.IFramesOnly)
		require.Equal(t, p.IndependentSegments, decodedPlaylist.IndependentSegments)
		require.Equal(t, p.Live, decodedPlaylist.Live)
		require.Equal(t, p.Master, decodedPlaylist.Master)
		require.Equal(t, p.Target, decodedPlaylist.Target)
		require.Equal(t, p.Version, decodedPlaylist.Version)
		require.Equal(t, len(p.Items), len(decodedPlaylist.Items))

		for i := 0; i < len(p.Items); i++ {
			require.Equal(t, p.Items[i], decodedPlaylist.Items[i])
		}
	}
}

func TestRetentionAttributesForTags(t *testing.T) {
	randomAttributesString := "UNKNOWN=123,UNKNOWN2=\"123\",UNKNOWN3=11x33,UNKNOWN4=1.33,"
	expectedAttributesMap := map[string]string{
		"UNKNOWN":  "123",
		"UNKNOWN2": "\"123\"",
		"UNKNOWN3": "11x33",
		"UNKNOWN4": "1.33",
	}

	testCases := []Tag{
		NewSessionKeyItem(SessionKeyItemTag + ":" + randomAttributesString + "METHOD=1"),
		NewKeyItem(KeyItemTag + ":" + randomAttributesString + "METHOD=1"),
		NewDateRangeItem(DateRangeItemTag + ":" + randomAttributesString + "ID=1"),
		NewMapItem(MapItemTag + ":" + randomAttributesString + "URI=1"),
		NewSessionDataItem(SessionDataItemTag + ":" + randomAttributesString + "DATA-ID=1"),
		mustTag(NewPlaybackStart(PlaybackStartTag + ":" + randomAttributesString + "TIME-OFFSET=1")),
		NewMediaItem(MediaItemTag + ":" + randomAttributesString + "TYPE=\"123\",GROUP-ID=\"123\",NAME=\"123\""),
		NewPlaylistItem(
			PlaylistIframeTag+":"+
				randomAttributesString+
				"RESOLUTION=1x1,AVERAGE-BANDWIDTH=2,FRAME-RATE=12,BANDWIDTH=2,URI=123", true),
		NewDefineItem(DefineTag + ":" + randomAttributesString + "NAME=\"123\""),
	}

	for _, tc := range testCases {
		tagString := tc.String()
		require.Nil(t, tc.Validate())
		attributes := parser.ParseAttributes(tagString)

		for expectedKey, expectedValue := range expectedAttributesMap {
			require.Equal(t, expectedValue, attributes[expectedKey], tagString)
		}
	}

	streamInfTag := NewPlaylistItem(
		PlaylistItemTag+":"+
			randomAttributesString+
			"RESOLUTION=1x1,AVERAGE-BANDWIDTH=2,FRAME-RATE=12,BANDWIDTH=2,URI=123", false)
	streamInfTagString := streamInfTag.String()
	tagAndURI := strings.Split(streamInfTagString, "\n")
	require.Len(t, tagAndURI, 2)

	streamInfTagStringAttributes := parser.ParseAttributes(tagAndURI[0])
	for expectedKey, expectedValue := range expectedAttributesMap {
		require.Equal(t, expectedValue, streamInfTagStringAttributes[expectedKey], streamInfTagStringAttributes)
	}

}

func mustTag(tag Tag, err error) Tag {
	if err != nil {
		panic(err)
	}
	return tag
}
