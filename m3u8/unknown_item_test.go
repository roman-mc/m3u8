package m3u8

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnknownItem_String(t *testing.T) {
	testCases := []struct {
		line    string
		tagName string
	}{
		{
			line:    "RANDOMSTRING-NOT-VALID-BY-SPECIFICATION",
			tagName: "",
		},
		{
			line:    "#EXT-RANDOMSTRING-VALID-BY-SPECIFICATION",
			tagName: "#EXT-RANDOMSTRING-VALID-BY-SPECIFICATION",
		},
	}

	for _, tc := range testCases {
		item := NewUnknownItem(tc.line, nil)

		encodedItem := item.String()
		require.Equal(t, tc.line, encodedItem)
		require.Equal(t, tc.tagName, item.GetTagName())
	}
}
