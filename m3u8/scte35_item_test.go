package m3u8

import (
	"strings"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

func TestNewSCTE35Item(t *testing.T) {
	attributesKV := []string{
		"TYPE=0x50",
		"CUE=\"/DA9AAAAAAAAAP/wBQb+AAAAAAAnAiVDVUVJAAAAAX+/DhZQQ0tfUENLX1ZPRF85MDAwMDgyOTk3AQAA7sN72Q==\"",
		"DURATION=32400.0",
		"ELAPSED=32400.0",
		"ID=\"QAQ\"",
		"TIME=1448928000.000",
		"UPID=\"0x0e:0x50434b5f50434b5f564f445f39303030303832393937\"",
		"BLACKOUT=MAYBE",
		"CUE-OUT=YES",
		"CUE-IN=YES",
		"SEGNE=\"3:3\"",
		"RANDOM-ATTRIBUTE=\"RANDOM-VALUE\"",
		"RANDOM-ATTRIBUTE2=VALUE",
	}
	line := "#EXT-X-SCTE35:" + strings.Join(attributesKV, ",")

	item, err := NewSCTE35Item(line)
	require.NoError(t, err)
	require.Len(t, item.attributes, 2)

	require.Equal(t, pointer.ToInt(0x50), item.Type)
	require.Equal(t, "/DA9AAAAAAAAAP/wBQb+AAAAAAAnAiVDVUVJAAAAAX+/DhZQQ0tfUENLX1ZPRF85MDAwMDgyOTk3AQAA7sN72Q==", item.Cue)
	require.Equal(t, pointer.ToFloat64(32400.0), item.Duration)
	require.Equal(t, pointer.ToFloat64(32400.0), item.Elapsed)
	require.Equal(t, pointer.ToString("QAQ"), item.ID)
	require.Equal(t, pointer.ToFloat64(1448928000.000), item.Time)
	require.Equal(t, pointer.ToString("0x0e:0x50434b5f50434b5f564f445f39303030303832393937"), item.UPID)
	require.Equal(t, pointer.ToString("MAYBE"), item.Blackout)
	require.Equal(t, pointer.ToString("YES"), item.CueOut)
	require.Equal(t, pointer.ToString("YES"), item.CueIn)
	require.Equal(t, pointer.ToString("3:3"), item.Segne)
	require.Equal(t, "\"RANDOM-VALUE\"", item.attributes["RANDOM-ATTRIBUTE"])

	itemEncoded := item.String()
	for _, attributeKV := range attributesKV {
		require.True(t, strings.Contains(itemEncoded, attributeKV), attributeKV, itemEncoded)
	}
}
