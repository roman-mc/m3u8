package m3u8

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/NBCUDTC/midnight-hls-go-parser-src/m3u8/parser"
)

type state struct {
	open        bool
	currentItem Item
	master      bool
}

// ReadString parses a text string and returns a playlist
func ReadString(text string) (*Playlist, error) {
	return Read(strings.NewReader(text))
}

// ReadFile reads text from a file and returns a playlist
func ReadFile(path string) (*Playlist, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Read(bytes.NewReader(f))
}

// Read reads text from an io.Reader and returns a playlist
func Read(reader io.Reader) (*Playlist, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	pl := NewPlaylist()
	st := &state{}
	header := true
	eof := false

	for !eof {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, err
			}

			eof = true
		}

		value := strings.TrimSpace(line)
		if len(value) == 0 {
			continue
		}
		if header && value != HeaderTag {
			return nil, ErrPlaylistInvalid
		}

		if value == HeaderTag {
			if header == false {
				return nil, ErrPlaylistInvalid
			}
			header = false
			continue
		}

		if err := parseLine(value, pl, st); err != nil {
			return nil, err
		}
	}

	if pl.Version == nil {
		version := 1
		pl.Version = &version
	}

	if !pl.IsValid() {
		return nil, ErrPlaylistInvalidType

	}

	return pl, nil
}

// parseLine parses all tags and attributes (implemented by this lib)
func parseLine(line string, pl *Playlist, st *state) error {
	lineIsParsed := false

	for tag, tagValue := range tagsMap {
		if !matchTag(line, tag) {
			continue
		}
		lineIsParsed = true

		err := tagValue.ReadLine(line, pl, st)
		if err != nil {
			pl.Items = append(pl.Items, NewUnknownItem(line, nil))
		}
		break
	}

	if lineIsParsed {
		return nil
	}
	if st.currentItem != nil && st.open {
		return parseNextLine(line, pl, st)
	}

	pl.Items = append(pl.Items, NewUnknownItem(line, nil))

	return nil
}

func parseNextLine(line string, pl *Playlist, st *state) error {
	value := strings.Replace(line, "\n", "", -1)
	value = strings.Replace(value, "\r", "", -1)
	if st.master {
		// PlaylistItem
		it, ok := st.currentItem.(*PlaylistItem)
		if !ok {
			return parseError(line, ErrPlaylistItemInvalid)
		}
		it.URI = value
		pl.Items = append(pl.Items, it)
	} else {
		// SegmentItem
		it, ok := st.currentItem.(*SegmentItem)
		if !ok {
			return parseError(line, ErrSegmentItemInvalid)
		}
		it.Segment = value
		pl.Items = append(pl.Items, it)
	}

	st.open = false

	return nil
}

func matchTag(line, tag string) bool {
	return strings.HasPrefix(line, tag) && !strings.HasPrefix(line, tag+"-")
}

func parseIntValue(line string, tag string) (int, error) {
	var v int
	_, err := fmt.Sscanf(line, tag+":%d", &v)
	return v, err
}

func parseIntPtr(line string, tag string) (*int, error) {
	var ptr int
	_, err := fmt.Sscanf(line, tag+":%d", &ptr)
	return &ptr, err
}

func parseStringPtr(line string, tag string) *string {
	value := strings.Replace(line, tag+":", "", -1)
	if value == "" {
		return nil
	}
	return &value
}

func ParseYesNoPtr(line string, tag string) *bool {
	value := strings.Replace(line, tag+":", "", -1)
	var b bool
	if value == parser.YesValue {
		b = true
	} else {
		b = false
	}

	return &b
}

func parseError(line string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("error: %v when parsing playlist error for line: %s", err, line)
}
