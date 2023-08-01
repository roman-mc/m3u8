package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAttributes(t *testing.T) {
	line := "TEST-ID=\"Help\",URI=\"http://test\",ID=33\n"
	mapAttr := ParseAttributes(line)

	assert.NotNil(t, mapAttr)
	assert.Equal(t, "\"Help\"", mapAttr["TEST-ID"])
	assert.Equal(t, "\"http://test\"", mapAttr["URI"])
	assert.Equal(t, "33", mapAttr["ID"])
}

func TestSanitizeAttributeValue(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "123",
			output: "123",
		},
		{
			input:  "\"123\"",
			output: "123",
		},
		{
			input:  "\"123",
			output: "123",
		},
		{
			input:  "123\"",
			output: "123",
		},
		{
			input:  "'123'",
			output: "'123'",
		},
	}
	for _, tc := range testCases {
		got := SanitizeAttributeValue(tc.input)
		assert.Equal(t, tc.output, got)

	}
}

func TestParseTagName(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "#EXT-DEFINE:NAME=\"name\",VALUE=123",
			output: "#EXT-DEFINE",
		},
		{
			input:  "#EXTM3U",
			output: "#EXTM3U",
		},
		{
			input:  "#EXTINF:NAME=\"name\",VALUE=123",
			output: "#EXTINF",
		},
		{
			input:  " #EXTINF:NAME=\"name\",VALUE=123",
			output: "",
		},
		{
			input:  "#EXXTINF:NAME=\"name\",VALUE=123",
			output: "",
		},
		{
			input:  "EXTINF:NAME=\"name\",VALUE=123",
			output: "",
		},
	}

	for _, tc := range testCases {
		got := ParseTagName(tc.input)

		assert.Equal(t, tc.output, got)
	}
}

func TestParseFloat(t *testing.T) {
	attributes := map[string]string{
		"float":   "123.3",
		"int":     "123",
		"quotes":  "\"123.3\"",
		"invalid": "1pp2",
	}
	valFloat := 123.3
	valInt := float64(123)

	testCases := []struct {
		key      string
		expect   *float64
		hasError bool
	}{
		{
			key:    "float",
			expect: &valFloat,
		},
		{
			key:    "int",
			expect: &valInt,
		},
		{
			key:    "quotes",
			expect: &valFloat,
		},
		{
			key:      "invalid",
			expect:   nil,
			hasError: true,
		},
		{
			key:      "non-existing",
			expect:   nil,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		result, err := ParseFloat(attributes, tc.key)
		if tc.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, tc.expect, result)
	}
}

func TestParseInt(t *testing.T) {
	attributes := map[string]string{
		"int":     "123",
		"quotes":  "\"123\"",
		"invalid": "1pp2",
	}
	valInt := 123

	testCases := []struct {
		key      string
		expect   *int
		hasError bool
	}{
		{
			key:    "int",
			expect: &valInt,
		},
		{
			key:    "quotes",
			expect: &valInt,
		},
		{
			key:      "invalid",
			expect:   nil,
			hasError: true,
		},
		{
			key:      "non-existing",
			expect:   nil,
			hasError: false,
		},
	}

	for _, tc := range testCases {
		result, err := ParseInt(attributes, tc.key)
		if tc.hasError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		require.Equal(t, tc.expect, result)
	}
}

func TestParseYesNo(t *testing.T) {
	attributes := map[string]string{
		"true":          "YES",
		"quotes":        "\"YES\"",
		"false-invalid": "1pp2",
		"false-valid":   "NO",
	}
	valTrue := true
	valFalse := false

	testCases := []struct {
		key string
		val *bool
	}{
		{
			key: "true",
			val: &valTrue,
		},
		{
			key: "quotes",
			val: &valTrue,
		},
		{
			key: "false-invalid",
			val: &valFalse,
		},
		{
			key: "false-valid",
			val: &valFalse,
		},
		{
			key: "non-existing",
			val: nil,
		},
	}

	for _, tc := range testCases {
		result := ParseYesNo(attributes, tc.key)
		require.Equal(t, tc.val, result)
	}

}

func TestFormatYesNo(t *testing.T) {
	yes := FormatYesNo(true)
	no := FormatYesNo(false)
	require.Equal(t, "YES", yes)
	require.Equal(t, "NO", no)
}

func TestPointerTo(t *testing.T) {
	attributes := map[string]string{
		"string":       "YES",
		"stringQuotes": "\"YES\"",
	}
	val := "YES"

	testCases := []struct {
		key string
		val *string
	}{
		{
			key: "string",
			val: &val,
		},
		{
			key: "stringQuotes",
			val: &val,
		},
		{
			key: "non-existing",
			val: nil,
		},
	}

	for _, tc := range testCases {
		result := PointerTo(attributes, tc.key)
		require.Equal(t, tc.val, result)
	}

}

func TestAttributeExists(t *testing.T) {
	attributes := map[string]string{
		"exists": "1",
	}

	exists := AttributeExists("exists", attributes)
	notExists := AttributeExists("random", attributes)
	require.True(t, exists)
	require.False(t, notExists)
}

func TestPointerToFloat(t *testing.T) {
	attributes := map[string]string{
		"exists": "1.1",
	}

	val := 1.1
	exists := PointerToFloat(attributes, "exists")
	notExists := PointerToFloat(attributes, "notExists")
	require.Equal(t, &val, exists)
	require.Equal(t, (*float64)(nil), notExists)
}

func TestPointerToInt(t *testing.T) {
	attributes := map[string]string{
		"exists": "1111",
	}

	val := 1111
	exists := PointerToInt(attributes, "exists")
	notExists := PointerToInt(attributes, "notExists")
	require.Equal(t, &val, exists)
	require.Equal(t, (*int)(nil), notExists)
}

func TestParseResolution(t *testing.T) {
	attributes := map[string]string{
		"exists":        "33x11",
		"exists-quotes": "\"33x11\"",
		"left-wrong":    "vvx11",
		"right-wrong":   "33xvv",
		"both-wrong":    "ppxvv",
		"left-exists":   "11x",
		"right-exists":  "x11",
		"none-exists":   "x",
		"none-exists2":  "111",
	}

	expected := &Resolution{
		Width:  33,
		Height: 11,
	}

	testCases := []struct {
		key      string
		isErr    bool
		expected *Resolution
		toString string
	}{
		{
			key:      "exists",
			isErr:    false,
			expected: expected,
			toString: "33x11",
		},
		{
			key:      "exists-quotes",
			isErr:    false,
			expected: expected,
			toString: "33x11",
		},
		{
			key:      "left-wrong",
			isErr:    true,
			expected: nil,
		},
		{
			key:      "right-wrong",
			isErr:    true,
			expected: nil,
		},
		{
			key:      "both-wrong",
			isErr:    true,
			expected: nil,
		},
		{
			key:      "left-exists",
			isErr:    true,
			expected: nil,
		},
		{
			key:      "right-exists",
			isErr:    true,
			expected: nil,
		},
		{
			key:      "none-exists",
			isErr:    true,
			expected: nil,
		},
		{
			key:      "none-exists2",
			isErr:    true,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			got, err := ParseResolution(attributes, tc.key)
			if tc.isErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tc.expected, got)
			require.Equal(t, tc.toString, got.String())
		})
	}
}
