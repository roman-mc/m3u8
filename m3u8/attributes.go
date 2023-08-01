package m3u8

import (
	"fmt"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

// deleteKeys used on attributes map to delete all known attributes to leave only unknown
func deleteKeys(m map[string]string, keys ...string) {
	for _, key := range keys {
		delete(m, key)
	}
}

// attributesJoinMap utility function to unfold map's values into a slice
func attributesJoinMap(attributes []string, attributesMap map[string]string) []string {
	for attributeKey, attribute := range attributesMap {
		attributes = append(attributes, fmt.Sprintf(parser.FormatString, attributeKey, attribute))
	}

	return attributes
}
