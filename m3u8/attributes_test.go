package m3u8

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttributesJoinMap(t *testing.T) {
	attrs := []string{"attr1=v1", "attr2=v2"}
	attributesMap := map[string]string{
		"k1": "v1",
		"k2": "\"v1\"",
		"k3": "v3",
		"k4": "\"v4\"",
	}

	expectedAttrsMap := map[string]string{
		"k1":    "v1",
		"k2":    "\"v1\"",
		"k3":    "v3",
		"k4":    "\"v4\"",
		"attr1": "v1",
		"attr2": "v2",
	}
	gotAttrsMap := make(map[string]string, len(expectedAttrsMap))

	attrs = attributesJoinMap(attrs, attributesMap)
	for _, attribute := range attrs {
		kv := strings.Split(attribute, "=")
		require.Len(t, kv, 2)

		gotAttrsMap[kv[0]] = kv[1]
	}

	require.Equal(t, expectedAttrsMap, gotAttrsMap)
}
