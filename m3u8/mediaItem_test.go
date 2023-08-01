package m3u8

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMediaItem_Parse(t *testing.T) {
	line := `#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio-lo",LANGUAGE="fre",
ASSOC-LANGUAGE="spoken",NAME="Francais",AUTOSELECT=YES,
INSTREAM-ID="SERVICE3",CHARACTERISTICS="public.html",
CHANNELS="6",
"DEFAULT=NO,URI="frelo/prog_index.m3u8",STABLE-RENDITION-ID="1234",FORCED=YES
"`

	mi := NewMediaItem(line)
	assert.Equal(t, "AUDIO", mi.Type)
	assert.Equal(t, "audio-lo", mi.GroupID)
	assert.Equal(t, "Francais", mi.Name)

	assertNotNilEqual(t, "1234", mi.StableRenditionId)
	assertNotNilEqual(t, "fre", mi.Language)
	assertNotNilEqual(t, "spoken", mi.AssocLanguage)
	assertNotNilEqual(t, true, mi.AutoSelect)
	assertNotNilEqual(t, false, mi.Default)
	assertNotNilEqual(t, "frelo/prog_index.m3u8", mi.URI)
	assertNotNilEqual(t, true, mi.Forced)
	assertNotNilEqual(t, "SERVICE3", mi.InStreamID)
	assertNotNilEqual(t, "public.html", mi.Characteristics)
	assertNotNilEqual(t, "6", mi.Channels)

	expected := "#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID=\"audio-lo\",LANGUAGE=\"fre\",ASSOC-LANGUAGE=\"spoken\",NAME=\"Francais\",AUTOSELECT=YES,DEFAULT=NO,URI=\"frelo/prog_index.m3u8\",FORCED=YES,INSTREAM-ID=\"SERVICE3\",CHARACTERISTICS=\"public.html\",CHANNELS=\"6\",STABLE-RENDITION-ID=\"1234\""
	assertToString(t, expected, mi)
}
