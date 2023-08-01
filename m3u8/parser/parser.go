package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

const (
	QuotedFormatString    = `%s="%v"`
	FormatString          = `%s=%v`
	FormatHexString       = `%s=%x`
	FrameRateFormatString = `%s=%.3f`
	// Values

	NoneValue = "NONE"
	YesValue  = "YES"
	NoValue   = "NO"
)

var (
	parseRegex        = regexp.MustCompile(`([A-z0-9-]+)\s*=\s*("[^"]*"|[^,]*)`)
	ParseTagNameRegex = regexp.MustCompile(`^#EXT[A-Z-0-9]+:?`)
)

var (
	// ErrBandwidthMissing represents error when a segment does not have bandwidth
	ErrBandwidthMissing = errors.New("missing bandwidth")

	// ErrBandwidthInvalid represents error when a bandwidth is invalid
	ErrBandwidthInvalid = errors.New("invalid bandwidth")
)

// ParseAttributes parses a text line in playlist and returns an attributes map
func ParseAttributes(text string) map[string]string {
	res := make(map[string]string)
	value := strings.Replace(text, "\n", "", -1)
	matches := parseRegex.FindAllStringSubmatch(value, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			key := match[1]
			// don't remove quotes, so we preserve attribute values in their original state for unknown attributes,
			// but we remove them while extracting in tags (knowing how to decode them back)
			attributeValue := match[2]

			res[key] = attributeValue
		}
	}

	return res
}

func ParseTagName(text string) string {
	str := ParseTagNameRegex.FindString(text)
	if str == "" {
		return ""
	}

	if str[len(str)-1] == ':' {
		return str[:len(str)-1]
	}

	return str
}

func ParseFloat(attributes map[string]string, key string) (*float64, error) {
	stringValue, ok := attributes[key]
	if !ok {
		return nil, nil
	}
	stringValue = SanitizeAttributeValue(stringValue)

	value, err := strconv.ParseFloat(stringValue, 64)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

func ParseInt(attributes map[string]string, key string) (*int, error) {
	stringValue, ok := attributes[key]
	if !ok {
		return nil, nil
	}
	stringValue = SanitizeAttributeValue(stringValue)

	int64Value, err := strconv.ParseInt(stringValue, 0, 0)
	if err != nil {
		return nil, err
	}

	value := int(int64Value)

	return &value, nil
}

func ParseYesNo(attributes map[string]string, key string) *bool {
	stringValue, ok := attributes[key]

	if !ok {
		return nil
	}
	stringValue = SanitizeAttributeValue(stringValue)

	val := false

	if stringValue == YesValue {
		val = true
	}

	return &val
}

func FormatYesNo(value bool) string {
	if value {
		return YesValue
	}

	return NoValue
}

func ParseBandwidth(attributes map[string]string, key string) (int, error) {
	bw, ok := attributes[key]
	if !ok {
		return 0, ErrBandwidthMissing
	}
	bw = SanitizeAttributeValue(bw)

	bandwidth, err := strconv.ParseInt(bw, 0, 0)
	if err != nil {
		return 0, ErrBandwidthInvalid
	}

	return int(bandwidth), nil
}

func ParseResolution(attributes map[string]string, key string) (*Resolution, error) {
	resolution, ok := attributes[key]
	if !ok {
		return nil, nil
	}
	resolution = SanitizeAttributeValue(resolution)

	return NewResolution(resolution)
}

func AttributeExists(key string, attributes map[string]string) bool {
	_, ok := attributes[key]
	return ok
}

func PointerTo(attributes map[string]string, key string) *string {
	value, ok := attributes[key]

	if !ok {
		return nil
	}

	value = SanitizeAttributeValue(value)
	return &value
}

func PointerToFloat(attributes map[string]string, key string) *float64 {
	result, _ := ParseFloat(attributes, key)
	return result
}

func PointerToInt(attributes map[string]string, key string) *int {
	result, _ := ParseInt(attributes, key)
	return result
}

// SanitizeAttributeValue sanitizes attribute value
//
//	We cannot sanitize all values during parsing attributes (which seems handy),
//	because some (e.g. unknown) attributes should preserve their original state, which is hard to handle at the moment
func SanitizeAttributeValue(s string) string {
	return strings.Replace(s, `"`, "", -1)
}
