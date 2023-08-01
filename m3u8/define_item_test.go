package m3u8

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDefineItem(t *testing.T) {
	value := "\"val1\""
	imp := "\"imp1\""
	queryParam := "\"queryParam1\""
	encodedName := fmt.Sprintf("%s=%s", AttributeName, "\"name\"")
	encodedValue := fmt.Sprintf("%s=%s", AttributeValue, value)
	encodedImp := fmt.Sprintf("%s=%s", AttributeImport, imp)
	encodedQueryParam := fmt.Sprintf("%s=%s", AttributeQueryParam, queryParam)

	expectedValue := "val1"
	expectedImp := "imp1"
	expectedQueryParam := "queryParam1"

	testCases := []struct {
		line                      string
		expectDecoded             *DefineItem
		expectedAttributesEncoded []string
	}{
		{
			line: fmt.Sprintf("%s:%s,%s,%s,%s",
				DefineTag,
				encodedName,
				encodedValue,
				encodedImp,
				encodedQueryParam,
			),
			expectDecoded: &DefineItem{
				Name:       "name",
				Value:      &expectedValue,
				Import:     &expectedImp,
				QueryParam: &expectedQueryParam,
				attributes: make(map[string]string),
			},
			expectedAttributesEncoded: []string{
				encodedName,
				encodedValue,
				encodedImp,
				encodedQueryParam,
			},
		},
		{
			line: fmt.Sprintf("%s:%s,%s,%s",
				DefineTag,
				encodedName,
				encodedValue,
				encodedImp,
			),
			expectDecoded: &DefineItem{
				Name:       "name",
				Value:      &expectedValue,
				Import:     &expectedImp,
				attributes: make(map[string]string),
			},
			expectedAttributesEncoded: []string{
				encodedName,
				encodedValue,
				encodedImp,
			},
		},
		{
			line: fmt.Sprintf("%s:%s,%s",
				DefineTag,
				encodedName,
				encodedValue,
			),
			expectDecoded: &DefineItem{
				Name:       "name",
				Value:      &expectedValue,
				attributes: make(map[string]string),
			},
			expectedAttributesEncoded: []string{
				encodedName,
				encodedValue,
			},
		},
		{
			line: fmt.Sprintf("%s:%s",
				DefineTag,
				encodedName,
			),
			expectDecoded: &DefineItem{
				Name:       "name",
				attributes: make(map[string]string),
			},
			expectedAttributesEncoded: []string{
				encodedName,
			},
		},
	}

	for _, tc := range testCases {
		defineTag := NewDefineItem(tc.line)
		require.Equal(t, tc.expectDecoded, defineTag)

		encodedTag := defineTag.String()

		for _, attributeKV := range tc.expectedAttributesEncoded {
			require.Contains(t, encodedTag, attributeKV)
		}
	}
}
